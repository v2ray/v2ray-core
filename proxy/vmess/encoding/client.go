package encoding

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"hash/fnv"
	"io"

	"golang.org/x/crypto/chacha20poly1305"

	"v2ray.com/core/app/log"
	"v2ray.com/core/common"
	"v2ray.com/core/common/bitmask"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/crypto"
	"v2ray.com/core/common/dice"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/protocol"
	"v2ray.com/core/common/serial"
	"v2ray.com/core/proxy/vmess"
)

func hashTimestamp(t protocol.Timestamp) []byte {
	bytes := make([]byte, 0, 32)
	bytes = t.Bytes(bytes)
	bytes = t.Bytes(bytes)
	bytes = t.Bytes(bytes)
	bytes = t.Bytes(bytes)
	return bytes
}

// ClientSession stores connection session info for VMess client.
type ClientSession struct {
	requestBodyKey  []byte
	requestBodyIV   []byte
	responseHeader  byte
	responseBodyKey []byte
	responseBodyIV  []byte
	responseReader  io.Reader
	idHash          protocol.IDHash
}

// NewClientSession creates a new ClientSession.
func NewClientSession(idHash protocol.IDHash) *ClientSession {
	randomBytes := make([]byte, 33) // 16 + 16 + 1
	common.Must2(rand.Read(randomBytes))

	session := &ClientSession{}
	session.requestBodyKey = randomBytes[:16]
	session.requestBodyIV = randomBytes[16:32]
	session.responseHeader = randomBytes[32]
	responseBodyKey := md5.Sum(session.requestBodyKey)
	responseBodyIV := md5.Sum(session.requestBodyIV)
	session.responseBodyKey = responseBodyKey[:]
	session.responseBodyIV = responseBodyIV[:]
	session.idHash = idHash

	return session
}

func (c *ClientSession) EncodeRequestHeader(header *protocol.RequestHeader, writer io.Writer) {
	timestamp := protocol.NewTimestampGenerator(protocol.NowTime(), 30)()
	account, err := header.User.GetTypedAccount()
	if err != nil {
		log.Trace(newError("failed to get user account: ", err).AtError())
		return
	}
	idHash := c.idHash(account.(*vmess.InternalAccount).AnyValidID().Bytes())
	common.Must2(idHash.Write(timestamp.Bytes(nil)))
	common.Must2(writer.Write(idHash.Sum(nil)))

	buffer := make([]byte, 0, 512)
	buffer = append(buffer, Version)
	buffer = append(buffer, c.requestBodyIV...)
	buffer = append(buffer, c.requestBodyKey...)
	buffer = append(buffer, c.responseHeader, byte(header.Option))
	padingLen := dice.Roll(16)
	if header.Security.Is(protocol.SecurityType_LEGACY) {
		// Disable padding in legacy mode for a smooth transition.
		padingLen = 0
	}
	security := byte(padingLen<<4) | byte(header.Security)
	buffer = append(buffer, security, byte(0), byte(header.Command))

	if header.Command != protocol.RequestCommandMux {
		buffer = header.Port.Bytes(buffer)

		switch header.Address.Family() {
		case net.AddressFamilyIPv4:
			buffer = append(buffer, AddrTypeIPv4)
			buffer = append(buffer, header.Address.IP()...)
		case net.AddressFamilyIPv6:
			buffer = append(buffer, AddrTypeIPv6)
			buffer = append(buffer, header.Address.IP()...)
		case net.AddressFamilyDomain:
			buffer = append(buffer, AddrTypeDomain, byte(len(header.Address.Domain())))
			buffer = append(buffer, header.Address.Domain()...)
		}
	}

	if padingLen > 0 {
		pading := make([]byte, padingLen)
		common.Must2(rand.Read(pading))
		buffer = append(buffer, pading...)
	}

	fnv1a := fnv.New32a()
	common.Must2(fnv1a.Write(buffer))

	buffer = fnv1a.Sum(buffer)

	timestampHash := md5.New()
	common.Must2(timestampHash.Write(hashTimestamp(timestamp)))
	iv := timestampHash.Sum(nil)
	aesStream := crypto.NewAesEncryptionStream(account.(*vmess.InternalAccount).ID.CmdKey(), iv)
	aesStream.XORKeyStream(buffer, buffer)
	common.Must2(writer.Write(buffer))
}

func (c *ClientSession) EncodeRequestBody(request *protocol.RequestHeader, writer io.Writer) buf.Writer {
	var sizeParser crypto.ChunkSizeEncoder = crypto.PlainChunkSizeParser{}
	if request.Option.Has(protocol.RequestOptionChunkMasking) {
		sizeParser = NewShakeSizeParser(c.requestBodyIV)
	}
	if request.Security.Is(protocol.SecurityType_NONE) {
		if request.Option.Has(protocol.RequestOptionChunkStream) {
			if request.Command == protocol.RequestCommandTCP {
				return crypto.NewChunkStreamWriter(sizeParser, writer)
			}
			auth := &crypto.AEADAuthenticator{
				AEAD:                    new(NoOpAuthenticator),
				NonceGenerator:          crypto.NoOpBytesGenerator{},
				AdditionalDataGenerator: crypto.NoOpBytesGenerator{},
			}
			return crypto.NewAuthenticationWriter(auth, sizeParser, writer, protocol.TransferTypePacket)
		}

		return buf.NewWriter(writer)
	}

	if request.Security.Is(protocol.SecurityType_LEGACY) {
		aesStream := crypto.NewAesEncryptionStream(c.requestBodyKey, c.requestBodyIV)
		cryptionWriter := crypto.NewCryptionWriter(aesStream, writer)
		if request.Option.Has(protocol.RequestOptionChunkStream) {
			auth := &crypto.AEADAuthenticator{
				AEAD:                    new(FnvAuthenticator),
				NonceGenerator:          crypto.NoOpBytesGenerator{},
				AdditionalDataGenerator: crypto.NoOpBytesGenerator{},
			}
			return crypto.NewAuthenticationWriter(auth, sizeParser, cryptionWriter, request.Command.TransferType())
		}

		return buf.NewWriter(cryptionWriter)
	}

	if request.Security.Is(protocol.SecurityType_AES128_GCM) {
		block, _ := aes.NewCipher(c.requestBodyKey)
		aead, _ := cipher.NewGCM(block)

		auth := &crypto.AEADAuthenticator{
			AEAD: aead,
			NonceGenerator: &ChunkNonceGenerator{
				Nonce: append([]byte(nil), c.requestBodyIV...),
				Size:  aead.NonceSize(),
			},
			AdditionalDataGenerator: crypto.NoOpBytesGenerator{},
		}
		return crypto.NewAuthenticationWriter(auth, sizeParser, writer, request.Command.TransferType())
	}

	if request.Security.Is(protocol.SecurityType_CHACHA20_POLY1305) {
		aead, _ := chacha20poly1305.New(GenerateChacha20Poly1305Key(c.requestBodyKey))

		auth := &crypto.AEADAuthenticator{
			AEAD: aead,
			NonceGenerator: &ChunkNonceGenerator{
				Nonce: append([]byte(nil), c.requestBodyIV...),
				Size:  aead.NonceSize(),
			},
			AdditionalDataGenerator: crypto.NoOpBytesGenerator{},
		}
		return crypto.NewAuthenticationWriter(auth, sizeParser, writer, request.Command.TransferType())
	}

	panic("Unknown security type.")
}

func (c *ClientSession) DecodeResponseHeader(reader io.Reader) (*protocol.ResponseHeader, error) {
	aesStream := crypto.NewAesDecryptionStream(c.responseBodyKey, c.responseBodyIV)
	c.responseReader = crypto.NewCryptionReader(aesStream, reader)

	buffer := make([]byte, 256)

	_, err := io.ReadFull(c.responseReader, buffer[:4])
	if err != nil {
		log.Trace(newError("failed to read response header").Base(err))
		return nil, err
	}

	if buffer[0] != c.responseHeader {
		return nil, newError("unexpected response header. Expecting ", int(c.responseHeader), " but actually ", int(buffer[0]))
	}

	header := &protocol.ResponseHeader{
		Option: bitmask.Byte(buffer[1]),
	}

	if buffer[2] != 0 {
		cmdID := buffer[2]
		dataLen := int(buffer[3])
		_, err := io.ReadFull(c.responseReader, buffer[:dataLen])
		if err != nil {
			log.Trace(newError("failed to read response command").Base(err))
			return nil, err
		}
		data := buffer[:dataLen]
		command, err := UnmarshalCommand(cmdID, data)
		if err == nil {
			header.Command = command
		}
	}

	return header, nil
}

func (c *ClientSession) DecodeResponseBody(request *protocol.RequestHeader, reader io.Reader) buf.Reader {
	var sizeParser crypto.ChunkSizeDecoder = crypto.PlainChunkSizeParser{}
	if request.Option.Has(protocol.RequestOptionChunkMasking) {
		sizeParser = NewShakeSizeParser(c.responseBodyIV)
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
		if request.Option.Has(protocol.RequestOptionChunkStream) {
			auth := &crypto.AEADAuthenticator{
				AEAD:                    new(FnvAuthenticator),
				NonceGenerator:          crypto.NoOpBytesGenerator{},
				AdditionalDataGenerator: crypto.NoOpBytesGenerator{},
			}
			return crypto.NewAuthenticationReader(auth, sizeParser, c.responseReader, request.Command.TransferType())
		}

		return buf.NewReader(c.responseReader)
	}

	if request.Security.Is(protocol.SecurityType_AES128_GCM) {
		block, _ := aes.NewCipher(c.responseBodyKey)
		aead, _ := cipher.NewGCM(block)

		auth := &crypto.AEADAuthenticator{
			AEAD: aead,
			NonceGenerator: &ChunkNonceGenerator{
				Nonce: append([]byte(nil), c.responseBodyIV...),
				Size:  aead.NonceSize(),
			},
			AdditionalDataGenerator: crypto.NoOpBytesGenerator{},
		}
		return crypto.NewAuthenticationReader(auth, sizeParser, reader, request.Command.TransferType())
	}

	if request.Security.Is(protocol.SecurityType_CHACHA20_POLY1305) {
		aead, _ := chacha20poly1305.New(GenerateChacha20Poly1305Key(c.responseBodyKey))

		auth := &crypto.AEADAuthenticator{
			AEAD: aead,
			NonceGenerator: &ChunkNonceGenerator{
				Nonce: append([]byte(nil), c.responseBodyIV...),
				Size:  aead.NonceSize(),
			},
			AdditionalDataGenerator: crypto.NoOpBytesGenerator{},
		}
		return crypto.NewAuthenticationReader(auth, sizeParser, reader, request.Command.TransferType())
	}

	panic("Unknown security type.")
}

type ChunkNonceGenerator struct {
	Nonce []byte
	Size  int
	count uint16
}

func (g *ChunkNonceGenerator) Next() []byte {
	serial.Uint16ToBytes(g.count, g.Nonce[:0])
	g.count++
	return g.Nonce[:g.Size]
}
