package crypto

import (
	"crypto/aes"
	"crypto/cipher"
)

func NewAesDecryptionStream(key []byte, iv []byte) cipher.Stream {
	aesBlock, _ := aes.NewCipher(key)
	return cipher.NewCFBDecrypter(aesBlock, iv)
}

func NewAesEncryptionStream(key []byte, iv []byte) cipher.Stream {
	aesBlock, _ := aes.NewCipher(key)
	return cipher.NewCFBEncrypter(aesBlock, iv)
}
