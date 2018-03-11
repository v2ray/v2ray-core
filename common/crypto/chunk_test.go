package crypto_test

import (
	"io"
	"testing"

	"v2ray.com/core/common/buf"
	. "v2ray.com/core/common/crypto"
	. "v2ray.com/ext/assert"
)

func TestChunkStreamIO(t *testing.T) {
	assert := With(t)

	cache := buf.NewSize(8192)

	writer := NewChunkStreamWriter(PlainChunkSizeParser{}, cache)
	reader := NewChunkStreamReader(PlainChunkSizeParser{}, cache)

	b := buf.New()
	b.AppendBytes('a', 'b', 'c', 'd')
	assert(writer.WriteMultiBuffer(buf.NewMultiBufferValue(b)), IsNil)

	b = buf.New()
	b.AppendBytes('e', 'f', 'g')
	assert(writer.WriteMultiBuffer(buf.NewMultiBufferValue(b)), IsNil)

	assert(writer.WriteMultiBuffer(buf.MultiBuffer{}), IsNil)

	assert(cache.Len(), Equals, 13)

	mb, err := reader.ReadMultiBuffer()
	assert(err, IsNil)
	assert(mb.Len(), Equals, 4)
	assert(mb[0].Bytes(), Equals, []byte("abcd"))

	mb, err = reader.ReadMultiBuffer()
	assert(err, IsNil)
	assert(mb.Len(), Equals, 3)
	assert(mb[0].Bytes(), Equals, []byte("efg"))

	_, err = reader.ReadMultiBuffer()
	assert(err, Equals, io.EOF)
}
