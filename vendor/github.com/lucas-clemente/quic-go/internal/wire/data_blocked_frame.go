package wire

import (
	"bytes"

	"github.com/lucas-clemente/quic-go/internal/protocol"
	"github.com/lucas-clemente/quic-go/internal/utils"
)

// A DataBlockedFrame is a DATA_BLOCKED frame
type DataBlockedFrame struct {
	DataLimit protocol.ByteCount
}

func parseDataBlockedFrame(r *bytes.Reader, _ protocol.VersionNumber) (*DataBlockedFrame, error) {
	if _, err := r.ReadByte(); err != nil {
		return nil, err
	}
	offset, err := utils.ReadVarInt(r)
	if err != nil {
		return nil, err
	}
	return &DataBlockedFrame{
		DataLimit: protocol.ByteCount(offset),
	}, nil
}

func (f *DataBlockedFrame) Write(b *bytes.Buffer, version protocol.VersionNumber) error {
	typeByte := uint8(0x14)
	b.WriteByte(typeByte)
	utils.WriteVarInt(b, uint64(f.DataLimit))
	return nil
}

// Length of a written frame
func (f *DataBlockedFrame) Length(version protocol.VersionNumber) protocol.ByteCount {
	return 1 + utils.VarIntLen(uint64(f.DataLimit))
}
