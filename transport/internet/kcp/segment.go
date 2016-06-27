package kcp

import (
	"github.com/v2ray/v2ray-core/common"
	"github.com/v2ray/v2ray-core/common/alloc"
	"github.com/v2ray/v2ray-core/common/serial"
)

type SegmentCommand byte

const (
	SegmentCommandACK        SegmentCommand = 0
	SegmentCommandData       SegmentCommand = 1
	SegmentCommandTerminated SegmentCommand = 2
)

type SegmentOption byte

const (
	SegmentOptionClose SegmentOption = 1
)

type ISegment interface {
	common.Releasable
	ByteSize() int
	Bytes([]byte) []byte
}

type DataSegment struct {
	Conv            uint16
	Opt             SegmentOption
	ReceivingWindow uint32
	Timestamp       uint32
	Number          uint32
	Unacknowledged  uint32
	Data            *alloc.Buffer

	timeout    uint32
	ackSkipped uint32
	transmit   uint32
}

func (this *DataSegment) Bytes(b []byte) []byte {
	b = serial.Uint16ToBytes(this.Conv, b)
	b = append(b, byte(SegmentCommandData), byte(this.Opt))
	b = serial.Uint32ToBytes(this.ReceivingWindow, b)
	b = serial.Uint32ToBytes(this.Timestamp, b)
	b = serial.Uint32ToBytes(this.Number, b)
	b = serial.Uint32ToBytes(this.Unacknowledged, b)
	b = serial.Uint16ToBytes(uint16(this.Data.Len()), b)
	b = append(b, this.Data.Value...)
	return b
}

func (this *DataSegment) ByteSize() int {
	return 2 + 1 + 1 + 4 + 4 + 4 + 4 + 2 + this.Data.Len()
}

func (this *DataSegment) Release() {
	this.Data.Release()
}

type ACKSegment struct {
	Conv            uint16
	Opt             SegmentOption
	ReceivingWindow uint32
	Unacknowledged  uint32
	Count           byte
	NumberList      []uint32
	TimestampList   []uint32
}

func (this *ACKSegment) ByteSize() int {
	return 2 + 1 + 1 + 4 + 4 + 1 + len(this.NumberList)*4 + len(this.TimestampList)*4
}

func (this *ACKSegment) Bytes(b []byte) []byte {
	b = serial.Uint16ToBytes(this.Conv, b)
	b = append(b, byte(SegmentCommandACK), byte(this.Opt))
	b = serial.Uint32ToBytes(this.ReceivingWindow, b)
	b = serial.Uint32ToBytes(this.Unacknowledged, b)
	b = append(b, this.Count)
	for i := byte(0); i < this.Count; i++ {
		b = serial.Uint32ToBytes(this.NumberList[i], b)
		b = serial.Uint32ToBytes(this.TimestampList[i], b)
	}
	return b
}

func (this *ACKSegment) Release() {}

type TerminationSegment struct {
	Conv uint16
	Opt  SegmentOption
}

func (this *TerminationSegment) ByteSize() int {
	return 2 + 1 + 1
}

func (this *TerminationSegment) Bytes(b []byte) []byte {
	b = serial.Uint16ToBytes(this.Conv, b)
	b = append(b, byte(SegmentCommandTerminated), byte(this.Opt))
	return b
}

func (this *TerminationSegment) Release() {}

func ReadSegment(buf []byte) (ISegment, []byte) {
	if len(buf) <= 12 {
		return nil, nil
	}

	conv := serial.BytesToUint16(buf)
	buf = buf[2:]

	cmd := SegmentCommand(buf[0])
	opt := SegmentOption(buf[1])
	buf = buf[2:]

	if cmd == SegmentCommandData {
		seg := &DataSegment{
			Conv: conv,
			Opt:  opt,
		}
		seg.ReceivingWindow = serial.BytesToUint32(buf)
		buf = buf[4:]

		seg.Timestamp = serial.BytesToUint32(buf)
		buf = buf[4:]

		seg.Number = serial.BytesToUint32(buf)
		buf = buf[4:]

		seg.Unacknowledged = serial.BytesToUint32(buf)
		buf = buf[4:]

		len := serial.BytesToUint16(buf)
		buf = buf[2:]

		seg.Data = alloc.NewSmallBuffer().Clear().Append(buf[:len])
		buf = buf[len:]

		return seg, buf
	}

	if cmd == SegmentCommandACK {
		seg := &ACKSegment{
			Conv: conv,
			Opt:  opt,
		}
		seg.ReceivingWindow = serial.BytesToUint32(buf)
		buf = buf[4:]

		seg.Unacknowledged = serial.BytesToUint32(buf)
		buf = buf[4:]

		seg.Count = buf[0]
		buf = buf[1:]

		seg.NumberList = make([]uint32, 0, seg.Count)
		seg.TimestampList = make([]uint32, 0, seg.Count)

		for i := 0; i < int(seg.Count); i++ {
			seg.NumberList = append(seg.NumberList, serial.BytesToUint32(buf))
			seg.TimestampList = append(seg.TimestampList, serial.BytesToUint32(buf[4:]))
			buf = buf[8:]
		}

		return seg, buf
	}

	if cmd == SegmentCommandTerminated {
		return &TerminationSegment{
			Conv: conv,
			Opt:  opt,
		}, buf
	}

	return nil, nil
}
