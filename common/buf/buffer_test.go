package buf_test

import (
	"testing"

	. "v2ray.com/core/common/buf"
	"v2ray.com/core/common/serial"
	"v2ray.com/core/testing/assert"
)

func TestBufferClear(t *testing.T) {
	assert := assert.On(t)

	buffer := New()
	defer buffer.Release()

	payload := "Bytes"
	buffer.Append([]byte(payload))
	assert.Int(buffer.Len()).Equals(len(payload))

	buffer.Clear()
	assert.Int(buffer.Len()).Equals(0)
}

func TestBufferIsEmpty(t *testing.T) {
	assert := assert.On(t)

	buffer := New()
	defer buffer.Release()

	assert.Bool(buffer.IsEmpty()).IsTrue()
}

func TestBufferString(t *testing.T) {
	assert := assert.On(t)

	buffer := New()
	defer buffer.Release()

	assert.Error(buffer.AppendSupplier(serial.WriteString("Test String"))).IsNil()
	assert.String(buffer.String()).Equals("Test String")
}

func TestBufferWrite(t *testing.T) {
	assert := assert.On(t)

	buffer := NewLocal(8)
	nBytes, err := buffer.Write([]byte("abcd"))
	assert.Error(err).IsNil()
	assert.Int(nBytes).Equals(4)
	nBytes, err = buffer.Write([]byte("abcde"))
	assert.Error(err).IsNil()
	assert.Int(nBytes).Equals(4)
	assert.String(buffer.String()).Equals("abcdabcd")
}

func BenchmarkNewBuffer8192(b *testing.B) {
	for i := 0; i < b.N; i++ {
		buffer := New()
		buffer.Release()
	}
}

func BenchmarkNewLocalBuffer8192(b *testing.B) {
	for i := 0; i < b.N; i++ {
		buffer := NewLocal(8192)
		buffer.Release()
	}
}
