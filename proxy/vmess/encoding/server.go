package encoding

import (
	"context"
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
	"v2ray.com/core/common/signal"
	"v2ray.com/core/proxy/vmess"
)

type sessionId struct {
	user  [16]byte
	key   [16]byte
	nonce [16]byte
}

type SessionHistory struct {
	sync.RWMutex
	cache map[sessionId]time.Time
	token *signal.Semaphore
	ctx   context.Context
}

func NewSessionHistory(ctx context.Context) *SessionHistory {
	h := &SessionHistory{
		cache: make(map[sessionId]time.Time, 128),
		token: signal.NewSemaphore(1),
		ctx:   ctx,
	}
	return h
}

func (h *SessionHistory) add(session sessionId) {
	h.Lock()
	h.cache[session] = time.Now().Add(time.Minute * 3)

	select {
	case <-h.token.Wait():
		go h.run()
	default:
	}

	h.Unlock()
}

func (h *SessionHistory) has(session sessionId) bool {
	h.RLock()
	defer h.RUnlock()

	if expire, found := h.cache[session]; found {
		return expire.After(time.Now())
	}
	return false
}

func (h *SessionHistory) run() {
	defer h.token.Signal()

	for {
		select {
		case <-h.ctx.Done():
			return
		case <-time.After(time.Second * 30):
		}
		session2Remove := make([]sessionId, 0, 16)
		now := time.Now()
		h.Lock()
		if len(h.cache) == 0 {
			h.Unlock()
			return
		}
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

func getSizeMask(b []byte) crypto.Uint16Generator {
	return crypto.NewShakeUint16Generator(b)
}

type ServerSession struct {
	userValidator   protocol.UserValidator
	sessionHistory  *SessionHistory
	requestBodyKey  []byte
	requestBodyIV   []byte
	responseBodyKey []byte
	responseBodyIV  []byte
	responseHeader  byte
	responseWriter  io.Writer
}

// NewServerSession creates a new ServerSession, using the given UserValidator.
// The ServerSession instance doesn't take ownership of the validator.
func NewServerSession(validator protocol.UserValidator, sessionHistory *SessionHistory) *ServerSession {
	return &ServerSession{
		userValidator:  validator,
		sessionHistory: sessionHistory,
	}
}

func (v *ServerSession) DecodeRequestHeader(reader io.Reader) (*protocol.RequestHeader, error) {
	buffer := make([]byte, 512)

	_, err := io.ReadFull(reader, buffer[:protocol.IDBytesLen])
	if err != nil {
		return nil, errors.New("failed to read request header").Base(err).Path("VMess", "Server")
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
		return nil, errors.New("failed to get user account").Base(err).Path("VMess", "Server")
	}
	vmessAccount := account.(*vmess.InternalAccount)

	aesStream := crypto.NewAesDecryptionStream(vmessAccount.ID.CmdKey(), iv)
	decryptor := crypto.NewCryptionReader(aesStream, reader)

	nBytes, err := io.ReadFull(decryptor, buffer[:41])
	if err != nil {
		return nil, errors.New("failed to read request header").Base(err).Path("VMess", "Server")
	}
	bufferLen := nBytes

	request := &protocol.RequestHeader{
		User:    user,
		Version: buffer[0],
	}

	if request.Version != Version {
		return nil, errors.New("invalid protocol version ", request.Version).Path("VMess", "Server")
	}

	v.requestBodyIV = append([]byte(nil), buffer[1:17]...)   // 16 bytes
	v.requestBodyKey = append([]byte(nil), buffer[17:33]...) // 16 bytes
	var sid sessionId
	copy(sid.user[:], vmessAccount.ID.Bytes())
	copy(sid.key[:], v.requestBodyKey)
	copy(sid.nonce[:], v.requestBodyIV)
	if v.sessionHistory.has(sid) {
		return nil, errors.New("duplicated session id, possibly under replay attack").Path("VMess", "Server")
	}
	v.sessionHistory.add(sid)

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
			return nil, errors.New("failed to read IPv4 address").Base(err).Path("VMess", "Server")
		}
		request.Address = net.IPAddress(buffer[41:45])
	case AddrTypeIPv6:
		_, err = io.ReadFull(decryptor, buffer[41:57]) // 16 bytes
		bufferLen += 16
		if err != nil {
			return nil, errors.New("failed to read IPv6 address").Base(err).Path("VMess", "Server")
		}
		request.Address = net.IPAddress(buffer[41:57])
	case AddrTypeDomain:
		_, err = io.ReadFull(decryptor, buffer[41:42])
		if err != nil {
			return nil, errors.New("failed to read domain address").Base(err).Path("VMess", "Server")
		}
		domainLength := int(buffer[41])
		if domainLength == 0 {
			return nil, errors.New("zero length domain").Base(err).Path("VMess", "Server")
		}
		_, err = io.ReadFull(decryptor, buffer[42:42+domainLength])
		if err != nil {
			return nil, errors.New("failed to read domain address").Base(err).Path("VMess", "Server")
		}
		bufferLen += 1 + domainLength
		request.Address = net.DomainAddress(string(buffer[42 : 42+domainLength]))
	}

	if padingLen > 0 {
		_, err = io.ReadFull(decryptor, buffer[bufferLen:bufferLen+padingLen])
		if err != nil {
			return nil, errors.New("failed to read padding").Base(err).Path("VMess", "Server")
		}
		bufferLen += padingLen
	}

	_, err = io.ReadFull(decryptor, buffer[bufferLen:bufferLen+4])
	if err != nil {
		return nil, errors.New("failed to read checksum").Base(err).Path("VMess", "Server")
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
	var sizeMask crypto.Uint16Generator = crypto.StaticUint16Generator(0)
	if request.Option.Has(protocol.RequestOptionChunkMasking) {
		sizeMask = getSizeMask(v.requestBodyIV)
	}
	if request.Security.Is(protocol.SecurityType_NONE) {
		if request.Option.Has(protocol.RequestOptionChunkStream) {
			auth := &crypto.AEADAuthenticator{
				AEAD:                    NoOpAuthenticator{},
				NonceGenerator:          crypto.NoOpBytesGenerator{},
				AdditionalDataGenerator: crypto.NoOpBytesGenerator{},
			}
			authReader = crypto.NewAuthenticationReader(auth, reader, sizeMask)
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
			authReader = crypto.NewAuthenticationReader(auth, cryptionReader, sizeMask)
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
		authReader = crypto.NewAuthenticationReader(auth, reader, sizeMask)
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
		authReader = crypto.NewAuthenticationReader(auth, reader, sizeMask)
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
	var sizeMask crypto.Uint16Generator = crypto.StaticUint16Generator(0)
	if request.Option.Has(protocol.RequestOptionChunkMasking) {
		sizeMask = getSizeMask(v.responseBodyIV)
	}
	if request.Security.Is(protocol.SecurityType_NONE) {
		if request.Option.Has(protocol.RequestOptionChunkStream) {
			auth := &crypto.AEADAuthenticator{
				AEAD:                    NoOpAuthenticator{},
				NonceGenerator:          crypto.NoOpBytesGenerator{},
				AdditionalDataGenerator: crypto.NoOpBytesGenerator{},
			}
			authWriter = crypto.NewAuthenticationWriter(auth, writer, sizeMask)
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
			authWriter = crypto.NewAuthenticationWriter(auth, v.responseWriter, sizeMask)
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
		authWriter = crypto.NewAuthenticationWriter(auth, writer, sizeMask)
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
		authWriter = crypto.NewAuthenticationWriter(auth, writer, sizeMask)
	}

	return buf.NewWriter(authWriter)
}
