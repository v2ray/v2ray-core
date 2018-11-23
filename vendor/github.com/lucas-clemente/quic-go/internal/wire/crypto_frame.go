package wire

import (
	"bytes"
	"io"

	"github.com/lucas-clemente/quic-go/internal/protocol"
	"github.com/lucas-clemente/quic-go/internal/utils"
)

// A CryptoFrame is a CRYPTO frame
type CryptoFrame struct {
	Offset protocol.ByteCount
	Data   []byte
}

func parseCryptoFrame(r *bytes.Reader, _ protocol.VersionNumber) (*CryptoFrame, error) {
	if _, err := r.ReadByte(); err != nil {
		return nil, err
	}

	frame := &CryptoFrame{}
	offset, err := utils.ReadVarInt(r)
	if err != nil {
		return nil, err
	}
	frame.Offset = protocol.ByteCount(offset)
	dataLen, err := utils.ReadVarInt(r)
	if err != nil {
		return nil, err
	}
	if dataLen > uint64(r.Len()) {
		return nil, io.EOF
	}
	if dataLen != 0 {
		frame.Data = make([]byte, dataLen)
		if _, err := io.ReadFull(r, frame.Data); err != nil {
			// this should never happen, since we already checked the dataLen earlier
			return nil, err
		}
	}
	return frame, nil
}

func (f *CryptoFrame) Write(b *bytes.Buffer, _ protocol.VersionNumber) error {
	b.WriteByte(0x6)
	utils.WriteVarInt(b, uint64(f.Offset))
	utils.WriteVarInt(b, uint64(len(f.Data)))
	b.Write(f.Data)
	return nil
}

// Length of a written frame
func (f *CryptoFrame) Length(_ protocol.VersionNumber) protocol.ByteCount {
	return 1 + utils.VarIntLen(uint64(f.Offset)) + utils.VarIntLen(uint64(len(f.Data))) + protocol.ByteCount(len(f.Data))
}

// MaxDataLen returns the maximum data length
func (f *CryptoFrame) MaxDataLen(maxSize protocol.ByteCount) protocol.ByteCount {
	// pretend that the data size will be 1 bytes
	// if it turns out that varint encoding the length will consume 2 bytes, we need to adjust the data length afterwards
	headerLen := 1 + utils.VarIntLen(uint64(f.Offset)) + 1
	if headerLen > maxSize {
		return 0
	}
	maxDataLen := maxSize - headerLen
	if utils.VarIntLen(uint64(maxDataLen)) != 1 {
		maxDataLen--
	}
	return maxDataLen
}
