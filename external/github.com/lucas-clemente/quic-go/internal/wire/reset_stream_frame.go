package wire

import (
	"bytes"

	"v2ray.com/core/external/github.com/lucas-clemente/quic-go/internal/protocol"
	"v2ray.com/core/external/github.com/lucas-clemente/quic-go/internal/utils"
)

// A ResetStreamFrame is a RESET_STREAM frame in QUIC
type ResetStreamFrame struct {
	StreamID   protocol.StreamID
	ErrorCode  protocol.ApplicationErrorCode
	ByteOffset protocol.ByteCount
}

func parseResetStreamFrame(r *bytes.Reader, version protocol.VersionNumber) (*ResetStreamFrame, error) {
	if _, err := r.ReadByte(); err != nil { // read the TypeByte
		return nil, err
	}

	var streamID protocol.StreamID
	var errorCode uint16
	var byteOffset protocol.ByteCount
	sid, err := utils.ReadVarInt(r)
	if err != nil {
		return nil, err
	}
	streamID = protocol.StreamID(sid)
	errorCode, err = utils.BigEndian.ReadUint16(r)
	if err != nil {
		return nil, err
	}
	bo, err := utils.ReadVarInt(r)
	if err != nil {
		return nil, err
	}
	byteOffset = protocol.ByteCount(bo)

	return &ResetStreamFrame{
		StreamID:   streamID,
		ErrorCode:  protocol.ApplicationErrorCode(errorCode),
		ByteOffset: byteOffset,
	}, nil
}

func (f *ResetStreamFrame) Write(b *bytes.Buffer, version protocol.VersionNumber) error {
	b.WriteByte(0x4)
	utils.WriteVarInt(b, uint64(f.StreamID))
	utils.BigEndian.WriteUint16(b, uint16(f.ErrorCode))
	utils.WriteVarInt(b, uint64(f.ByteOffset))
	return nil
}

// Length of a written frame
func (f *ResetStreamFrame) Length(version protocol.VersionNumber) protocol.ByteCount {
	return 1 + utils.VarIntLen(uint64(f.StreamID)) + 2 + utils.VarIntLen(uint64(f.ByteOffset))
}
