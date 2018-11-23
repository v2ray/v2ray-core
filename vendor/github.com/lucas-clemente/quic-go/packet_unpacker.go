package quic

import (
	"bytes"
	"fmt"

	"github.com/lucas-clemente/quic-go/internal/protocol"
	"github.com/lucas-clemente/quic-go/internal/qerr"
	"github.com/lucas-clemente/quic-go/internal/wire"
)

type unpackedPacket struct {
	encryptionLevel protocol.EncryptionLevel
	frames          []wire.Frame
}

type quicAEAD interface {
	OpenInitial(dst, src []byte, pn protocol.PacketNumber, ad []byte) ([]byte, error)
	OpenHandshake(dst, src []byte, pn protocol.PacketNumber, ad []byte) ([]byte, error)
	Open1RTT(dst, src []byte, pn protocol.PacketNumber, ad []byte) ([]byte, error)
}

// The packetUnpacker unpacks QUIC packets.
type packetUnpacker struct {
	aead    quicAEAD
	version protocol.VersionNumber
}

var _ unpacker = &packetUnpacker{}

func newPacketUnpacker(aead quicAEAD, version protocol.VersionNumber) unpacker {
	return &packetUnpacker{
		aead:    aead,
		version: version,
	}
}

func (u *packetUnpacker) Unpack(headerBinary []byte, hdr *wire.Header, data []byte) (*unpackedPacket, error) {
	buf := *getPacketBuffer()
	buf = buf[:0]
	defer putPacketBuffer(&buf)

	var decrypted []byte
	var encryptionLevel protocol.EncryptionLevel
	var err error
	switch hdr.Type {
	case protocol.PacketTypeInitial:
		decrypted, err = u.aead.OpenInitial(buf, data, hdr.PacketNumber, headerBinary)
		encryptionLevel = protocol.EncryptionInitial
	case protocol.PacketTypeHandshake:
		decrypted, err = u.aead.OpenHandshake(buf, data, hdr.PacketNumber, headerBinary)
		encryptionLevel = protocol.EncryptionHandshake
	default:
		if hdr.IsLongHeader {
			return nil, fmt.Errorf("unknown packet type: %s", hdr.Type)
		}
		decrypted, err = u.aead.Open1RTT(buf, data, hdr.PacketNumber, headerBinary)
		encryptionLevel = protocol.Encryption1RTT
	}
	if err != nil {
		// Wrap err in quicError so that public reset is sent by session
		return nil, qerr.Error(qerr.DecryptionFailure, err.Error())
	}

	fs, err := u.parseFrames(decrypted)
	if err != nil {
		return nil, err
	}

	return &unpackedPacket{
		encryptionLevel: encryptionLevel,
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
