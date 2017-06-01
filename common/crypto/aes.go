package crypto

import (
	"crypto/aes"
	"crypto/cipher"

	"v2ray.com/core/common"
)

// NewAesDecryptionStream creates a new AES encryption stream based on given key and IV.
// Caller must ensure the length of key and IV is either 16, 24 or 32 bytes.
func NewAesDecryptionStream(key []byte, iv []byte) cipher.Stream {
	aesBlock, err := aes.NewCipher(key)
	common.Must(err)
	return cipher.NewCFBDecrypter(aesBlock, iv)
}

// NewAesEncryptionStream creates a new AES description stream based on given key and IV.
// Caller must ensure the length of key and IV is either 16, 24 or 32 bytes.
func NewAesEncryptionStream(key []byte, iv []byte) cipher.Stream {
	aesBlock, err := aes.NewCipher(key)
	common.Must(err)
	return cipher.NewCFBEncrypter(aesBlock, iv)
}
