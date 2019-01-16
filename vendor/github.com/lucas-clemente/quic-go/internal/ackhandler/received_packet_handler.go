package ackhandler

import (
	"fmt"
	"time"

	"github.com/lucas-clemente/quic-go/internal/congestion"
	"github.com/lucas-clemente/quic-go/internal/protocol"
	"github.com/lucas-clemente/quic-go/internal/utils"
	"github.com/lucas-clemente/quic-go/internal/wire"
)

const (
	// maximum delay that can be applied to an ACK for a retransmittable packet
	ackSendDelay = 25 * time.Millisecond
	// initial maximum number of retransmittable packets received before sending an ack.
	initialRetransmittablePacketsBeforeAck = 2
	// number of retransmittable that an ACK is sent for
	retransmittablePacketsBeforeAck = 10
	// 1/5 RTT delay when doing ack decimation
	ackDecimationDelay = 1.0 / 4
	// 1/8 RTT delay when doing ack decimation
	shortAckDecimationDelay = 1.0 / 8
	// Minimum number of packets received before ack decimation is enabled.
	// This intends to avoid the beginning of slow start, when CWNDs may be
	// rapidly increasing.
	minReceivedBeforeAckDecimation = 100
	// Maximum number of packets to ack immediately after a missing packet for
	// fast retransmission to kick in at the sender.  This limit is created to
	// reduce the number of acks sent that have no benefit for fast retransmission.
	// Set to the number of nacks needed for fast retransmit plus one for protection
	// against an ack loss
	maxPacketsAfterNewMissing = 4
)

type receivedPacketHandler struct {
	initialPackets   *receivedPacketTracker
	handshakePackets *receivedPacketTracker
	oneRTTPackets    *receivedPacketTracker
}

var _ ReceivedPacketHandler = &receivedPacketHandler{}

// NewReceivedPacketHandler creates a new receivedPacketHandler
func NewReceivedPacketHandler(
	rttStats *congestion.RTTStats,
	logger utils.Logger,
	version protocol.VersionNumber,
) ReceivedPacketHandler {
	return &receivedPacketHandler{
		initialPackets:   newReceivedPacketTracker(rttStats, logger, version),
		handshakePackets: newReceivedPacketTracker(rttStats, logger, version),
		oneRTTPackets:    newReceivedPacketTracker(rttStats, logger, version),
	}
}

func (h *receivedPacketHandler) ReceivedPacket(
	pn protocol.PacketNumber,
	encLevel protocol.EncryptionLevel,
	rcvTime time.Time,
	shouldInstigateAck bool,
) error {
	switch encLevel {
	case protocol.EncryptionInitial:
		return h.initialPackets.ReceivedPacket(pn, rcvTime, shouldInstigateAck)
	case protocol.EncryptionHandshake:
		return h.handshakePackets.ReceivedPacket(pn, rcvTime, shouldInstigateAck)
	case protocol.Encryption1RTT:
		return h.oneRTTPackets.ReceivedPacket(pn, rcvTime, shouldInstigateAck)
	default:
		return fmt.Errorf("received packet with unknown encryption level: %s", encLevel)
	}
}

// only to be used with 1-RTT packets
func (h *receivedPacketHandler) IgnoreBelow(pn protocol.PacketNumber) {
	h.oneRTTPackets.IgnoreBelow(pn)
}

func (h *receivedPacketHandler) GetAlarmTimeout() time.Time {
	initialAlarm := h.initialPackets.GetAlarmTimeout()
	handshakeAlarm := h.handshakePackets.GetAlarmTimeout()
	oneRTTAlarm := h.oneRTTPackets.GetAlarmTimeout()
	return utils.MinNonZeroTime(utils.MinNonZeroTime(initialAlarm, handshakeAlarm), oneRTTAlarm)
}

func (h *receivedPacketHandler) GetAckFrame(encLevel protocol.EncryptionLevel) *wire.AckFrame {
	switch encLevel {
	case protocol.EncryptionInitial:
		return h.initialPackets.GetAckFrame()
	case protocol.EncryptionHandshake:
		return h.handshakePackets.GetAckFrame()
	case protocol.Encryption1RTT:
		return h.oneRTTPackets.GetAckFrame()
	default:
		return nil
	}
}
