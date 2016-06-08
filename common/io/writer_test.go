package io_test

import (
	"bytes"
	"crypto/rand"
	"testing"

	"github.com/v2ray/v2ray-core/common/alloc"
	. "github.com/v2ray/v2ray-core/common/io"
	"github.com/v2ray/v2ray-core/testing/assert"
)

func TestAdaptiveWriter(t *testing.T) {
	assert := assert.On(t)

	lb := alloc.NewLargeBuffer()
	rand.Read(lb.Value)

	writeBuffer := make([]byte, 0, 1024*1024)

	writer := NewAdaptiveWriter(NewBufferedWriter(bytes.NewBuffer(writeBuffer)))
	err := writer.Write(lb)
	assert.Error(err).IsNil()
	assert.Bytes(lb.Bytes()).Equals(writeBuffer)
}
