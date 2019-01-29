package wire

import (
	"bytes"

	"v2ray.com/core/external/github.com/lucas-clemente/quic-go/internal/protocol"
	"v2ray.com/core/external/github.com/lucas-clemente/quic-go/internal/utils"
)

// A MaxStreamDataFrame is a MAX_STREAM_DATA frame
type MaxStreamDataFrame struct {
	StreamID   protocol.StreamID
	ByteOffset protocol.ByteCount
}

func parseMaxStreamDataFrame(r *bytes.Reader, version protocol.VersionNumber) (*MaxStreamDataFrame, error) {
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

	return &MaxStreamDataFrame{
		StreamID:   protocol.StreamID(sid),
		ByteOffset: protocol.ByteCount(offset),
	}, nil
}

func (f *MaxStreamDataFrame) Write(b *bytes.Buffer, version protocol.VersionNumber) error {
	b.WriteByte(0x11)
	utils.WriteVarInt(b, uint64(f.StreamID))
	utils.WriteVarInt(b, uint64(f.ByteOffset))
	return nil
}

// Length of a written frame
func (f *MaxStreamDataFrame) Length(version protocol.VersionNumber) protocol.ByteCount {
	return 1 + utils.VarIntLen(uint64(f.StreamID)) + utils.VarIntLen(uint64(f.ByteOffset))
}
