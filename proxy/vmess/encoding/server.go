package encoding

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"hash/fnv"
	"io"
	"sync"
	"time"

	"golang.org/x/crypto/chacha20poly1305"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/crypto"
	"v2ray.com/core/common/errors"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/protocol"
	"v2ray.com/core/common/serial"
	"v2ray.com/core/proxy/vmess"
)

type sessionId struct {
	user  [16]byte
	key   [16]byte
	nonce [16]byte
}

type sessionHistory struct {
	sync.RWMutex
	cache map[sessionId]time.Time
}

func newSessionHistory() *sessionHistory {
	h := &sessionHistory{
		cache: make(map[sessionId]time.Time, 128),
	}
	go h.run()
	return h
}

func (h *sessionHistory) Add(session sessionId) {
	h.Lock()
	h.cache[session] = time.Now().Add(time.Minute * 3)
	h.Unlock()
}

func (h *sessionHistory) Has(session sessionId) bool {
	h.RLock()
	defer h.RUnlock()

	if expire, found := h.cache[session]; found {
		return expire.After(time.Now())
	}
	return false
}

func (h *sessionHistory) run() {
	for {
		time.Sleep(time.Second * 30)
		session2Remove := make([]sessionId, 0, 16)
		now := time.Now()
		h.Lock()
		for session, expire := range h.cache {
			if expire.Before(now) {
				session2Remove = append(session2Remove, session)
			}
		}
		for _, session := range session2Remove {
			delete(h.cache, session)
		}
		h.Unlock()
	}
}

var (
	globalSessionHistory = newSessionHistory()
)

type ServerSession struct {
	userValidator   protocol.UserValidator
	requestBodyKey  []byte
	requestBodyIV   []byte
	responseBodyKey []byte
	responseBodyIV  []byte
	responseHeader  byte
	responseWriter  io.Writer
}

// NewServerSession creates a new ServerSession, using the given UserValidator.
// The ServerSession instance doesn't take ownership of the validator.
func NewServerSession(validator protocol.UserValidator) *ServerSession {
	return &ServerSession{
		userValidator: validator,
	}
}

func (v *ServerSession) DecodeRequestHeader(reader io.Reader) (*protocol.RequestHeader, error) {
	buffer := make([]byte, 512)

	_, err := io.ReadFull(reader, buffer[:protocol.IDBytesLen])
	if err != nil {
		return nil, errors.Base(err).Message("VMess|Server: Failed to read request header.")
	}

	user, timestamp, valid := v.userValidator.Get(buffer[:protocol.IDBytesLen])
	if !valid {
		return nil, errors.New("VMess|Server: Invalid user.")
	}

	timestampHash := md5.New()
	timestampHash.Write(hashTimestamp(timestamp))
	iv := timestampHash.Sum(nil)
	account, err := user.GetTypedAccount()
	if err != nil {
		return nil, errors.Base(err).Message("VMess|Server: Failed to get user account.")
	}
	vmessAccount := account.(*vmess.InternalAccount)

	aesStream := crypto.NewAesDecryptionStream(vmessAccount.ID.CmdKey(), iv)
	decryptor := crypto.NewCryptionReader(aesStream, reader)

	nBytes, err := io.ReadFull(decryptor, buffer[:41])
	if err != nil {
		return nil, errors.Base(err).Message("VMess|Server: Failed to read request header.")
	}
	bufferLen := nBytes

	request := &protocol.RequestHeader{
		User:    user,
		Version: buffer[0],
	}

	if request.Version != Version {
		return nil, errors.New("VMess|Server: Invalid protocol version ", request.Version)
	}

	v.requestBodyIV = append([]byte(nil), buffer[1:17]...)   // 16 bytes
	v.requestBodyKey = append([]byte(nil), buffer[17:33]...) // 16 bytes
	var sid sessionId
	copy(sid.user[:], vmessAccount.ID.Bytes())
	copy(sid.key[:], v.requestBodyKey)
	copy(sid.nonce[:], v.requestBodyIV)
	if globalSessionHistory.Has(sid) {
		return nil, errors.New("VMess|Server: Duplicated session id. Possibly under reply attack.")
	}
	globalSessionHistory.Add(sid)

	v.responseHeader = buffer[33]                       // 1 byte
	request.Option = protocol.RequestOption(buffer[34]) // 1 byte
	padingLen := int(buffer[35] >> 4)
	request.Security = protocol.NormSecurity(protocol.Security(buffer[35] & 0x0F))
	// 1 bytes reserved
	request.Command = protocol.RequestCommand(buffer[37])

	request.Port = net.PortFromBytes(buffer[38:40])

	switch buffer[40] {
	case AddrTypeIPv4:
		_, err = io.ReadFull(decryptor, buffer[41:45]) // 4 bytes
		bufferLen += 4
		if err != nil {
			return nil, errors.Base(err).Message("VMess|Server: Failed to read IPv4.")
		}
		request.Address = net.IPAddress(buffer[41:45])
	case AddrTypeIPv6:
		_, err = io.ReadFull(decryptor, buffer[41:57]) // 16 bytes
		bufferLen += 16
		if err != nil {
			return nil, errors.Base(err).Message("VMess|Server: Failed to read IPv6 address.")
		}
		request.Address = net.IPAddress(buffer[41:57])
	case AddrTypeDomain:
		_, err = io.ReadFull(decryptor, buffer[41:42])
		if err != nil {
			return nil, errors.Base(err).Message("VMess:Server: Failed to read domain.")
		}
		domainLength := int(buffer[41])
		if domainLength == 0 {
			return nil, errors.New("VMess|Server: Zero length domain.")
		}
		_, err = io.ReadFull(decryptor, buffer[42:42+domainLength])
		if err != nil {
			return nil, errors.Base(err).Message("VMess|Server: Failed to read domain.")
		}
		bufferLen += 1 + domainLength
		request.Address = net.DomainAddress(string(buffer[42 : 42+domainLength]))
	}

	if padingLen > 0 {
		_, err = io.ReadFull(decryptor, buffer[bufferLen:bufferLen+padingLen])
		if err != nil {
			return nil, errors.Base(err).Message("VMess|Server: Failed to read padding.")
		}
		bufferLen += padingLen
	}

	_, err = io.ReadFull(decryptor, buffer[bufferLen:bufferLen+4])
	if err != nil {
		return nil, errors.Base(err).Message("VMess|Server: Failed to read checksum.")
	}

	fnv1a := fnv.New32a()
	fnv1a.Write(buffer[:bufferLen])
	actualHash := fnv1a.Sum32()
	expectedHash := serial.BytesToUint32(buffer[bufferLen : bufferLen+4])

	if actualHash != expectedHash {
		return nil, errors.New("VMess|Server: Invalid auth.")
	}

	if request.Address == nil {
		return nil, errors.New("VMess|Server: Invalid remote address.")
	}

	return request, nil
}

func (v *ServerSession) DecodeRequestBody(request *protocol.RequestHeader, reader io.Reader) buf.Reader {
	var authReader io.Reader
	if request.Security.Is(protocol.SecurityType_NONE) {
		if request.Option.Has(protocol.RequestOptionChunkStream) {
			auth := &crypto.AEADAuthenticator{
				AEAD:                    NoOpAuthenticator{},
				NonceGenerator:          crypto.NoOpBytesGenerator{},
				AdditionalDataGenerator: crypto.NoOpBytesGenerator{},
			}
			authReader = crypto.NewAuthenticationReader(auth, reader)
		} else {
			authReader = reader
		}
	} else if request.Security.Is(protocol.SecurityType_LEGACY) {
		aesStream := crypto.NewAesDecryptionStream(v.requestBodyKey, v.requestBodyIV)
		cryptionReader := crypto.NewCryptionReader(aesStream, reader)
		if request.Option.Has(protocol.RequestOptionChunkStream) {
			auth := &crypto.AEADAuthenticator{
				AEAD:                    new(FnvAuthenticator),
				NonceGenerator:          crypto.NoOpBytesGenerator{},
				AdditionalDataGenerator: crypto.NoOpBytesGenerator{},
			}
			authReader = crypto.NewAuthenticationReader(auth, cryptionReader)
		} else {
			authReader = cryptionReader
		}
	} else if request.Security.Is(protocol.SecurityType_AES128_GCM) {
		block, _ := aes.NewCipher(v.requestBodyKey)
		aead, _ := cipher.NewGCM(block)

		auth := &crypto.AEADAuthenticator{
			AEAD: aead,
			NonceGenerator: &ChunkNonceGenerator{
				Nonce: append([]byte(nil), v.requestBodyIV...),
				Size:  aead.NonceSize(),
			},
			AdditionalDataGenerator: crypto.NoOpBytesGenerator{},
		}
		authReader = crypto.NewAuthenticationReader(auth, reader)
	} else if request.Security.Is(protocol.SecurityType_CHACHA20_POLY1305) {
		aead, _ := chacha20poly1305.New(GenerateChacha20Poly1305Key(v.requestBodyKey))

		auth := &crypto.AEADAuthenticator{
			AEAD: aead,
			NonceGenerator: &ChunkNonceGenerator{
				Nonce: append([]byte(nil), v.requestBodyIV...),
				Size:  aead.NonceSize(),
			},
			AdditionalDataGenerator: crypto.NoOpBytesGenerator{},
		}
		authReader = crypto.NewAuthenticationReader(auth, reader)
	}

	return buf.NewReader(authReader)
}

func (v *ServerSession) EncodeResponseHeader(header *protocol.ResponseHeader, writer io.Writer) {
	responseBodyKey := md5.Sum(v.requestBodyKey)
	responseBodyIV := md5.Sum(v.requestBodyIV)
	v.responseBodyKey = responseBodyKey[:]
	v.responseBodyIV = responseBodyIV[:]

	aesStream := crypto.NewAesEncryptionStream(v.responseBodyKey, v.responseBodyIV)
	encryptionWriter := crypto.NewCryptionWriter(aesStream, writer)
	v.responseWriter = encryptionWriter

	encryptionWriter.Write([]byte{v.responseHeader, byte(header.Option)})
	err := MarshalCommand(header.Command, encryptionWriter)
	if err != nil {
		encryptionWriter.Write([]byte{0x00, 0x00})
	}
}

func (v *ServerSession) EncodeResponseBody(request *protocol.RequestHeader, writer io.Writer) buf.Writer {
	var authWriter io.Writer
	if request.Security.Is(protocol.SecurityType_NONE) {
		if request.Option.Has(protocol.RequestOptionChunkStream) {
			auth := &crypto.AEADAuthenticator{
				AEAD:                    new(FnvAuthenticator),
				NonceGenerator:          crypto.NoOpBytesGenerator{},
				AdditionalDataGenerator: crypto.NoOpBytesGenerator{},
			}
			authWriter = crypto.NewAuthenticationWriter(auth, writer)
		} else {
			authWriter = writer
		}
	} else if request.Security.Is(protocol.SecurityType_LEGACY) {
		if request.Option.Has(protocol.RequestOptionChunkStream) {
			auth := &crypto.AEADAuthenticator{
				AEAD:                    new(FnvAuthenticator),
				NonceGenerator:          crypto.NoOpBytesGenerator{},
				AdditionalDataGenerator: crypto.NoOpBytesGenerator{},
			}
			authWriter = crypto.NewAuthenticationWriter(auth, v.responseWriter)
		} else {
			authWriter = v.responseWriter
		}
	} else if request.Security.Is(protocol.SecurityType_AES128_GCM) {
		block, _ := aes.NewCipher(v.responseBodyKey)
		aead, _ := cipher.NewGCM(block)

		auth := &crypto.AEADAuthenticator{
			AEAD: aead,
			NonceGenerator: &ChunkNonceGenerator{
				Nonce: append([]byte(nil), v.responseBodyIV...),
				Size:  aead.NonceSize(),
			},
			AdditionalDataGenerator: crypto.NoOpBytesGenerator{},
		}
		authWriter = crypto.NewAuthenticationWriter(auth, writer)
	} else if request.Security.Is(protocol.SecurityType_CHACHA20_POLY1305) {
		aead, _ := chacha20poly1305.New(GenerateChacha20Poly1305Key(v.responseBodyKey))

		auth := &crypto.AEADAuthenticator{
			AEAD: aead,
			NonceGenerator: &ChunkNonceGenerator{
				Nonce: append([]byte(nil), v.responseBodyIV...),
				Size:  aead.NonceSize(),
			},
			AdditionalDataGenerator: crypto.NoOpBytesGenerator{},
		}
		authWriter = crypto.NewAuthenticationWriter(auth, writer)
	}

	return buf.NewWriter(authWriter)
}
