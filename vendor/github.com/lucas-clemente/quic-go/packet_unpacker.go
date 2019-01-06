package quic

import (
	"bytes"
	"fmt"

	"github.com/lucas-clemente/quic-go/internal/handshake"
	"github.com/lucas-clemente/quic-go/internal/protocol"
	"github.com/lucas-clemente/quic-go/internal/qerr"
	"github.com/lucas-clemente/quic-go/internal/utils"
	"github.com/lucas-clemente/quic-go/internal/wire"
)

type unpackedPacket struct {
	packetNumber    protocol.PacketNumber // the decoded packet number
	hdr             *wire.ExtendedHeader
	encryptionLevel protocol.EncryptionLevel
	frames          []wire.Frame
}

// The packetUnpacker unpacks QUIC packets.
type packetUnpacker struct {
	cs handshake.CryptoSetup

	largestRcvdPacketNumber protocol.PacketNumber

	version protocol.VersionNumber
}

var _ unpacker = &packetUnpacker{}

func newPacketUnpacker(cs handshake.CryptoSetup, version protocol.VersionNumber) unpacker {
	return &packetUnpacker{
		cs:      cs,
		version: version,
	}
}

func (u *packetUnpacker) Unpack(hdr *wire.Header, data []byte) (*unpackedPacket, error) {
	r := bytes.NewReader(data)

	var encLevel protocol.EncryptionLevel
	switch hdr.Type {
	case protocol.PacketTypeInitial:
		encLevel = protocol.EncryptionInitial
	case protocol.PacketTypeHandshake:
		encLevel = protocol.EncryptionHandshake
	default:
		if hdr.IsLongHeader {
			return nil, fmt.Errorf("unknown packet type: %s", hdr.Type)
		}
		encLevel = protocol.Encryption1RTT
	}
	opener, err := u.cs.GetOpener(encLevel)
	if err != nil {
		return nil, err
	}
	hdrLen := int(hdr.ParsedLen())
	// The packet number can be up to 4 bytes long, but we won't know the length until we decrypt it.
	// 1. save a copy of the 4 bytes
	origPNBytes := make([]byte, 4)
	copy(origPNBytes, data[hdrLen:hdrLen+4])
	// 2. decrypt the header, assuming a 4 byte packet number
	opener.DecryptHeader(
		data[hdrLen+4:hdrLen+4+16],
		&data[0],
		data[hdrLen:hdrLen+4],
	)

	// 3. parse the header (and learn the actual length of the packet number)
	extHdr, err := hdr.ParseExtended(r, u.version)
	if err != nil {
		return nil, fmt.Errorf("error parsing extended header: %s", err)
	}
	extHdr.Raw = data[:hdrLen+int(extHdr.PacketNumberLen)]
	// 4. if the packet number is shorter than 4 bytes, replace the remaining bytes with the copy we saved earlier
	if extHdr.PacketNumberLen != protocol.PacketNumberLen4 {
		copy(data[hdrLen+int(extHdr.PacketNumberLen):hdrLen+4], origPNBytes[int(extHdr.PacketNumberLen):])
	}
	data = data[hdrLen+int(extHdr.PacketNumberLen):]

	pn := protocol.DecodePacketNumber(
		extHdr.PacketNumberLen,
		u.largestRcvdPacketNumber,
		extHdr.PacketNumber,
	)

	decrypted, err := opener.Open(data[:0], data, pn, extHdr.Raw)
	if err != nil {
		return nil, err
	}

	// Only do this after decrypting, so we are sure the packet is not attacker-controlled
	u.largestRcvdPacketNumber = utils.MaxPacketNumber(u.largestRcvdPacketNumber, pn)

	fs, err := u.parseFrames(decrypted)
	if err != nil {
		return nil, err
	}

	return &unpackedPacket{
		hdr:             extHdr,
		packetNumber:    pn,
		encryptionLevel: encLevel,
		frames:          fs,
	}, nil
}

func (u *packetUnpacker) parseFrames(decrypted []byte) ([]wire.Frame, error) {
	r := bytes.NewReader(decrypted)
	if r.Len() == 0 {
		return nil, qerr.MissingPayload
	}

	fs := make([]wire.Frame, 0, 2)
	// Read all frames in the packet
	for {
		frame, err := wire.ParseNextFrame(r, u.version)
		if err != nil {
			return nil, err
		}
		if frame == nil {
			break
		}
		fs = append(fs, frame)
	}
	return fs, nil
}
