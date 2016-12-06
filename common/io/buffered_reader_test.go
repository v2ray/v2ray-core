package io_test

import (
	"testing"

	"crypto/rand"
	"v2ray.com/core/common/alloc"
	. "v2ray.com/core/common/io"
	"v2ray.com/core/testing/assert"
)

func TestBufferedReader(t *testing.T) {
	assert := assert.On(t)

	content := alloc.NewBuffer()
	content.FillFrom(rand.Reader)

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
}
