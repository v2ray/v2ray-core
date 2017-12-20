package event

import "sync"

type Event uint16

type Handler func(data interface{}) error

type Registry interface {
	On(Event, Handler)
}

type Listener struct {
	sync.RWMutex
	events map[Event][]Handler
}

func (l *Listener) On(e Event, h Handler) {
	l.Lock()
	defer l.Unlock()

	if l.events == nil {
		l.events = make(map[Event][]Handler)
	}

	handlers := l.events[e]
	handlers = append(handlers, h)
	l.events[e] = handlers
}

func (l *Listener) Fire(e Event, data interface{}) error {
	l.RLock()
	defer l.RUnlock()

	if l.events == nil {
		return nil
	}

	for _, h := range l.events[e] {
		if err := h(data); err != nil {
			return err
		}
	}

	return nil
}
