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
	b1.WriteString("ab")

	b2 := New()
	b2.WriteString("cd")
	mb := MultiBuffer{b1, b2}

	bs := make([]byte, 32)
	_, nBytes := SplitBytes(mb, bs)
	assert(nBytes, Equals, 4)
	assert(bs[:nBytes], Equals, []byte("abcd"))
}

func TestMultiBufferAppend(t *testing.T) {
	assert := With(t)

	var mb MultiBuffer
	b := New()
	b.WriteString("ab")
	mb = append(mb, b)
	assert(mb.Len(), Equals, int32(2))
}

func TestMultiBufferSliceBySizeLarge(t *testing.T) {
	assert := With(t)

	lb := make([]byte, 8*1024)
	common.Must2(io.ReadFull(rand.Reader, lb))

	mb := MergeBytes(nil, lb)

	mb, mb2 := SplitSize(mb, 1024)
	assert(mb2.Len(), Equals, int32(1024))
}
