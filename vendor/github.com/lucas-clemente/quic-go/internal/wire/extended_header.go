package wire

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	"github.com/lucas-clemente/quic-go/internal/protocol"
	"github.com/lucas-clemente/quic-go/internal/utils"
)

// ExtendedHeader is the header of a QUIC packet.
type ExtendedHeader struct {
	Header

	typeByte byte

	PacketNumberLen protocol.PacketNumberLen
	PacketNumber    protocol.PacketNumber

	KeyPhase int
}

func (h *ExtendedHeader) parse(b *bytes.Reader, v protocol.VersionNumber) (*ExtendedHeader, error) {
	// read the (now unencrypted) first byte
	var err error
	h.typeByte, err = b.ReadByte()
	if err != nil {
		return nil, err
	}
	if _, err := b.Seek(int64(h.ParsedLen())-1, io.SeekCurrent); err != nil {
		return nil, err
	}
	if h.IsLongHeader {
		return h.parseLongHeader(b, v)
	}
	return h.parseShortHeader(b, v)
}

func (h *ExtendedHeader) parseLongHeader(b *bytes.Reader, v protocol.VersionNumber) (*ExtendedHeader, error) {
	if h.typeByte&0xc != 0 {
		return nil, errors.New("5th and 6th bit must be 0")
	}
	if err := h.readPacketNumber(b); err != nil {
		return nil, err
	}
	return h, nil
}

func (h *ExtendedHeader) parseShortHeader(b *bytes.Reader, v protocol.VersionNumber) (*ExtendedHeader, error) {
	if h.typeByte&0x18 != 0 {
		return nil, errors.New("4th and 5th bit must be 0")
	}

	h.KeyPhase = int(h.typeByte&0x4) >> 2

	if err := h.readPacketNumber(b); err != nil {
		return nil, err
	}
	return h, nil
}

func (h *ExtendedHeader) readPacketNumber(b *bytes.Reader) error {
	h.PacketNumberLen = protocol.PacketNumberLen(h.typeByte&0x3) + 1
	pn, err := utils.BigEndian.ReadUintN(b, uint8(h.PacketNumberLen))
	if err != nil {
		return err
	}
	h.PacketNumber = protocol.PacketNumber(pn)
	return nil
}

// Write writes the Header.
func (h *ExtendedHeader) Write(b *bytes.Buffer, ver protocol.VersionNumber) error {
	if h.IsLongHeader {
		return h.writeLongHeader(b, ver)
	}
	return h.writeShortHeader(b, ver)
}

func (h *ExtendedHeader) writeLongHeader(b *bytes.Buffer, v protocol.VersionNumber) error {
	var packetType uint8
	switch h.Type {
	case protocol.PacketTypeInitial:
		packetType = 0x0
	case protocol.PacketType0RTT:
		packetType = 0x1
	case protocol.PacketTypeHandshake:
		packetType = 0x2
	case protocol.PacketTypeRetry:
		packetType = 0x3
	}
	firstByte := 0xc0 | packetType<<4
	if h.Type == protocol.PacketTypeRetry {
		odcil, err := encodeSingleConnIDLen(h.OrigDestConnectionID)
		if err != nil {
			return err
		}
		firstByte |= odcil
	} else { // Retry packets don't have a packet number
		firstByte |= uint8(h.PacketNumberLen - 1)
	}

	b.WriteByte(firstByte)
	utils.BigEndian.WriteUint32(b, uint32(h.Version))
	connIDLen, err := encodeConnIDLen(h.DestConnectionID, h.SrcConnectionID)
	if err != nil {
		return err
	}
	b.WriteByte(connIDLen)
	b.Write(h.DestConnectionID.Bytes())
	b.Write(h.SrcConnectionID.Bytes())

	switch h.Type {
	case protocol.PacketTypeRetry:
		b.Write(h.OrigDestConnectionID.Bytes())
		b.Write(h.Token)
		return nil
	case protocol.PacketTypeInitial:
		utils.WriteVarInt(b, uint64(len(h.Token)))
		b.Write(h.Token)
	}

	utils.WriteVarInt(b, uint64(h.Length))
	return h.writePacketNumber(b)
}

// TODO: add support for the key phase
func (h *ExtendedHeader) writeShortHeader(b *bytes.Buffer, v protocol.VersionNumber) error {
	typeByte := 0x40 | uint8(h.PacketNumberLen-1)
	typeByte |= byte(h.KeyPhase << 2)

	b.WriteByte(typeByte)
	b.Write(h.DestConnectionID.Bytes())
	return h.writePacketNumber(b)
}

func (h *ExtendedHeader) writePacketNumber(b *bytes.Buffer) error {
	if h.PacketNumberLen == protocol.PacketNumberLenInvalid || h.PacketNumberLen > protocol.PacketNumberLen4 {
		return fmt.Errorf("invalid packet number length: %d", h.PacketNumberLen)
	}
	utils.BigEndian.WriteUintN(b, uint8(h.PacketNumberLen), uint64(h.PacketNumber))
	return nil
}

// GetLength determines the length of the Header.
func (h *ExtendedHeader) GetLength(v protocol.VersionNumber) protocol.ByteCount {
	if h.IsLongHeader {
		length := 1 /* type byte */ + 4 /* version */ + 1 /* conn id len byte */ + protocol.ByteCount(h.DestConnectionID.Len()+h.SrcConnectionID.Len()) + protocol.ByteCount(h.PacketNumberLen) + utils.VarIntLen(uint64(h.Length))
		if h.Type == protocol.PacketTypeInitial {
			length += utils.VarIntLen(uint64(len(h.Token))) + protocol.ByteCount(len(h.Token))
		}
		return length
	}

	length := protocol.ByteCount(1 /* type byte */ + h.DestConnectionID.Len())
	length += protocol.ByteCount(h.PacketNumberLen)
	return length
}

// Log logs the Header
func (h *ExtendedHeader) Log(logger utils.Logger) {
	if h.IsLongHeader {
		var token string
		if h.Type == protocol.PacketTypeInitial || h.Type == protocol.PacketTypeRetry {
			if len(h.Token) == 0 {
				token = "Token: (empty), "
			} else {
				token = fmt.Sprintf("Token: %#x, ", h.Token)
			}
			if h.Type == protocol.PacketTypeRetry {
				logger.Debugf("\tLong Header{Type: %s, DestConnectionID: %s, SrcConnectionID: %s, %sOrigDestConnectionID: %s, Version: %s}", h.Type, h.DestConnectionID, h.SrcConnectionID, token, h.OrigDestConnectionID, h.Version)
				return
			}
		}
		logger.Debugf("\tLong Header{Type: %s, DestConnectionID: %s, SrcConnectionID: %s, %sPacketNumber: %#x, PacketNumberLen: %d, Length: %d, Version: %s}", h.Type, h.DestConnectionID, h.SrcConnectionID, token, h.PacketNumber, h.PacketNumberLen, h.Length, h.Version)
	} else {
		logger.Debugf("\tShort Header{DestConnectionID: %s, PacketNumber: %#x, PacketNumberLen: %d, KeyPhase: %d}", h.DestConnectionID, h.PacketNumber, h.PacketNumberLen, h.KeyPhase)
	}
}

func encodeConnIDLen(dest, src protocol.ConnectionID) (byte, error) {
	dcil, err := encodeSingleConnIDLen(dest)
	if err != nil {
		return 0, err
	}
	scil, err := encodeSingleConnIDLen(src)
	if err != nil {
		return 0, err
	}
	return scil | dcil<<4, nil
}

func encodeSingleConnIDLen(id protocol.ConnectionID) (byte, error) {
	len := id.Len()
	if len == 0 {
		return 0, nil
	}
	if len < 4 || len > 18 {
		return 0, fmt.Errorf("invalid connection ID length: %d bytes", len)
	}
	return byte(len - 3), nil
}
