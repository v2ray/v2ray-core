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
	h.Unlock()

	select {
	case <-h.token.Wait():
		go h.run()
	default:
	}
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

func (s *ServerSession) DecodeRequestHeader(reader io.Reader) (*protocol.RequestHeader, error) {
	buffer := make([]byte, 512)

	_, err := io.ReadFull(reader, buffer[:protocol.IDBytesLen])
	if err != nil {
		return nil, newError("failed to read request header").Base(err)
	}

	user, timestamp, valid := s.userValidator.Get(buffer[:protocol.IDBytesLen])
	if !valid {
		return nil, newError("invalid user")
	}

	timestampHash := md5.New()
	timestampHash.Write(hashTimestamp(timestamp))
	iv := timestampHash.Sum(nil)
	account, err := user.GetTypedAccount()
	if err != nil {
		return nil, newError("failed to get user account").Base(err)
	}
	vmessAccount := account.(*vmess.InternalAccount)

	aesStream := crypto.NewAesDecryptionStream(vmessAccount.ID.CmdKey(), iv)
	decryptor := crypto.NewCryptionReader(aesStream, reader)

	nBytes, err := io.ReadFull(decryptor, buffer[:41])
	if err != nil {
		return nil, newError("failed to read request header").Base(err)
	}
	bufferLen := nBytes

	request := &protocol.RequestHeader{
		User:    user,
		Version: buffer[0],
	}

	if request.Version != Version {
		return nil, newError("invalid protocol version ", request.Version)
	}

	s.requestBodyIV = append([]byte(nil), buffer[1:17]...)   // 16 bytes
	s.requestBodyKey = append([]byte(nil), buffer[17:33]...) // 16 bytes
	var sid sessionId
	copy(sid.user[:], vmessAccount.ID.Bytes())
	copy(sid.key[:], s.requestBodyKey)
	copy(sid.nonce[:], s.requestBodyIV)
	if s.sessionHistory.has(sid) {
		return nil, newError("duplicated session id, possibly under replay attack")
	}
	s.sessionHistory.add(sid)

	s.responseHeader = buffer[33]                       // 1 byte
	request.Option = protocol.RequestOption(buffer[34]) // 1 byte
	padingLen := int(buffer[35] >> 4)
	request.Security = protocol.NormSecurity(protocol.Security(buffer[35] & 0x0F))
	// 1 bytes reserved
	request.Command = protocol.RequestCommand(buffer[37])

	if request.Command != protocol.RequestCommandMux {
		request.Port = net.PortFromBytes(buffer[38:40])

		switch buffer[40] {
		case AddrTypeIPv4:
			_, err = io.ReadFull(decryptor, buffer[41:45]) // 4 bytes
			bufferLen += 4
			if err != nil {
				return nil, newError("failed to read IPv4 address").Base(err)
			}
			request.Address = net.IPAddress(buffer[41:45])
		case AddrTypeIPv6:
			_, err = io.ReadFull(decryptor, buffer[41:57]) // 16 bytes
			bufferLen += 16
			if err != nil {
				return nil, newError("failed to read IPv6 address").Base(err)
			}
			request.Address = net.IPAddress(buffer[41:57])
		case AddrTypeDomain:
			_, err = io.ReadFull(decryptor, buffer[41:42])
			if err != nil {
				return nil, newError("failed to read domain address").Base(err)
			}
			domainLength := int(buffer[41])
			if domainLength == 0 {
				return nil, newError("zero length domain").Base(err)
			}
			_, err = io.ReadFull(decryptor, buffer[42:42+domainLength])
			if err != nil {
				return nil, newError("failed to read domain address").Base(err)
			}
			bufferLen += 1 + domainLength
			request.Address = net.DomainAddress(string(buffer[42 : 42+domainLength]))
		}
	}

	if padingLen > 0 {
		_, err = io.ReadFull(decryptor, buffer[bufferLen:bufferLen+padingLen])
		if err != nil {
			return nil, newError("failed to read padding").Base(err)
		}
		bufferLen += padingLen
	}

	_, err = io.ReadFull(decryptor, buffer[bufferLen:bufferLen+4])
	if err != nil {
		return nil, newError("failed to read checksum").Base(err)
	}

	fnv1a := fnv.New32a()
	fnv1a.Write(buffer[:bufferLen])
	actualHash := fnv1a.Sum32()
	expectedHash := serial.BytesToUint32(buffer[bufferLen : bufferLen+4])

	if actualHash != expectedHash {
		return nil, newError("invalid auth")
	}

	if request.Address == nil {
		return nil, newError("invalid remote address")
	}

	return request, nil
}

func (s *ServerSession) DecodeRequestBody(request *protocol.RequestHeader, reader io.Reader) buf.Reader {
	var sizeParser crypto.ChunkSizeDecoder = crypto.PlainChunkSizeParser{}
	if request.Option.Has(protocol.RequestOptionChunkMasking) {
		sizeParser = NewShakeSizeParser(s.requestBodyIV)
	}
	if request.Security.Is(protocol.SecurityType_NONE) {
		if request.Option.Has(protocol.RequestOptionChunkStream) {
			if request.Command == protocol.RequestCommandTCP {
				return crypto.NewChunkStreamReader(sizeParser, reader)
			}

			auth := &crypto.AEADAuthenticator{
				AEAD:                    new(NoOpAuthenticator),
				NonceGenerator:          crypto.NoOpBytesGenerator{},
				AdditionalDataGenerator: crypto.NoOpBytesGenerator{},
			}
			return crypto.NewAuthenticationReader(auth, sizeParser, reader, protocol.TransferTypePacket)
		}

		return buf.NewReader(reader)
	}

	if request.Security.Is(protocol.SecurityType_LEGACY) {
		aesStream := crypto.NewAesDecryptionStream(s.requestBodyKey, s.requestBodyIV)
		cryptionReader := crypto.NewCryptionReader(aesStream, reader)
		if request.Option.Has(protocol.RequestOptionChunkStream) {
			auth := &crypto.AEADAuthenticator{
				AEAD:                    new(FnvAuthenticator),
				NonceGenerator:          crypto.NoOpBytesGenerator{},
				AdditionalDataGenerator: crypto.NoOpBytesGenerator{},
			}
			return crypto.NewAuthenticationReader(auth, sizeParser, cryptionReader, request.Command.TransferType())
		}

		return buf.NewReader(cryptionReader)
	}

	if request.Security.Is(protocol.SecurityType_AES128_GCM) {
		block, _ := aes.NewCipher(s.requestBodyKey)
		aead, _ := cipher.NewGCM(block)

		auth := &crypto.AEADAuthenticator{
			AEAD: aead,
			NonceGenerator: &ChunkNonceGenerator{
				Nonce: append([]byte(nil), s.requestBodyIV...),
				Size:  aead.NonceSize(),
			},
			AdditionalDataGenerator: crypto.NoOpBytesGenerator{},
		}
		return crypto.NewAuthenticationReader(auth, sizeParser, reader, request.Command.TransferType())
	}

	if request.Security.Is(protocol.SecurityType_CHACHA20_POLY1305) {
		aead, _ := chacha20poly1305.New(GenerateChacha20Poly1305Key(s.requestBodyKey))

		auth := &crypto.AEADAuthenticator{
			AEAD: aead,
			NonceGenerator: &ChunkNonceGenerator{
				Nonce: append([]byte(nil), s.requestBodyIV...),
				Size:  aead.NonceSize(),
			},
			AdditionalDataGenerator: crypto.NoOpBytesGenerator{},
		}
		return crypto.NewAuthenticationReader(auth, sizeParser, reader, request.Command.TransferType())
	}

	panic("Unknown security type.")
}

func (s *ServerSession) EncodeResponseHeader(header *protocol.ResponseHeader, writer io.Writer) {
	responseBodyKey := md5.Sum(s.requestBodyKey)
	responseBodyIV := md5.Sum(s.requestBodyIV)
	s.responseBodyKey = responseBodyKey[:]
	s.responseBodyIV = responseBodyIV[:]

	aesStream := crypto.NewAesEncryptionStream(s.responseBodyKey, s.responseBodyIV)
	encryptionWriter := crypto.NewCryptionWriter(aesStream, writer)
	s.responseWriter = encryptionWriter

	encryptionWriter.Write([]byte{s.responseHeader, byte(header.Option)})
	err := MarshalCommand(header.Command, encryptionWriter)
	if err != nil {
		encryptionWriter.Write([]byte{0x00, 0x00})
	}
}

func (s *ServerSession) EncodeResponseBody(request *protocol.RequestHeader, writer io.Writer) buf.Writer {
	var sizeParser crypto.ChunkSizeEncoder = crypto.PlainChunkSizeParser{}
	if request.Option.Has(protocol.RequestOptionChunkMasking) {
		sizeParser = NewShakeSizeParser(s.responseBodyIV)
	}
	if request.Security.Is(protocol.SecurityType_NONE) {
		if request.Option.Has(protocol.RequestOptionChunkStream) {
			if request.Command == protocol.RequestCommandTCP {
				return crypto.NewChunkStreamWriter(sizeParser, writer)
			}

			auth := &crypto.AEADAuthenticator{
				AEAD:                    new(NoOpAuthenticator),
				NonceGenerator:          &crypto.NoOpBytesGenerator{},
				AdditionalDataGenerator: crypto.NoOpBytesGenerator{},
			}
			return crypto.NewAuthenticationWriter(auth, sizeParser, writer, protocol.TransferTypePacket)
		}

		return buf.NewWriter(writer)
	}

	if request.Security.Is(protocol.SecurityType_LEGACY) {
		if request.Option.Has(protocol.RequestOptionChunkStream) {
			auth := &crypto.AEADAuthenticator{
				AEAD:                    new(FnvAuthenticator),
				NonceGenerator:          crypto.NoOpBytesGenerator{},
				AdditionalDataGenerator: crypto.NoOpBytesGenerator{},
			}
			return crypto.NewAuthenticationWriter(auth, sizeParser, s.responseWriter, request.Command.TransferType())
		}

		return buf.NewWriter(s.responseWriter)
	}

	if request.Security.Is(protocol.SecurityType_AES128_GCM) {
		block, _ := aes.NewCipher(s.responseBodyKey)
		aead, _ := cipher.NewGCM(block)

		auth := &crypto.AEADAuthenticator{
			AEAD: aead,
			NonceGenerator: &ChunkNonceGenerator{
				Nonce: append([]byte(nil), s.responseBodyIV...),
				Size:  aead.NonceSize(),
			},
			AdditionalDataGenerator: crypto.NoOpBytesGenerator{},
		}
		return crypto.NewAuthenticationWriter(auth, sizeParser, writer, request.Command.TransferType())
	}

	if request.Security.Is(protocol.SecurityType_CHACHA20_POLY1305) {
		aead, _ := chacha20poly1305.New(GenerateChacha20Poly1305Key(s.responseBodyKey))

		auth := &crypto.AEADAuthenticator{
			AEAD: aead,
			NonceGenerator: &ChunkNonceGenerator{
				Nonce: append([]byte(nil), s.responseBodyIV...),
				Size:  aead.NonceSize(),
			},
			AdditionalDataGenerator: crypto.NoOpBytesGenerator{},
		}
		return crypto.NewAuthenticationWriter(auth, sizeParser, writer, request.Command.TransferType())
	}

	panic("Unknown security type.")
}
