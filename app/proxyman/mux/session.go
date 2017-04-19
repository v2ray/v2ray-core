package mux

import (
	"sync"

	"v2ray.com/core/transport/ray"
)

type SessionManager struct {
	sync.RWMutex
	count    uint16
	sessions map[uint16]*Session
	closed   bool
}

func NewSessionManager() *SessionManager {
	return &SessionManager{
		count:    0,
		sessions: make(map[uint16]*Session, 32),
	}
}

func (m *SessionManager) Size() int {
	m.RLock()
	defer m.RUnlock()

	return len(m.sessions)
}

func (m *SessionManager) Allocate() *Session {
	m.Lock()
	defer m.Unlock()

	if m.closed {
		return nil
	}

	m.count++
	s := &Session{
		ID:     m.count,
		parent: m,
	}
	m.sessions[s.ID] = s
	return s
}

func (m *SessionManager) Add(s *Session) {
	m.Lock()
	defer m.Unlock()

	m.sessions[s.ID] = s
}

func (m *SessionManager) Remove(id uint16) {
	m.Lock()
	defer m.Unlock()

	delete(m.sessions, id)
}

func (m *SessionManager) Get(id uint16) (*Session, bool) {
	m.RLock()
	defer m.RUnlock()

	if m.closed {
		return nil, false
	}

	s, found := m.sessions[id]
	return s, found
}

func (m *SessionManager) CloseIfNoSession() bool {
	m.RLock()
	defer m.RUnlock()

	if m.closed {
		return true
	}

	if len(m.sessions) != 0 {
		return false
	}

	m.closed = true
	return true
}

func (m *SessionManager) Close() {
	m.RLock()
	defer m.RUnlock()

	if m.closed {
		return
	}

	m.closed = true

	for _, s := range m.sessions {
		s.input.Close()
		s.output.Close()
	}

	m.sessions = make(map[uint16]*Session)
}

type Session struct {
	sync.Mutex
	input          ray.InputStream
	output         ray.OutputStream
	parent         *SessionManager
	ID             uint16
	uplinkClosed   bool
	downlinkClosed bool
}

func (s *Session) CloseUplink() {
	var allDone bool
	s.Lock()
	s.uplinkClosed = true
	allDone = s.uplinkClosed && s.downlinkClosed
	s.Unlock()
	if allDone {
		s.parent.Remove(s.ID)
	}
}

func (s *Session) CloseDownlink() {
	var allDone bool
	s.Lock()
	s.downlinkClosed = true
	allDone = s.uplinkClosed && s.downlinkClosed
	s.Unlock()
	if allDone {
		s.parent.Remove(s.ID)
	}
}
