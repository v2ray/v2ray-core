package encoding

import (
	"crypto/md5"
	"hash/fnv"

	"v2ray.com/core/common"
	"v2ray.com/core/common/serial"

	"golang.org/x/crypto/sha3"
)

// Authenticate authenticates a byte array using Fnv hash.
func Authenticate(b []byte) uint32 {
	fnv1hash := fnv.New32a()
	common.Must2(fnv1hash.Write(b))
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
func (*FnvAuthenticator) NonceSize() int {
	return 0
}

// Overhead impelements AEAD.Overhead().
func (*FnvAuthenticator) Overhead() int {
	return 4
}

// Seal implements AEAD.Seal().
func (*FnvAuthenticator) Seal(dst, nonce, plaintext, additionalData []byte) []byte {
	dst = serial.Uint32ToBytes(Authenticate(plaintext), dst)
	return append(dst, plaintext...)
}

// Open implements AEAD.Open().
func (*FnvAuthenticator) Open(dst, nonce, ciphertext, additionalData []byte) ([]byte, error) {
	if serial.BytesToUint32(ciphertext[:4]) != Authenticate(ciphertext[4:]) {
		return dst, newError("invalid authentication")
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

type ShakeSizeParser struct {
	shake  sha3.ShakeHash
	buffer [2]byte
}

func NewShakeSizeParser(nonce []byte) *ShakeSizeParser {
	shake := sha3.NewShake128()
	common.Must2(shake.Write(nonce))
	return &ShakeSizeParser{
		shake: shake,
	}
}

func (*ShakeSizeParser) SizeBytes() int {
	return 2
}

func (s *ShakeSizeParser) next() uint16 {
	common.Must2(s.shake.Read(s.buffer[:]))
	return serial.BytesToUint16(s.buffer[:])
}

func (s *ShakeSizeParser) Decode(b []byte) (uint16, error) {
	mask := s.next()
	size := serial.BytesToUint16(b)
	return mask ^ size, nil
}

func (s *ShakeSizeParser) Encode(size uint16, b []byte) []byte {
	mask := s.next()
	return serial.Uint16ToBytes(mask^size, b[:0])
}
