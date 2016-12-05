package kcp_test

import (
	"testing"

	"v2ray.com/core/common/alloc"
	"v2ray.com/core/testing/assert"
	. "v2ray.com/core/transport/internet/kcp"
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
		Conv:        1,
		Timestamp:   3,
		Number:      4,
		SendingNext: 5,
		Data:        alloc.NewLocalBuffer(512).Clear().Append([]byte{'a', 'b', 'c', 'd'}),
	}

	nBytes := seg.ByteSize()
	bytes := seg.Bytes(nil)

	assert.Int(len(bytes)).Equals(nBytes)

	iseg, _ := ReadSegment(bytes)
	seg2 := iseg.(*DataSegment)
	assert.Uint16(seg2.Conv).Equals(seg.Conv)
	assert.Uint32(seg2.Timestamp).Equals(seg.Timestamp)
	assert.Uint32(seg2.SendingNext).Equals(seg.SendingNext)
	assert.Uint32(seg2.Number).Equals(seg.Number)
	assert.Bytes(seg2.Data.Bytes()).Equals(seg.Data.Bytes())
}

func TestACKSegment(t *testing.T) {
	assert := assert.On(t)

	seg := &AckSegment{
		Conv:            1,
		ReceivingWindow: 2,
		ReceivingNext:   3,
		Timestamp:       10,
		Count:           5,
		NumberList:      []uint32{1, 3, 5, 7, 9},
	}

	nBytes := seg.ByteSize()
	bytes := seg.Bytes(nil)

	assert.Int(len(bytes)).Equals(nBytes)

	iseg, _ := ReadSegment(bytes)
	seg2 := iseg.(*AckSegment)
	assert.Uint16(seg2.Conv).Equals(seg.Conv)
	assert.Uint32(seg2.ReceivingWindow).Equals(seg.ReceivingWindow)
	assert.Uint32(seg2.ReceivingNext).Equals(seg.ReceivingNext)
	assert.Byte(seg2.Count).Equals(seg.Count)
	assert.Uint32(seg2.Timestamp).Equals(seg.Timestamp)
	for i := byte(0); i < seg2.Count; i++ {
		assert.Uint32(seg2.NumberList[i]).Equals(seg.NumberList[i])
	}
}

func TestCmdSegment(t *testing.T) {
	assert := assert.On(t)

	seg := &CmdOnlySegment{
		Conv:         1,
		Command:      CommandPing,
		Option:       SegmentOptionClose,
		SendingNext:  11,
		ReceivinNext: 13,
		PeerRTO:      15,
	}

	nBytes := seg.ByteSize()
	bytes := seg.Bytes(nil)

	assert.Int(len(bytes)).Equals(nBytes)

	iseg, _ := ReadSegment(bytes)
	seg2 := iseg.(*CmdOnlySegment)
	assert.Uint16(seg2.Conv).Equals(seg.Conv)
	assert.Byte(byte(seg2.Command)).Equals(byte(seg.Command))
	assert.Byte(byte(seg2.Option)).Equals(byte(seg.Option))
	assert.Uint32(seg2.SendingNext).Equals(seg.SendingNext)
	assert.Uint32(seg2.ReceivinNext).Equals(seg.ReceivinNext)
	assert.Uint32(seg2.PeerRTO).Equals(seg.PeerRTO)
}
