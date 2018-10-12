package core

import (
	"sync"

	"v2ray.com/core/features/stats"
)

type syncStatManager struct {
	sync.RWMutex
	stats.Manager
}

func (s *syncStatManager) Start() error {
	s.RLock()
	defer s.RUnlock()

	if s.Manager == nil {
		return nil
	}

	return s.Manager.Start()
}

func (s *syncStatManager) Close() error {
	s.RLock()
	defer s.RUnlock()

	if s.Manager == nil {
		return nil
	}
	return s.Manager.Close()
}

func (s *syncStatManager) RegisterCounter(name string) (stats.Counter, error) {
	s.RLock()
	defer s.RUnlock()

	if s.Manager == nil {
		return nil, newError("StatManager not set.")
	}
	return s.Manager.RegisterCounter(name)
}

func (s *syncStatManager) GetCounter(name string) stats.Counter {
	s.RLock()
	defer s.RUnlock()

	if s.Manager == nil {
		return nil
	}
	return s.Manager.GetCounter(name)
}

func (s *syncStatManager) Set(m stats.Manager) {
	if m == nil {
		return
	}
	s.Lock()
	defer s.Unlock()

	if s.Manager != nil {
		s.Manager.Close() // nolint: errcheck
	}
	s.Manager = m
}
