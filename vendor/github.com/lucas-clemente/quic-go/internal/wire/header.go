package wire

import (
	"bytes"
	"crypto/rand"
	"fmt"

	"github.com/lucas-clemente/quic-go/internal/protocol"
	"github.com/lucas-clemente/quic-go/internal/utils"
)

// Header is the header of a QUIC packet.
type Header struct {
	Raw []byte

	Version protocol.VersionNumber

	DestConnectionID     protocol.ConnectionID
	SrcConnectionID      protocol.ConnectionID
	OrigDestConnectionID protocol.ConnectionID // only needed in the Retry packet

	PacketNumberLen protocol.PacketNumberLen
	PacketNumber    protocol.PacketNumber

	IsVersionNegotiation bool
	SupportedVersions    []protocol.VersionNumber // Version Number sent in a Version Negotiation Packet by the server

	Type         protocol.PacketType
	IsLongHeader bool
	KeyPhase     int
	PayloadLen   protocol.ByteCount
	Token        []byte
}

// Write writes the Header.
func (h *Header) Write(b *bytes.Buffer, pers protocol.Perspective, ver protocol.VersionNumber) error {
	if h.IsLongHeader {
		return h.writeLongHeader(b, ver)
	}
	return h.writeShortHeader(b, ver)
}

// TODO: add support for the key phase
func (h *Header) writeLongHeader(b *bytes.Buffer, v protocol.VersionNumber) error {
	b.WriteByte(byte(0x80 | h.Type))
	utils.BigEndian.WriteUint32(b, uint32(h.Version))
	connIDLen, err := encodeConnIDLen(h.DestConnectionID, h.SrcConnectionID)
	if err != nil {
		return err
	}
	b.WriteByte(connIDLen)
	b.Write(h.DestConnectionID.Bytes())
	b.Write(h.SrcConnectionID.Bytes())

	if h.Type == protocol.PacketTypeInitial {
		utils.WriteVarInt(b, uint64(len(h.Token)))
		b.Write(h.Token)
	}

	if h.Type == protocol.PacketTypeRetry {
		odcil, err := encodeSingleConnIDLen(h.OrigDestConnectionID)
		if err != nil {
			return err
		}
		// randomize the first 4 bits
		odcilByte := make([]byte, 1)
		_, _ = rand.Read(odcilByte) // it's safe to ignore the error here
		odcilByte[0] = (odcilByte[0] & 0xf0) | odcil
		b.Write(odcilByte)
		b.Write(h.OrigDestConnectionID.Bytes())
		b.Write(h.Token)
		return nil
	}

	utils.WriteVarInt(b, uint64(h.PayloadLen))
	return utils.WriteVarIntPacketNumber(b, h.PacketNumber, h.PacketNumberLen)
}

func (h *Header) writeShortHeader(b *bytes.Buffer, v protocol.VersionNumber) error {
	typeByte := byte(0x30)
	typeByte |= byte(h.KeyPhase << 6)

	b.WriteByte(typeByte)
	b.Write(h.DestConnectionID.Bytes())
	return utils.WriteVarIntPacketNumber(b, h.PacketNumber, h.PacketNumberLen)
}

// GetLength determines the length of the Header.
func (h *Header) GetLength(v protocol.VersionNumber) protocol.ByteCount {
	if h.IsLongHeader {
		length := 1 /* type byte */ + 4 /* version */ + 1 /* conn id len byte */ + protocol.ByteCount(h.DestConnectionID.Len()+h.SrcConnectionID.Len()) + protocol.ByteCount(h.PacketNumberLen) + utils.VarIntLen(uint64(h.PayloadLen))
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
func (h *Header) Log(logger utils.Logger) {
	if h.IsLongHeader {
		if h.Version == 0 {
			logger.Debugf("\tVersionNegotiationPacket{DestConnectionID: %s, SrcConnectionID: %s, SupportedVersions: %s}", h.DestConnectionID, h.SrcConnectionID, h.SupportedVersions)
		} else {
			var token string
			if h.Type == protocol.PacketTypeInitial || h.Type == protocol.PacketTypeRetry {
				if len(h.Token) == 0 {
					token = "Token: (empty), "
				} else {
					token = fmt.Sprintf("Token: %#x, ", h.Token)
				}
			}
			if h.Type == protocol.PacketTypeRetry {
				logger.Debugf("\tLong Header{Type: %s, DestConnectionID: %s, SrcConnectionID: %s, %sOrigDestConnectionID: %s, Version: %s}", h.Type, h.DestConnectionID, h.SrcConnectionID, token, h.OrigDestConnectionID, h.Version)
				return
			}
			logger.Debugf("\tLong Header{Type: %s, DestConnectionID: %s, SrcConnectionID: %s, %sPacketNumber: %#x, PacketNumberLen: %d, PayloadLen: %d, Version: %s}", h.Type, h.DestConnectionID, h.SrcConnectionID, token, h.PacketNumber, h.PacketNumberLen, h.PayloadLen, h.Version)
		}
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

func decodeConnIDLen(enc byte) (int /*dest conn id len*/, int /*src conn id len*/) {
	return decodeSingleConnIDLen(enc >> 4), decodeSingleConnIDLen(enc & 0xf)
}

func decodeSingleConnIDLen(enc uint8) int {
	if enc == 0 {
		return 0
	}
	return int(enc) + 3
}
