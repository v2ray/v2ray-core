package kcp_test

import (
	"testing"

	. "v2ray.com/core/transport/internet/kcp"
	. "v2ray.com/ext/assert"
)

func TestBadSegment(t *testing.T) {
	assert := With(t)

	seg, buf := ReadSegment(nil)
	assert(seg, IsNil)
	assert(len(buf), Equals, 0)
}

func TestDataSegment(t *testing.T) {
	assert := With(t)

	seg := &DataSegment{
		Conv:        1,
		Timestamp:   3,
		Number:      4,
		SendingNext: 5,
	}
	seg.Data().Append([]byte{'a', 'b', 'c', 'd'})

	nBytes := seg.ByteSize()
	bytes := make([]byte, nBytes)
	seg.Bytes()(bytes)

	assert(len(bytes), Equals, nBytes)

	iseg, _ := ReadSegment(bytes)
	seg2 := iseg.(*DataSegment)
	assert(seg2.Conv, Equals, seg.Conv)
	assert(seg2.Timestamp, Equals, seg.Timestamp)
	assert(seg2.SendingNext, Equals, seg.SendingNext)
	assert(seg2.Number, Equals, seg.Number)
	assert(seg2.Data().Bytes(), Equals, seg.Data().Bytes())
}

func Test1ByteDataSegment(t *testing.T) {
	assert := With(t)

	seg := &DataSegment{
		Conv:        1,
		Timestamp:   3,
		Number:      4,
		SendingNext: 5,
	}
	seg.Data().AppendBytes('a')

	nBytes := seg.ByteSize()
	bytes := make([]byte, nBytes)
	seg.Bytes()(bytes)

	assert(len(bytes), Equals, nBytes)

	iseg, _ := ReadSegment(bytes)
	seg2 := iseg.(*DataSegment)
	assert(seg2.Conv, Equals, seg.Conv)
	assert(seg2.Timestamp, Equals, seg.Timestamp)
	assert(seg2.SendingNext, Equals, seg.SendingNext)
	assert(seg2.Number, Equals, seg.Number)
	assert(seg2.Data().Bytes(), Equals, seg.Data().Bytes())
}

func TestACKSegment(t *testing.T) {
	assert := With(t)

	seg := &AckSegment{
		Conv:            1,
		ReceivingWindow: 2,
		ReceivingNext:   3,
		Timestamp:       10,
		NumberList:      []uint32{1, 3, 5, 7, 9},
	}

	nBytes := seg.ByteSize()
	bytes := make([]byte, nBytes)
	seg.Bytes()(bytes)

	assert(len(bytes), Equals, nBytes)

	iseg, _ := ReadSegment(bytes)
	seg2 := iseg.(*AckSegment)
	assert(seg2.Conv, Equals, seg.Conv)
	assert(seg2.ReceivingWindow, Equals, seg.ReceivingWindow)
	assert(seg2.ReceivingNext, Equals, seg.ReceivingNext)
	assert(len(seg2.NumberList), Equals, len(seg.NumberList))
	assert(seg2.Timestamp, Equals, seg.Timestamp)
	for i, number := range seg2.NumberList {
		assert(number, Equals, seg.NumberList[i])
	}
}

func TestCmdSegment(t *testing.T) {
	assert := With(t)

	seg := &CmdOnlySegment{
		Conv:          1,
		Cmd:           CommandPing,
		Option:        SegmentOptionClose,
		SendingNext:   11,
		ReceivingNext: 13,
		PeerRTO:       15,
	}

	nBytes := seg.ByteSize()
	bytes := make([]byte, nBytes)
	seg.Bytes()(bytes)

	assert(len(bytes), Equals, nBytes)

	iseg, _ := ReadSegment(bytes)
	seg2 := iseg.(*CmdOnlySegment)
	assert(seg2.Conv, Equals, seg.Conv)
	assert(byte(seg2.Command()), Equals, byte(seg.Command()))
	assert(byte(seg2.Option), Equals, byte(seg.Option))
	assert(seg2.SendingNext, Equals, seg.SendingNext)
	assert(seg2.ReceivingNext, Equals, seg.ReceivingNext)
	assert(seg2.PeerRTO, Equals, seg.PeerRTO)
}
