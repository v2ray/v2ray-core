package user

import (
	"crypto/hmac"
	"crypto/md5"
)

type CounterHash interface {
	Hash(key []byte, counter int64) []byte
}

type StringHash interface {
	Hash(key []byte, data []byte) []byte
}

type TimeHash struct {
	baseHash StringHash
}

func NewTimeHash(baseHash StringHash) CounterHash {
	return TimeHash{
		baseHash: baseHash,
	}
}

func (h TimeHash) Hash(key []byte, counter int64) []byte {
	counterBytes := int64ToBytes(counter)
	return h.baseHash.Hash(key, counterBytes)
}

type HMACHash struct {
}

func (h HMACHash) Hash(key []byte, data []byte) []byte {
	hash := hmac.New(md5.New, key)
	hash.Write(data)
	return hash.Sum(nil)
}

func int64ToBytes(value int64) []byte {
	return []byte{
		byte(value >> 56),
		byte(value >> 48),
		byte(value >> 40),
		byte(value >> 32),
		byte(value >> 24),
		byte(value >> 16),
		byte(value >> 8),
		byte(value),
	}
}
