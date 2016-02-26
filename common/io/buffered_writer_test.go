package io_test

import (
	"testing"

	"github.com/v2ray/v2ray-core/common/alloc"
	. "github.com/v2ray/v2ray-core/common/io"
	v2testing "github.com/v2ray/v2ray-core/testing"
	"github.com/v2ray/v2ray-core/testing/assert"
)

func TestBufferedWriter(t *testing.T) {
	v2testing.Current(t)

	content := alloc.NewLargeBuffer().Clear()

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
