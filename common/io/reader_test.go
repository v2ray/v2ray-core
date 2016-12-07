package io_test

import (
	"bytes"
	"testing"

	"v2ray.com/core/common/alloc"
	. "v2ray.com/core/common/io"
	"v2ray.com/core/testing/assert"
)

func TestAdaptiveReader(t *testing.T) {
	assert := assert.On(t)

	rawContent := make([]byte, 1024*1024)
	buffer := bytes.NewBuffer(rawContent)

	reader := NewAdaptiveReader(buffer)
	b1, err := reader.Read()
	assert.Error(err).IsNil()
	assert.Bool(b1.IsFull()).IsTrue()
	assert.Int(b1.Len()).Equals(alloc.BufferSize)
	assert.Int(buffer.Len()).Equals(cap(rawContent) - alloc.BufferSize)

	b2, err := reader.Read()
	assert.Error(err).IsNil()
	assert.Bool(b2.IsFull()).IsTrue()
	assert.Int(buffer.Len()).Equals(1007648)
}
