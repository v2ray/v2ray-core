package ackhandler

import (
	"crypto/rand"
	"math"

	"github.com/lucas-clemente/quic-go/internal/protocol"
	"github.com/lucas-clemente/quic-go/internal/wire"
)

// The packetNumberGenerator generates the packet number for the next packet
// it randomly skips a packet number every averagePeriod packets (on average)
// it is guarantued to never skip two consecutive packet numbers
type packetNumberGenerator struct {
	averagePeriod protocol.PacketNumber

	next       protocol.PacketNumber
	nextToSkip protocol.PacketNumber

	history []protocol.PacketNumber
}

func newPacketNumberGenerator(initial, averagePeriod protocol.PacketNumber) *packetNumberGenerator {
	g := &packetNumberGenerator{
		next:          initial,
		averagePeriod: averagePeriod,
	}
	g.generateNewSkip()
	return g
}

func (p *packetNumberGenerator) Peek() protocol.PacketNumber {
	return p.next
}

func (p *packetNumberGenerator) Pop() protocol.PacketNumber {
	next := p.next

	// generate a new packet number for the next packet
	p.next++

	if p.next == p.nextToSkip {
		if len(p.history)+1 > protocol.MaxTrackedSkippedPackets {
			p.history = p.history[1:]
		}
		p.history = append(p.history, p.next)
		p.next++
		p.generateNewSkip()
	}

	return next
}

func (p *packetNumberGenerator) generateNewSkip() {
	num := p.getRandomNumber()
	skip := protocol.PacketNumber(num) * (p.averagePeriod - 1) / (math.MaxUint16 / 2)
	// make sure that there are never two consecutive packet numbers that are skipped
	p.nextToSkip = p.next + 2 + skip
}

// getRandomNumber() generates a cryptographically secure random number between 0 and MaxUint16 (= 65535)
// The expectation value is 65535/2
func (p *packetNumberGenerator) getRandomNumber() uint16 {
	b := make([]byte, 2)
	rand.Read(b) // ignore the error here

	num := uint16(b[0])<<8 + uint16(b[1])
	return num
}

func (p *packetNumberGenerator) Validate(ack *wire.AckFrame) bool {
	for _, pn := range p.history {
		if ack.AcksPacket(pn) {
			return false
		}
	}
	return true
}
