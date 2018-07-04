package mtproto

import (
	"crypto/rand"
	"crypto/sha256"
	"io"
	"sync"

	"v2ray.com/core/common"
)

const (
	HeaderSize = 64
)

type Authentication struct {
	Header        [HeaderSize]byte
	DecodingKey   [32]byte
	EncodingKey   [32]byte
	DecodingNonce [16]byte
	EncodingNonce [16]byte
}

func (a *Authentication) DataCenterID() uint16 {
	return (uint16(a.Header[61]) << 8) | uint16(a.Header[60])
}

func (a *Authentication) ApplySecret(b []byte) {
	a.DecodingKey = sha256.Sum256(append(a.DecodingKey[:], b...))
	a.EncodingKey = sha256.Sum256(append(a.EncodingKey[:], b...))
}

func generateRandomBytes(random []byte) {
	for {
		common.Must2(rand.Read(random[:]))

		if random[0] == 0xef {
			continue
		}

		val := (uint32(random[3]) << 24) | (uint32(random[2]) << 16) | (uint32(random[1]) << 8) | uint32(random[0])
		if val == 0x44414548 || val == 0x54534f50 || val == 0x20544547 || val == 0x4954504f || val == 0xeeeeeeee {
			continue
		}

		if 0x00000000 == (uint32(random[7])<<24)|(uint32(random[6])<<16)|(uint32(random[5])<<8)|uint32(random[4]) {
			continue
		}

		return
	}
}

func NewAuthentication() *Authentication {
	auth := getAuthenticationObject()
	random := auth.Header[:]
	generateRandomBytes(random)
	copy(auth.EncodingKey[:], random[8:])
	copy(auth.EncodingNonce[:], random[8+32:])
	keyivInverse := Inverse(random[8 : 8+32+16])
	copy(auth.DecodingKey[:], keyivInverse)
	copy(auth.DecodingNonce[:], keyivInverse[32:])
	return auth
}

func ReadAuthentication(reader io.Reader) (*Authentication, error) {
	auth := getAuthenticationObject()

	if _, err := io.ReadFull(reader, auth.Header[:]); err != nil {
		putAuthenticationObject(auth)
		return nil, err
	}

	copy(auth.DecodingKey[:], auth.Header[8:])
	copy(auth.DecodingNonce[:], auth.Header[8+32:])
	keyivInverse := Inverse(auth.Header[8 : 8+32+16])
	copy(auth.EncodingKey[:], keyivInverse)
	copy(auth.EncodingNonce[:], keyivInverse[32:])

	return auth, nil
}

// Inverse returns a new byte array. It is a sequence of bytes when the input is read from end to beginning.Inverse
// Visible for testing only.
func Inverse(b []byte) []byte {
	lenb := len(b)
	b2 := make([]byte, lenb)
	for i, v := range b {
		b2[lenb-i-1] = v
	}
	return b2
}

var (
	authPool = sync.Pool{
		New: func() interface{} {
			return new(Authentication)
		},
	}
)

func getAuthenticationObject() *Authentication {
	return authPool.Get().(*Authentication)
}

func putAuthenticationObject(auth *Authentication) {
	authPool.Put(auth)
}
