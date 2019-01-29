package wire

import (
	"bytes"
	"io"

	"v2ray.com/core/external/github.com/lucas-clemente/quic-go/internal/protocol"
	"v2ray.com/core/external/github.com/lucas-clemente/quic-go/internal/qerr"
	"v2ray.com/core/external/github.com/lucas-clemente/quic-go/internal/utils"
)

// A ConnectionCloseFrame is a CONNECTION_CLOSE frame
type ConnectionCloseFrame struct {
	IsApplicationError bool
	ErrorCode          qerr.ErrorCode
	ReasonPhrase       string
}

func parseConnectionCloseFrame(r *bytes.Reader, version protocol.VersionNumber) (*ConnectionCloseFrame, error) {
	typeByte, err := r.ReadByte()
	if err != nil {
		return nil, err
	}

	f := &ConnectionCloseFrame{IsApplicationError: typeByte == 0x1d}
	ec, err := utils.BigEndian.ReadUint16(r)
	if err != nil {
		return nil, err
	}
	f.ErrorCode = qerr.ErrorCode(ec)
	// read the Frame Type, if this is not an application error
	if !f.IsApplicationError {
		if _, err := utils.ReadVarInt(r); err != nil {
			return nil, err
		}
	}
	var reasonPhraseLen uint64
	reasonPhraseLen, err = utils.ReadVarInt(r)
	if err != nil {
		return nil, err
	}
	// shortcut to prevent the unnecessary allocation of dataLen bytes
	// if the dataLen is larger than the remaining length of the packet
	// reading the whole reason phrase would result in EOF when attempting to READ
	if int(reasonPhraseLen) > r.Len() {
		return nil, io.EOF
	}

	reasonPhrase := make([]byte, reasonPhraseLen)
	if _, err := io.ReadFull(r, reasonPhrase); err != nil {
		// this should never happen, since we already checked the reasonPhraseLen earlier
		return nil, err
	}
	f.ReasonPhrase = string(reasonPhrase)
	return f, nil
}

// Length of a written frame
func (f *ConnectionCloseFrame) Length(version protocol.VersionNumber) protocol.ByteCount {
	length := 1 + 2 + utils.VarIntLen(uint64(len(f.ReasonPhrase))) + protocol.ByteCount(len(f.ReasonPhrase))
	if !f.IsApplicationError {
		length++ // for the frame type
	}
	return length
}

func (f *ConnectionCloseFrame) Write(b *bytes.Buffer, version protocol.VersionNumber) error {
	if f.IsApplicationError {
		b.WriteByte(0x1d)
	} else {
		b.WriteByte(0x1c)
	}

	utils.BigEndian.WriteUint16(b, uint16(f.ErrorCode))
	if !f.IsApplicationError {
		utils.WriteVarInt(b, 0)
	}
	utils.WriteVarInt(b, uint64(len(f.ReasonPhrase)))
	b.WriteString(f.ReasonPhrase)
	return nil
}
