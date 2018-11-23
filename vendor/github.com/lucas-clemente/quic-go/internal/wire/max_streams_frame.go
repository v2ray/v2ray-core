package wire

import (
	"bytes"

	"github.com/lucas-clemente/quic-go/internal/protocol"
	"github.com/lucas-clemente/quic-go/internal/utils"
)

// A MaxStreamsFrame is a MAX_STREAMS frame
type MaxStreamsFrame struct {
	Type       protocol.StreamType
	MaxStreams uint64
}

func parseMaxStreamsFrame(r *bytes.Reader, _ protocol.VersionNumber) (*MaxStreamsFrame, error) {
	typeByte, err := r.ReadByte()
	if err != nil {
		return nil, err
	}

	f := &MaxStreamsFrame{}
	switch typeByte {
	case 0x12:
		f.Type = protocol.StreamTypeBidi
	case 0x13:
		f.Type = protocol.StreamTypeUni
	}
	streamID, err := utils.ReadVarInt(r)
	if err != nil {
		return nil, err
	}
	f.MaxStreams = streamID
	return f, nil
}

func (f *MaxStreamsFrame) Write(b *bytes.Buffer, _ protocol.VersionNumber) error {
	switch f.Type {
	case protocol.StreamTypeBidi:
		b.WriteByte(0x12)
	case protocol.StreamTypeUni:
		b.WriteByte(0x13)
	}
	utils.WriteVarInt(b, f.MaxStreams)
	return nil
}

// Length of a written frame
func (f *MaxStreamsFrame) Length(protocol.VersionNumber) protocol.ByteCount {
	return 1 + utils.VarIntLen(f.MaxStreams)
}
