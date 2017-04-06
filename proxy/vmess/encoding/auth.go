package encoding

import (
	"crypto/md5"
	"hash/fnv"

	"v2ray.com/core/common/errors"
	"v2ray.com/core/common/serial"
)

// Authenticate authenticates a byte array using Fnv hash.
func Authenticate(b []byte) uint32 {
	fnv1hash := fnv.New32a()
	fnv1hash.Write(b)
	return fnv1hash.Sum32()
}

type NoOpAuthenticator struct{}

func (NoOpAuthenticator) NonceSize() int {
	return 0
}

func (NoOpAuthenticator) Overhead() int {
	return 0
}

// Seal implements AEAD.Seal().
func (NoOpAuthenticator) Seal(dst, nonce, plaintext, additionalData []byte) []byte {
	return append(dst[:0], plaintext...)
}

// Open implements AEAD.Open().
func (NoOpAuthenticator) Open(dst, nonce, ciphertext, additionalData []byte) ([]byte, error) {
	return append(dst[:0], ciphertext...), nil
}

// FnvAuthenticator is an AEAD based on Fnv hash.
type FnvAuthenticator struct {
}

// NonceSize implements AEAD.NonceSize().
func (v *FnvAuthenticator) NonceSize() int {
	return 0
}

// Overhead impelements AEAD.Overhead().
func (v *FnvAuthenticator) Overhead() int {
	return 4
}

// Seal implements AEAD.Seal().
func (v *FnvAuthenticator) Seal(dst, nonce, plaintext, additionalData []byte) []byte {
	dst = serial.Uint32ToBytes(Authenticate(plaintext), dst)
	return append(dst, plaintext...)
}

// Open implements AEAD.Open().
func (v *FnvAuthenticator) Open(dst, nonce, ciphertext, additionalData []byte) ([]byte, error) {
	if serial.BytesToUint32(ciphertext[:4]) != Authenticate(ciphertext[4:]) {
		return dst, errors.New("invalid authentication").Path("Proxy", "VMess", "Encoding", "FnvAuthenticator")
	}
	return append(dst, ciphertext[4:]...), nil
}

// GenerateChacha20Poly1305Key generates a 32-byte key from a given 16-byte array.
func GenerateChacha20Poly1305Key(b []byte) []byte {
	key := make([]byte, 32)
	t := md5.Sum(b)
	copy(key, t[:])
	t = md5.Sum(key[:16])
	copy(key[16:], t[:])
	return key
}
