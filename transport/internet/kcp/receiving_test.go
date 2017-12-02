package kcp_test

import (
	"testing"

	. "v2ray.com/core/transport/internet/kcp"
	. "v2ray.com/ext/assert"
)

func TestRecivingWindow(t *testing.T) {
	assert := With(t)

	window := NewReceivingWindow(3)

	seg0 := &DataSegment{}
	seg1 := &DataSegment{}
	seg2 := &DataSegment{}
	seg3 := &DataSegment{}

	assert(window.Set(0, seg0), IsTrue)
	assert(window.RemoveFirst(), Equals, seg0)
	e := window.RemoveFirst()
	assert(e, IsNil)

	assert(window.Set(1, seg1), IsTrue)
	assert(window.Set(2, seg2), IsTrue)

	window.Advance()
	assert(window.Set(2, seg3), IsTrue)

	assert(window.RemoveFirst(), Equals, seg1)
	assert(window.Remove(1), Equals, seg2)
	assert(window.Remove(2), Equals, seg3)
}
