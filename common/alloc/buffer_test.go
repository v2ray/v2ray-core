package alloc

import (
	"testing"

	v2testing "github.com/v2ray/v2ray-core/testing"
	"github.com/v2ray/v2ray-core/testing/assert"
)

func TestBufferClear(t *testing.T) {
	v2testing.Current(t)

	buffer := NewBuffer().Clear()
	defer buffer.Release()

	payload := "Bytes"
	buffer.Append([]byte(payload))
	assert.Int(buffer.Len()).Equals(len(payload))

	buffer.Clear()
	assert.Int(buffer.Len()).Equals(0)
}

func TestBufferIsFull(t *testing.T) {
	v2testing.Current(t)

	buffer := NewBuffer()
	defer buffer.Release()

	assert.Bool(buffer.IsFull()).IsTrue()

	buffer.Clear()
	assert.Bool(buffer.IsFull()).IsFalse()
}

func TestBufferPrepend(t *testing.T) {
	v2testing.Current(t)

	buffer := NewBuffer().Clear()
	defer buffer.Release()

	buffer.Append([]byte{'a', 'b', 'c'})
	buffer.Prepend([]byte{'x', 'y', 'z'})

	assert.Int(buffer.Len()).Equals(6)
	assert.Bytes(buffer.Value).Equals([]byte("xyzabc"))

	buffer.Prepend([]byte{'u', 'v', 'w'})
	assert.Bytes(buffer.Value).Equals([]byte("uvwxyzabc"))
}
