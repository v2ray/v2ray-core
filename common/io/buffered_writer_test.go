package io_test

import (
	"testing"

	"v2ray.com/core/common/alloc"
	. "v2ray.com/core/common/io"
	"v2ray.com/core/testing/assert"
)

func TestBufferedWriter(t *testing.T) {
	assert := assert.On(t)

	content := alloc.NewBuffer().Clear()

	writer := NewBufferedWriter(content)
	assert.Bool(writer.Cached()).IsTrue()

	payload := make([]byte, 16)

	nBytes, err := writer.Write(payload)
	assert.Int(nBytes).Equals(16)
	assert.Error(err).IsNil()

	assert.Bool(content.IsEmpty()).IsTrue()

	writer.SetCached(false)
	assert.Int(content.Len()).Equals(16)
}
