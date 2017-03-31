package mux

import (
	"context"
	"sync"
	"time"

	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/signal"
	"v2ray.com/core/proxy"
	"v2ray.com/core/transport/ray"
)

const (
	maxParallel = 8
	maxTotal    = 128
)

type clientSession struct {
	sync.Mutex
	outboundRay    ray.OutboundRay
	parent         *Client
	id             uint16
	uplinkClosed   bool
	downlinkClosed bool
}

func (s *clientSession) checkAndRemove() {
	s.Lock()
	if s.uplinkClosed && s.downlinkClosed {
		s.parent.remove(s.id)
	}
	s.Unlock()
}

func (s *clientSession) closeUplink() {
	s.Lock()
	s.uplinkClosed = true
	s.Unlock()
	s.checkAndRemove()
}

func (s *clientSession) closeDownlink() {
	s.Lock()
	s.downlinkClosed = true
	s.Unlock()
	s.checkAndRemove()
}

type Client struct {
	access     sync.RWMutex
	count      uint16
	sessions   map[uint16]*clientSession
	inboundRay ray.InboundRay
}

func (m *Client) IsFullyOccupied() bool {
	m.access.RLock()
	defer m.access.RUnlock()

	return len(m.sessions) >= maxParallel
}

func (m *Client) IsFullyUsed() bool {
	m.access.RLock()
	defer m.access.RUnlock()

	return m.count >= maxTotal
}

func (m *Client) remove(id uint16) {
	m.access.Lock()
	delete(m.sessions, id)
	m.access.Unlock()
}

func (m *Client) fetchInput(ctx context.Context, session *clientSession) {
	dest, _ := proxy.TargetFromContext(ctx)
	writer := &muxWriter{
		dest:   dest,
		id:     session.id,
		writer: m.inboundRay.InboundInput(),
	}
	_, timer := signal.CancelAfterInactivity(ctx, time.Minute*5)
	buf.PipeUntilEOF(timer, session.outboundRay.OutboundInput(), writer)
	writer.Close()
	session.closeUplink()
}

func (m *Client) Dispatch(ctx context.Context, outboundRay ray.OutboundRay) {
	m.access.Lock()
	defer m.access.Unlock()

	m.count++
	id := m.count
	session := &clientSession{
		outboundRay: outboundRay,
		parent:      m,
		id:          id,
	}
	m.sessions[id] = session
	go m.fetchInput(ctx, session)
}

func (m *Client) fetchOutput() {
	reader := NewReader(m.inboundRay.InboundOutput())
	for {
		meta, err := reader.ReadMetadata()
		if err != nil {
			break
		}
		m.access.RLock()
		session, found := m.sessions[meta.SessionID]
		m.access.RUnlock()
		if found && meta.SessionStatus == SessionStatusEnd {
			session.closeDownlink()
		}
		if !meta.Option.Has(OptionData) {
			continue
		}

		for {
			data, more, err := reader.Read()
			if err != nil {
				break
			}
			if found {
				if err := session.outboundRay.OutboundOutput().Write(data); err != nil {
					break
				}
			}
			if !more {
				break
			}
		}
	}
}
