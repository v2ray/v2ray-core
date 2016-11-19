package io_test

import (
	"bytes"
	"crypto/rand"
	"testing"

	"v2ray.com/core/common/alloc"
	. "v2ray.com/core/common/io"
	"v2ray.com/core/testing/assert"
)

func TestAdaptiveWriter(t *testing.T) {
	assert := assert.On(t)

	lb := alloc.NewBuffer()
	rand.Read(lb.Value)

	writeBuffer := make([]byte, 0, 1024*1024)

	writer := NewAdaptiveWriter(NewBufferedWriter(bytes.NewBuffer(writeBuffer)))
	err := writer.Write(lb)
	assert.Error(err).IsNil()
	assert.Bytes(lb.Bytes()).Equals(writeBuffer)
}
