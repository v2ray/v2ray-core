package kcp

import (
	"sync"

	"github.com/v2ray/v2ray-core/common"
	"github.com/v2ray/v2ray-core/common/alloc"
	_ "github.com/v2ray/v2ray-core/common/log"
	"github.com/v2ray/v2ray-core/common/serial"
)

var (
	dataSegmentPool = &sync.Pool{
		New: func() interface{} { return new(DataSegment) },
	}
	ackSegmentPool = &sync.Pool{
		New: func() interface{} { return new(AckSegment) },
	}
	cmdSegmentPool = &sync.Pool{
		New: func() interface{} { return new(CmdOnlySegment) },
	}
)

type SegmentCommand byte

const (
	SegmentCommandACK        SegmentCommand = 0
	SegmentCommandData       SegmentCommand = 1
	SegmentCommandTerminated SegmentCommand = 2
	SegmentCommandPing       SegmentCommand = 3
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
	Opt         SegmentOption
	Timestamp   uint32
	Number      uint32
	SendingNext uint32
	Data        *alloc.Buffer

	timeout    uint32
	ackSkipped uint32
	transmit   uint32
}

func NewDataSegment() *DataSegment {
	return dataSegmentPool.Get().(*DataSegment)
}

func (this *DataSegment) Bytes(b []byte) []byte {
	b = serial.Uint16ToBytes(this.Conv, b)
	b = append(b, byte(SegmentCommandData), byte(this.Opt))
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
	this.Opt = 0
	this.timeout = 0
	this.ackSkipped = 0
	this.transmit = 0
	dataSegmentPool.Put(this)
}

type AckSegment struct {
	Conv            uint16
	Opt             SegmentOption
	ReceivingWindow uint32
	ReceivingNext   uint32
	Count           byte
	NumberList      []uint32
	TimestampList   []uint32
}

func NewAckSegment() *AckSegment {
	seg := ackSegmentPool.Get().(*AckSegment)
	if seg.NumberList == nil {
		seg.NumberList = make([]uint32, 0, 128)
	}
	if seg.TimestampList == nil {
		seg.TimestampList = make([]uint32, 0, 128)
	}
	return seg
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
	b = append(b, byte(SegmentCommandACK), byte(this.Opt))
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
	this.Opt = 0
	this.Count = 0
	this.NumberList = this.NumberList[:0]
	this.TimestampList = this.TimestampList[:0]
	ackSegmentPool.Put(this)
}

type CmdOnlySegment struct {
	Conv         uint16
	Cmd          SegmentCommand
	Opt          SegmentOption
	SendingNext  uint32
	ReceivinNext uint32
}

func NewCmdOnlySegment() *CmdOnlySegment {
	return cmdSegmentPool.Get().(*CmdOnlySegment)
}

func (this *CmdOnlySegment) ByteSize() int {
	return 2 + 1 + 1 + 4 + 4
}

func (this *CmdOnlySegment) Bytes(b []byte) []byte {
	b = serial.Uint16ToBytes(this.Conv, b)
	b = append(b, byte(this.Cmd), byte(this.Opt))
	b = serial.Uint32ToBytes(this.SendingNext, b)
	b = serial.Uint32ToBytes(this.ReceivinNext, b)
	return b
}

func (this *CmdOnlySegment) Release() {
	this.Opt = 0
	cmdSegmentPool.Put(this)
}

func ReadSegment(buf []byte) (Segment, []byte) {
	if len(buf) <= 4 {
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
		seg.Data = alloc.NewSmallBuffer().Clear().Append(buf[:dataLen])
		buf = buf[dataLen:]

		return seg, buf
	}

	if cmd == SegmentCommandACK {
		seg := &AckSegment{
			Conv: conv,
			Opt:  opt,
		}
		if len(buf) < 9 {
			return nil, nil
		}

		seg.ReceivingWindow = serial.BytesToUint32(buf)
		buf = buf[4:]

		seg.ReceivingNext = serial.BytesToUint32(buf)
		buf = buf[4:]

		seg.Count = buf[0]
		buf = buf[1:]

		seg.NumberList = make([]uint32, 0, seg.Count)
		seg.TimestampList = make([]uint32, 0, seg.Count)

		if len(buf) < int(seg.Count)*8 {
			return nil, nil
		}
		for i := 0; i < int(seg.Count); i++ {
			seg.NumberList = append(seg.NumberList, serial.BytesToUint32(buf))
			seg.TimestampList = append(seg.TimestampList, serial.BytesToUint32(buf[4:]))
			buf = buf[8:]
		}

		return seg, buf
	}

	seg := &CmdOnlySegment{
		Conv: conv,
		Cmd:  cmd,
		Opt:  opt,
	}

	if len(buf) < 8 {
		return nil, nil
	}

	seg.SendingNext = serial.BytesToUint32(buf)
	buf = buf[4:]

	seg.ReceivinNext = serial.BytesToUint32(buf)
	buf = buf[4:]

	return seg, buf
}
