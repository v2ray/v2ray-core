package quic

import (
	"errors"
	"fmt"
	"net"

	"v2ray.com/core/external/github.com/lucas-clemente/quic-go/internal/flowcontrol"
	"v2ray.com/core/external/github.com/lucas-clemente/quic-go/internal/handshake"
	"v2ray.com/core/external/github.com/lucas-clemente/quic-go/internal/protocol"
	"v2ray.com/core/external/github.com/lucas-clemente/quic-go/internal/wire"
)

type streamOpenErr struct{ error }

var _ net.Error = &streamOpenErr{}

func (e streamOpenErr) Temporary() bool { return e.error == errTooManyOpenStreams }
func (streamOpenErr) Timeout() bool     { return false }

// errTooManyOpenStreams is used internally by the outgoing streams maps.
var errTooManyOpenStreams = errors.New("too many open streams")

type streamsMap struct {
	perspective protocol.Perspective

	sender            streamSender
	newFlowController func(protocol.StreamID) flowcontrol.StreamFlowController

	outgoingBidiStreams *outgoingBidiStreamsMap
	outgoingUniStreams  *outgoingUniStreamsMap
	incomingBidiStreams *incomingBidiStreamsMap
	incomingUniStreams  *incomingUniStreamsMap
}

var _ streamManager = &streamsMap{}

func newStreamsMap(
	sender streamSender,
	newFlowController func(protocol.StreamID) flowcontrol.StreamFlowController,
	maxIncomingStreams uint64,
	maxIncomingUniStreams uint64,
	perspective protocol.Perspective,
	version protocol.VersionNumber,
) streamManager {
	m := &streamsMap{
		perspective:       perspective,
		newFlowController: newFlowController,
		sender:            sender,
	}
	newBidiStream := func(id protocol.StreamID) streamI {
		return newStream(id, m.sender, m.newFlowController(id), version)
	}
	newUniSendStream := func(id protocol.StreamID) sendStreamI {
		return newSendStream(id, m.sender, m.newFlowController(id), version)
	}
	newUniReceiveStream := func(id protocol.StreamID) receiveStreamI {
		return newReceiveStream(id, m.sender, m.newFlowController(id), version)
	}
	m.outgoingBidiStreams = newOutgoingBidiStreamsMap(
		protocol.FirstStream(protocol.StreamTypeBidi, perspective),
		newBidiStream,
		sender.queueControlFrame,
	)
	m.incomingBidiStreams = newIncomingBidiStreamsMap(
		protocol.FirstStream(protocol.StreamTypeBidi, perspective.Opposite()),
		protocol.MaxStreamID(protocol.StreamTypeBidi, maxIncomingStreams, perspective.Opposite()),
		maxIncomingStreams,
		sender.queueControlFrame,
		newBidiStream,
	)
	m.outgoingUniStreams = newOutgoingUniStreamsMap(
		protocol.FirstStream(protocol.StreamTypeUni, perspective),
		newUniSendStream,
		sender.queueControlFrame,
	)
	m.incomingUniStreams = newIncomingUniStreamsMap(
		protocol.FirstStream(protocol.StreamTypeUni, perspective.Opposite()),
		protocol.MaxStreamID(protocol.StreamTypeUni, maxIncomingUniStreams, perspective.Opposite()),
		maxIncomingUniStreams,
		sender.queueControlFrame,
		newUniReceiveStream,
	)
	return m
}

func (m *streamsMap) OpenStream() (Stream, error) {
	return m.outgoingBidiStreams.OpenStream()
}

func (m *streamsMap) OpenStreamSync() (Stream, error) {
	return m.outgoingBidiStreams.OpenStreamSync()
}

func (m *streamsMap) OpenUniStream() (SendStream, error) {
	return m.outgoingUniStreams.OpenStream()
}

func (m *streamsMap) OpenUniStreamSync() (SendStream, error) {
	return m.outgoingUniStreams.OpenStreamSync()
}

func (m *streamsMap) AcceptStream() (Stream, error) {
	return m.incomingBidiStreams.AcceptStream()
}

func (m *streamsMap) AcceptUniStream() (ReceiveStream, error) {
	return m.incomingUniStreams.AcceptStream()
}

func (m *streamsMap) DeleteStream(id protocol.StreamID) error {
	switch id.Type() {
	case protocol.StreamTypeUni:
		if id.InitiatedBy() == m.perspective {
			return m.outgoingUniStreams.DeleteStream(id)
		}
		return m.incomingUniStreams.DeleteStream(id)
	case protocol.StreamTypeBidi:
		if id.InitiatedBy() == m.perspective {
			return m.outgoingBidiStreams.DeleteStream(id)
		}
		return m.incomingBidiStreams.DeleteStream(id)
	}
	panic("")
}

func (m *streamsMap) GetOrOpenReceiveStream(id protocol.StreamID) (receiveStreamI, error) {
	switch id.Type() {
	case protocol.StreamTypeUni:
		if id.InitiatedBy() == m.perspective {
			// an outgoing unidirectional stream is a send stream, not a receive stream
			return nil, fmt.Errorf("peer attempted to open receive stream %d", id)
		}
		return m.incomingUniStreams.GetOrOpenStream(id)
	case protocol.StreamTypeBidi:
		if id.InitiatedBy() == m.perspective {
			return m.outgoingBidiStreams.GetStream(id)
		}
		return m.incomingBidiStreams.GetOrOpenStream(id)
	}
	panic("")
}

func (m *streamsMap) GetOrOpenSendStream(id protocol.StreamID) (sendStreamI, error) {
	switch id.Type() {
	case protocol.StreamTypeUni:
		if id.InitiatedBy() == m.perspective {
			return m.outgoingUniStreams.GetStream(id)
		}
		// an incoming unidirectional stream is a receive stream, not a send stream
		return nil, fmt.Errorf("peer attempted to open send stream %d", id)
	case protocol.StreamTypeBidi:
		if id.InitiatedBy() == m.perspective {
			return m.outgoingBidiStreams.GetStream(id)
		}
		return m.incomingBidiStreams.GetOrOpenStream(id)
	}
	panic("")
}

func (m *streamsMap) HandleMaxStreamsFrame(f *wire.MaxStreamsFrame) error {
	id := protocol.MaxStreamID(f.Type, f.MaxStreams, m.perspective)
	switch id.Type() {
	case protocol.StreamTypeUni:
		m.outgoingUniStreams.SetMaxStream(id)
	case protocol.StreamTypeBidi:
		fmt.Printf("")
		m.outgoingBidiStreams.SetMaxStream(id)
	}
	return nil
}

func (m *streamsMap) UpdateLimits(p *handshake.TransportParameters) {
	// Max{Uni,Bidi}StreamID returns the highest stream ID that the peer is allowed to open.
	m.outgoingBidiStreams.SetMaxStream(protocol.MaxStreamID(protocol.StreamTypeBidi, p.MaxBidiStreams, m.perspective))
	m.outgoingUniStreams.SetMaxStream(protocol.MaxStreamID(protocol.StreamTypeUni, p.MaxUniStreams, m.perspective))
}

func (m *streamsMap) CloseWithError(err error) {
	m.outgoingBidiStreams.CloseWithError(err)
	m.outgoingUniStreams.CloseWithError(err)
	m.incomingBidiStreams.CloseWithError(err)
	m.incomingUniStreams.CloseWithError(err)
}
