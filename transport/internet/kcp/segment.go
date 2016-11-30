package kcp

import (
	"v2ray.com/core/common"
	"v2ray.com/core/common/alloc"
	_ "v2ray.com/core/common/log"
	"v2ray.com/core/common/serial"
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
	Conversation() uint16
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

	timeout  uint32
	transmit uint32
}

func NewDataSegment() *DataSegment {
	return new(DataSegment)
}

func (v *DataSegment) Conversation() uint16 {
	return v.Conv
}

func (v *DataSegment) SetData(b []byte) {
	if v.Data == nil {
		v.Data = alloc.NewSmallBuffer()
	}
	v.Data.Clear().Append(b)
}

func (v *DataSegment) Bytes(b []byte) []byte {
	b = serial.Uint16ToBytes(v.Conv, b)
	b = append(b, byte(CommandData), byte(v.Option))
	b = serial.Uint32ToBytes(v.Timestamp, b)
	b = serial.Uint32ToBytes(v.Number, b)
	b = serial.Uint32ToBytes(v.SendingNext, b)
	b = serial.Uint16ToBytes(uint16(v.Data.Len()), b)
	b = append(b, v.Data.Value...)
	return b
}

func (v *DataSegment) ByteSize() int {
	return 2 + 1 + 1 + 4 + 4 + 4 + 2 + v.Data.Len()
}

func (v *DataSegment) Release() {
	v.Data.Release()
	v.Data = nil
}

type AckSegment struct {
	Conv            uint16
	Option          SegmentOption
	ReceivingWindow uint32
	ReceivingNext   uint32
	Timestamp       uint32
	Count           byte
	NumberList      []uint32
}

const ackNumberLimit = 128

func NewAckSegment() *AckSegment {
	return &AckSegment{
		NumberList: make([]uint32, 0, ackNumberLimit),
	}
}

func (v *AckSegment) Conversation() uint16 {
	return v.Conv
}

func (v *AckSegment) PutTimestamp(timestamp uint32) {
	if timestamp-v.Timestamp < 0x7FFFFFFF {
		v.Timestamp = timestamp
	}
}

func (v *AckSegment) PutNumber(number uint32) {
	v.Count++
	v.NumberList = append(v.NumberList, number)
}

func (v *AckSegment) IsFull() bool {
	return v.Count == ackNumberLimit
}

func (v *AckSegment) ByteSize() int {
	return 2 + 1 + 1 + 4 + 4 + 4 + 1 + int(v.Count)*4
}

func (v *AckSegment) Bytes(b []byte) []byte {
	b = serial.Uint16ToBytes(v.Conv, b)
	b = append(b, byte(CommandACK), byte(v.Option))
	b = serial.Uint32ToBytes(v.ReceivingWindow, b)
	b = serial.Uint32ToBytes(v.ReceivingNext, b)
	b = serial.Uint32ToBytes(v.Timestamp, b)
	b = append(b, v.Count)
	for i := byte(0); i < v.Count; i++ {
		b = serial.Uint32ToBytes(v.NumberList[i], b)
	}
	return b
}

func (v *AckSegment) Release() {
	v.NumberList = nil
}

type CmdOnlySegment struct {
	Conv         uint16
	Command      Command
	Option       SegmentOption
	SendingNext  uint32
	ReceivinNext uint32
	PeerRTO      uint32
}

func NewCmdOnlySegment() *CmdOnlySegment {
	return new(CmdOnlySegment)
}

func (v *CmdOnlySegment) Conversation() uint16 {
	return v.Conv
}

func (v *CmdOnlySegment) ByteSize() int {
	return 2 + 1 + 1 + 4 + 4 + 4
}

func (v *CmdOnlySegment) Bytes(b []byte) []byte {
	b = serial.Uint16ToBytes(v.Conv, b)
	b = append(b, byte(v.Command), byte(v.Option))
	b = serial.Uint32ToBytes(v.SendingNext, b)
	b = serial.Uint32ToBytes(v.ReceivinNext, b)
	b = serial.Uint32ToBytes(v.PeerRTO, b)
	return b
}

func (v *CmdOnlySegment) Release() {
}

func ReadSegment(buf []byte) (Segment, []byte) {
	if len(buf) < 4 {
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
		seg.SetData(buf[:dataLen])
		buf = buf[dataLen:]

		return seg, buf
	}

	if cmd == CommandACK {
		seg := NewAckSegment()
		seg.Conv = conv
		seg.Option = opt
		if len(buf) < 13 {
			return nil, nil
		}

		seg.ReceivingWindow = serial.BytesToUint32(buf)
		buf = buf[4:]

		seg.ReceivingNext = serial.BytesToUint32(buf)
		buf = buf[4:]

		seg.Timestamp = serial.BytesToUint32(buf)
		buf = buf[4:]

		count := int(buf[0])
		buf = buf[1:]

		if len(buf) < count*4 {
			return nil, nil
		}
		for i := 0; i < count; i++ {
			seg.PutNumber(serial.BytesToUint32(buf))
			buf = buf[4:]
		}

		return seg, buf
	}

	seg := NewCmdOnlySegment()
	seg.Conv = conv
	seg.Command = cmd
	seg.Option = opt

	if len(buf) < 12 {
		return nil, nil
	}

	seg.SendingNext = serial.BytesToUint32(buf)
	buf = buf[4:]

	seg.ReceivinNext = serial.BytesToUint32(buf)
	buf = buf[4:]

	seg.PeerRTO = serial.BytesToUint32(buf)
	buf = buf[4:]

	return seg, buf
}
