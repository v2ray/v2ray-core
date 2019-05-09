package kcp_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	. "v2ray.com/core/transport/internet/kcp"
)

func TestBadSegment(t *testing.T) {
	seg, buf := ReadSegment(nil)
	if seg != nil {
		t.Error("non-nil seg")
	}
	if len(buf) != 0 {
		t.Error("buf len: ", len(buf))
	}
}

func TestDataSegment(t *testing.T) {
	seg := &DataSegment{
		Conv:        1,
		Timestamp:   3,
		Number:      4,
		SendingNext: 5,
	}
	seg.Data().Write([]byte{'a', 'b', 'c', 'd'})

	nBytes := seg.ByteSize()
	bytes := make([]byte, nBytes)
	seg.Serialize(bytes)

	iseg, _ := ReadSegment(bytes)
	seg2 := iseg.(*DataSegment)
	if r := cmp.Diff(seg2, seg, cmpopts.IgnoreUnexported(DataSegment{})); r != "" {
		t.Error(r)
	}
	if r := cmp.Diff(seg2.Data().Bytes(), seg.Data().Bytes()); r != "" {
		t.Error(r)
	}
}

func Test1ByteDataSegment(t *testing.T) {
	seg := &DataSegment{
		Conv:        1,
		Timestamp:   3,
		Number:      4,
		SendingNext: 5,
	}
	seg.Data().WriteByte('a')

	nBytes := seg.ByteSize()
	bytes := make([]byte, nBytes)
	seg.Serialize(bytes)

	iseg, _ := ReadSegment(bytes)
	seg2 := iseg.(*DataSegment)
	if r := cmp.Diff(seg2, seg, cmpopts.IgnoreUnexported(DataSegment{})); r != "" {
		t.Error(r)
	}
	if r := cmp.Diff(seg2.Data().Bytes(), seg.Data().Bytes()); r != "" {
		t.Error(r)
	}
}

func TestACKSegment(t *testing.T) {
	seg := &AckSegment{
		Conv:            1,
		ReceivingWindow: 2,
		ReceivingNext:   3,
		Timestamp:       10,
		NumberList:      []uint32{1, 3, 5, 7, 9},
	}

	nBytes := seg.ByteSize()
	bytes := make([]byte, nBytes)
	seg.Serialize(bytes)

	iseg, _ := ReadSegment(bytes)
	seg2 := iseg.(*AckSegment)
	if r := cmp.Diff(seg2, seg); r != "" {
		t.Error(r)
	}
}

func TestCmdSegment(t *testing.T) {
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
	seg.Serialize(bytes)

	iseg, _ := ReadSegment(bytes)
	seg2 := iseg.(*CmdOnlySegment)
	if r := cmp.Diff(seg2, seg); r != "" {
		t.Error(r)
	}
}
