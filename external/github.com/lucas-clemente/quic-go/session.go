package quic

import (
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net"
	"reflect"
	"sync"
	"time"

	"v2ray.com/core/external/github.com/lucas-clemente/quic-go/internal/ackhandler"
	"v2ray.com/core/external/github.com/lucas-clemente/quic-go/internal/congestion"
	"v2ray.com/core/external/github.com/lucas-clemente/quic-go/internal/flowcontrol"
	"v2ray.com/core/external/github.com/lucas-clemente/quic-go/internal/handshake"
	"v2ray.com/core/external/github.com/lucas-clemente/quic-go/internal/protocol"
	"v2ray.com/core/external/github.com/lucas-clemente/quic-go/internal/qerr"
	"v2ray.com/core/external/github.com/lucas-clemente/quic-go/internal/utils"
	"v2ray.com/core/external/github.com/lucas-clemente/quic-go/internal/wire"
)

type unpacker interface {
	Unpack(hdr *wire.Header, data []byte) (*unpackedPacket, error)
}

type streamGetter interface {
	GetOrOpenReceiveStream(protocol.StreamID) (receiveStreamI, error)
	GetOrOpenSendStream(protocol.StreamID) (sendStreamI, error)
}

type streamManager interface {
	GetOrOpenSendStream(protocol.StreamID) (sendStreamI, error)
	GetOrOpenReceiveStream(protocol.StreamID) (receiveStreamI, error)
	OpenStream() (Stream, error)
	OpenUniStream() (SendStream, error)
	OpenStreamSync() (Stream, error)
	OpenUniStreamSync() (SendStream, error)
	AcceptStream() (Stream, error)
	AcceptUniStream() (ReceiveStream, error)
	DeleteStream(protocol.StreamID) error
	UpdateLimits(*handshake.TransportParameters)
	HandleMaxStreamsFrame(*wire.MaxStreamsFrame) error
	CloseWithError(error)
}

type cryptoStreamHandler interface {
	RunHandshake() error
	io.Closer
	ConnectionState() handshake.ConnectionState
}

type receivedPacket struct {
	remoteAddr net.Addr
	hdr        *wire.Header
	rcvTime    time.Time
	data       []byte

	buffer *packetBuffer
}

type closeError struct {
	err       error
	remote    bool
	sendClose bool
}

var errCloseForRecreating = errors.New("closing session in order to recreate it")

// A Session is a QUIC session
type session struct {
	sessionRunner sessionRunner

	destConnID protocol.ConnectionID
	srcConnID  protocol.ConnectionID

	perspective protocol.Perspective
	version     protocol.VersionNumber
	config      *Config

	conn connection

	streamsMap streamManager

	rttStats *congestion.RTTStats

	cryptoStreamManager   *cryptoStreamManager
	sentPacketHandler     ackhandler.SentPacketHandler
	receivedPacketHandler ackhandler.ReceivedPacketHandler
	framer                framer
	windowUpdateQueue     *windowUpdateQueue
	connFlowController    flowcontrol.ConnectionFlowController

	unpacker unpacker
	packer   packer

	cryptoStreamHandler cryptoStreamHandler

	receivedPackets  chan *receivedPacket
	sendingScheduled chan struct{}

	closeOnce sync.Once
	closed    utils.AtomicBool
	// closeChan is used to notify the run loop that it should terminate
	closeChan                 chan closeError
	connectionClosePacket     *packedPacket
	packetsReceivedAfterClose int

	ctx       context.Context
	ctxCancel context.CancelFunc

	undecryptablePackets []*receivedPacket

	clientHelloWritten    <-chan struct{}
	handshakeCompleteChan chan struct{} // is closed when the handshake completes
	handshakeComplete     bool

	receivedFirstPacket              bool
	receivedFirstForwardSecurePacket bool

	sessionCreationTime     time.Time
	lastNetworkActivityTime time.Time
	// pacingDeadline is the time when the next packet should be sent
	pacingDeadline time.Time

	peerParams *handshake.TransportParameters

	timer *utils.Timer
	// keepAlivePingSent stores whether a Ping frame was sent to the peer or not
	// it is reset as soon as we receive a packet from the peer
	keepAlivePingSent bool

	logger utils.Logger
}

var _ Session = &session{}
var _ streamSender = &session{}

var newSession = func(
	conn connection,
	runner sessionRunner,
	clientDestConnID protocol.ConnectionID,
	destConnID protocol.ConnectionID,
	srcConnID protocol.ConnectionID,
	conf *Config,
	tlsConf *tls.Config,
	params *handshake.TransportParameters,
	logger utils.Logger,
	v protocol.VersionNumber,
) (quicSession, error) {
	s := &session{
		conn:                  conn,
		sessionRunner:         runner,
		config:                conf,
		srcConnID:             srcConnID,
		destConnID:            destConnID,
		perspective:           protocol.PerspectiveServer,
		handshakeCompleteChan: make(chan struct{}),
		logger:                logger,
		version:               v,
	}
	s.preSetup()
	s.sentPacketHandler = ackhandler.NewSentPacketHandler(0, s.rttStats, s.logger)
	initialStream := newCryptoStream()
	handshakeStream := newCryptoStream()
	s.streamsMap = newStreamsMap(
		s,
		s.newFlowController,
		uint64(s.config.MaxIncomingStreams),
		uint64(s.config.MaxIncomingUniStreams),
		s.perspective,
		s.version,
	)
	s.framer = newFramer(s.streamsMap, s.version)
	cs, err := handshake.NewCryptoSetupServer(
		initialStream,
		handshakeStream,
		clientDestConnID,
		params,
		s.processTransportParameters,
		tlsConf,
		conf.Versions,
		v,
		logger,
		protocol.PerspectiveServer,
	)
	if err != nil {
		return nil, err
	}
	s.cryptoStreamHandler = cs
	s.packer = newPacketPacker(
		s.destConnID,
		s.srcConnID,
		initialStream,
		handshakeStream,
		s.sentPacketHandler,
		s.RemoteAddr(),
		nil, // no token
		cs,
		s.framer,
		s.receivedPacketHandler,
		s.perspective,
		s.version,
	)
	s.cryptoStreamManager = newCryptoStreamManager(cs, initialStream, handshakeStream)

	if err := s.postSetup(); err != nil {
		return nil, err
	}
	s.unpacker = newPacketUnpacker(cs, s.version)
	return s, nil
}

// declare this as a variable, such that we can it mock it in the tests
var newClientSession = func(
	conn connection,
	runner sessionRunner,
	token []byte,
	origDestConnID protocol.ConnectionID,
	destConnID protocol.ConnectionID,
	srcConnID protocol.ConnectionID,
	conf *Config,
	tlsConf *tls.Config,
	initialPacketNumber protocol.PacketNumber,
	params *handshake.TransportParameters,
	initialVersion protocol.VersionNumber,
	logger utils.Logger,
	v protocol.VersionNumber,
) (quicSession, error) {
	s := &session{
		conn:                  conn,
		sessionRunner:         runner,
		config:                conf,
		srcConnID:             srcConnID,
		destConnID:            destConnID,
		perspective:           protocol.PerspectiveClient,
		handshakeCompleteChan: make(chan struct{}),
		logger:                logger,
		version:               v,
	}
	s.preSetup()
	s.sentPacketHandler = ackhandler.NewSentPacketHandler(initialPacketNumber, s.rttStats, s.logger)
	initialStream := newCryptoStream()
	handshakeStream := newCryptoStream()
	cs, clientHelloWritten, err := handshake.NewCryptoSetupClient(
		initialStream,
		handshakeStream,
		origDestConnID,
		s.destConnID,
		params,
		s.processTransportParameters,
		tlsConf,
		initialVersion,
		conf.Versions,
		v,
		logger,
		protocol.PerspectiveClient,
	)
	if err != nil {
		return nil, err
	}
	s.clientHelloWritten = clientHelloWritten
	s.cryptoStreamHandler = cs
	s.cryptoStreamManager = newCryptoStreamManager(cs, initialStream, handshakeStream)
	s.unpacker = newPacketUnpacker(cs, s.version)
	s.streamsMap = newStreamsMap(
		s,
		s.newFlowController,
		uint64(s.config.MaxIncomingStreams),
		uint64(s.config.MaxIncomingUniStreams),
		s.perspective,
		s.version,
	)
	s.framer = newFramer(s.streamsMap, s.version)
	s.packer = newPacketPacker(
		s.destConnID,
		s.srcConnID,
		initialStream,
		handshakeStream,
		s.sentPacketHandler,
		s.RemoteAddr(),
		token,
		cs,
		s.framer,
		s.receivedPacketHandler,
		s.perspective,
		s.version,
	)
	return s, s.postSetup()
}

func (s *session) preSetup() {
	s.rttStats = &congestion.RTTStats{}
	s.receivedPacketHandler = ackhandler.NewReceivedPacketHandler(s.rttStats, s.logger, s.version)
	s.connFlowController = flowcontrol.NewConnectionFlowController(
		protocol.InitialMaxData,
		protocol.ByteCount(s.config.MaxReceiveConnectionFlowControlWindow),
		s.onHasConnectionWindowUpdate,
		s.rttStats,
		s.logger,
	)
}

func (s *session) postSetup() error {
	s.receivedPackets = make(chan *receivedPacket, protocol.MaxSessionUnprocessedPackets)
	s.closeChan = make(chan closeError, 1)
	s.sendingScheduled = make(chan struct{}, 1)
	s.undecryptablePackets = make([]*receivedPacket, 0, protocol.MaxUndecryptablePackets)
	s.ctx, s.ctxCancel = context.WithCancel(context.Background())

	s.timer = utils.NewTimer()
	now := time.Now()
	s.lastNetworkActivityTime = now
	s.sessionCreationTime = now

	s.windowUpdateQueue = newWindowUpdateQueue(s.streamsMap, s.connFlowController, s.framer.QueueControlFrame)
	return nil
}

// run the session main loop
func (s *session) run() error {
	defer s.ctxCancel()

	go func() {
		if err := s.cryptoStreamHandler.RunHandshake(); err != nil {
			s.closeLocal(err)
			return
		}
		close(s.handshakeCompleteChan)
	}()
	if s.perspective == protocol.PerspectiveClient {
		select {
		case <-s.clientHelloWritten:
			s.scheduleSending()
		case closeErr := <-s.closeChan:
			// put the close error back into the channel, so that the run loop can receive it
			s.closeChan <- closeErr
		}
	}

	var closeErr closeError

runLoop:
	for {
		// Close immediately if requested
		select {
		case closeErr = <-s.closeChan:
			break runLoop
		case <-s.handshakeCompleteChan:
			s.handleHandshakeComplete()
		default:
		}

		s.maybeResetTimer()

		select {
		case closeErr = <-s.closeChan:
			break runLoop
		case <-s.timer.Chan():
			s.timer.SetRead()
			// We do all the interesting stuff after the switch statement, so
			// nothing to see here.
		case <-s.sendingScheduled:
			// We do all the interesting stuff after the switch statement, so
			// nothing to see here.
		case p := <-s.receivedPackets:
			// Only reset the timers if this packet was actually processed.
			// This avoids modifying any state when handling undecryptable packets,
			// which could be injected by an attacker.
			if wasProcessed := s.handlePacketImpl(p); !wasProcessed {
				continue
			}
		case <-s.handshakeCompleteChan:
			s.handleHandshakeComplete()
		}

		now := time.Now()
		if timeout := s.sentPacketHandler.GetAlarmTimeout(); !timeout.IsZero() && timeout.Before(now) {
			// This could cause packets to be retransmitted.
			// Check it before trying to send packets.
			if err := s.sentPacketHandler.OnAlarm(); err != nil {
				s.closeLocal(err)
			}
		}

		var pacingDeadline time.Time
		if s.pacingDeadline.IsZero() { // the timer didn't have a pacing deadline set
			pacingDeadline = s.sentPacketHandler.TimeUntilSend()
		}
		if s.config.KeepAlive && !s.keepAlivePingSent && s.handshakeComplete && time.Since(s.lastNetworkActivityTime) >= s.peerParams.IdleTimeout/2 {
			// send a PING frame since there is no activity in the session
			s.logger.Debugf("Sending a keep-alive ping to keep the connection alive.")
			s.framer.QueueControlFrame(&wire.PingFrame{})
			s.keepAlivePingSent = true
		} else if !pacingDeadline.IsZero() && now.Before(pacingDeadline) {
			// If we get to this point before the pacing deadline, we should wait until that deadline.
			// This can happen when scheduleSending is called, or a packet is received.
			// Set the timer and restart the run loop.
			s.pacingDeadline = pacingDeadline
			continue
		}

		if !s.handshakeComplete && now.Sub(s.sessionCreationTime) >= s.config.HandshakeTimeout {
			s.closeLocal(qerr.Error(qerr.HandshakeTimeout, "Crypto handshake did not complete in time."))
			continue
		}
		if s.handshakeComplete && now.Sub(s.lastNetworkActivityTime) >= s.config.IdleTimeout {
			s.closeLocal(qerr.Error(qerr.NetworkIdleTimeout, "No recent network activity."))
			continue
		}

		if err := s.sendPackets(); err != nil {
			s.closeLocal(err)
		}
	}

	if err := s.handleCloseError(closeErr); err != nil {
		s.logger.Infof("Handling close error failed: %s", err)
	}
	s.closed.Set(true)
	s.logger.Infof("Connection %s closed.", s.srcConnID)
	s.cryptoStreamHandler.Close()
	return closeErr.err
}

func (s *session) Context() context.Context {
	return s.ctx
}

func (s *session) ConnectionState() ConnectionState {
	return s.cryptoStreamHandler.ConnectionState()
}

func (s *session) maybeResetTimer() {
	var deadline time.Time
	if s.config.KeepAlive && s.handshakeComplete && !s.keepAlivePingSent {
		deadline = s.lastNetworkActivityTime.Add(s.peerParams.IdleTimeout / 2)
	} else {
		deadline = s.lastNetworkActivityTime.Add(s.config.IdleTimeout)
	}

	if ackAlarm := s.receivedPacketHandler.GetAlarmTimeout(); !ackAlarm.IsZero() {
		deadline = utils.MinTime(deadline, ackAlarm)
	}
	if lossTime := s.sentPacketHandler.GetAlarmTimeout(); !lossTime.IsZero() {
		deadline = utils.MinTime(deadline, lossTime)
	}
	if !s.handshakeComplete {
		handshakeDeadline := s.sessionCreationTime.Add(s.config.HandshakeTimeout)
		deadline = utils.MinTime(deadline, handshakeDeadline)
	}
	if !s.pacingDeadline.IsZero() {
		deadline = utils.MinTime(deadline, s.pacingDeadline)
	}

	s.timer.Reset(deadline)
}

func (s *session) handleHandshakeComplete() {
	s.handshakeComplete = true
	s.handshakeCompleteChan = nil // prevent this case from ever being selected again
	s.sessionRunner.onHandshakeComplete(s)

	// The client completes the handshake first (after sending the CFIN).
	// We need to make sure they learn about the peer completing the handshake,
	// in order to stop retransmitting handshake packets.
	// They will stop retransmitting handshake packets when receiving the first forward-secure packet.
	// We need to make sure that a retransmittable forward-secure packet is sent,
	// independent from the application protocol.
	if s.perspective == protocol.PerspectiveServer {
		s.queueControlFrame(&wire.PingFrame{})
		s.sentPacketHandler.SetHandshakeComplete()
	}
}

func (s *session) handlePacketImpl(p *receivedPacket) bool /* was the packet successfully processed */ {
	var wasQueued bool

	defer func() {
		// Put back the packet buffer if the packet wasn't queued for later decryption.
		if !wasQueued {
			p.buffer.Release()
		}
	}()

	// The server can change the source connection ID with the first Handshake packet.
	// After this, all packets with a different source connection have to be ignored.
	if s.receivedFirstPacket && p.hdr.IsLongHeader && !p.hdr.SrcConnectionID.Equal(s.destConnID) {
		s.logger.Debugf("Dropping packet with unexpected source connection ID: %s (expected %s)", p.hdr.SrcConnectionID, s.destConnID)
		return false
	}
	// drop 0-RTT packets
	if p.hdr.Type == protocol.PacketType0RTT {
		return false
	}

	packet, err := s.unpacker.Unpack(p.hdr, p.data)
	if err != nil {
		if err == handshake.ErrOpenerNotYetAvailable {
			// Sealer for this encryption level not yet available.
			// Try again later.
			wasQueued = true
			s.tryQueueingUndecryptablePacket(p)
			return false
		}
		// This might be a packet injected by an attacker.
		// Drop it.
		s.logger.Debugf("Dropping packet that could not be unpacked. Unpack error: %s", err)
		return false
	}

	if s.logger.Debug() {
		s.logger.Debugf("<- Reading packet %#x (%d bytes) for connection %s, %s", packet.packetNumber, len(p.data), p.hdr.DestConnectionID, packet.encryptionLevel)
		packet.hdr.Log(s.logger)
	}

	if err := s.handleUnpackedPacket(packet, p.rcvTime); err != nil {
		s.closeLocal(err)
		return false
	}
	return true
}

func (s *session) handleUnpackedPacket(packet *unpackedPacket, rcvTime time.Time) error {
	if len(packet.data) == 0 {
		return qerr.MissingPayload
	}

	// The server can change the source connection ID with the first Handshake packet.
	if s.perspective == protocol.PerspectiveClient && !s.receivedFirstPacket && packet.hdr.IsLongHeader && !packet.hdr.SrcConnectionID.Equal(s.destConnID) {
		s.logger.Debugf("Received first packet. Switching destination connection ID to: %s", packet.hdr.SrcConnectionID)
		s.destConnID = packet.hdr.SrcConnectionID
		s.packer.ChangeDestConnectionID(s.destConnID)
	}

	s.receivedFirstPacket = true
	s.lastNetworkActivityTime = rcvTime
	s.keepAlivePingSent = false

	// The client completes the handshake first (after sending the CFIN).
	// We know that the server completed the handshake as soon as we receive a forward-secure packet.
	if s.perspective == protocol.PerspectiveClient {
		if !s.receivedFirstForwardSecurePacket && packet.encryptionLevel == protocol.Encryption1RTT {
			s.receivedFirstForwardSecurePacket = true
			s.sentPacketHandler.SetHandshakeComplete()
		}
	}

	r := bytes.NewReader(packet.data)
	var isRetransmittable bool
	for {
		frame, err := wire.ParseNextFrame(r, s.version)
		if err != nil {
			return err
		}
		if frame == nil {
			break
		}
		if ackhandler.IsFrameRetransmittable(frame) {
			isRetransmittable = true
		}
		if err := s.handleFrame(frame, packet.packetNumber, packet.encryptionLevel); err != nil {
			return err
		}
	}

	if err := s.receivedPacketHandler.ReceivedPacket(packet.packetNumber, packet.encryptionLevel, rcvTime, isRetransmittable); err != nil {
		return err
	}
	return nil
}

func (s *session) handleFrame(f wire.Frame, pn protocol.PacketNumber, encLevel protocol.EncryptionLevel) error {
	var err error
	wire.LogFrame(s.logger, f, false)
	switch frame := f.(type) {
	case *wire.CryptoFrame:
		err = s.handleCryptoFrame(frame, encLevel)
	case *wire.StreamFrame:
		err = s.handleStreamFrame(frame, encLevel)
	case *wire.AckFrame:
		err = s.handleAckFrame(frame, pn, encLevel)
	case *wire.ConnectionCloseFrame:
		s.closeRemote(qerr.Error(frame.ErrorCode, frame.ReasonPhrase))
	case *wire.ResetStreamFrame:
		err = s.handleResetStreamFrame(frame)
	case *wire.MaxDataFrame:
		s.handleMaxDataFrame(frame)
	case *wire.MaxStreamDataFrame:
		err = s.handleMaxStreamDataFrame(frame)
	case *wire.MaxStreamsFrame:
		err = s.handleMaxStreamsFrame(frame)
	case *wire.DataBlockedFrame:
	case *wire.StreamDataBlockedFrame:
	case *wire.StreamsBlockedFrame:
	case *wire.StopSendingFrame:
		err = s.handleStopSendingFrame(frame)
	case *wire.PingFrame:
	case *wire.PathChallengeFrame:
		s.handlePathChallengeFrame(frame)
	case *wire.PathResponseFrame:
		// since we don't send PATH_CHALLENGEs, we don't expect PATH_RESPONSEs
		err = errors.New("unexpected PATH_RESPONSE frame")
	case *wire.NewTokenFrame:
	case *wire.NewConnectionIDFrame:
	case *wire.RetireConnectionIDFrame:
		// since we don't send new connection IDs, we don't expect retirements
		err = errors.New("unexpected RETIRE_CONNECTION_ID frame")
	default:
		err = fmt.Errorf("unexpected frame type: %s", reflect.ValueOf(&frame).Elem().Type().Name())
	}
	return err
}

// handlePacket is called by the server with a new packet
func (s *session) handlePacket(p *receivedPacket) {
	if s.closed.Get() {
		s.handlePacketAfterClosed(p)
	}
	// Discard packets once the amount of queued packets is larger than
	// the channel size, protocol.MaxSessionUnprocessedPackets
	select {
	case s.receivedPackets <- p:
	default:
	}
}

func (s *session) handlePacketAfterClosed(p *receivedPacket) {
	s.packetsReceivedAfterClose++
	if s.connectionClosePacket == nil {
		return
	}
	// exponential backoff
	// only send a CONNECTION_CLOSE for the 1st, 2nd, 4th, 8th, 16th, ... packet arriving
	for n := s.packetsReceivedAfterClose; n > 1; n = n / 2 {
		if n%2 != 0 {
			return
		}
	}
	s.logger.Debugf("Received %d packets after sending CONNECTION_CLOSE. Retransmitting.", s.packetsReceivedAfterClose)
	if err := s.conn.Write(s.connectionClosePacket.raw); err != nil {
		s.logger.Debugf("Error retransmitting CONNECTION_CLOSE: %s", err)
	}
}

func (s *session) handleCryptoFrame(frame *wire.CryptoFrame, encLevel protocol.EncryptionLevel) error {
	encLevelChanged, err := s.cryptoStreamManager.HandleCryptoFrame(frame, encLevel)
	if err != nil {
		return err
	}
	if encLevelChanged {
		s.tryDecryptingQueuedPackets()
	}
	return nil
}

func (s *session) handleStreamFrame(frame *wire.StreamFrame, encLevel protocol.EncryptionLevel) error {
	if encLevel < protocol.Encryption1RTT {
		return qerr.Error(qerr.UnencryptedStreamData, fmt.Sprintf("received unencrypted stream data on stream %d", frame.StreamID))
	}
	str, err := s.streamsMap.GetOrOpenReceiveStream(frame.StreamID)
	if err != nil {
		return err
	}
	if str == nil {
		// Stream is closed and already garbage collected
		// ignore this StreamFrame
		return nil
	}
	return str.handleStreamFrame(frame)
}

func (s *session) handleMaxDataFrame(frame *wire.MaxDataFrame) {
	s.connFlowController.UpdateSendWindow(frame.ByteOffset)
}

func (s *session) handleMaxStreamDataFrame(frame *wire.MaxStreamDataFrame) error {
	str, err := s.streamsMap.GetOrOpenSendStream(frame.StreamID)
	if err != nil {
		return err
	}
	if str == nil {
		// stream is closed and already garbage collected
		return nil
	}
	str.handleMaxStreamDataFrame(frame)
	return nil
}

func (s *session) handleMaxStreamsFrame(frame *wire.MaxStreamsFrame) error {
	return s.streamsMap.HandleMaxStreamsFrame(frame)
}

func (s *session) handleResetStreamFrame(frame *wire.ResetStreamFrame) error {
	str, err := s.streamsMap.GetOrOpenReceiveStream(frame.StreamID)
	if err != nil {
		return err
	}
	if str == nil {
		// stream is closed and already garbage collected
		return nil
	}
	return str.handleResetStreamFrame(frame)
}

func (s *session) handleStopSendingFrame(frame *wire.StopSendingFrame) error {
	str, err := s.streamsMap.GetOrOpenSendStream(frame.StreamID)
	if err != nil {
		return err
	}
	if str == nil {
		// stream is closed and already garbage collected
		return nil
	}
	str.handleStopSendingFrame(frame)
	return nil
}

func (s *session) handlePathChallengeFrame(frame *wire.PathChallengeFrame) {
	s.queueControlFrame(&wire.PathResponseFrame{Data: frame.Data})
}

func (s *session) handleAckFrame(frame *wire.AckFrame, pn protocol.PacketNumber, encLevel protocol.EncryptionLevel) error {
	if err := s.sentPacketHandler.ReceivedAck(frame, pn, encLevel, s.lastNetworkActivityTime); err != nil {
		return err
	}
	if encLevel == protocol.Encryption1RTT {
		s.receivedPacketHandler.IgnoreBelow(s.sentPacketHandler.GetLowestPacketNotConfirmedAcked())
	}
	return nil
}

// closeLocal closes the session and send a CONNECTION_CLOSE containing the error
func (s *session) closeLocal(e error) {
	s.closeOnce.Do(func() {
		s.sessionRunner.retireConnectionID(s.srcConnID)
		s.closeChan <- closeError{err: e, sendClose: true, remote: false}
	})
}

// destroy closes the session without sending the error on the wire
func (s *session) destroy(e error) {
	s.closeOnce.Do(func() {
		s.sessionRunner.removeConnectionID(s.srcConnID)
		s.closeChan <- closeError{err: e, sendClose: false, remote: false}
	})
}

// closeForRecreating closes the session in order to recreate it immediately afterwards
// It returns the first packet number that should be used in the new session.
func (s *session) closeForRecreating() protocol.PacketNumber {
	s.destroy(errCloseForRecreating)
	nextPN, _ := s.sentPacketHandler.PeekPacketNumber()
	return nextPN
}

func (s *session) closeRemote(e error) {
	s.closeOnce.Do(func() {
		s.sessionRunner.removeConnectionID(s.srcConnID)
		s.closeChan <- closeError{err: e, remote: true}
	})
}

// Close the connection. It sends a qerr.PeerGoingAway.
// It waits until the run loop has stopped before returning
func (s *session) Close() error {
	s.closeLocal(nil)
	<-s.ctx.Done()
	return nil
}

func (s *session) CloseWithError(code protocol.ApplicationErrorCode, e error) error {
	s.closeLocal(qerr.Error(qerr.ErrorCode(code), e.Error()))
	<-s.ctx.Done()
	return nil
}

func (s *session) handleCloseError(closeErr closeError) error {
	if closeErr.err == nil {
		closeErr.err = qerr.PeerGoingAway
	}

	var quicErr *qerr.QuicError
	var ok bool
	if quicErr, ok = closeErr.err.(*qerr.QuicError); !ok {
		quicErr = qerr.ToQuicError(closeErr.err)
	}
	// Don't log 'normal' reasons
	if quicErr.ErrorCode == qerr.PeerGoingAway || quicErr.ErrorCode == qerr.NetworkIdleTimeout {
		s.logger.Infof("Closing connection %s.", s.srcConnID)
	} else {
		s.logger.Errorf("Closing session with error: %s", closeErr.err.Error())
	}

	s.streamsMap.CloseWithError(quicErr)

	if !closeErr.sendClose {
		return nil
	}

	// If this is a remote close we're done here
	if closeErr.remote {
		return nil
	}
	// otherwise send a CONNECTION_CLOSE
	return s.sendConnectionClose(quicErr)
}

func (s *session) processTransportParameters(params *handshake.TransportParameters) {
	s.peerParams = params
	s.streamsMap.UpdateLimits(params)
	s.packer.HandleTransportParameters(params)
	s.connFlowController.UpdateSendWindow(params.InitialMaxData)
	// the crypto stream is the only open stream at this moment
	// so we don't need to update stream flow control windows
}

func (s *session) sendPackets() error {
	s.pacingDeadline = time.Time{}

	sendMode := s.sentPacketHandler.SendMode()
	if sendMode == ackhandler.SendNone { // shortcut: return immediately if there's nothing to send
		return nil
	}

	numPackets := s.sentPacketHandler.ShouldSendNumPackets()
	var numPacketsSent int
sendLoop:
	for {
		switch sendMode {
		case ackhandler.SendNone:
			break sendLoop
		case ackhandler.SendAck:
			// If we already sent packets, and the send mode switches to SendAck,
			// we've just become congestion limited.
			// There's no need to try to send an ACK at this moment.
			if numPacketsSent > 0 {
				return nil
			}
			// We can at most send a single ACK only packet.
			// There will only be a new ACK after receiving new packets.
			// SendAck is only returned when we're congestion limited, so we don't need to set the pacingt timer.
			return s.maybeSendAckOnlyPacket()
		case ackhandler.SendPTO:
			if err := s.sendProbePacket(); err != nil {
				return err
			}
			numPacketsSent++
		case ackhandler.SendRetransmission:
			sentPacket, err := s.maybeSendRetransmission()
			if err != nil {
				return err
			}
			if sentPacket {
				numPacketsSent++
				// This can happen if a retransmission queued, but it wasn't necessary to send it.
				// e.g. when an Initial is queued, but we already received a packet from the server.
			}
		case ackhandler.SendAny:
			sentPacket, err := s.sendPacket()
			if err != nil {
				return err
			}
			if !sentPacket {
				break sendLoop
			}
			numPacketsSent++
		default:
			return fmt.Errorf("BUG: invalid send mode %d", sendMode)
		}
		if numPacketsSent >= numPackets {
			break
		}
		sendMode = s.sentPacketHandler.SendMode()
	}
	// Only start the pacing timer if we sent as many packets as we were allowed.
	// There will probably be more to send when calling sendPacket again.
	if numPacketsSent == numPackets {
		s.pacingDeadline = s.sentPacketHandler.TimeUntilSend()
	}
	return nil
}

func (s *session) maybeSendAckOnlyPacket() error {
	packet, err := s.packer.MaybePackAckPacket()
	if err != nil {
		return err
	}
	if packet == nil {
		return nil
	}
	s.sentPacketHandler.SentPacket(packet.ToAckHandlerPacket())
	return s.sendPackedPacket(packet)
}

// maybeSendRetransmission sends retransmissions for at most one packet.
// It takes care that Initials aren't retransmitted, if a packet from the server was already received.
func (s *session) maybeSendRetransmission() (bool, error) {
	var retransmitPacket *ackhandler.Packet
	for {
		retransmitPacket = s.sentPacketHandler.DequeuePacketForRetransmission()
		if retransmitPacket == nil {
			return false, nil
		}

		// Don't retransmit Initial packets if we already received a response.
		// An Initial might have been retransmitted multiple times before we receive a response.
		// As soon as we receive one response, we don't need to send any more Initials.
		if s.perspective == protocol.PerspectiveClient && s.receivedFirstPacket && retransmitPacket.PacketType == protocol.PacketTypeInitial {
			s.logger.Debugf("Skipping retransmission of packet %d. Already received a response to an Initial.", retransmitPacket.PacketNumber)
			continue
		}
		break
	}

	s.logger.Debugf("Dequeueing retransmission for packet 0x%x", retransmitPacket.PacketNumber)
	packets, err := s.packer.PackRetransmission(retransmitPacket)
	if err != nil {
		return false, err
	}
	ackhandlerPackets := make([]*ackhandler.Packet, len(packets))
	for i, packet := range packets {
		ackhandlerPackets[i] = packet.ToAckHandlerPacket()
	}
	s.sentPacketHandler.SentPacketsAsRetransmission(ackhandlerPackets, retransmitPacket.PacketNumber)
	for _, packet := range packets {
		if err := s.sendPackedPacket(packet); err != nil {
			return false, err
		}
	}
	return true, nil
}

func (s *session) sendProbePacket() error {
	p, err := s.sentPacketHandler.DequeueProbePacket()
	if err != nil {
		return err
	}
	s.logger.Debugf("Sending a retransmission for %#x as a probe packet.", p.PacketNumber)

	packets, err := s.packer.PackRetransmission(p)
	if err != nil {
		return err
	}
	ackhandlerPackets := make([]*ackhandler.Packet, len(packets))
	for i, packet := range packets {
		ackhandlerPackets[i] = packet.ToAckHandlerPacket()
	}
	s.sentPacketHandler.SentPacketsAsRetransmission(ackhandlerPackets, p.PacketNumber)
	for _, packet := range packets {
		if err := s.sendPackedPacket(packet); err != nil {
			return err
		}
	}
	return nil
}

func (s *session) sendPacket() (bool, error) {
	if isBlocked, offset := s.connFlowController.IsNewlyBlocked(); isBlocked {
		s.framer.QueueControlFrame(&wire.DataBlockedFrame{DataLimit: offset})
	}
	s.windowUpdateQueue.QueueAll()

	packet, err := s.packer.PackPacket()
	if err != nil || packet == nil {
		return false, err
	}
	s.sentPacketHandler.SentPacket(packet.ToAckHandlerPacket())
	if err := s.sendPackedPacket(packet); err != nil {
		return false, err
	}
	return true, nil
}

func (s *session) sendPackedPacket(packet *packedPacket) error {
	defer packet.buffer.Release()
	s.logPacket(packet)
	return s.conn.Write(packet.raw)
}

func (s *session) sendConnectionClose(quicErr *qerr.QuicError) error {
	packet, err := s.packer.PackConnectionClose(&wire.ConnectionCloseFrame{
		ErrorCode:    quicErr.ErrorCode,
		ReasonPhrase: quicErr.ErrorMessage,
	})
	if err != nil {
		return err
	}
	s.connectionClosePacket = packet
	s.logPacket(packet)
	return s.conn.Write(packet.raw)
}

func (s *session) logPacket(packet *packedPacket) {
	if !s.logger.Debug() {
		// We don't need to allocate the slices for calling the format functions
		return
	}
	s.logger.Debugf("-> Sending packet 0x%x (%d bytes) for connection %s, %s", packet.header.PacketNumber, len(packet.raw), s.srcConnID, packet.EncryptionLevel())
	packet.header.Log(s.logger)
	for _, frame := range packet.frames {
		wire.LogFrame(s.logger, frame, true)
	}
}

// GetOrOpenStream either returns an existing stream, a newly opened stream, or nil if a stream with the provided ID is already closed.
// It is *only* needed for gQUIC's H2.
// It will be removed as soon as gQUIC moves towards the IETF H2/QUIC stream mapping.
func (s *session) GetOrOpenStream(id protocol.StreamID) (Stream, error) {
	str, err := s.streamsMap.GetOrOpenSendStream(id)
	if str != nil {
		if bstr, ok := str.(Stream); ok {
			return bstr, err
		}
		return nil, fmt.Errorf("Stream %d is not a bidirectional stream", id)
	}
	// make sure to return an actual nil value here, not an Stream with value nil
	return nil, err
}

// AcceptStream returns the next stream openend by the peer
func (s *session) AcceptStream() (Stream, error) {
	return s.streamsMap.AcceptStream()
}

func (s *session) AcceptUniStream() (ReceiveStream, error) {
	return s.streamsMap.AcceptUniStream()
}

// OpenStream opens a stream
func (s *session) OpenStream() (Stream, error) {
	return s.streamsMap.OpenStream()
}

func (s *session) OpenStreamSync() (Stream, error) {
	return s.streamsMap.OpenStreamSync()
}

func (s *session) OpenUniStream() (SendStream, error) {
	return s.streamsMap.OpenUniStream()
}

func (s *session) OpenUniStreamSync() (SendStream, error) {
	return s.streamsMap.OpenUniStreamSync()
}

func (s *session) newStream(id protocol.StreamID) streamI {
	flowController := s.newFlowController(id)
	return newStream(id, s, flowController, s.version)
}

func (s *session) newFlowController(id protocol.StreamID) flowcontrol.StreamFlowController {
	var initialSendWindow protocol.ByteCount
	if s.peerParams != nil {
		if id.Type() == protocol.StreamTypeUni {
			initialSendWindow = s.peerParams.InitialMaxStreamDataUni
		} else {
			if id.InitiatedBy() == s.perspective {
				initialSendWindow = s.peerParams.InitialMaxStreamDataBidiLocal
			} else {
				initialSendWindow = s.peerParams.InitialMaxStreamDataBidiRemote
			}
		}
	}
	return flowcontrol.NewStreamFlowController(
		id,
		s.connFlowController,
		protocol.InitialMaxStreamData,
		protocol.ByteCount(s.config.MaxReceiveStreamFlowControlWindow),
		initialSendWindow,
		s.onHasStreamWindowUpdate,
		s.rttStats,
		s.logger,
	)
}

// scheduleSending signals that we have data for sending
func (s *session) scheduleSending() {
	select {
	case s.sendingScheduled <- struct{}{}:
	default:
	}
}

func (s *session) tryQueueingUndecryptablePacket(p *receivedPacket) {
	if s.handshakeComplete {
		s.logger.Debugf("Received undecryptable packet from %s after the handshake (%d bytes)", p.remoteAddr.String(), len(p.data))
		return
	}
	if len(s.undecryptablePackets)+1 > protocol.MaxUndecryptablePackets {
		s.logger.Infof("Dropping undecrytable packet (%d bytes). Undecryptable packet queue full.", len(p.data))
		return
	}
	s.logger.Infof("Queueing packet (%d bytes) for later decryption", len(p.data))
	s.undecryptablePackets = append(s.undecryptablePackets, p)
}

func (s *session) tryDecryptingQueuedPackets() {
	for _, p := range s.undecryptablePackets {
		s.handlePacket(p)
	}
	s.undecryptablePackets = s.undecryptablePackets[:0]
}

func (s *session) queueControlFrame(f wire.Frame) {
	s.framer.QueueControlFrame(f)
	s.scheduleSending()
}

func (s *session) onHasStreamWindowUpdate(id protocol.StreamID) {
	s.windowUpdateQueue.AddStream(id)
	s.scheduleSending()
}

func (s *session) onHasConnectionWindowUpdate() {
	s.windowUpdateQueue.AddConnection()
	s.scheduleSending()
}

func (s *session) onHasStreamData(id protocol.StreamID) {
	s.framer.AddActiveStream(id)
	s.scheduleSending()
}

func (s *session) onStreamCompleted(id protocol.StreamID) {
	if err := s.streamsMap.DeleteStream(id); err != nil {
		s.closeLocal(err)
	}
}

func (s *session) LocalAddr() net.Addr {
	return s.conn.LocalAddr()
}

func (s *session) RemoteAddr() net.Addr {
	return s.conn.RemoteAddr()
}

func (s *session) GetVersion() protocol.VersionNumber {
	return s.version
}
