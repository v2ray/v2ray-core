package stats

//go:generate go run $GOPATH/src/v2ray.com/core/common/errors/errorgen/main.go -pkg stats -path App,Stats

import (
	"context"
	"sync"
	"sync/atomic"

	"v2ray.com/core"
)

type Counter struct {
	value int64
}

func (c *Counter) Value() int64 {
	return atomic.LoadInt64(&c.value)
}

func (c *Counter) Exchange(newValue int64) int64 {
	return atomic.SwapInt64(&c.value, newValue)
}

func (c *Counter) Add(delta int64) int64 {
	return atomic.AddInt64(&c.value, delta)
}

type Manager struct {
	access   sync.RWMutex
	counters map[string]*Counter
}

func NewManager(ctx context.Context, config *Config) (*Manager, error) {
	return &Manager{
		counters: make(map[string]*Counter),
	}, nil
}

func (m *Manager) RegisterCounter(name string) (core.StatCounter, error) {
	m.access.Lock()
	defer m.access.Unlock()

	if _, found := m.counters[name]; found {
		return nil, newError("Counter ", name, " already registered.")
	}
	c := new(Counter)
	m.counters[name] = c
	return c, nil
}

func (m *Manager) GetCounter(name string) core.StatCounter {
	m.access.RLock()
	defer m.access.RUnlock()

	if c, found := m.counters[name]; found {
		return c
	}
	return nil
}

func (m *Manager) Start() error {
	return nil
}

func (m *Manager) Close() error {
	return nil
}
