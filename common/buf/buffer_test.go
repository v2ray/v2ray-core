package buf_test

import (
	"testing"

	"v2ray.com/core/common"
	"v2ray.com/core/common/compare"

	. "v2ray.com/core/common/buf"
	"v2ray.com/core/common/serial"
	. "v2ray.com/ext/assert"
)

func TestBufferClear(t *testing.T) {
	assert := With(t)

	buffer := New()
	defer buffer.Release()

	payload := "Bytes"
	buffer.Write([]byte(payload))
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

func TestBufferSlice(t *testing.T) {
	{
		b := New()
		common.Must2(b.Write([]byte("abcd")))
		bytes := b.BytesFrom(-2)
		if err := compare.BytesEqualWithDetail(bytes, []byte{'c', 'd'}); err != nil {
			t.Error(err)
		}
	}

	{
		b := New()
		common.Must2(b.Write([]byte("abcd")))
		bytes := b.BytesTo(-2)
		if err := compare.BytesEqualWithDetail(bytes, []byte{'a', 'b'}); err != nil {
			t.Error(err)
		}
	}

	{
		b := New()
		common.Must2(b.Write([]byte("abcd")))
		bytes := b.BytesRange(-3, -1)
		if err := compare.BytesEqualWithDetail(bytes, []byte{'b', 'c'}); err != nil {
			t.Error(err)
		}
	}
}

func BenchmarkNewBuffer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		buffer := New()
		buffer.Release()
	}
}
