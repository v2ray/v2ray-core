package crypto_test

import (
	"crypto/cipher"
	"testing"

	. "v2ray.com/core/common/crypto"
)

const benchSize = 1024 * 1024

func benchmarkStream(b *testing.B, c cipher.Stream) {
	b.SetBytes(benchSize)
	input := make([]byte, benchSize)
	output := make([]byte, benchSize)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.XORKeyStream(output, input)
	}
}

func BenchmarkChaCha20(b *testing.B) {
	key := make([]byte, 32)
	nonce := make([]byte, 8)
	c := NewChaCha20Stream(key, nonce)
	benchmarkStream(b, c)
}

func BenchmarkChaCha20IETF(b *testing.B) {
	key := make([]byte, 32)
	nonce := make([]byte, 12)
	c := NewChaCha20Stream(key, nonce)
	benchmarkStream(b, c)
}

func BenchmarkAESEncryption(b *testing.B) {
	key := make([]byte, 32)
	iv := make([]byte, 16)
	c := NewAesEncryptionStream(key, iv)

	benchmarkStream(b, c)
}

func BenchmarkAESDecryption(b *testing.B) {
	key := make([]byte, 32)
	iv := make([]byte, 16)
	c := NewAesDecryptionStream(key, iv)

	benchmarkStream(b, c)
}
