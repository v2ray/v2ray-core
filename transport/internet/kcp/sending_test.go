package kcp_test

import (
	"testing"

	"v2ray.com/core/testing/assert"
	. "v2ray.com/core/transport/internet/kcp"
)

func TestSendingWindow(t *testing.T) {
	assert := assert.On(t)

	window := NewSendingWindow(5, nil, nil)
	window.Push(0, []byte{})
	window.Push(1, []byte{})
	window.Push(2, []byte{})
	assert.Int(window.Len()).Equals(3)

	window.Remove(1)
	assert.Int(window.Len()).Equals(3)
	assert.Uint32(window.FirstNumber()).Equals(0)

	window.Remove(0)
	assert.Int(window.Len()).Equals(1)
	assert.Uint32(window.FirstNumber()).Equals(2)

	window.Remove(0)
	assert.Int(window.Len()).Equals(0)

	window.Push(4, []byte{})
	assert.Int(window.Len()).Equals(1)
	assert.Uint32(window.FirstNumber()).Equals(4)

	window.Push(5, []byte{})
	assert.Int(window.Len()).Equals(2)

	window.Remove(1)
	assert.Int(window.Len()).Equals(2)

	window.Remove(0)
	assert.Int(window.Len()).Equals(0)
}
