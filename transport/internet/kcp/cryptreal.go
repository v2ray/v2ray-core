package kcp

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"v2ray.com/core/common"
)

func NewAEADAESGCMBasedOnSeed(seed string) cipher.AEAD {
	HashedSeed := sha256.Sum256([]byte(seed))
	aesBlock := common.Must2(aes.NewCipher(HashedSeed[:16])).(cipher.Block)
	return common.Must2(cipher.NewGCM(aesBlock)).(cipher.AEAD)
}
