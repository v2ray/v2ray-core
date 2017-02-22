package uuid

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"

	"v2ray.com/core/common/errors"
)

var (
	byteGroups = []int{8, 4, 4, 4, 12}
)

type UUID [16]byte

// String returns the string representation of this UUID.
func (u *UUID) String() string {
	bytes := u.Bytes()
	result := hex.EncodeToString(bytes[0 : byteGroups[0]/2])
	start := byteGroups[0] / 2
	for i := 1; i < len(byteGroups); i++ {
		nBytes := byteGroups[i] / 2
		result += "-"
		result += hex.EncodeToString(bytes[start : start+nBytes])
		start += nBytes
	}
	return result
}

// Bytes returns the bytes representation of this UUID.
func (u *UUID) Bytes() []byte {
	return u[:]
}

// Equals returns true if this UUID equals another UUID by value.
func (u *UUID) Equals(another *UUID) bool {
	if u == nil && another == nil {
		return true
	}
	if u == nil || another == nil {
		return false
	}
	return bytes.Equal(u.Bytes(), another.Bytes())
}

// Next generates a deterministic random UUID based on this UUID.
func (u *UUID) Next() *UUID {
	md5hash := md5.New()
	md5hash.Write(u.Bytes())
	md5hash.Write([]byte("16167dc8-16b6-4e6d-b8bb-65dd68113a81"))
	newid := new(UUID)
	for {
		md5hash.Sum(newid[:0])
		if !newid.Equals(u) {
			return newid
		}
		md5hash.Write([]byte("533eff8a-4113-4b10-b5ce-0f5d76b98cd2"))
	}
}

// New creates an UUID with random value.
func New() *UUID {
	uuid := new(UUID)
	rand.Read(uuid.Bytes())
	return uuid
}

// ParseBytes converts an UUID in byte form to object.
func ParseBytes(b []byte) (*UUID, error) {
	if len(b) != 16 {
		return nil, errors.New("Invalid UUID: ", b)
	}
	uuid := new(UUID)
	copy(uuid[:], b)
	return uuid, nil
}

// ParseString converts an UUID in string form to object.
func ParseString(str string) (*UUID, error) {
	text := []byte(str)
	if len(text) < 32 {
		return nil, errors.New("Invalid UUID: ", str)
	}

	uuid := new(UUID)
	b := uuid.Bytes()

	for _, byteGroup := range byteGroups {
		if text[0] == '-' {
			text = text[1:]
		}

		_, err := hex.Decode(b[:byteGroup/2], text[:byteGroup])

		if err != nil {
			return nil, err
		}

		text = text[byteGroup:]
		b = b[byteGroup/2:]
	}

	return uuid, nil
}
