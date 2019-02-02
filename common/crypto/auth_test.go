package crypto_test

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"testing"

	"v2ray.com/core/common"
	"v2ray.com/core/common/buf"
	. "v2ray.com/core/common/crypto"
	"v2ray.com/core/common/protocol"
	. "v2ray.com/ext/assert"
)

func TestAuthenticationReaderWriter(t *testing.T) {
	assert := With(t)

	key := make([]byte, 16)
	rand.Read(key)
	block, err := aes.NewCipher(key)
	common.Must(err)

	aead, err := cipher.NewGCM(block)
	common.Must(err)

	const payloadSize = 1024 * 80
	rawPayload := make([]byte, payloadSize)
	rand.Read(rawPayload)

	payload := buf.MergeBytes(nil, rawPayload)
	assert(payload.Len(), Equals, int32(payloadSize))

	cache := bytes.NewBuffer(nil)
	iv := make([]byte, 12)
	rand.Read(iv)

	writer := NewAuthenticationWriter(&AEADAuthenticator{
		AEAD:                    aead,
		NonceGenerator:          GenerateStaticBytes(iv),
		AdditionalDataGenerator: GenerateEmptyBytes(),
	}, PlainChunkSizeParser{}, cache, protocol.TransferTypeStream, nil)

	assert(writer.WriteMultiBuffer(payload), IsNil)
	assert(cache.Len(), Equals, int(82658))
	assert(writer.WriteMultiBuffer(buf.MultiBuffer{}), IsNil)

	reader := NewAuthenticationReader(&AEADAuthenticator{
		AEAD:                    aead,
		NonceGenerator:          GenerateStaticBytes(iv),
		AdditionalDataGenerator: GenerateEmptyBytes(),
	}, PlainChunkSizeParser{}, cache, protocol.TransferTypeStream, nil)

	var mb buf.MultiBuffer

	for mb.Len() < payloadSize {
		mb2, err := reader.ReadMultiBuffer()
		common.Must(err)

		mb, _ = buf.MergeMulti(mb, mb2)
	}

	assert(mb.Len(), Equals, int32(payloadSize))

	mbContent := make([]byte, payloadSize)
	buf.SplitBytes(mb, mbContent)
	assert(mbContent, Equals, rawPayload)

	_, err = reader.ReadMultiBuffer()
	assert(err, Equals, io.EOF)
}

func TestAuthenticationReaderWriterPacket(t *testing.T) {
	assert := With(t)

	key := make([]byte, 16)
	common.Must2(rand.Read(key))
	block, err := aes.NewCipher(key)
	common.Must(err)

	aead, err := cipher.NewGCM(block)
	common.Must(err)

	cache := buf.New()
	iv := make([]byte, 12)
	rand.Read(iv)

	writer := NewAuthenticationWriter(&AEADAuthenticator{
		AEAD:                    aead,
		NonceGenerator:          GenerateStaticBytes(iv),
		AdditionalDataGenerator: GenerateEmptyBytes(),
	}, PlainChunkSizeParser{}, cache, protocol.TransferTypePacket, nil)

	var payload buf.MultiBuffer
	pb1 := buf.New()
	pb1.Write([]byte("abcd"))
	payload = append(payload, pb1)

	pb2 := buf.New()
	pb2.Write([]byte("efgh"))
	payload = append(payload, pb2)

	assert(writer.WriteMultiBuffer(payload), IsNil)
	assert(cache.Len(), GreaterThan, int32(0))
	assert(writer.WriteMultiBuffer(buf.MultiBuffer{}), IsNil)
	common.Must(err)

	reader := NewAuthenticationReader(&AEADAuthenticator{
		AEAD:                    aead,
		NonceGenerator:          GenerateStaticBytes(iv),
		AdditionalDataGenerator: GenerateEmptyBytes(),
	}, PlainChunkSizeParser{}, cache, protocol.TransferTypePacket, nil)

	mb, err := reader.ReadMultiBuffer()
	common.Must(err)

	mb, b1 := buf.SplitFirst(mb)
	assert(b1.String(), Equals, "abcd")

	mb, b2 := buf.SplitFirst(mb)
	assert(b2.String(), Equals, "efgh")

	assert(mb.IsEmpty(), IsTrue)

	_, err = reader.ReadMultiBuffer()
	assert(err, Equals, io.EOF)
}
