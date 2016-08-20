package io_test

import (
	"testing"

	"v2ray.com/core/common/alloc"
	. "v2ray.com/core/common/io"
	"v2ray.com/core/testing/assert"
)

func TestBufferedReader(t *testing.T) {
	assert := assert.On(t)

	content := alloc.NewLargeBuffer()
	len := content.Len()

	reader := NewBufferedReader(content)
	assert.Bool(reader.Cached()).IsTrue()

	payload := make([]byte, 16)

	nBytes, err := reader.Read(payload)
	assert.Int(nBytes).Equals(16)
	assert.Error(err).IsNil()

	len2 := content.Len()
	assert.Int(len - len2).GreaterThan(16)

	nBytes, err = reader.Read(payload)
	assert.Int(nBytes).Equals(16)
	assert.Error(err).IsNil()

	assert.Int(content.Len()).Equals(len2)
	reader.SetCached(false)

	payload2 := alloc.NewBuffer()
	reader.Read(payload2.Value)

	assert.Int(content.Len()).Equals(len2)

	reader.Read(payload2.Value)
	assert.Int(content.Len()).LessThan(len2)
}
