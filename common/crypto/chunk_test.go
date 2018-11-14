package crypto_test

import (
	"bytes"
	"io"
	"testing"

	"v2ray.com/core/common"
	"v2ray.com/core/common/buf"
	. "v2ray.com/core/common/crypto"
	. "v2ray.com/ext/assert"
)

func TestChunkStreamIO(t *testing.T) {
	assert := With(t)

	cache := bytes.NewBuffer(make([]byte, 0, 8192))

	writer := NewChunkStreamWriter(PlainChunkSizeParser{}, cache)
	reader := NewChunkStreamReader(PlainChunkSizeParser{}, cache)

	b := buf.New()
	b.WriteString("abcd")
	common.Must(writer.WriteMultiBuffer(buf.NewMultiBufferValue(b)))

	b = buf.New()
	b.WriteString("efg")
	common.Must(writer.WriteMultiBuffer(buf.NewMultiBufferValue(b)))

	common.Must(writer.WriteMultiBuffer(buf.MultiBuffer{}))

	if cache.Len() != 13 {
		t.Fatalf("Cache length is %d, want 13", cache.Len())
	}

	mb, err := reader.ReadMultiBuffer()
	common.Must(err)

	assert(mb.Len(), Equals, int32(4))
	assert(mb[0].Bytes(), Equals, []byte("abcd"))

	mb, err = reader.ReadMultiBuffer()
	assert(err, IsNil)
	assert(mb.Len(), Equals, int32(3))
	assert(mb[0].Bytes(), Equals, []byte("efg"))

	_, err = reader.ReadMultiBuffer()
	assert(err, Equals, io.EOF)
}
