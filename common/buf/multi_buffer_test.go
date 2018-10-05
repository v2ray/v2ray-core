package buf_test

import (
	"crypto/rand"
	"io"
	"testing"

	"v2ray.com/core/common"
	. "v2ray.com/core/common/buf"
	. "v2ray.com/ext/assert"
)

func TestMultiBufferRead(t *testing.T) {
	assert := With(t)

	b1 := New()
	b1.WriteBytes('a', 'b')

	b2 := New()
	b2.WriteBytes('c', 'd')
	mb := NewMultiBufferValue(b1, b2)

	bs := make([]byte, 32)
	nBytes, err := mb.Read(bs)
	assert(err, IsNil)
	assert(nBytes, Equals, 4)
	assert(bs[:nBytes], Equals, []byte("abcd"))
}

func TestMultiBufferAppend(t *testing.T) {
	assert := With(t)

	var mb MultiBuffer
	b := New()
	b.WriteBytes('a', 'b')
	mb.Append(b)
	assert(mb.Len(), Equals, int32(2))
}

func TestMultiBufferSliceBySizeLarge(t *testing.T) {
	assert := With(t)

	lb := make([]byte, 8*1024)
	common.Must2(io.ReadFull(rand.Reader, lb))

	var mb MultiBuffer
	common.Must2(mb.Write(lb))

	mb2 := mb.SliceBySize(1024)
	assert(mb2.Len(), Equals, int32(1024))
}

func TestInterface(t *testing.T) {
	assert := With(t)

	assert((*MultiBuffer)(nil), Implements, (*io.WriterTo)(nil))
	assert((*MultiBuffer)(nil), Implements, (*io.ReaderFrom)(nil))
}
