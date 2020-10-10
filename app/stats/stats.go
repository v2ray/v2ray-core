// +build !confonly

package stats

//go:generate go run v2ray.com/core/common/errors/errorgen

import (
	"context"
	"sync"

	"v2ray.com/core/common"
	"v2ray.com/core/common/errors"
	"v2ray.com/core/features/stats"
)

// Manager is an implementation of stats.Manager.
type Manager struct {
	access   sync.RWMutex
	counters map[string]*Counter
	channels map[string]*Channel
	running  bool
}

// NewManager creates an instance of Statistics Manager.
func NewManager(ctx context.Context, config *Config) (*Manager, error) {
	m := &Manager{
		counters: make(map[string]*Counter),
		channels: make(map[string]*Channel),
	}

	return m, nil
}

// Type implements common.HasType.
func (*Manager) Type() interface{} {
	return stats.ManagerType()
}

// RegisterCounter implements stats.Manager.
func (m *Manager) RegisterCounter(name string) (stats.Counter, error) {
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

// UnregisterCounter implements stats.Manager.
func (m *Manager) UnregisterCounter(name string) error {
	m.access.Lock()
	defer m.access.Unlock()

	if _, found := m.counters[name]; found {
		newError("remove counter ", name).AtDebug().WriteToLog()
		delete(m.counters, name)
	}
	return nil
}

// GetCounter implements stats.Manager.
func (m *Manager) GetCounter(name string) stats.Counter {
	m.access.RLock()
	defer m.access.RUnlock()

	if c, found := m.counters[name]; found {
		return c
	}
	return nil
}

// VisitCounters calls visitor function on all managed counters.
func (m *Manager) VisitCounters(visitor func(string, stats.Counter) bool) {
	m.access.RLock()
	defer m.access.RUnlock()

	for name, c := range m.counters {
		if !visitor(name, c) {
			break
		}
	}
}

// RegisterChannel implements stats.Manager.
func (m *Manager) RegisterChannel(name string) (stats.Channel, error) {
	m.access.Lock()
	defer m.access.Unlock()

	if _, found := m.channels[name]; found {
		return nil, newError("Channel ", name, " already registered.")
	}
	newError("create new channel ", name).AtDebug().WriteToLog()
	c := NewChannel(&ChannelConfig{BufferSize: 64, Blocking: false})
	m.channels[name] = c
	if m.running {
		return c, c.Start()
	}
	return c, nil
}

// UnregisterChannel implements stats.Manager.
func (m *Manager) UnregisterChannel(name string) error {
	m.access.Lock()
	defer m.access.Unlock()

	if c, found := m.channels[name]; found {
		newError("remove channel ", name).AtDebug().WriteToLog()
		delete(m.channels, name)
		return c.Close()
	}
	return nil
}

// GetChannel implements stats.Manager.
func (m *Manager) GetChannel(name string) stats.Channel {
	m.access.RLock()
	defer m.access.RUnlock()

	if c, found := m.channels[name]; found {
		return c
	}
	return nil
}

// Start implements common.Runnable.
func (m *Manager) Start() error {
	m.access.Lock()
	defer m.access.Unlock()
	m.running = true
	errs := []error{}
	for _, channel := range m.channels {
		if err := channel.Start(); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) != 0 {
		return errors.Combine(errs...)
	}
	return nil
}

// Close implement common.Closable.
func (m *Manager) Close() error {
	m.access.Lock()
	defer m.access.Unlock()
	m.running = false
	errs := []error{}
	for name, channel := range m.channels {
		newError("remove channel ", name).AtDebug().WriteToLog()
		delete(m.channels, name)
		if err := channel.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) != 0 {
		return errors.Combine(errs...)
	}
	return nil
}

func init() {
	common.Must(common.RegisterConfig((*Config)(nil), func(ctx context.Context, config interface{}) (interface{}, error) {
		return NewManager(ctx, config.(*Config))
	}))
}