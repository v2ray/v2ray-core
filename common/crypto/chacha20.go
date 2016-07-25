package crypto

import (
	"crypto/cipher"

	"github.com/v2ray/v2ray-core/common/crypto/internal"
)

func NewChaCha20Stream(key []byte, iv []byte) cipher.Stream {
	return internal.NewChaCha20Stream(key, iv, 20)
}
