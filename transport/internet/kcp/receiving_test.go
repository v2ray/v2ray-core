package kcp_test

import (
	"io"
	"testing"
	"time"

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
	assert.Bool(queue.Put(alloc.NewSmallBuffer().Clear().AppendString("abcd"))).IsTrue()
	assert.Bool(queue.Put(alloc.NewSmallBuffer().Clear().AppendString("efg"))).IsTrue()
	assert.Bool(queue.Put(alloc.NewSmallBuffer().Clear().AppendString("more content"))).IsFalse()

	b := make([]byte, 1024)
	nBytes, err := queue.Read(b)
	assert.Error(err).IsNil()
	assert.Int(nBytes).Equals(7)
	assert.String(string(b[:nBytes])).Equals("abcdefg")

	assert.Bool(queue.Put(alloc.NewSmallBuffer().Clear().AppendString("1"))).IsTrue()
	queue.Close()
	nBytes, err = queue.Read(b)
	assert.Error(err).Equals(io.EOF)
}

func TestRecivingQueueTimeout(t *testing.T) {
	assert := assert.On(t)

	queue := NewReceivingQueue(2)
	assert.Bool(queue.Put(alloc.NewSmallBuffer().Clear().AppendString("abcd"))).IsTrue()
	queue.SetReadDeadline(time.Now().Add(time.Second))

	b := make([]byte, 1024)
	nBytes, err := queue.Read(b)
	assert.Error(err).IsNil()
	assert.Int(nBytes).Equals(4)
	assert.String(string(b[:nBytes])).Equals("abcd")

	nBytes, err = queue.Read(b)
	assert.Error(err).IsNotNil()
}
