package mux

import (
	"sync"

	"v2ray.com/core/transport/ray"
)

type SessionManager struct {
	sync.RWMutex
	count    uint16
	sessions map[uint16]*Session
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

func (m *SessionManager) Allocate(s *Session) {
	m.Lock()
	defer m.Unlock()

	m.count++
	s.ID = m.count
	m.sessions[s.ID] = s
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

	s, found := m.sessions[id]
	return s, found
}

func (m *SessionManager) Close() {
	m.RLock()
	defer m.RUnlock()

	for _, s := range m.sessions {
		s.output.CloseError()
	}
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
