package crypto

import (
	"crypto/cipher"

	"github.com/aead/chacha20"
)

// NewChaCha20Stream creates a new Chacha/20 cipher stream. Caller must ensure that key is 32-bytes long and iv is either 8 or 12 bytes.
func NewChaCha20Stream(key []byte, iv []byte) cipher.Stream {
	var keyArray [32]byte
	var nonce [12]byte
	copy(keyArray[:], key)
	switch len(iv) {
	case 8:
		copy(nonce[4:], iv)
	case 12:
		copy(nonce[:], iv)
	}
	return chacha20.NewCipher(&nonce, &keyArray)
}
