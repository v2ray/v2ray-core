package kcp_test

import (
	"testing"

	"github.com/v2ray/v2ray-core/common/alloc"
	"github.com/v2ray/v2ray-core/testing/assert"
	. "github.com/v2ray/v2ray-core/transport/internet/kcp"
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

func TestRecivingQueue(t *testing.T) {
	assert := assert.On(t)

	queue := NewReceivingQueue(2)
	queue.Put(alloc.NewSmallBuffer().Clear().AppendString("abcd"))
	queue.Put(alloc.NewSmallBuffer().Clear().AppendString("efg"))
	assert.Bool(queue.IsFull()).IsTrue()

	b := make([]byte, 1024)
	nBytes := queue.Read(b)
	assert.Int(nBytes).Equals(7)
	assert.String(string(b[:nBytes])).Equals("abcdefg")

	queue.Put(alloc.NewSmallBuffer().Clear().AppendString("1"))
	queue.Close()
	nBytes = queue.Read(b)
	assert.Int(nBytes).Equals(0)
}
