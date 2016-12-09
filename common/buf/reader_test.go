package buf_test

import (
	"bytes"
	"testing"

	. "v2ray.com/core/common/buf"
	"v2ray.com/core/testing/assert"
)

func TestAdaptiveReader(t *testing.T) {
	assert := assert.On(t)

	rawContent := make([]byte, 1024*1024)
	buffer := bytes.NewBuffer(rawContent)

	reader := NewReader(buffer)
	b1, err := reader.Read()
	assert.Error(err).IsNil()
	assert.Bool(b1.IsFull()).IsTrue()
	assert.Int(b1.Len()).Equals(Size)
	assert.Int(buffer.Len()).Equals(cap(rawContent) - Size)

	b2, err := reader.Read()
	assert.Error(err).IsNil()
	assert.Bool(b2.IsFull()).IsTrue()
	assert.Int(buffer.Len()).Equals(1007616)
}
