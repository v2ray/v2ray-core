package crypto_test

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"testing"

	"v2ray.com/core/common/buf"
	. "v2ray.com/core/common/crypto"
	"v2ray.com/core/testing/assert"
)

func TestAuthenticationReaderWriter(t *testing.T) {
	assert := assert.On(t)

	key := make([]byte, 16)
	rand.Read(key)
	block, err := aes.NewCipher(key)
	assert.Error(err).IsNil()

	aead, err := cipher.NewGCM(block)
	assert.Error(err).IsNil()

	rawPayload := make([]byte, 8192*10)
	rand.Read(rawPayload)

	payload := buf.NewLocal(8192 * 10)
	payload.Append(rawPayload)

	cache := buf.NewLocal(160 * 1024)
	iv := make([]byte, 12)
	rand.Read(iv)

	writer := NewAuthenticationWriter(&AEADAuthenticator{
		AEAD: aead,
		NonceGenerator: &StaticBytesGenerator{
			Content: iv,
		},
		AdditionalDataGenerator: &NoOpBytesGenerator{},
	}, PlainChunkSizeParser{}, cache, ModeStream)

	assert.Error(writer.Write(buf.NewMultiBufferValue(payload))).IsNil()
	assert.Int(cache.Len()).Equals(83360)
	assert.Error(writer.Write(buf.NewMultiBuffer())).IsNil()
	assert.Error(err).IsNil()

	reader := NewAuthenticationReader(&AEADAuthenticator{
		AEAD: aead,
		NonceGenerator: &StaticBytesGenerator{
			Content: iv,
		},
		AdditionalDataGenerator: &NoOpBytesGenerator{},
	}, PlainChunkSizeParser{}, cache, ModeStream)

	mb := buf.NewMultiBuffer()

	for mb.Len() < len(rawPayload) {
		mb2, err := reader.Read()
		assert.Error(err).IsNil()

		mb.AppendMulti(mb2)
	}

	mbContent := make([]byte, 8192*10)
	mb.Read(mbContent)
	assert.Bytes(mbContent).Equals(rawPayload)

	_, err = reader.Read()
	assert.Error(err).Equals(io.EOF)
}

func TestAuthenticationReaderWriterPacket(t *testing.T) {
	assert := assert.On(t)

	key := make([]byte, 16)
	rand.Read(key)
	block, err := aes.NewCipher(key)
	assert.Error(err).IsNil()

	aead, err := cipher.NewGCM(block)
	assert.Error(err).IsNil()

	cache := buf.NewLocal(1024)
	iv := make([]byte, 12)
	rand.Read(iv)

	writer := NewAuthenticationWriter(&AEADAuthenticator{
		AEAD: aead,
		NonceGenerator: &StaticBytesGenerator{
			Content: iv,
		},
		AdditionalDataGenerator: &NoOpBytesGenerator{},
	}, PlainChunkSizeParser{}, cache, ModePacket)

	payload := buf.NewMultiBuffer()
	pb1 := buf.New()
	pb1.Append([]byte("abcd"))
	payload.Append(pb1)

	pb2 := buf.New()
	pb2.Append([]byte("efgh"))
	payload.Append(pb2)

	assert.Error(writer.Write(payload)).IsNil()
	assert.Int(cache.Len()).GreaterThan(0)
	assert.Error(writer.Write(buf.NewMultiBuffer())).IsNil()
	assert.Error(err).IsNil()

	reader := NewAuthenticationReader(&AEADAuthenticator{
		AEAD: aead,
		NonceGenerator: &StaticBytesGenerator{
			Content: iv,
		},
		AdditionalDataGenerator: &NoOpBytesGenerator{},
	}, PlainChunkSizeParser{}, cache, ModePacket)

	mb, err := reader.Read()
	assert.Error(err).IsNil()

	b1 := mb.SplitFirst()
	assert.String(b1.String()).Equals("abcd")
	b2 := mb.SplitFirst()
	assert.String(b2.String()).Equals("efgh")
	assert.Bool(mb.IsEmpty()).IsTrue()

	_, err = reader.Read()
	assert.Error(err).Equals(io.EOF)
}
