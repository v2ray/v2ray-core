package io

import (
	"crypto/aes"
	"crypto/cipher"
	"io"
)

func NewAesDecryptReader(key []byte, iv []byte, reader io.Reader) (*CryptionReader, error) {
	aesBlock, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesStream := cipher.NewCFBDecrypter(aesBlock, iv)
	return NewCryptionReader(aesStream, reader), nil
}

func NewAesEncryptWriter(key []byte, iv []byte, writer io.Writer) (*CryptionWriter, error) {
	aesBlock, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesStream := cipher.NewCFBEncrypter(aesBlock, iv)
	return NewCryptionWriter(aesStream, writer), nil
}
