package kcp_test

import (
	"testing"

	"github.com/v2ray/v2ray-core/common/alloc"
	"github.com/v2ray/v2ray-core/testing/assert"
	. "github.com/v2ray/v2ray-core/transport/internet/kcp"
)

func TestSimpleAuthenticator(t *testing.T) {
	assert := assert.On(t)

	buffer := alloc.NewBuffer().Clear()
	buffer.AppendBytes('a', 'b', 'c', 'd', 'e', 'f', 'g')

	auth := NewSimpleAuthenticator()
	auth.Seal(buffer)

	assert.Bool(auth.Open(buffer)).IsTrue()
	assert.String(buffer.String()).Equals("abcdefg")
}
