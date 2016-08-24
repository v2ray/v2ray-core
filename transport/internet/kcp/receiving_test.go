package kcp_test

import (
	"testing"

	"v2ray.com/core/testing/assert"
	. "v2ray.com/core/transport/internet/kcp"
)

func TestRecivingWindow(t *testing.T) {
	assert := assert.On(t)

	window := NewReceivingWindow(3)

	seg0 := &DataSegment{}
	seg1 := &DataSegment{}
	seg2 := &DataSegment{}
	seg3 := &DataSegment{}

	assert.Bool(window.Set(0, seg0)).IsTrue()
	assert.Pointer(window.RemoveFirst()).Equals(seg0)
	e := window.RemoveFirst()
	if e != nil {
		assert.Fail("Expecting nil.")
	}

	assert.Bool(window.Set(1, seg1)).IsTrue()
	assert.Bool(window.Set(2, seg2)).IsTrue()

	window.Advance()
	assert.Bool(window.Set(2, seg3)).IsTrue()

	assert.Pointer(window.RemoveFirst()).Equals(seg1)
	assert.Pointer(window.Remove(1)).Equals(seg2)
	assert.Pointer(window.Remove(2)).Equals(seg3)
}
