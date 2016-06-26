package kcp_test

import (
	"testing"

	"github.com/v2ray/v2ray-core/testing/assert"
	. "github.com/v2ray/v2ray-core/transport/internet/kcp"
)

func TestSendingQueue(t *testing.T) {
	assert := assert.On(t)

	queue := NewSendingQueue(3)

	seg0 := &Segment{}
	seg1 := &Segment{}
	seg2 := &Segment{}
	seg3 := &Segment{}

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

	seg0 := &Segment{}
	seg1 := &Segment{}
	seg2 := &Segment{}
	seg3 := &Segment{}

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
