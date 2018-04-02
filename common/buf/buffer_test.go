package buf_test

import (
	"testing"

	. "v2ray.com/core/common/buf"
	"v2ray.com/core/common/serial"
	. "v2ray.com/ext/assert"
)

func TestBufferClear(t *testing.T) {
	assert := With(t)

	buffer := New()
	defer buffer.Release()

	payload := "Bytes"
	buffer.Append([]byte(payload))
	assert(buffer.Len(), Equals, int32(len(payload)))

	buffer.Clear()
	assert(buffer.Len(), Equals, int32(0))
}

func TestBufferIsEmpty(t *testing.T) {
	assert := With(t)

	buffer := New()
	defer buffer.Release()

	assert(buffer.IsEmpty(), IsTrue)
}

func TestBufferString(t *testing.T) {
	assert := With(t)

	buffer := New()
	defer buffer.Release()

	assert(buffer.AppendSupplier(serial.WriteString("Test String")), IsNil)
	assert(buffer.String(), Equals, "Test String")
}

func BenchmarkNewBuffer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		buffer := New()
		buffer.Release()
	}
}

func BenchmarkNewLocalBuffer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		buffer := NewSize(Size)
		buffer.Release()
	}
}
