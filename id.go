package core

import (
	"crypto/hmac"
	"crypto/md5"
	"encoding/hex"
	"hash"
	mrand "math/rand"
	"time"

	"github.com/v2ray/v2ray-core/log"
)

const (
	IDBytesLen = 16
)

// The ID of en entity, in the form of an UUID.
type ID struct {
	String string
	Bytes  []byte
	hasher hash.Hash
}

func NewID(id string) (ID, error) {
	idBytes, err := UUIDToID(id)
	if err != nil {
		return ID{}, log.Error("Failed to parse id %s", id)
	}
	hasher := hmac.New(md5.New, idBytes)
	return ID{id, idBytes, hasher}, nil
}

func (v ID) TimeRangeHash(rangeSec int) []byte {
	nowSec := time.Now().UTC().Unix()
	delta := mrand.Intn(rangeSec*2) - rangeSec

	targetSec := nowSec + int64(delta)
	return v.TimeHash(targetSec)
}

func (v ID) TimeHash(timeSec int64) []byte {
	buffer := []byte{
		byte(timeSec >> 56),
		byte(timeSec >> 48),
		byte(timeSec >> 40),
		byte(timeSec >> 32),
		byte(timeSec >> 24),
		byte(timeSec >> 16),
		byte(timeSec >> 8),
		byte(timeSec),
	}
	return v.Hash(buffer)
}

func (v ID) Hash(data []byte) []byte {
	return v.hasher.Sum(data)
}

var byteGroups = []int{8, 4, 4, 4, 12}

// TODO: leverage a full functional UUID library
func UUIDToID(uuid string) (v []byte, err error) {
	v = make([]byte, 16)

	text := []byte(uuid)
	if len(text) < 32 {
		err = log.Error("uuid: invalid UUID string: %s", text)
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
