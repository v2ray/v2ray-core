package core

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

// The ID of en entity, in the form of an UUID.
type VID [16]byte

// Hash generates a MD5 hash based on current VID and a suffix string.
func (v VID) Hash(suffix []byte) []byte {
	md5 := md5.New()
	md5.Write(v[:])
	md5.Write(suffix)
	return md5.Sum(nil)
}

var byteGroups = []int{8, 4, 4, 4, 12}

// TODO: leverage a full functional UUID library
func UUIDToVID(uuid string) (v VID, err error) {
	text := []byte(uuid)
	if len(text) < 32 {
		err = fmt.Errorf("uuid: invalid UUID string: %s", text)
		return
	}

	b := v[:]

	for _, byteGroup := range byteGroups {
		if text[0] == '-' {
			text = text[1:]
		}

		_, err = hex.Decode(b[:byteGroup/2], text[:byteGroup])

		if err != nil {
			return
		}

		text = text[byteGroup:]
		b = b[byteGroup/2:]
	}

	return
}
