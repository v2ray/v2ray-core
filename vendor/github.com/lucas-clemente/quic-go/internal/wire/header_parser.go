package wire

import (
	"bytes"
	"fmt"
	"io"

	"github.com/lucas-clemente/quic-go/internal/protocol"
	"github.com/lucas-clemente/quic-go/internal/qerr"
	"github.com/lucas-clemente/quic-go/internal/utils"
)

// The InvariantHeader is the version independent part of the header
type InvariantHeader struct {
	IsLongHeader     bool
	Version          protocol.VersionNumber
	SrcConnectionID  protocol.ConnectionID
	DestConnectionID protocol.ConnectionID

	typeByte byte
}

// ParseInvariantHeader parses the version independent part of the header
func ParseInvariantHeader(b *bytes.Reader, shortHeaderConnIDLen int) (*InvariantHeader, error) {
	typeByte, err := b.ReadByte()
	if err != nil {
		return nil, err
	}

	h := &InvariantHeader{typeByte: typeByte}
	h.IsLongHeader = typeByte&0x80 > 0

	// If this is not a Long Header, it could either be a Public Header or a Short Header.
	if !h.IsLongHeader {
		var err error
		h.DestConnectionID, err = protocol.ReadConnectionID(b, shortHeaderConnIDLen)
		if err != nil {
			return nil, err
		}
		return h, nil
	}
	// Long Header
	v, err := utils.BigEndian.ReadUint32(b)
	if err != nil {
		return nil, err
	}
	h.Version = protocol.VersionNumber(v)
	connIDLenByte, err := b.ReadByte()
	if err != nil {
		return nil, err
	}
	dcil, scil := decodeConnIDLen(connIDLenByte)
	h.DestConnectionID, err = protocol.ReadConnectionID(b, dcil)
	if err != nil {
		return nil, err
	}
	h.SrcConnectionID, err = protocol.ReadConnectionID(b, scil)
	if err != nil {
		return nil, err
	}
	return h, nil
}

// Parse parses the version dependent part of the header
func (iv *InvariantHeader) Parse(b *bytes.Reader, sentBy protocol.Perspective, ver protocol.VersionNumber) (*Header, error) {
	if iv.IsLongHeader {
		if iv.Version == 0 { // Version Negotiation Packet
			return iv.parseVersionNegotiationPacket(b)
		}
		return iv.parseLongHeader(b, sentBy, ver)
	}
	return iv.parseShortHeader(b, ver)
}

func (iv *InvariantHeader) toHeader() *Header {
	return &Header{
		IsLongHeader:     iv.IsLongHeader,
		DestConnectionID: iv.DestConnectionID,
		SrcConnectionID:  iv.SrcConnectionID,
		Version:          iv.Version,
	}
}

func (iv *InvariantHeader) parseVersionNegotiationPacket(b *bytes.Reader) (*Header, error) {
	h := iv.toHeader()
	if b.Len() == 0 {
		return nil, qerr.Error(qerr.InvalidVersionNegotiationPacket, "empty version list")
	}
	h.IsVersionNegotiation = true
	h.SupportedVersions = make([]protocol.VersionNumber, b.Len()/4)
	for i := 0; b.Len() > 0; i++ {
		v, err := utils.BigEndian.ReadUint32(b)
		if err != nil {
			return nil, qerr.InvalidVersionNegotiationPacket
		}
		h.SupportedVersions[i] = protocol.VersionNumber(v)
	}
	return h, nil
}

func (iv *InvariantHeader) parseLongHeader(b *bytes.Reader, sentBy protocol.Perspective, v protocol.VersionNumber) (*Header, error) {
	h := iv.toHeader()
	h.Type = protocol.PacketType(iv.typeByte & 0x7f)

	if h.Type != protocol.PacketTypeInitial && h.Type != protocol.PacketTypeRetry && h.Type != protocol.PacketType0RTT && h.Type != protocol.PacketTypeHandshake {
		return nil, qerr.Error(qerr.InvalidPacketHeader, fmt.Sprintf("Received packet with invalid packet type: %d", h.Type))
	}

	if h.Type == protocol.PacketTypeRetry {
		odcilByte, err := b.ReadByte()
		if err != nil {
			return nil, err
		}
		odcil := decodeSingleConnIDLen(odcilByte & 0xf)
		h.OrigDestConnectionID, err = protocol.ReadConnectionID(b, odcil)
		if err != nil {
			return nil, err
		}
		h.Token = make([]byte, b.Len())
		if _, err := io.ReadFull(b, h.Token); err != nil {
			return nil, err
		}
		return h, nil
	}

	if h.Type == protocol.PacketTypeInitial {
		tokenLen, err := utils.ReadVarInt(b)
		if err != nil {
			return nil, err
		}
		if tokenLen > uint64(b.Len()) {
			return nil, io.EOF
		}
		h.Token = make([]byte, tokenLen)
		if _, err := io.ReadFull(b, h.Token); err != nil {
			return nil, err
		}
	}

	pl, err := utils.ReadVarInt(b)
	if err != nil {
		return nil, err
	}
	h.Length = protocol.ByteCount(pl)
	pn, pnLen, err := utils.ReadVarIntPacketNumber(b)
	if err != nil {
		return nil, err
	}
	h.PacketNumber = pn
	h.PacketNumberLen = pnLen

	return h, nil
}

func (iv *InvariantHeader) parseShortHeader(b *bytes.Reader, v protocol.VersionNumber) (*Header, error) {
	h := iv.toHeader()
	h.KeyPhase = int(iv.typeByte&0x40) >> 6

	pn, pnLen, err := utils.ReadVarIntPacketNumber(b)
	if err != nil {
		return nil, err
	}
	h.PacketNumber = pn
	h.PacketNumberLen = pnLen

	return h, nil
}
