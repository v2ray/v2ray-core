package buf_test

import (
	"testing"

	. "v2ray.com/core/common/buf"
	"v2ray.com/core/testing/assert"
)

func TestMultiBufferRead(t *testing.T) {
	assert := assert.On(t)

	b1 := New()
	b1.AppendBytes('a', 'b')

	b2 := New()
	b2.AppendBytes('c', 'd')
	mb := NewMultiBufferValue(b1, b2)

	bs := make([]byte, 32)
	nBytes, err := mb.Read(bs)
	assert.Error(err).IsNil()
	assert.Int(nBytes).Equals(4)
	assert.Bytes(bs[:nBytes]).Equals([]byte("abcd"))
}

func TestMultiBufferAppend(t *testing.T) {
	assert := assert.On(t)

	var mb MultiBuffer
	b := New()
	b.AppendBytes('a', 'b')
	mb.Append(b)
	assert.Int(mb.Len()).Equals(2)
}
