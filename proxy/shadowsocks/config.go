package shadowsocks

import (
	"bytes"
	"crypto/cipher"
	"crypto/md5"
	"v2ray.com/core/common/crypto"
	"v2ray.com/core/common/errors"
	"v2ray.com/core/common/protocol"
)

type ShadowsocksAccount struct {
	Cipher      Cipher
	Key         []byte
	OneTimeAuth Account_OneTimeAuth
}

func (v *ShadowsocksAccount) Equals(another protocol.Account) bool {
	if account, ok := another.(*ShadowsocksAccount); ok {
		return bytes.Equal(v.Key, account.Key)
	}
	return false
}

func (v *Account) GetCipher() (Cipher, error) {
	switch v.CipherType {
	case CipherType_AES_128_CFB:
		return &AesCfb{KeyBytes: 16}, nil
	case CipherType_AES_256_CFB:
		return &AesCfb{KeyBytes: 32}, nil
	case CipherType_CHACHA20:
		return &ChaCha20{IVBytes: 8}, nil
	case CipherType_CHACHA20_IEFT:
		return &ChaCha20{IVBytes: 12}, nil
	default:
		return nil, errors.New("Unsupported cipher.")
	}
}

func (v *Account) AsAccount() (protocol.Account, error) {
	cipher, err := v.GetCipher()
	if err != nil {
		return nil, errors.Base(err).Message("Shadowsocks|Account: Failed to get cipher.")
	}
	return &ShadowsocksAccount{
		Cipher:      cipher,
		Key:         v.GetCipherKey(),
		OneTimeAuth: v.Ota,
	}, nil
}

func (v *Account) GetCipherKey() []byte {
	ct, err := v.GetCipher()
	if err != nil {
		return nil
	}
	return PasswordToCipherKey(v.Password, ct.KeySize())
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

func (v *AesCfb) KeySize() int {
	return v.KeyBytes
}

func (v *AesCfb) IVSize() int {
	return 16
}

func (v *AesCfb) NewEncodingStream(key []byte, iv []byte) (cipher.Stream, error) {
	stream := crypto.NewAesEncryptionStream(key, iv)
	return stream, nil
}

func (v *AesCfb) NewDecodingStream(key []byte, iv []byte) (cipher.Stream, error) {
	stream := crypto.NewAesDecryptionStream(key, iv)
	return stream, nil
}

type ChaCha20 struct {
	IVBytes int
}

func (v *ChaCha20) KeySize() int {
	return 32
}

func (v *ChaCha20) IVSize() int {
	return v.IVBytes
}

func (v *ChaCha20) NewEncodingStream(key []byte, iv []byte) (cipher.Stream, error) {
	return crypto.NewChaCha20Stream(key, iv), nil
}

func (v *ChaCha20) NewDecodingStream(key []byte, iv []byte) (cipher.Stream, error) {
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
