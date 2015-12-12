package uuid

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
)

var (
	byteGroups = []int{8, 4, 4, 4, 12}

	InvalidID = errors.New("Invalid ID.")
)

type UUID struct {
	byteValue   []byte
	stringValue string
}

func (this *UUID) String() string {
	return this.stringValue
}

func (this *UUID) Bytes() []byte {
	return this.byteValue[:]
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
	bytes := make([]byte, 16)
	rand.Read(bytes)
	uuid, _ := ParseBytes(bytes)
	return uuid
}

func ParseBytes(bytes []byte) (*UUID, error) {
	if len(bytes) != 16 {
		return nil, InvalidID
	}
	return &UUID{
		byteValue:   bytes,
		stringValue: bytesToString(bytes),
	}, nil
}

func ParseString(str string) (*UUID, error) {
	text := []byte(str)
	if len(text) < 32 {
		return nil, InvalidID
	}

	uuid := &UUID{
		byteValue:   make([]byte, 16),
		stringValue: str,
	}
	b := uuid.byteValue[:]

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
