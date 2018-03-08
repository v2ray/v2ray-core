package buf_test

import (
	"crypto/rand"
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
	assert(buffer.Len(), Equals, len(payload))

	buffer.Clear()
	assert(buffer.Len(), Equals, 0)
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

func TestBufferWrite(t *testing.T) {
	assert := With(t)

	buffer := NewLocal(8)
	nBytes, err := buffer.Write([]byte("abcd"))
	assert(err, IsNil)
	assert(nBytes, Equals, 4)
	nBytes, err = buffer.Write([]byte("abcde"))
	assert(err, IsNil)
	assert(nBytes, Equals, 4)
	assert(buffer.String(), Equals, "abcdabcd")
}

func TestSyncPool(t *testing.T) {
	assert := With(t)

	p := NewPool(32)
	b := p.Allocate()
	assert(b.Len(), Equals, 0)

	assert(b.AppendSupplier(ReadFrom(rand.Reader)), IsNil)
	assert(b.Len(), Equals, 32)

	b.Release()
}

func BenchmarkNewBuffer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		buffer := New()
		buffer.Release()
	}
}

func BenchmarkNewLocalBuffer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		buffer := NewLocal(Size)
		buffer.Release()
	}
}
