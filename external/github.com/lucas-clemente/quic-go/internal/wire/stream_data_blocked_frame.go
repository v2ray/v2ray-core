package wire

import (
	"bytes"

	"v2ray.com/core/external/github.com/lucas-clemente/quic-go/internal/protocol"
	"v2ray.com/core/external/github.com/lucas-clemente/quic-go/internal/utils"
)

// A StreamDataBlockedFrame is a STREAM_DATA_BLOCKED frame
type StreamDataBlockedFrame struct {
	StreamID  protocol.StreamID
	DataLimit protocol.ByteCount
}

func parseStreamDataBlockedFrame(r *bytes.Reader, _ protocol.VersionNumber) (*StreamDataBlockedFrame, error) {
	if _, err := r.ReadByte(); err != nil {
		return nil, err
	}

	sid, err := utils.ReadVarInt(r)
	if err != nil {
		return nil, err
	}
	offset, err := utils.ReadVarInt(r)
	if err != nil {
		return nil, err
	}

	return &StreamDataBlockedFrame{
		StreamID:  protocol.StreamID(sid),
		DataLimit: protocol.ByteCount(offset),
	}, nil
}

func (f *StreamDataBlockedFrame) Write(b *bytes.Buffer, version protocol.VersionNumber) error {
	b.WriteByte(0x15)
	utils.WriteVarInt(b, uint64(f.StreamID))
	utils.WriteVarInt(b, uint64(f.DataLimit))
	return nil
}

// Length of a written frame
func (f *StreamDataBlockedFrame) Length(version protocol.VersionNumber) protocol.ByteCount {
	return 1 + utils.VarIntLen(uint64(f.StreamID)) + utils.VarIntLen(uint64(f.DataLimit))
}
