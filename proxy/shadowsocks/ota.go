package shadowsocks

import (
	"crypto/hmac"
	"crypto/sha1"
	"io"

	"github.com/v2ray/v2ray-core/common/alloc"
	"github.com/v2ray/v2ray-core/common/log"
	"github.com/v2ray/v2ray-core/common/serial"
	"github.com/v2ray/v2ray-core/transport"
)

const (
	AuthSize = 10
)

type KeyGenerator func() []byte

type Authenticator struct {
	key KeyGenerator
}

func NewAuthenticator(keygen KeyGenerator) *Authenticator {
	return &Authenticator{
		key: keygen,
	}
}

func (this *Authenticator) Authenticate(auth []byte, data []byte) []byte {
	hasher := hmac.New(sha1.New, this.key())
	hasher.Write(data)
	res := hasher.Sum(nil)
	return append(auth, res[:AuthSize]...)
}

func HeaderKeyGenerator(key []byte, iv []byte) func() []byte {
	return func() []byte {
		newKey := make([]byte, 0, len(key)+len(iv))
		newKey = append(newKey, key...)
		newKey = append(newKey, iv...)
		return newKey
	}
}

func ChunkKeyGenerator(iv []byte) func() []byte {
	chunkId := 0
	return func() []byte {
		newKey := make([]byte, 0, len(iv)+4)
		newKey = append(newKey, iv...)
		newKey = append(newKey, serial.IntLiteral(chunkId).Bytes()...)
		chunkId++
		return newKey
	}
}

type ChunkReader struct {
	reader io.Reader
	auth   *Authenticator
}

func NewChunkReader(reader io.Reader, auth *Authenticator) *ChunkReader {
	return &ChunkReader{
		reader: reader,
		auth:   auth,
	}
}

func (this *ChunkReader) Read() (*alloc.Buffer, error) {
	buffer := alloc.NewLargeBuffer()
	if _, err := io.ReadFull(this.reader, buffer.Value[:2]); err != nil {
		alloc.Release(buffer)
		return nil, err
	}
	// There is a potential buffer overflow here. Large buffer is 64K bytes,
	// while uin16 + 10 will be more than that
	length := serial.BytesLiteral(buffer.Value[:2]).Uint16Value() + AuthSize
	if _, err := io.ReadFull(this.reader, buffer.Value[:length]); err != nil {
		alloc.Release(buffer)
		return nil, err
	}
	buffer.Slice(0, int(length))

	authBytes := buffer.Value[:AuthSize]
	payload := buffer.Value[AuthSize:]

	actualAuthBytes := this.auth.Authenticate(nil, payload)
	if !serial.BytesLiteral(authBytes).Equals(serial.BytesLiteral(actualAuthBytes)) {
		alloc.Release(buffer)
		log.Debug("AuthenticationReader: Unexpected auth: ", authBytes)
		return nil, transport.CorruptedPacket
	}
	buffer.Value = payload

	return buffer, nil
}
