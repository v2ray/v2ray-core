package protocol

import (
	"crypto/hmac"
	"crypto/md5"
	"hash"
)

func TimestampHash() hash.Hash {
	return md5.New()
}

func IDHash(key []byte) hash.Hash {
	return hmac.New(md5.New, key)
}
