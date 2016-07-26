package crypto

import (
	"crypto/cipher"

	"github.com/v2ray/v2ray-core/common/crypto/internal"
)

// NewChaCha20Stream creates a new Chacha20 encryption/descryption stream based on give key and IV.
// Caller must ensure the length of key is 32 bytes, and length of IV is either 8 or 12 bytes.
func NewChaCha20Stream(key []byte, iv []byte) cipher.Stream {
	return internal.NewChaCha20Stream(key, iv, 20)
}
