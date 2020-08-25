package quic

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"time"

	"v2ray.com/core/external/github.com/lucas-clemente/quic-go/internal/ackhandler"
	"v2ray.com/core/external/github.com/lucas-clemente/quic-go/internal/handshake"
	"v2ray.com/core/external/github.com/lucas-clemente/quic-go/internal/protocol"
	"v2ray.com/core/external/github.com/lucas-clemente/quic-go/internal/utils"
	"v2ray.com/core/external/github.com/lucas-clemente/quic-go/internal/wire"
)

type packer interface {
	PackPacket() (*packedPacket, error)
	MaybePackAckPacket() (*packedPacket, error)
	PackRetransmission(packet *ackhandler.Packet) ([]*packedPacket, error)
	PackConnectionClose(*wire.ConnectionCloseFrame) (*packedPacket, error)

	HandleTransportParameters(*handshake.TransportParameters)
	ChangeDestConnectionID(protocol.ConnectionID)
}

type packedPacket struct {
	header *wire.ExtendedHeader
	raw    []byte
	frames []wire.Frame

	buffer *packetBuffer
}

func (p *packedPacket) EncryptionLevel() protocol.EncryptionLevel {
	if !p.header.IsLongHeader {
		return protocol.Encryption1RTT
	}
	switch p.header.Type {
	case protocol.PacketTypeInitial:
		return protocol.EncryptionInitial
	case protocol.PacketTypeHandshake:
		return protocol.EncryptionHandshake
	default:
		return protocol.EncryptionUnspecified
	}
}

func (p *packedPacket) ToAckHandlerPacket() *ackhandler.Packet {
	return &ackhandler.Packet{
		PacketNumber:    p.header.PacketNumber,
		PacketType:      p.header.Type,
		Frames:          p.frames,
		Length:          protocol.ByteCount(len(p.raw)),
		EncryptionLevel: p.EncryptionLevel(),
		SendTime:        time.Now(),
	}
}

func getMaxPacketSize(addr net.Addr) protocol.ByteCount {
	maxSize := protocol.ByteCount(protocol.MinInitialPacketSize)
	// If this is not a UDP address, we don't know anything about the MTU.
	// Use the minimum size of an Initial packet as the max packet size.
	if udpAddr, ok := addr.(*net.UDPAddr); ok {
		// If ip is not an IPv4 address, To4 returns nil.
		// Note that there might be some corner cases, where this is not correct.
		// See https://stackoverflow.com/questions/22751035/golang-distinguish-ipv4-ipv6.
		if udpAddr.IP.To4() == nil {
			maxSize = protocol.MaxPacketSizeIPv6
		} else {
			maxSize = protocol.MaxPacketSizeIPv4
		}
	}
	return maxSize
}

type packetNumberManager interface {
	PeekPacketNumber() (protocol.PacketNumber, protocol.PacketNumberLen)
	PopPacketNumber() protocol.PacketNumber
}

type sealingManager interface {
	GetSealer() (protocol.EncryptionLevel, handshake.Sealer)
	GetSealerWithEncryptionLevel(protocol.EncryptionLevel) (handshake.Sealer, error)
}

type frameSource interface {
	AppendStreamFrames([]wire.Frame, protocol.ByteCount) []wire.Frame
	AppendControlFrames([]wire.Frame, protocol.ByteCount) ([]wire.Frame, protocol.ByteCount)
}

type ackFrameSource interface {
	GetAckFrame(protocol.EncryptionLevel) *wire.AckFrame
}

type packetPacker struct {
	destConnID protocol.ConnectionID
	srcConnID  protocol.ConnectionID

	perspective protocol.Perspective
	version     protocol.VersionNumber
	cryptoSetup sealingManager

	initialStream   cryptoStream
	handshakeStream cryptoStream

	token []byte

	pnManager packetNumberManager
	framer    frameSource
	acks      ackFrameSource

	maxPacketSize             protocol.ByteCount
	numNonRetransmittableAcks int
}

var _ packer = &packetPacker{}

func newPacketPacker(
	destConnID protocol.ConnectionID,
	srcConnID protocol.ConnectionID,
	initialStream cryptoStream,
	handshakeStream cryptoStream,
	packetNumberManager packetNumberManager,
	remoteAddr net.Addr, // only used for determining the max packet size
	token []byte,
	cryptoSetup sealingManager,
	framer frameSource,
	acks ackFrameSource,
	perspective protocol.Perspective,
	version protocol.VersionNumber,
) *packetPacker {
	return &packetPacker{
		cryptoSetup:     cryptoSetup,
		token:           token,
		destConnID:      destConnID,
		srcConnID:       srcConnID,
		initialStream:   initialStream,
		handshakeStream: handshakeStream,
		perspective:     perspective,
		version:         version,
		framer:          framer,
		acks:            acks,
		pnManager:       packetNumberManager,
		maxPacketSize:   getMaxPacketSize(remoteAddr),
	}
}

// PackConnectionClose packs a packet that ONLY contains a ConnectionCloseFrame
func (p *packetPacker) PackConnectionClose(ccf *wire.ConnectionCloseFrame) (*packedPacket, error) {
	frames := []wire.Frame{ccf}
	encLevel, sealer := p.cryptoSetup.GetSealer()
	header := p.getHeader(encLevel)
	return p.writeAndSealPacket(header, frames, sealer)
}

func (p *packetPacker) MaybePackAckPacket() (*packedPacket, error) {
	ack := p.acks.GetAckFrame(protocol.Encryption1RTT)
	if ack == nil {
		return nil, nil
	}
	// TODO(#1534): only pack ACKs with the right encryption level
	encLevel, sealer := p.cryptoSetup.GetSealer()
	header := p.getHeader(encLevel)
	frames := []wire.Frame{ack}
	return p.writeAndSealPacket(header, frames, sealer)
}

// PackRetransmission packs a retransmission
// For packets sent after completion of the handshake, it might happen that 2 packets have to be sent.
// This can happen e.g. when a longer packet number is used in the header.
func (p *packetPacker) PackRetransmission(packet *ackhandler.Packet) ([]*packedPacket, error) {
	var controlFrames []wire.Frame
	var streamFrames []*wire.StreamFrame
	for _, f := range packet.Frames {
		// CRYPTO frames are treated as control frames here.
		// Since we're making sure that the header can never be larger for a retransmission,
		// we never have to split CRYPTO frames.
		if sf, ok := f.(*wire.StreamFrame); ok {
			sf.DataLenPresent = true
			streamFrames = append(streamFrames, sf)
		} else {
			controlFrames = append(controlFrames, f)
		}
	}

	var packets []*packedPacket
	encLevel := packet.EncryptionLevel
	sealer, err := p.cryptoSetup.GetSealerWithEncryptionLevel(encLevel)
	if err != nil {
		return nil, err
	}
	for len(controlFrames) > 0 || len(streamFrames) > 0 {
		var frames []wire.Frame
		var length protocol.ByteCount

		header := p.getHeader(encLevel)
		headerLen := header.GetLength(p.version)
		maxSize := p.maxPacketSize - protocol.ByteCount(sealer.Overhead()) - headerLen

		for len(controlFrames) > 0 {
			frame := controlFrames[0]
			frameLen := frame.Length(p.version)
			if length+frameLen > maxSize {
				break
			}
			length += frameLen
			frames = append(frames, frame)
			controlFrames = controlFrames[1:]
		}

		for len(streamFrames) > 0 && length+protocol.MinStreamFrameSize < maxSize {
			frame := streamFrames[0]
			frame.DataLenPresent = false
			frameToAdd := frame

			sf, err := frame.MaybeSplitOffFrame(maxSize-length, p.version)
			if err != nil {
				return nil, err
			}
			if sf != nil {
				frameToAdd = sf
			} else {
				streamFrames = streamFrames[1:]
			}
			frame.DataLenPresent = true
			length += frameToAdd.Length(p.version)
			frames = append(frames, frameToAdd)
		}
		if sf, ok := frames[len(frames)-1].(*wire.StreamFrame); ok {
			sf.DataLenPresent = false
		}
		p, err := p.writeAndSealPacket(header, frames, sealer)
		if err != nil {
			return nil, err
		}
		packets = append(packets, p)
	}
	return packets, nil
}

// PackPacket packs a new packet
// the other controlFrames are sent in the next packet, but might be queued and sent in the next packet if the packet would overflow MaxPacketSize otherwise
func (p *packetPacker) PackPacket() (*packedPacket, error) {
	packet, err := p.maybePackCryptoPacket()
	if err != nil {
		return nil, err
	}
	if packet != nil {
		return packet, nil
	}

	encLevel, sealer := p.cryptoSetup.GetSealer()
	header := p.getHeader(encLevel)
	headerLen := header.GetLength(p.version)
	if err != nil {
		return nil, err
	}

	maxSize := p.maxPacketSize - protocol.ByteCount(sealer.Overhead()) - headerLen
	frames, err := p.composeNextPacket(maxSize)
	if err != nil {
		return nil, err
	}

	// Check if we have enough frames to send
	if len(frames) == 0 {
		return nil, nil
	}
	// check if this packet only contains an ACK
	if !ackhandler.HasRetransmittableFrames(frames) {
		if p.numNonRetransmittableAcks >= protocol.MaxNonRetransmittableAcks {
			frames = append(frames, &wire.PingFrame{})
			p.numNonRetransmittableAcks = 0
		} else {
			p.numNonRetransmittableAcks++
		}
	} else {
		p.numNonRetransmittableAcks = 0
	}

	return p.writeAndSealPacket(header, frames, sealer)
}

func (p *packetPacker) maybePackCryptoPacket() (*packedPacket, error) {
	var s cryptoStream
	var encLevel protocol.EncryptionLevel

	hasData := p.initialStream.HasData()
	ack := p.acks.GetAckFrame(protocol.EncryptionInitial)
	if hasData || ack != nil {
		s = p.initialStream
		encLevel = protocol.EncryptionInitial
	} else {
		hasData = p.handshakeStream.HasData()
		ack = p.acks.GetAckFrame(protocol.EncryptionHandshake)
		if hasData || ack != nil {
			s = p.handshakeStream
			encLevel = protocol.EncryptionHandshake
		}
	}
	if s == nil {
		return nil, nil
	}
	sealer, err := p.cryptoSetup.GetSealerWithEncryptionLevel(encLevel)
	if err != nil {
		// The sealer
		return nil, err
	}

	hdr := p.getHeader(encLevel)
	hdrLen := hdr.GetLength(p.version)
	var length protocol.ByteCount
	frames := make([]wire.Frame, 0, 2)
	if ack != nil {
		frames = append(frames, ack)
		length += ack.Length(p.version)
	}
	if hasData {
		cf := s.PopCryptoFrame(p.maxPacketSize - hdrLen - protocol.ByteCount(sealer.Overhead()) - length)
		frames = append(frames, cf)
	}
	return p.writeAndSealPacket(hdr, frames, sealer)
}

func (p *packetPacker) composeNextPacket(maxFrameSize protocol.ByteCount) ([]wire.Frame, error) {
	var length protocol.ByteCount
	var frames []wire.Frame

	// ACKs need to go first, so that the sentPacketHandler will recognize them
	if ack := p.acks.GetAckFrame(protocol.Encryption1RTT); ack != nil {
		frames = append(frames, ack)
		length += ack.Length(p.version)
	}

	var lengthAdded protocol.ByteCount
	frames, lengthAdded = p.framer.AppendControlFrames(frames, maxFrameSize-length)
	length += lengthAdded

	// temporarily increase the maxFrameSize by the (minimum) length of the DataLen field
	// this leads to a properly sized packet in all cases, since we do all the packet length calculations with STREAM frames that have the DataLen set
	// however, for the last STREAM frame in the packet, we can omit the DataLen, thus yielding a packet of exactly the correct size
	// the length is encoded to either 1 or 2 bytes
	maxFrameSize++

	frames = p.framer.AppendStreamFrames(frames, maxFrameSize-length)
	if len(frames) > 0 {
		lastFrame := frames[len(frames)-1]
		if sf, ok := lastFrame.(*wire.StreamFrame); ok {
			sf.DataLenPresent = false
		}
	}
	return frames, nil
}

func (p *packetPacker) getHeader(encLevel protocol.EncryptionLevel) *wire.ExtendedHeader {
	pn, pnLen := p.pnManager.PeekPacketNumber()
	header := &wire.ExtendedHeader{}
	header.PacketNumber = pn
	header.PacketNumberLen = pnLen
	header.Version = p.version
	header.DestConnectionID = p.destConnID

	if encLevel != protocol.Encryption1RTT {
		header.IsLongHeader = true
		// Always send Initial and Handshake packets with the maximum packet number length.
		// This simplifies retransmissions: Since the header can't get any larger,
		// we don't need to split CRYPTO frames.
		header.PacketNumberLen = protocol.PacketNumberLen4
		header.SrcConnectionID = p.srcConnID
		// Set the length to the maximum packet size.
		// Since it is encoded as a varint, this guarantees us that the header will end up at most as big as GetLength() returns.
		header.Length = p.maxPacketSize
		switch encLevel {
		case protocol.EncryptionInitial:
			header.Type = protocol.PacketTypeInitial
		case protocol.EncryptionHandshake:
			header.Type = protocol.PacketTypeHandshake
		}
	}

	return header
}

func (p *packetPacker) writeAndSealPacket(
	header *wire.ExtendedHeader,
	frames []wire.Frame,
	sealer handshake.Sealer,
) (*packedPacket, error) {
	packetBuffer := getPacketBuffer()
	buffer := bytes.NewBuffer(packetBuffer.Slice[:0])

	addPaddingForInitial := p.perspective == protocol.PerspectiveClient && header.Type == protocol.PacketTypeInitial

	if header.IsLongHeader {
		if p.perspective == protocol.PerspectiveClient && header.Type == protocol.PacketTypeInitial {
			header.Token = p.token
		}
		if addPaddingForInitial {
			headerLen := header.GetLength(p.version)
			header.Length = protocol.ByteCount(header.PacketNumberLen) + protocol.MinInitialPacketSize - headerLen
		} else {
			// long header packets always use 4 byte packet number, so we never need to pad short payloads
			length := protocol.ByteCount(sealer.Overhead()) + protocol.ByteCount(header.PacketNumberLen)
			for _, frame := range frames {
				length += frame.Length(p.version)
			}
			header.Length = length
		}
	}

	if err := header.Write(buffer, p.version); err != nil {
		return nil, err
	}
	payloadOffset := buffer.Len()

	// write all frames but the last one
	for _, frame := range frames[:len(frames)-1] {
		if err := frame.Write(buffer, p.version); err != nil {
			return nil, err
		}
	}
	lastFrame := frames[len(frames)-1]
	if addPaddingForInitial {
		// when appending padding, we need to make sure that the last STREAM frames has the data length set
		if sf, ok := lastFrame.(*wire.StreamFrame); ok {
			sf.DataLenPresent = true
		}
	} else {
		payloadLen := buffer.Len() - payloadOffset + int(lastFrame.Length(p.version))
		if paddingLen := 4 - int(header.PacketNumberLen) - payloadLen; paddingLen > 0 {
			// Pad the packet such that packet number length + payload length is 4 bytes.
			// This is needed to enable the peer to get a 16 byte sample for header protection.
			buffer.Write(bytes.Repeat([]byte{0}, paddingLen))
		}
	}
	if err := lastFrame.Write(buffer, p.version); err != nil {
		return nil, err
	}

	if addPaddingForInitial {
		paddingLen := protocol.MinInitialPacketSize - sealer.Overhead() - buffer.Len()
		if paddingLen > 0 {
			buffer.Write(bytes.Repeat([]byte{0}, paddingLen))
		}
	}

	if size := protocol.ByteCount(buffer.Len() + sealer.Overhead()); size > p.maxPacketSize {
		return nil, fmt.Errorf("PacketPacker BUG: packet too large (%d bytes, allowed %d bytes)", size, p.maxPacketSize)
	}

	raw := buffer.Bytes()
	_ = sealer.Seal(raw[payloadOffset:payloadOffset], raw[payloadOffset:], header.PacketNumber, raw[:payloadOffset])
	raw = raw[0 : buffer.Len()+sealer.Overhead()]

	pnOffset := payloadOffset - int(header.PacketNumberLen)
	sealer.EncryptHeader(
		raw[pnOffset+4:pnOffset+4+16],
		&raw[0],
		raw[pnOffset:payloadOffset],
	)

	num := p.pnManager.PopPacketNumber()
	if num != header.PacketNumber {
		return nil, errors.New("packetPacker BUG: Peeked and Popped packet numbers do not match")
	}
	return &packedPacket{
		header: header,
		raw:    raw,
		frames: frames,
		buffer: packetBuffer,
	}, nil
}

func (p *packetPacker) ChangeDestConnectionID(connID protocol.ConnectionID) {
	p.destConnID = connID
}

func (p *packetPacker) HandleTransportParameters(params *handshake.TransportParameters) {
	if params.MaxPacketSize != 0 {
		p.maxPacketSize = utils.MinByteCount(p.maxPacketSize, params.MaxPacketSize)
	}
}
