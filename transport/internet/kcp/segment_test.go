package kcp_test

import (
	"testing"

	"github.com/v2ray/v2ray-core/common/alloc"
	"github.com/v2ray/v2ray-core/testing/assert"
	. "github.com/v2ray/v2ray-core/transport/internet/kcp"
)

func TestBadSegment(t *testing.T) {
	assert := assert.On(t)

	seg, buf := ReadSegment(nil)
	assert.Pointer(seg).IsNil()
	assert.Int(len(buf)).Equals(0)
}

func TestDataSegment(t *testing.T) {
	assert := assert.On(t)

	seg := &DataSegment{
		Conv:            1,
		ReceivingWindow: 2,
		Timestamp:       3,
		Number:          4,
		Unacknowledged:  5,
		Data:            alloc.NewSmallBuffer().Clear().Append([]byte{'a', 'b', 'c', 'd'}),
	}

	nBytes := seg.ByteSize()
	bytes := seg.Bytes(nil)

	assert.Int(len(bytes)).Equals(nBytes)

	iseg, _ := ReadSegment(bytes)
	seg2 := iseg.(*DataSegment)
	assert.Uint16(seg2.Conv).Equals(seg.Conv)
	assert.Uint32(seg2.ReceivingWindow).Equals(seg.ReceivingWindow)
	assert.Uint32(seg2.Timestamp).Equals(seg.Timestamp)
	assert.Uint32(seg2.Unacknowledged).Equals(seg.Unacknowledged)
	assert.Uint32(seg2.Number).Equals(seg.Number)
	assert.Bytes(seg2.Data.Value).Equals(seg.Data.Value)
}

func TestACKSegment(t *testing.T) {
	assert := assert.On(t)

	seg := &ACKSegment{
		Conv:            1,
		ReceivingWindow: 2,
		Unacknowledged:  3,
		Count:           5,
		NumberList:      []uint32{1, 3, 5, 7, 9},
		TimestampList:   []uint32{2, 4, 6, 8, 10},
	}

	nBytes := seg.ByteSize()
	bytes := seg.Bytes(nil)

	assert.Int(len(bytes)).Equals(nBytes)

	iseg, _ := ReadSegment(bytes)
	seg2 := iseg.(*ACKSegment)
	assert.Uint16(seg2.Conv).Equals(seg.Conv)
	assert.Uint32(seg2.ReceivingWindow).Equals(seg.ReceivingWindow)
	assert.Uint32(seg2.Unacknowledged).Equals(seg.Unacknowledged)
	assert.Byte(seg2.Count).Equals(seg.Count)
	for i := byte(0); i < seg2.Count; i++ {
		assert.Uint32(seg2.TimestampList[i]).Equals(seg.TimestampList[i])
		assert.Uint32(seg2.NumberList[i]).Equals(seg.NumberList[i])
	}
}
