package alloc_test

import (
	"testing"

	. "v2ray.com/core/common/alloc"
	"v2ray.com/core/testing/assert"
)

func TestBufferClear(t *testing.T) {
	assert := assert.On(t)

	buffer := NewBuffer().Clear()
	defer buffer.Release()

	payload := "Bytes"
	buffer.Append([]byte(payload))
	assert.Int(buffer.Len()).Equals(len(payload))

	buffer.Clear()
	assert.Int(buffer.Len()).Equals(0)
}

func TestBufferIsFull(t *testing.T) {
	assert := assert.On(t)

	buffer := NewBuffer()
	defer buffer.Release()

	assert.Bool(buffer.IsFull()).IsTrue()

	buffer.Clear()
	assert.Bool(buffer.IsFull()).IsFalse()
}

func TestBufferPrepend(t *testing.T) {
	assert := assert.On(t)

	buffer := NewBuffer().Clear()
	defer buffer.Release()

	buffer.Append([]byte{'a', 'b', 'c'})
	buffer.Prepend([]byte{'x', 'y', 'z'})

	assert.Int(buffer.Len()).Equals(6)
	assert.Bytes(buffer.Value).Equals([]byte("xyzabc"))

	buffer.Prepend([]byte{'u', 'v', 'w'})
	assert.Bytes(buffer.Value).Equals([]byte("uvwxyzabc"))
}

func TestBufferString(t *testing.T) {
	assert := assert.On(t)

	buffer := NewBuffer().Clear()
	defer buffer.Release()

	buffer.AppendString("Test String")
	assert.String(buffer.String()).Equals("Test String")
}

func BenchmarkNewBuffer8192(b *testing.B) {
	for i := 0; i < b.N; i++ {
		buffer := NewBuffer()
		buffer.Release()
	}
}

func BenchmarkNewLocalBuffer8192(b *testing.B) {
	for i := 0; i < b.N; i++ {
		buffer := NewLocalBuffer(8192)
		buffer.Release()
	}
}

func BenchmarkNewBuffer2048(b *testing.B) {
	for i := 0; i < b.N; i++ {
		buffer := NewSmallBuffer()
		buffer.Release()
	}
}

func BenchmarkNewLocalBuffer2048(b *testing.B) {
	for i := 0; i < b.N; i++ {
		buffer := NewLocalBuffer(2048)
		buffer.Release()
	}
}
