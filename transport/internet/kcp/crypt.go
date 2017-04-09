package kcp

import (
	"crypto/cipher"
	"hash/fnv"

	"v2ray.com/core/common/serial"
)

// SimpleAuthenticator is a legacy AEAD used for KCP encryption.
type SimpleAuthenticator struct{}

// NewSimpleAuthenticator creates a new SimpleAuthenticator
func NewSimpleAuthenticator() cipher.AEAD {
	return &SimpleAuthenticator{}
}

// NonceSize implements cipher.AEAD.NonceSize().
func (v *SimpleAuthenticator) NonceSize() int {
	return 0
}

// Overhead implements cipher.AEAD.NonceSize().
func (v *SimpleAuthenticator) Overhead() int {
	return 6
}

// Seal implements cipher.AEAD.Seal().
func (v *SimpleAuthenticator) Seal(dst, nonce, plain, extra []byte) []byte {
	dst = append(dst, 0, 0, 0, 0)
	dst = serial.Uint16ToBytes(uint16(len(plain)), dst)
	dst = append(dst, plain...)

	fnvHash := fnv.New32a()
	fnvHash.Write(dst[4:])
	fnvHash.Sum(dst[:0])

	len := len(dst)
	xtra := 4 - len%4
	if xtra != 4 {
		dst = append(dst, make([]byte, xtra)...)
	}
	xorfwd(dst)
	if xtra != 4 {
		dst = dst[:len]
	}
	return dst
}

// Open implements cipher.AEAD.Open().
func (v *SimpleAuthenticator) Open(dst, nonce, cipherText, extra []byte) ([]byte, error) {
	dst = append(dst, cipherText...)
	dstLen := len(dst)
	xtra := 4 - dstLen%4
	if xtra != 4 {
		dst = append(dst, make([]byte, xtra)...)
	}
	xorbkd(dst)
	if xtra != 4 {
		dst = dst[:dstLen]
	}

	fnvHash := fnv.New32a()
	fnvHash.Write(dst[4:])
	if serial.BytesToUint32(dst[:4]) != fnvHash.Sum32() {
		return nil, newError("invalid auth")
	}

	length := serial.BytesToUint16(dst[4:6])
	if len(dst)-6 != int(length) {
		return nil, newError("invalid auth")
	}

	return dst[6:], nil
}
