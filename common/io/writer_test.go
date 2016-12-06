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
	lb.FillFrom(rand.Reader)

	expectedBytes := append([]byte(nil), lb.Bytes()...)

	writeBuffer := bytes.NewBuffer(make([]byte, 0, 1024*1024))

	writer := NewAdaptiveWriter(NewBufferedWriter(writeBuffer))
	err := writer.Write(lb)
	assert.Error(err).IsNil()
	assert.Bytes(expectedBytes).Equals(writeBuffer.Bytes())
}
