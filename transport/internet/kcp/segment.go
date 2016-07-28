package kcp

import (
	"github.com/v2ray/v2ray-core/common"
	"github.com/v2ray/v2ray-core/common/alloc"
	_ "github.com/v2ray/v2ray-core/common/log"
	"github.com/v2ray/v2ray-core/common/serial"
)

type Command byte

const (
	CommandACK       Command = 0
	CommandData      Command = 1
	CommandTerminate Command = 2
	CommandPing      Command = 3
)

type SegmentOption byte

const (
	SegmentOptionClose SegmentOption = 1
)

type Segment interface {
	common.Releasable
	ByteSize() int
	Bytes([]byte) []byte
}

const (
	DataSegmentOverhead = 18
)

type DataSegment struct {
	Conv        uint16
	Option      SegmentOption
	Timestamp   uint32
	Number      uint32
	SendingNext uint32
	Data        *alloc.Buffer

	timeout    uint32
	ackSkipped uint32
	transmit   uint32
}

func NewDataSegment() *DataSegment {
	return new(DataSegment)
}

func (this *DataSegment) Bytes(b []byte) []byte {
	b = serial.Uint16ToBytes(this.Conv, b)
	b = append(b, byte(CommandData), byte(this.Option))
	b = serial.Uint32ToBytes(this.Timestamp, b)
	b = serial.Uint32ToBytes(this.Number, b)
	b = serial.Uint32ToBytes(this.SendingNext, b)
	b = serial.Uint16ToBytes(uint16(this.Data.Len()), b)
	b = append(b, this.Data.Value...)
	return b
}

func (this *DataSegment) ByteSize() int {
	return 2 + 1 + 1 + 4 + 4 + 4 + 2 + this.Data.Len()
}

func (this *DataSegment) Release() {
	this.Data.Release()
	this.Data = nil
}

type AckSegment struct {
	Conv            uint16
	Option          SegmentOption
	ReceivingWindow uint32
	ReceivingNext   uint32
	Count           byte
	NumberList      []uint32
	TimestampList   []uint32
}

func NewAckSegment() *AckSegment {
	return new(AckSegment)
}

func (this *AckSegment) PutNumber(number uint32, timestamp uint32) {
	this.Count++
	this.NumberList = append(this.NumberList, number)
	this.TimestampList = append(this.TimestampList, timestamp)
}

func (this *AckSegment) IsFull() bool {
	return this.Count == 128
}

func (this *AckSegment) ByteSize() int {
	return 2 + 1 + 1 + 4 + 4 + 1 + int(this.Count)*4 + int(this.Count)*4
}

func (this *AckSegment) Bytes(b []byte) []byte {
	b = serial.Uint16ToBytes(this.Conv, b)
	b = append(b, byte(CommandACK), byte(this.Option))
	b = serial.Uint32ToBytes(this.ReceivingWindow, b)
	b = serial.Uint32ToBytes(this.ReceivingNext, b)
	b = append(b, this.Count)
	for i := byte(0); i < this.Count; i++ {
		b = serial.Uint32ToBytes(this.NumberList[i], b)
		b = serial.Uint32ToBytes(this.TimestampList[i], b)
	}
	return b
}

func (this *AckSegment) Release() {
	this.NumberList = nil
	this.TimestampList = nil
}

type CmdOnlySegment struct {
	Conv         uint16
	Command      Command
	Option       SegmentOption
	SendingNext  uint32
	ReceivinNext uint32
}

func NewCmdOnlySegment() *CmdOnlySegment {
	return new(CmdOnlySegment)
}

func (this *CmdOnlySegment) ByteSize() int {
	return 2 + 1 + 1 + 4 + 4
}

func (this *CmdOnlySegment) Bytes(b []byte) []byte {
	b = serial.Uint16ToBytes(this.Conv, b)
	b = append(b, byte(this.Command), byte(this.Option))
	b = serial.Uint32ToBytes(this.SendingNext, b)
	b = serial.Uint32ToBytes(this.ReceivinNext, b)
	return b
}

func (this *CmdOnlySegment) Release() {
}

func ReadSegment(buf []byte) (Segment, []byte) {
	if len(buf) <= 4 {
		return nil, nil
	}

	conv := serial.BytesToUint16(buf)
	buf = buf[2:]

	cmd := Command(buf[0])
	opt := SegmentOption(buf[1])
	buf = buf[2:]

	if cmd == CommandData {
		seg := NewDataSegment()
		seg.Conv = conv
		seg.Option = opt
		if len(buf) < 16 {
			return nil, nil
		}
		seg.Timestamp = serial.BytesToUint32(buf)
		buf = buf[4:]

		seg.Number = serial.BytesToUint32(buf)
		buf = buf[4:]

		seg.SendingNext = serial.BytesToUint32(buf)
		buf = buf[4:]

		dataLen := int(serial.BytesToUint16(buf))
		buf = buf[2:]

		if len(buf) < dataLen {
			return nil, nil
		}
		seg.Data = AllocateBuffer().Clear().Append(buf[:dataLen])
		buf = buf[dataLen:]

		return seg, buf
	}

	if cmd == CommandACK {
		seg := NewAckSegment()
		seg.Conv = conv
		seg.Option = opt
		if len(buf) < 9 {
			return nil, nil
		}

		seg.ReceivingWindow = serial.BytesToUint32(buf)
		buf = buf[4:]

		seg.ReceivingNext = serial.BytesToUint32(buf)
		buf = buf[4:]

		count := int(buf[0])
		buf = buf[1:]

		if len(buf) < count*8 {
			return nil, nil
		}
		for i := 0; i < count; i++ {
			seg.PutNumber(serial.BytesToUint32(buf), serial.BytesToUint32(buf[4:]))
			buf = buf[8:]
		}

		return seg, buf
	}

	seg := NewCmdOnlySegment()
	seg.Conv = conv
	seg.Command = cmd
	seg.Option = opt

	if len(buf) < 8 {
		return nil, nil
	}

	seg.SendingNext = serial.BytesToUint32(buf)
	buf = buf[4:]

	seg.ReceivinNext = serial.BytesToUint32(buf)
	buf = buf[4:]

	return seg, buf
}
