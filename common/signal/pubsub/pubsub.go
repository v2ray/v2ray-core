package pubsub

import (
	"sync"
	"time"

	"v2ray.com/core/common"
	"v2ray.com/core/common/task"
)

type Subscriber struct {
	name    string
	buffer  chan interface{}
	removed chan struct{}
}

func (s *Subscriber) push(msg interface{}) {
	select {
	case s.buffer <- msg:
	default:
	}
}

func (s *Subscriber) Wait() <-chan interface{} {
	return s.buffer
}

func (s *Subscriber) Close() {
	close(s.removed)
}

func (s *Subscriber) IsClosed() bool {
	select {
	case <-s.removed:
		return true
	default:
		return false
	}
}

type Service struct {
	sync.RWMutex
	subs  []*Subscriber
	ctask *task.Periodic
}

func NewService() *Service {
	s := &Service{}
	s.ctask = &task.Periodic{
		Execute:  s.cleanup,
		Interval: time.Second * 30,
	}
	common.Must(s.ctask.Start())
	return s
}

func (s *Service) cleanup() error {
	s.Lock()
	defer s.Unlock()

	if len(s.subs) < 16 {
		return nil
	}

	newSub := make([]*Subscriber, 0, len(s.subs))
	for _, sub := range s.subs {
		if !sub.IsClosed() {
			newSub = append(newSub, sub)
		}
	}

	s.subs = newSub
	return nil
}

func (s *Service) Subscribe(name string) *Subscriber {
	sub := &Subscriber{
		name:    name,
		buffer:  make(chan interface{}, 16),
		removed: make(chan struct{}),
	}
	s.Lock()
	s.subs = append(s.subs, sub)
	s.Unlock()
	return sub
}

func (s *Service) Publish(name string, message interface{}) {
	s.RLock()
	defer s.RUnlock()

	for _, sub := range s.subs {
		if sub.name == name && !sub.IsClosed() {
			sub.push(message)
		}
	}
}
