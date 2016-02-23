package shadowsocks

import (
	"crypto/cipher"
	"crypto/md5"

	"github.com/v2ray/v2ray-core/common/crypto"
	"github.com/v2ray/v2ray-core/common/protocol"
)

type Cipher interface {
	KeySize() int
	IVSize() int
	NewEncodingStream(key []byte, iv []byte) (cipher.Stream, error)
	NewDecodingStream(key []byte, iv []byte) (cipher.Stream, error)
}

type AesCfb struct {
	KeyBytes int
}

func (this *AesCfb) KeySize() int {
	return this.KeyBytes
}

func (this *AesCfb) IVSize() int {
	return 16
}

func (this *AesCfb) NewEncodingStream(key []byte, iv []byte) (cipher.Stream, error) {
	stream, err := crypto.NewAesEncryptionStream(key, iv)
	if err != nil {
		return nil, err
	}
	return stream, nil
}

func (this *AesCfb) NewDecodingStream(key []byte, iv []byte) (cipher.Stream, error) {
	stream, err := crypto.NewAesDecryptionStream(key, iv)
	if err != nil {
		return nil, err
	}
	return stream, nil
}

type ChaCha20 struct {
	IVBytes int
}

func (this *ChaCha20) KeySize() int {
	return 32
}

func (this *ChaCha20) IVSize() int {
	return this.IVBytes
}

func (this *ChaCha20) NewEncodingStream(key []byte, iv []byte) (cipher.Stream, error) {
	return crypto.NewChaCha20Stream(key, iv), nil
}

func (this *ChaCha20) NewDecodingStream(key []byte, iv []byte) (cipher.Stream, error) {
	return crypto.NewChaCha20Stream(key, iv), nil
}

type Config struct {
	Cipher Cipher
	Key    []byte
	UDP    bool
	Level  protocol.UserLevel
}

func PasswordToCipherKey(password string, keySize int) []byte {
	pwdBytes := []byte(password)
	key := make([]byte, 0, keySize)

	md5Sum := md5.Sum(pwdBytes)
	key = append(key, md5Sum[:]...)

	for len(key) < keySize {
		md5Hash := md5.New()
		md5Hash.Write(md5Sum[:])
		md5Hash.Write(pwdBytes)
		md5Hash.Sum(md5Sum[:0])

		key = append(key, md5Sum[:]...)
	}
	return key
}
