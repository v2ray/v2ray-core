package encoding

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"hash/fnv"
	"io"

	"golang.org/x/crypto/chacha20poly1305"

	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/crypto"
	"v2ray.com/core/common/dice"
	"v2ray.com/core/common/errors"
	"v2ray.com/core/common/log"
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
	rand.Read(randomBytes)

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

func (v *ClientSession) EncodeRequestHeader(header *protocol.RequestHeader, writer io.Writer) {
	timestamp := protocol.NewTimestampGenerator(protocol.NowTime(), 30)()
	account, err := header.User.GetTypedAccount()
	if err != nil {
		log.Error("VMess: Failed to get user account: ", err)
		return
	}
	idHash := v.idHash(account.(*vmess.InternalAccount).AnyValidID().Bytes())
	idHash.Write(timestamp.Bytes(nil))
	writer.Write(idHash.Sum(nil))

	buffer := make([]byte, 0, 512)
	buffer = append(buffer, Version)
	buffer = append(buffer, v.requestBodyIV...)
	buffer = append(buffer, v.requestBodyKey...)
	buffer = append(buffer, v.responseHeader, byte(header.Option))
	padingLen := dice.Roll(16)
	if header.Security.Is(protocol.SecurityType_LEGACY) {
		// Disable padding in legacy mode for a smooth transition.
		padingLen = 0
	}
	security := byte(padingLen<<4) | byte(header.Security)
	buffer = append(buffer, security, byte(0), byte(header.Command))
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

	if padingLen > 0 {
		pading := make([]byte, padingLen)
		rand.Read(pading)
		buffer = append(buffer, pading...)
	}

	fnv1a := fnv.New32a()
	fnv1a.Write(buffer)

	buffer = fnv1a.Sum(buffer)

	timestampHash := md5.New()
	timestampHash.Write(hashTimestamp(timestamp))
	iv := timestampHash.Sum(nil)
	aesStream := crypto.NewAesEncryptionStream(account.(*vmess.InternalAccount).ID.CmdKey(), iv)
	aesStream.XORKeyStream(buffer, buffer)
	writer.Write(buffer)

	return
}

func (v *ClientSession) EncodeRequestBody(request *protocol.RequestHeader, writer io.Writer) buf.Writer {
	var authWriter io.Writer
	if request.Security.Is(protocol.SecurityType_NONE) {
		if request.Option.Has(protocol.RequestOptionChunkStream) {
			auth := &crypto.AEADAuthenticator{
				AEAD:                    NoOpAuthenticator{},
				NonceGenerator:          crypto.NoOpBytesGenerator{},
				AdditionalDataGenerator: crypto.NoOpBytesGenerator{},
			}
			authWriter = crypto.NewAuthenticationWriter(auth, writer)
		} else {
			authWriter = writer
		}
	} else if request.Security.Is(protocol.SecurityType_LEGACY) {
		aesStream := crypto.NewAesEncryptionStream(v.requestBodyKey, v.requestBodyIV)
		cryptionWriter := crypto.NewCryptionWriter(aesStream, writer)
		if request.Option.Has(protocol.RequestOptionChunkStream) {
			auth := &crypto.AEADAuthenticator{
				AEAD:                    new(FnvAuthenticator),
				NonceGenerator:          crypto.NoOpBytesGenerator{},
				AdditionalDataGenerator: crypto.NoOpBytesGenerator{},
			}
			authWriter = crypto.NewAuthenticationWriter(auth, cryptionWriter)
		} else {
			authWriter = cryptionWriter
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
		authWriter = crypto.NewAuthenticationWriter(auth, writer)
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
		authWriter = crypto.NewAuthenticationWriter(auth, writer)
	}

	return buf.NewWriter(authWriter)

}

func (v *ClientSession) DecodeResponseHeader(reader io.Reader) (*protocol.ResponseHeader, error) {
	aesStream := crypto.NewAesDecryptionStream(v.responseBodyKey, v.responseBodyIV)
	v.responseReader = crypto.NewCryptionReader(aesStream, reader)

	buffer := make([]byte, 256)

	_, err := io.ReadFull(v.responseReader, buffer[:4])
	if err != nil {
		log.Info("VMess|Client: Failed to read response header: ", err)
		return nil, err
	}

	if buffer[0] != v.responseHeader {
		return nil, errors.Format("VMess|Client: Unexpected response header. Expecting %d but actually %d", v.responseHeader, buffer[0])
	}

	header := &protocol.ResponseHeader{
		Option: protocol.ResponseOption(buffer[1]),
	}

	if buffer[2] != 0 {
		cmdID := buffer[2]
		dataLen := int(buffer[3])
		_, err := io.ReadFull(v.responseReader, buffer[:dataLen])
		if err != nil {
			log.Info("VMess|Client: Failed to read response command: ", err)
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

func (v *ClientSession) DecodeResponseBody(request *protocol.RequestHeader, reader io.Reader) buf.Reader {
	var authReader io.Reader
	if request.Security.Is(protocol.SecurityType_NONE) {
		if request.Option.Has(protocol.RequestOptionChunkStream) {
			auth := &crypto.AEADAuthenticator{
				AEAD:                    new(FnvAuthenticator),
				NonceGenerator:          crypto.NoOpBytesGenerator{},
				AdditionalDataGenerator: crypto.NoOpBytesGenerator{},
			}
			authReader = crypto.NewAuthenticationReader(auth, reader)
		} else {
			authReader = reader
		}
	} else if request.Security.Is(protocol.SecurityType_LEGACY) {
		if request.Option.Has(protocol.RequestOptionChunkStream) {
			auth := &crypto.AEADAuthenticator{
				AEAD:                    new(FnvAuthenticator),
				NonceGenerator:          crypto.NoOpBytesGenerator{},
				AdditionalDataGenerator: crypto.NoOpBytesGenerator{},
			}
			authReader = crypto.NewAuthenticationReader(auth, v.responseReader)
		} else {
			authReader = v.responseReader
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
		authReader = crypto.NewAuthenticationReader(auth, reader)
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
		authReader = crypto.NewAuthenticationReader(auth, reader)
	}

	return buf.NewReader(authReader)
}

type ChunkNonceGenerator struct {
	Nonce []byte
	Size  int
	count uint16
}

func (v *ChunkNonceGenerator) Next() []byte {
	serial.Uint16ToBytes(v.count, v.Nonce[:0])
	v.count++
	return v.Nonce[:v.Size]
}
