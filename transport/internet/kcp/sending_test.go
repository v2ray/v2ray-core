package kcp_test

import (
	"testing"

	"v2ray.com/core/testing/assert"
	. "v2ray.com/core/transport/internet/kcp"
)

func TestSendingWindow(t *testing.T) {
	assert := assert.On(t)

	window := NewSendingWindow(5, nil, nil)
	window.Push(&DataSegment{
		Number: 0,
	})
	window.Push(&DataSegment{
		Number: 1,
	})
	window.Push(&DataSegment{
		Number: 2,
	})
	assert.Int(window.Len()).Equals(3)

	window.Remove(1)
	assert.Int(window.Len()).Equals(3)
	assert.Uint32(window.First().Number).Equals(0)

	window.Remove(0)
	assert.Int(window.Len()).Equals(1)
	assert.Uint32(window.First().Number).Equals(2)

	window.Remove(0)
	assert.Int(window.Len()).Equals(0)

	window.Push(&DataSegment{
		Number: 4,
	})
	assert.Int(window.Len()).Equals(1)
	assert.Uint32(window.First().Number).Equals(4)

	window.Push(&DataSegment{
		Number: 5,
	})
	assert.Int(window.Len()).Equals(2)

	window.Remove(1)
	assert.Int(window.Len()).Equals(2)

	window.Remove(0)
	assert.Int(window.Len()).Equals(0)
}
