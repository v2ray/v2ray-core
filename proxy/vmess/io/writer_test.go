package io_test

import (
	"testing"

	"github.com/v2ray/v2ray-core/common/alloc"
	. "github.com/v2ray/v2ray-core/proxy/vmess/io"
	v2testing "github.com/v2ray/v2ray-core/testing"
	"github.com/v2ray/v2ray-core/testing/assert"
)

func TestAuthenticate(t *testing.T) {
	v2testing.Current(t)

	buffer := alloc.NewBuffer().Clear()
	buffer.AppendBytes(1, 2, 3, 4)
	Authenticate(buffer)
	assert.Bytes(buffer.Value).Equals([]byte{0, 8, 87, 52, 168, 125, 1, 2, 3, 4})

	b2, err := NewAuthChunkReader(buffer).Read()
	assert.Error(err).IsNil()
	assert.Bytes(b2.Value).Equals([]byte{1, 2, 3, 4})
}
