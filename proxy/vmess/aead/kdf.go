package aead

import (
	"crypto/hmac"
	"crypto/sha256"
	"hash"
)

func KDF(key []byte, path ...string) []byte {
	hmacf := hmac.New(func() hash.Hash {
		return sha256.New()
	}, []byte(KDFSaltConst_VMessAEADKDF))

	for _, v := range path {
		hmacf = hmac.New(func() hash.Hash {
			return hmacf
		}, []byte(v))
	}
	hmacf.Write(key)
	return hmacf.Sum(nil)
}

func KDF16(key []byte, path ...string) []byte {
	r := KDF(key, path...)
	return r[:16]
}
