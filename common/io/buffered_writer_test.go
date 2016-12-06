package io_test

import (
	"crypto/rand"
	"testing"

	"v2ray.com/core/common/alloc"
	. "v2ray.com/core/common/io"
	"v2ray.com/core/testing/assert"
)

func TestBufferedWriter(t *testing.T) {
	assert := assert.On(t)

	content := alloc.NewBuffer()

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

func TestBufferedWriterLargePayload(t *testing.T) {
	assert := assert.On(t)

	content := alloc.NewLocalBuffer(128 * 1024)

	writer := NewBufferedWriter(content)
	assert.Bool(writer.Cached()).IsTrue()

	payload := make([]byte, 64*1024)
	rand.Read(payload)

	nBytes, err := writer.Write(payload[:1024])
	assert.Int(nBytes).Equals(1024)
	assert.Error(err).IsNil()

	assert.Bool(content.IsEmpty()).IsTrue()

	nBytes, err = writer.Write(payload[1024:])
	assert.Error(err).IsNil()
	assert.Int(nBytes).Equals(63 * 1024)
	assert.Bytes(content.Bytes()).Equals(payload)
}
