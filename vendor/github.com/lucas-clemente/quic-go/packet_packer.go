package quic

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/lucas-clemente/quic-go/internal/ackhandler"
	"github.com/lucas-clemente/quic-go/internal/handshake"
	"github.com/lucas-clemente/quic-go/internal/protocol"
	"github.com/lucas-clemente/quic-go/internal/utils"
	"github.com/lucas-clemente/quic-go/internal/wire"
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
	header          *wire.Header
	raw             []byte
	frames          []wire.Frame
	encryptionLevel protocol.EncryptionLevel
}

func (p *packedPacket) ToAckHandlerPacket() *ackhandler.Packet {
	return &ackhandler.Packet{
		PacketNumber:    p.header.PacketNumber,
		PacketType:      p.header.Type,
		Frames:          p.frames,
		Length:          protocol.ByteCount(len(p.raw)),
		EncryptionLevel: p.encryptionLevel,
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
	GetAckFrame() *wire.AckFrame
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
	hasSentPacket             bool // has the packetPacker already sent a packet
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
	raw, err := p.writeAndSealPacket(header, frames, sealer)
	return &packedPacket{
		header:          header,
		raw:             raw,
		frames:          frames,
		encryptionLevel: encLevel,
	}, err
}

func (p *packetPacker) MaybePackAckPacket() (*packedPacket, error) {
	ack := p.acks.GetAckFrame()
	if ack == nil {
		return nil, nil
	}
	// TODO(#1534): only pack ACKs with the right encryption level
	encLevel, sealer := p.cryptoSetup.GetSealer()
	header := p.getHeader(encLevel)
	frames := []wire.Frame{ack}
	raw, err := p.writeAndSealPacket(header, frames, sealer)
	return &packedPacket{
		header:          header,
		raw:             raw,
		frames:          frames,
		encryptionLevel: encLevel,
	}, err
}

// PackRetransmission packs a retransmission
// For packets sent after completion of the handshake, it might happen that 2 packets have to be sent.
// This can happen e.g. when a longer packet number is used in the header.
func (p *packetPacker) PackRetransmission(packet *ackhandler.Packet) ([]*packedPacket, error) {
	if packet.EncryptionLevel != protocol.Encryption1RTT {
		p, err := p.packHandshakeRetransmission(packet)
		return []*packedPacket{p}, err
	}

	var controlFrames []wire.Frame
	var streamFrames []*wire.StreamFrame
	for _, f := range packet.Frames {
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
		raw, err := p.writeAndSealPacket(header, frames, sealer)
		if err != nil {
			return nil, err
		}
		packets = append(packets, &packedPacket{
			header:          header,
			raw:             raw,
			frames:          frames,
			encryptionLevel: encLevel,
		})
	}
	return packets, nil
}

// packHandshakeRetransmission retransmits a handshake packet
func (p *packetPacker) packHandshakeRetransmission(packet *ackhandler.Packet) (*packedPacket, error) {
	sealer, err := p.cryptoSetup.GetSealerWithEncryptionLevel(packet.EncryptionLevel)
	if err != nil {
		return nil, err
	}
	// make sure that the retransmission for an Initial packet is sent as an Initial packet
	if packet.PacketType == protocol.PacketTypeInitial {
		p.hasSentPacket = false
	}
	header := p.getHeader(packet.EncryptionLevel)
	header.Type = packet.PacketType
	raw, err := p.writeAndSealPacket(header, packet.Frames, sealer)
	return &packedPacket{
		header:          header,
		raw:             raw,
		frames:          packet.Frames,
		encryptionLevel: packet.EncryptionLevel,
	}, err
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
	// if this is the first packet to be send, make sure it contains stream data
	if !p.hasSentPacket && packet == nil {
		return nil, nil
	}

	encLevel, sealer := p.cryptoSetup.GetSealer()
	header := p.getHeader(encLevel)
	headerLen := header.GetLength(p.version)
	if err != nil {
		return nil, err
	}

	maxSize := p.maxPacketSize - protocol.ByteCount(sealer.Overhead()) - headerLen
	frames, err := p.composeNextPacket(maxSize, p.canSendData(encLevel))
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

	raw, err := p.writeAndSealPacket(header, frames, sealer)
	if err != nil {
		return nil, err
	}
	return &packedPacket{
		header:          header,
		raw:             raw,
		frames:          frames,
		encryptionLevel: encLevel,
	}, nil
}

func (p *packetPacker) maybePackCryptoPacket() (*packedPacket, error) {
	var s cryptoStream
	var encLevel protocol.EncryptionLevel
	if p.initialStream.HasData() {
		s = p.initialStream
		encLevel = protocol.EncryptionInitial
	} else if p.handshakeStream.HasData() {
		s = p.handshakeStream
		encLevel = protocol.EncryptionHandshake
	}
	if s == nil {
		return nil, nil
	}
	hdr := p.getHeader(encLevel)
	hdrLen := hdr.GetLength(p.version)
	sealer, err := p.cryptoSetup.GetSealerWithEncryptionLevel(encLevel)
	if err != nil {
		return nil, err
	}
	var length protocol.ByteCount
	frames := make([]wire.Frame, 0, 2)
	if ack := p.acks.GetAckFrame(); ack != nil {
		frames = append(frames, ack)
		length += ack.Length(p.version)
	}
	cf := s.PopCryptoFrame(p.maxPacketSize - hdrLen - protocol.ByteCount(sealer.Overhead()) - length)
	frames = append(frames, cf)
	raw, err := p.writeAndSealPacket(hdr, frames, sealer)
	if err != nil {
		return nil, err
	}
	return &packedPacket{
		header:          hdr,
		raw:             raw,
		frames:          frames,
		encryptionLevel: encLevel,
	}, nil
}

func (p *packetPacker) composeNextPacket(
	maxFrameSize protocol.ByteCount,
	canSendStreamFrames bool,
) ([]wire.Frame, error) {
	var length protocol.ByteCount
	var frames []wire.Frame

	// ACKs need to go first, so that the sentPacketHandler will recognize them
	if ack := p.acks.GetAckFrame(); ack != nil {
		frames = append(frames, ack)
		length += ack.Length(p.version)
	}

	var lengthAdded protocol.ByteCount
	frames, lengthAdded = p.framer.AppendControlFrames(frames, maxFrameSize-length)
	length += lengthAdded

	if !canSendStreamFrames {
		return frames, nil
	}

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

func (p *packetPacker) getHeader(encLevel protocol.EncryptionLevel) *wire.Header {
	pn, pnLen := p.pnManager.PeekPacketNumber()
	header := &wire.Header{
		PacketNumber:     pn,
		PacketNumberLen:  pnLen,
		Version:          p.version,
		DestConnectionID: p.destConnID,
	}

	if encLevel != protocol.Encryption1RTT {
		header.IsLongHeader = true
		header.SrcConnectionID = p.srcConnID
		// Set the payload len to maximum size.
		// Since it is encoded as a varint, this guarantees us that the header will end up at most as big as GetLength() returns.
		header.PayloadLen = p.maxPacketSize
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
	header *wire.Header,
	frames []wire.Frame,
	sealer handshake.Sealer,
) ([]byte, error) {
	raw := *getPacketBuffer()
	buffer := bytes.NewBuffer(raw[:0])

	addPadding := p.perspective == protocol.PerspectiveClient && header.Type == protocol.PacketTypeInitial && !p.hasSentPacket

	// the payload length is only needed for Long Headers
	if header.IsLongHeader {
		if p.perspective == protocol.PerspectiveClient && header.Type == protocol.PacketTypeInitial {
			header.Token = p.token
		}
		if addPadding {
			headerLen := header.GetLength(p.version)
			header.PayloadLen = protocol.ByteCount(protocol.MinInitialPacketSize) - headerLen
		} else {
			payloadLen := protocol.ByteCount(sealer.Overhead())
			for _, frame := range frames {
				payloadLen += frame.Length(p.version)
			}
			header.PayloadLen = payloadLen
		}
	}

	if err := header.Write(buffer, p.perspective, p.version); err != nil {
		return nil, err
	}
	payloadStartIndex := buffer.Len()

	// the Initial packet needs to be padded, so the last STREAM frame must have the data length present
	if p.perspective == protocol.PerspectiveClient && header.Type == protocol.PacketTypeInitial {
		lastFrame := frames[len(frames)-1]
		if sf, ok := lastFrame.(*wire.StreamFrame); ok {
			sf.DataLenPresent = true
		}
	}
	for _, frame := range frames {
		if err := frame.Write(buffer, p.version); err != nil {
			return nil, err
		}
	}
	if addPadding {
		paddingLen := protocol.MinInitialPacketSize - sealer.Overhead() - buffer.Len()
		if paddingLen > 0 {
			buffer.Write(bytes.Repeat([]byte{0}, paddingLen))
		}
	}

	if size := protocol.ByteCount(buffer.Len() + sealer.Overhead()); size > p.maxPacketSize {
		return nil, fmt.Errorf("PacketPacker BUG: packet too large (%d bytes, allowed %d bytes)", size, p.maxPacketSize)
	}

	raw = raw[0:buffer.Len()]
	_ = sealer.Seal(raw[payloadStartIndex:payloadStartIndex], raw[payloadStartIndex:], header.PacketNumber, raw[:payloadStartIndex])
	raw = raw[0 : buffer.Len()+sealer.Overhead()]

	num := p.pnManager.PopPacketNumber()
	if num != header.PacketNumber {
		return nil, errors.New("packetPacker BUG: Peeked and Popped packet numbers do not match")
	}
	p.hasSentPacket = true
	return raw, nil
}

func (p *packetPacker) canSendData(encLevel protocol.EncryptionLevel) bool {
	return encLevel == protocol.Encryption1RTT
}

func (p *packetPacker) ChangeDestConnectionID(connID protocol.ConnectionID) {
	p.destConnID = connID
}

func (p *packetPacker) HandleTransportParameters(params *handshake.TransportParameters) {
	if params.MaxPacketSize != 0 {
		p.maxPacketSize = utils.MinByteCount(p.maxPacketSize, params.MaxPacketSize)
	}
}
