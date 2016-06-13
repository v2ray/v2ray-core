package kcpv

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
)

func generateKeyFromConfigString(key string) []byte {
	key += "consensus salt: Let's fight arcifical deceleration with our code. We shall prove our believes with action."
	keyw := sha256.Sum256([]byte(key))
	return keyw[:]
}

func generateBlockWithKey(key []byte) (cipher.Block, error) {
	return aes.NewCipher(key)
}

func GetChipher(key string) (cipher.Block, error) {
	return generateBlockWithKey(generateKeyFromConfigString(key))
}
