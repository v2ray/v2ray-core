package io_test

import (
	"bytes"
	"testing"

	"github.com/v2ray/v2ray-core/common/alloc"
	. "github.com/v2ray/v2ray-core/common/io"
	v2testing "github.com/v2ray/v2ray-core/testing"
	"github.com/v2ray/v2ray-core/testing/assert"
)

func TestAdaptiveReader(t *testing.T) {
	v2testing.Current(t)

	rawContent := make([]byte, 1024*1024)

	reader := NewAdaptiveReader(bytes.NewBuffer(rawContent))
	b1, err := reader.Read()
	assert.Error(err).IsNil()
	assert.Bool(b1.IsFull()).IsTrue()
	assert.Int(b1.Len()).Equals(alloc.BufferSize)

	b2, err := reader.Read()
	assert.Error(err).IsNil()
	assert.Bool(b2.IsFull()).IsTrue()
	assert.Int(b2.Len()).Equals(alloc.LargeBufferSize)
}
