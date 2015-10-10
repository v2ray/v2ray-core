package alloc

import (
	"testing"

	"github.com/v2ray/v2ray-core/testing/unit"
)

func TestBufferClear(t *testing.T) {
	assert := unit.Assert(t)

	buffer := NewBuffer().Clear()
	defer buffer.Release()

	payload := "Bytes"
	buffer.Append([]byte(payload))
	assert.Int(buffer.Len()).Equals(len(payload))

	buffer.Clear()
	assert.Int(buffer.Len()).Equals(0)
}

func TestBufferIsFull(t *testing.T) {
	assert := unit.Assert(t)

	buffer := NewBuffer()
	defer buffer.Release()

	assert.Bool(buffer.IsFull()).IsTrue()

	buffer.Clear()
	assert.Bool(buffer.IsFull()).IsFalse()
}
