package shadowsocks

import (
	"crypto/cipher"
	"crypto/md5"

	"v2ray.com/core/common/crypto"
	"v2ray.com/core/common/protocol"
)

func (this *Config) GetCipher() Cipher {
	switch this.Cipher {
	case Config_AES_128_CFB:
		return &AesCfb{KeyBytes: 16}
	case Config_AES_256_CFB:
		return &AesCfb{KeyBytes: 32}
	case Config_CHACHA20:
		return &ChaCha20{IVBytes: 8}
	case Config_CHACHA20_IEFT:
		return &ChaCha20{IVBytes: 12}
	}
	panic("Failed to create Cipher. Should not happen.")
}

func (this *Account) Equals(another protocol.Account) bool {
	if account, ok := another.(*Account); ok {
		return account.Password == this.Password
	}
	return false
}

func (this *Account) AsAccount() (protocol.Account, error) {
	return this, nil
}

func (this *Account) GetCipherKey(size int) []byte {
	return PasswordToCipherKey(this.Password, size)
}

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
	stream := crypto.NewAesEncryptionStream(key, iv)
	return stream, nil
}

func (this *AesCfb) NewDecodingStream(key []byte, iv []byte) (cipher.Stream, error) {
	stream := crypto.NewAesDecryptionStream(key, iv)
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
