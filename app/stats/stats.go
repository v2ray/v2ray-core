package stats

//go:generate go run $GOPATH/src/v2ray.com/core/common/errors/errorgen/main.go -pkg stats -path App,Stats

import (
	"context"
	"sync"
	"sync/atomic"

	"v2ray.com/core"
)

// Counter is an implementation of core.StatCounter.
type Counter struct {
	value int64
}

// Value implements core.StatCounter.
func (c *Counter) Value() int64 {
	return atomic.LoadInt64(&c.value)
}

// Set implements core.StatCounter.
func (c *Counter) Set(newValue int64) int64 {
	return atomic.SwapInt64(&c.value, newValue)
}

// Add implements core.StatCounter.
func (c *Counter) Add(delta int64) int64 {
	return atomic.AddInt64(&c.value, delta)
}

// Manager is an implementation of core.StatManager.
type Manager struct {
	access   sync.RWMutex
	counters map[string]*Counter
}

func NewManager(ctx context.Context, config *Config) (*Manager, error) {
	m := &Manager{
		counters: make(map[string]*Counter),
	}

	v := core.FromContext(ctx)
	if v != nil {
		if err := v.RegisterFeature((*core.StatManager)(nil), m); err != nil {
			return nil, newError("failed to register StatManager").Base(err)
		}
	}

	return m, nil
}

func (m *Manager) RegisterCounter(name string) (core.StatCounter, error) {
	m.access.Lock()
	defer m.access.Unlock()

	if _, found := m.counters[name]; found {
		return nil, newError("Counter ", name, " already registered.")
	}
	newError("create new counter ", name).AtDebug().WriteToLog()
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

// Start implements common.Runnable.
func (m *Manager) Start() error {
	return nil
}

// Close implement common.Closable.
func (m *Manager) Close() error {
	return nil
}
