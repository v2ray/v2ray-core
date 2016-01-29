package shadowsocks

import (
	"crypto/hmac"
	"crypto/sha1"

	"github.com/v2ray/v2ray-core/common/serial"
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

func (this *Authenticator) AuthSize() int {
	return AuthSize
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
