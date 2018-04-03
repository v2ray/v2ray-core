package core

import (
	"sync"
)

type StatCounter interface {
	Value() int64
	Set(int64) int64
	Add(int64) int64
}

type StatManager interface {
	Feature

	RegisterCounter(string) (StatCounter, error)
	GetCounter(string) StatCounter
}

type syncStatManager struct {
	sync.RWMutex
	StatManager
}

func (s *syncStatManager) Start() error {
	s.RLock()
	defer s.RUnlock()

	if s.StatManager == nil {
		return nil
	}

	return s.StatManager.Start()
}

func (s *syncStatManager) Close() error {
	s.RLock()
	defer s.RUnlock()

	if s.StatManager == nil {
		return nil
	}
	return s.StatManager.Close()
}

func (s *syncStatManager) RegisterCounter(name string) (StatCounter, error) {
	s.RLock()
	defer s.RUnlock()

	if s.StatManager == nil {
		return nil, newError("StatManager not set.")
	}
	return s.StatManager.RegisterCounter(name)
}

func (s *syncStatManager) GetCounter(name string) StatCounter {
	s.RLock()
	defer s.RUnlock()

	if s.StatManager == nil {
		return nil
	}
	return s.StatManager.GetCounter(name)
}

func (s *syncStatManager) Set(m StatManager) {
	if m == nil {
		return
	}
	s.Lock()
	defer s.Unlock()

	if s.StatManager != nil {
		s.StatManager.Close()
	}
	s.StatManager = m
}
