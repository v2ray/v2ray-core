package buf_test

import (
	"testing"

	. "v2ray.com/core/common/buf"
	. "v2ray.com/ext/assert"
)

func TestMultiBufferRead(t *testing.T) {
	assert := With(t)

	b1 := New()
	b1.AppendBytes('a', 'b')

	b2 := New()
	b2.AppendBytes('c', 'd')
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
	b.AppendBytes('a', 'b')
	mb.Append(b)
	assert(mb.Len(), Equals, 2)
}
