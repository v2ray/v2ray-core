package uuid

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"errors"
)

var (
	byteGroups = []int{8, 4, 4, 4, 12}

	ErrorInvalidID = errors.New("Invalid ID.")
)

type UUID [16]byte

func (this *UUID) String() string {
	return bytesToString(this.Bytes())
}

func (this *UUID) Bytes() []byte {
	return this[:]
}

func (this *UUID) Equals(another *UUID) bool {
	if this == nil && another == nil {
		return true
	}
	if this == nil || another == nil {
		return false
	}
	return bytes.Equal(this.Bytes(), another.Bytes())
}

// Next generates a deterministic random UUID based on this UUID.
func (this *UUID) Next() *UUID {
	md5hash := md5.New()
	md5hash.Write(this.Bytes())
	md5hash.Write([]byte("16167dc8-16b6-4e6d-b8bb-65dd68113a81"))
	newid := new(UUID)
	for {
		md5hash.Sum(newid[:0])
		if !newid.Equals(this) {
			return newid
		}
		md5hash.Write([]byte("533eff8a-4113-4b10-b5ce-0f5d76b98cd2"))
	}
}

func bytesToString(bytes []byte) string {
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

func New() *UUID {
	uuid := new(UUID)
	rand.Read(uuid.Bytes())
	return uuid
}

func ParseBytes(b []byte) (*UUID, error) {
	if len(b) != 16 {
		return nil, ErrorInvalidID
	}
	uuid := new(UUID)
	copy(uuid[:], b)
	return uuid, nil
}

func ParseString(str string) (*UUID, error) {
	text := []byte(str)
	if len(text) < 32 {
		return nil, ErrorInvalidID
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
