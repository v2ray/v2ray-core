package kcp_test

import (
	"testing"

	. "v2ray.com/core/transport/internet/kcp"
	. "v2ray.com/ext/assert"
)

func TestSendingWindow(t *testing.T) {
	assert := With(t)

	window := NewSendingWindow(5, nil, nil)
	window.Push(0, []byte{})
	window.Push(1, []byte{})
	window.Push(2, []byte{})
	assert(window.Len(), Equals, 3)

	window.Remove(1)
	assert(window.Len(), Equals, 3)
	assert(window.FirstNumber(), Equals, uint32(0))

	window.Remove(0)
	assert(window.Len(), Equals, 1)
	assert(window.FirstNumber(), Equals, uint32(2))

	window.Remove(0)
	assert(window.Len(), Equals, 0)

	window.Push(4, []byte{})
	assert(window.Len(), Equals, 1)
	assert(window.FirstNumber(), Equals, uint32(4))

	window.Push(5, []byte{})
	assert(window.Len(), Equals, 2)

	window.Remove(1)
	assert(window.Len(), Equals, 2)

	window.Remove(0)
	assert(window.Len(), Equals, 0)
}
