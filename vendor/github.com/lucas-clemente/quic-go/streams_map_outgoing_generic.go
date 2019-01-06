package quic

import (
	"fmt"
	"sync"

	"github.com/lucas-clemente/quic-go/internal/protocol"
	"github.com/lucas-clemente/quic-go/internal/qerr"
	"github.com/lucas-clemente/quic-go/internal/wire"
)

//go:generate genny -in $GOFILE -out streams_map_outgoing_bidi.go gen "item=streamI Item=BidiStream streamTypeGeneric=protocol.StreamTypeBidi"
//go:generate genny -in $GOFILE -out streams_map_outgoing_uni.go gen "item=sendStreamI Item=UniStream streamTypeGeneric=protocol.StreamTypeUni"
type outgoingItemsMap struct {
	mutex sync.RWMutex
	cond  sync.Cond

	streams map[protocol.StreamID]item

	nextStream   protocol.StreamID // stream ID of the stream returned by OpenStream(Sync)
	maxStream    protocol.StreamID // the maximum stream ID we're allowed to open
	maxStreamSet bool              // was maxStream set. If not, it's not possible to any stream (also works for stream 0)
	blockedSent  bool              // was a STREAMS_BLOCKED sent for the current maxStream

	newStream            func(protocol.StreamID) item
	queueStreamIDBlocked func(*wire.StreamsBlockedFrame)

	closeErr error
}

func newOutgoingItemsMap(
	nextStream protocol.StreamID,
	newStream func(protocol.StreamID) item,
	queueControlFrame func(wire.Frame),
) *outgoingItemsMap {
	m := &outgoingItemsMap{
		streams:              make(map[protocol.StreamID]item),
		nextStream:           nextStream,
		newStream:            newStream,
		queueStreamIDBlocked: func(f *wire.StreamsBlockedFrame) { queueControlFrame(f) },
	}
	m.cond.L = &m.mutex
	return m
}

func (m *outgoingItemsMap) OpenStream() (item, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	return m.openStreamImpl()
}

func (m *outgoingItemsMap) OpenStreamSync() (item, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for {
		str, err := m.openStreamImpl()
		if err == nil {
			return str, err
		}
		if err != nil && err != qerr.TooManyOpenStreams {
			return nil, err
		}
		m.cond.Wait()
	}
}

func (m *outgoingItemsMap) openStreamImpl() (item, error) {
	if m.closeErr != nil {
		return nil, m.closeErr
	}
	if !m.maxStreamSet || m.nextStream > m.maxStream {
		if !m.blockedSent {
			if m.maxStreamSet {
				m.queueStreamIDBlocked(&wire.StreamsBlockedFrame{
					Type:        streamTypeGeneric,
					StreamLimit: m.maxStream.StreamNum(),
				})
			} else {
				m.queueStreamIDBlocked(&wire.StreamsBlockedFrame{
					Type:        streamTypeGeneric,
					StreamLimit: 0,
				})
			}
			m.blockedSent = true
		}
		return nil, qerr.TooManyOpenStreams
	}
	s := m.newStream(m.nextStream)
	m.streams[m.nextStream] = s
	m.nextStream += 4
	return s, nil
}

func (m *outgoingItemsMap) GetStream(id protocol.StreamID) (item, error) {
	m.mutex.RLock()
	if id >= m.nextStream {
		m.mutex.RUnlock()
		return nil, qerr.Error(qerr.InvalidStreamID, fmt.Sprintf("peer attempted to open stream %d", id))
	}
	s := m.streams[id]
	m.mutex.RUnlock()
	return s, nil
}

func (m *outgoingItemsMap) DeleteStream(id protocol.StreamID) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if _, ok := m.streams[id]; !ok {
		return fmt.Errorf("Tried to delete unknown stream %d", id)
	}
	delete(m.streams, id)
	return nil
}

func (m *outgoingItemsMap) SetMaxStream(id protocol.StreamID) {
	m.mutex.Lock()
	if !m.maxStreamSet || id > m.maxStream {
		m.maxStream = id
		m.maxStreamSet = true
		m.blockedSent = false
		m.cond.Broadcast()
	}
	m.mutex.Unlock()
}

func (m *outgoingItemsMap) CloseWithError(err error) {
	m.mutex.Lock()
	m.closeErr = err
	for _, str := range m.streams {
		str.closeForShutdown(err)
	}
	m.cond.Broadcast()
	m.mutex.Unlock()
}
