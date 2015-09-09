package vmess

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	mrand "math/rand"
	"testing"
  
  "github.com/v2ray/v2ray-core/testing/unit"
)

func randomBytes(p []byte, t *testing.T) {
  assert := unit.Assert(t)
  
	nBytes, err := rand.Read(p)
  assert.Error(err).IsNil()
  assert.Int(nBytes).Named("# bytes of random buffer").Equals(len(p))
}

func TestNormalReading(t *testing.T) {
  assert := unit.Assert(t)
  
	testSize := 256
	plaintext := make([]byte, testSize)
	randomBytes(plaintext, t)

	keySize := 16
	key := make([]byte, keySize)
	randomBytes(key, t)
	iv := make([]byte, keySize)
	randomBytes(iv, t)

	aesBlock, err := aes.NewCipher(key)
  assert.Error(err).IsNil()
  
	aesMode := cipher.NewCBCEncrypter(aesBlock, iv)

	ciphertext := make([]byte, testSize)
	aesMode.CryptBlocks(ciphertext, plaintext)

	ciphertextcopy := make([]byte, testSize)
	copy(ciphertextcopy, ciphertext)

	reader, err := NewDecryptionReader(bytes.NewReader(ciphertextcopy), key, iv)
	assert.Error(err).IsNil()

	readtext := make([]byte, testSize)
	readSize := 0
	for readSize < testSize {
		nBytes := mrand.Intn(16) + 1
		if nBytes > testSize-readSize {
			nBytes = testSize - readSize
		}
		bytesRead, err := reader.Read(readtext[readSize : readSize+nBytes])
		assert.Error(err).IsNil()
    assert.Int(bytesRead).Equals(nBytes)
		readSize += nBytes
	}
  assert.Bytes(readtext).Named("Plaintext").Equals(plaintext)
}
