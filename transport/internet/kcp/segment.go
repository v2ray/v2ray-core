package kcp

import (
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/serial"
)

// Command is a KCP command that indicate the purpose of a Segment.
type Command byte

const (
	// CommandACK indicates an AckSegment.
	CommandACK Command = 0
	// CommandData indicates a DataSegment.
	CommandData Command = 1
	// CommandTerminate indicates that peer terminates the connection.
	CommandTerminate Command = 2
	// CommandPing indicates a ping.
	CommandPing Command = 3
)

type SegmentOption byte

const (
	SegmentOptionClose SegmentOption = 1
)

type Segment interface {
	Release()
	Conversation() uint16
	Command() Command
	ByteSize() int32
	Bytes() buf.Supplier
	parse(conv uint16, cmd Command, opt SegmentOption, buf []byte) (bool, []byte)
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

	payload  *buf.Buffer
	timeout  uint32
	transmit uint32
}

func NewDataSegment() *DataSegment {
	return new(DataSegment)
}

func (s *DataSegment) parse(conv uint16, cmd Command, opt SegmentOption, buf []byte) (bool, []byte) {
	s.Conv = conv
	s.Option = opt
	if len(buf) < 15 {
		return false, nil
	}
	s.Timestamp = serial.BytesToUint32(buf)
	buf = buf[4:]

	s.Number = serial.BytesToUint32(buf)
	buf = buf[4:]

	s.SendingNext = serial.BytesToUint32(buf)
	buf = buf[4:]

	dataLen := int(serial.BytesToUint16(buf))
	buf = buf[2:]

	if len(buf) < dataLen {
		return false, nil
	}
	s.Data().Clear()
	s.Data().Write(buf[:dataLen])
	buf = buf[dataLen:]

	return true, buf
}

func (s *DataSegment) Conversation() uint16 {
	return s.Conv
}

func (*DataSegment) Command() Command {
	return CommandData
}

func (s *DataSegment) Detach() *buf.Buffer {
	r := s.payload
	s.payload = nil
	return r
}

func (s *DataSegment) Data() *buf.Buffer {
	if s.payload == nil {
		s.payload = buf.New()
	}
	return s.payload
}

func (s *DataSegment) Bytes() buf.Supplier {
	return func(b []byte) (int, error) {
		b = serial.Uint16ToBytes(s.Conv, b[:0])
		b = append(b, byte(CommandData), byte(s.Option))
		b = serial.Uint32ToBytes(s.Timestamp, b)
		b = serial.Uint32ToBytes(s.Number, b)
		b = serial.Uint32ToBytes(s.SendingNext, b)
		b = serial.Uint16ToBytes(uint16(s.payload.Len()), b)
		b = append(b, s.payload.Bytes()...)
		return len(b), nil
	}
}

func (s *DataSegment) ByteSize() int32 {
	return 2 + 1 + 1 + 4 + 4 + 4 + 2 + s.payload.Len()
}

func (s *DataSegment) Release() {
	s.payload.Release()
	s.payload = nil
}

type AckSegment struct {
	Conv            uint16
	Option          SegmentOption
	ReceivingWindow uint32
	ReceivingNext   uint32
	Timestamp       uint32
	NumberList      []uint32
}

const ackNumberLimit = 128

func NewAckSegment() *AckSegment {
	return new(AckSegment)
}

func (s *AckSegment) parse(conv uint16, cmd Command, opt SegmentOption, buf []byte) (bool, []byte) {
	s.Conv = conv
	s.Option = opt
	if len(buf) < 13 {
		return false, nil
	}

	s.ReceivingWindow = serial.BytesToUint32(buf)
	buf = buf[4:]

	s.ReceivingNext = serial.BytesToUint32(buf)
	buf = buf[4:]

	s.Timestamp = serial.BytesToUint32(buf)
	buf = buf[4:]

	count := int(buf[0])
	buf = buf[1:]

	if len(buf) < count*4 {
		return false, nil
	}
	for i := 0; i < count; i++ {
		s.PutNumber(serial.BytesToUint32(buf))
		buf = buf[4:]
	}

	return true, buf
}

func (s *AckSegment) Conversation() uint16 {
	return s.Conv
}

func (*AckSegment) Command() Command {
	return CommandACK
}

func (s *AckSegment) PutTimestamp(timestamp uint32) {
	if timestamp-s.Timestamp < 0x7FFFFFFF {
		s.Timestamp = timestamp
	}
}

func (s *AckSegment) PutNumber(number uint32) {
	s.NumberList = append(s.NumberList, number)
}

func (s *AckSegment) IsFull() bool {
	return len(s.NumberList) == ackNumberLimit
}

func (s *AckSegment) IsEmpty() bool {
	return len(s.NumberList) == 0
}

func (s *AckSegment) ByteSize() int32 {
	return 2 + 1 + 1 + 4 + 4 + 4 + 1 + int32(len(s.NumberList)*4)
}

func (s *AckSegment) Bytes() buf.Supplier {
	return func(b []byte) (int, error) {
		b = serial.Uint16ToBytes(s.Conv, b[:0])
		b = append(b, byte(CommandACK), byte(s.Option))
		b = serial.Uint32ToBytes(s.ReceivingWindow, b)
		b = serial.Uint32ToBytes(s.ReceivingNext, b)
		b = serial.Uint32ToBytes(s.Timestamp, b)
		count := byte(len(s.NumberList))
		b = append(b, count)
		for _, number := range s.NumberList {
			b = serial.Uint32ToBytes(number, b)
		}
		return int(s.ByteSize()), nil
	}
}

func (s *AckSegment) Release() {}

type CmdOnlySegment struct {
	Conv          uint16
	Cmd           Command
	Option        SegmentOption
	SendingNext   uint32
	ReceivingNext uint32
	PeerRTO       uint32
}

func NewCmdOnlySegment() *CmdOnlySegment {
	return new(CmdOnlySegment)
}

func (s *CmdOnlySegment) parse(conv uint16, cmd Command, opt SegmentOption, buf []byte) (bool, []byte) {
	s.Conv = conv
	s.Cmd = cmd
	s.Option = opt

	if len(buf) < 12 {
		return false, nil
	}

	s.SendingNext = serial.BytesToUint32(buf)
	buf = buf[4:]

	s.ReceivingNext = serial.BytesToUint32(buf)
	buf = buf[4:]

	s.PeerRTO = serial.BytesToUint32(buf)
	buf = buf[4:]

	return true, buf
}

func (s *CmdOnlySegment) Conversation() uint16 {
	return s.Conv
}

func (s *CmdOnlySegment) Command() Command {
	return s.Cmd
}

func (*CmdOnlySegment) ByteSize() int32 {
	return 2 + 1 + 1 + 4 + 4 + 4
}

func (s *CmdOnlySegment) Bytes() buf.Supplier {
	return func(b []byte) (int, error) {
		b = serial.Uint16ToBytes(s.Conv, b[:0])
		b = append(b, byte(s.Cmd), byte(s.Option))
		b = serial.Uint32ToBytes(s.SendingNext, b)
		b = serial.Uint32ToBytes(s.ReceivingNext, b)
		b = serial.Uint32ToBytes(s.PeerRTO, b)
		return len(b), nil
	}
}

func (*CmdOnlySegment) Release() {}

func ReadSegment(buf []byte) (Segment, []byte) {
	if len(buf) < 4 {
		return nil, nil
	}

	conv := serial.BytesToUint16(buf)
	buf = buf[2:]

	cmd := Command(buf[0])
	opt := SegmentOption(buf[1])
	buf = buf[2:]

	var seg Segment
	switch cmd {
	case CommandData:
		seg = NewDataSegment()
	case CommandACK:
		seg = NewAckSegment()
	default:
		seg = NewCmdOnlySegment()
	}

	valid, extra := seg.parse(conv, cmd, opt, buf)
	if !valid {
		return nil, nil
	}
	return seg, extra
}
