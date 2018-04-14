package uuid // import "v2ray.com/core/common/uuid"

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"

	"v2ray.com/core/common"
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

// New creates a UUID with random value.
func New() UUID {
	var uuid UUID
	common.Must2(rand.Read(uuid.Bytes()))
	return uuid
}

// ParseBytes converts a UUID in byte form to object.
func ParseBytes(b []byte) (UUID, error) {
	var uuid UUID
	if len(b) != 16 {
		return uuid, errors.New("invalid UUID: ", b)
	}
	copy(uuid[:], b)
	return uuid, nil
}

// ParseString converts a UUID in string form to object.
func ParseString(str string) (UUID, error) {
	var uuid UUID

	text := []byte(str)
	if len(text) < 32 {
		return uuid, errors.New("invalid UUID: ", str)
	}

	b := uuid.Bytes()

	for _, byteGroup := range byteGroups {
		if text[0] == '-' {
			text = text[1:]
		}

		if _, err := hex.Decode(b[:byteGroup/2], text[:byteGroup]); err != nil {
			return uuid, err
		}

		text = text[byteGroup:]
		b = b[byteGroup/2:]
	}

	return uuid, nil
}
