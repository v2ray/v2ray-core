package kcp_test

import (
	"testing"

	"v2ray.com/core/common/alloc"
	"v2ray.com/core/testing/assert"
	. "v2ray.com/core/transport/internet/kcp"
)

func TestSendingQueue(t *testing.T) {
	assert := assert.On(t)

	queue := NewSendingQueue(3)

	seg0 := alloc.NewLocalBuffer(512)
	seg1 := alloc.NewLocalBuffer(512)
	seg2 := alloc.NewLocalBuffer(512)
	seg3 := alloc.NewLocalBuffer(512)

	assert.Bool(queue.IsEmpty()).IsTrue()
	assert.Bool(queue.IsFull()).IsFalse()

	queue.Push(seg0)
	assert.Bool(queue.IsEmpty()).IsFalse()

	queue.Push(seg1)
	queue.Push(seg2)

	assert.Bool(queue.IsFull()).IsTrue()

	assert.Pointer(queue.Pop()).Equals(seg0)

	queue.Push(seg3)
	assert.Bool(queue.IsFull()).IsTrue()

	assert.Pointer(queue.Pop()).Equals(seg1)
	assert.Pointer(queue.Pop()).Equals(seg2)
	assert.Pointer(queue.Pop()).Equals(seg3)
	assert.Int(int(queue.Len())).Equals(0)
}

func TestSendingQueueClear(t *testing.T) {
	assert := assert.On(t)

	queue := NewSendingQueue(3)

	seg0 := alloc.NewLocalBuffer(512)
	seg1 := alloc.NewLocalBuffer(512)
	seg2 := alloc.NewLocalBuffer(512)
	seg3 := alloc.NewLocalBuffer(512)

	queue.Push(seg0)
	assert.Bool(queue.IsEmpty()).IsFalse()

	queue.Clear()
	assert.Bool(queue.IsEmpty()).IsTrue()

	queue.Push(seg1)
	queue.Push(seg2)
	queue.Push(seg3)

	queue.Clear()
	assert.Bool(queue.IsEmpty()).IsTrue()
}

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
