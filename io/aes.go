package io

import (
	"crypto/aes"
	"crypto/cipher"
	"io"
)

func NewAesDecryptReader(key []byte, iv []byte, reader io.Reader) (io.Reader, error) {
	aesBlock, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesMode := cipher.NewCBCDecrypter(aesBlock, iv)
	return NewCryptionReader(aesMode, reader), nil
}

func NewAesEncryptWriter(key []byte, iv []byte, writer io.Writer) (io.Writer, error) {
	aesBlock, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesMode := cipher.NewCBCEncrypter(aesBlock, iv)
	return NewCryptionWriter(aesMode, writer), nil
}
