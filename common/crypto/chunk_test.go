package crypto_test

import (
	"io"
	"testing"

	"v2ray.com/core/common/buf"
	. "v2ray.com/core/common/crypto"
	"v2ray.com/core/testing/assert"
)

func TestChunkStreamIO(t *testing.T) {
	assert := assert.On(t)

	cache := buf.NewLocal(8192)

	writer := NewChunkStreamWriter(PlainChunkSizeParser{}, cache)
	reader := NewChunkStreamReader(PlainChunkSizeParser{}, cache)

	b := buf.New()
	b.AppendBytes('a', 'b', 'c', 'd')
	assert.Error(writer.Write(buf.NewMultiBufferValue(b))).IsNil()

	b = buf.New()
	b.AppendBytes('e', 'f', 'g')
	assert.Error(writer.Write(buf.NewMultiBufferValue(b))).IsNil()

	assert.Error(writer.Write(buf.NewMultiBuffer())).IsNil()

	assert.Int(cache.Len()).Equals(13)

	mb, err := reader.Read()
	assert.Error(err).IsNil()
	assert.Int(mb.Len()).Equals(4)
	assert.Bytes(mb[0].Bytes()).Equals([]byte("abcd"))

	mb, err = reader.Read()
	assert.Error(err).IsNil()
	assert.Int(mb.Len()).Equals(3)
	assert.Bytes(mb[0].Bytes()).Equals([]byte("efg"))

	_, err = reader.Read()
	assert.Error(err).Equals(io.EOF)
}
