package crypto

import (
	"crypto/cipher"

	"github.com/aead/chacha20"
)

func makeNonce(nonce *[chacha20.NonceSize]byte, iv []byte) {
	switch len(iv) {
	case 8:
		copy(nonce[4:], iv)
	case 12:
		copy(nonce[:], iv)
	default:
		panic("bad nonce length")
	}
}

// NewChaCha20Stream creates a new Chacha20 encryption/descryption stream based on give key and IV.
// Caller must ensure the length of key is 32 bytes, and length of IV is either 8 or 12 bytes.
func NewChaCha20Stream(key []byte, iv []byte) cipher.Stream {
	var Key [32]byte
	var Nonce [12]byte

	if len(key) != 32 {
		panic("bad key length")
	}

	copy(Key[:], key)
	makeNonce(&Nonce, iv)
	return chacha20.NewCipher(&Nonce, &Key)
}
