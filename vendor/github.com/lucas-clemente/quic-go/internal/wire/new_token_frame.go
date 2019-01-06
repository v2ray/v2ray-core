package wire

import (
	"bytes"
	"io"

	"github.com/lucas-clemente/quic-go/internal/protocol"
	"github.com/lucas-clemente/quic-go/internal/utils"
)

// A NewTokenFrame is a NEW_TOKEN frame
type NewTokenFrame struct {
	Token []byte
}

func parseNewTokenFrame(r *bytes.Reader, _ protocol.VersionNumber) (*NewTokenFrame, error) {
	if _, err := r.ReadByte(); err != nil {
		return nil, err
	}
	tokenLen, err := utils.ReadVarInt(r)
	if err != nil {
		return nil, err
	}
	if uint64(r.Len()) < tokenLen {
		return nil, io.EOF
	}
	token := make([]byte, int(tokenLen))
	if _, err := io.ReadFull(r, token); err != nil {
		return nil, err
	}
	return &NewTokenFrame{Token: token}, nil
}

func (f *NewTokenFrame) Write(b *bytes.Buffer, _ protocol.VersionNumber) error {
	b.WriteByte(0x7)
	utils.WriteVarInt(b, uint64(len(f.Token)))
	b.Write(f.Token)
	return nil
}

// Length of a written frame
func (f *NewTokenFrame) Length(protocol.VersionNumber) protocol.ByteCount {
	return 1 + utils.VarIntLen(uint64(len(f.Token))) + protocol.ByteCount(len(f.Token))
}
