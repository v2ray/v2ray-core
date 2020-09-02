package v2board

import (
	"sync"
)

type UserSet struct {
	m map[V2RayUser]bool
	sync.RWMutex
}

func NewUserSet() *UserSet {
	return &UserSet{
		m: map[V2RayUser]bool{},
	}
}

func (s *UserSet) Add(item V2RayUser) {
	s.Lock()
	defer s.Unlock()
	s.m[item] = true
}

func (s *UserSet) Remove(item V2RayUser) {
	s.Lock()
	s.Unlock()
	delete(s.m, item)
}

func (s *UserSet) Has(item V2RayUser) bool {
	s.RLock()
	defer s.RUnlock()
	_, ok := s.m[item]
	return ok
}

func (s *UserSet) Len() int {
	return len(s.List())
}

func (s *UserSet) Clear() {
	s.Lock()
	defer s.Unlock()
	s.m = map[V2RayUser]bool{}
}

func (s *UserSet) IsEmpty() bool {
	if s.Len() == 0 {
		return true
	}
	return false
}

func (s *UserSet) List() []V2RayUser {
	s.RLock()
	defer s.RUnlock()
	list := []V2RayUser{}
	for item := range s.m {
		list = append(list, item)
	}
	return list
}
