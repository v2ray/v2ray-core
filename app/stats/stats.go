// +build !confonly

package stats

//go:generate errorgen

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"v2ray.com/core/features/stats"
)

// Counter is an implementation of stats.Counter.
type Counter struct {
	value int64
}

// Value implements stats.Counter.
func (c *Counter) Value() int64 {
	return atomic.LoadInt64(&c.value)
}

// Set implements stats.Counter.
func (c *Counter) Set(newValue int64) int64 {
	return atomic.SwapInt64(&c.value, newValue)
}

// Add implements stats.Counter.
func (c *Counter) Add(delta int64) int64 {
	return atomic.AddInt64(&c.value, delta)
}

// Channel is an implementation of stats.Channel
type Channel struct {
	channel     chan interface{}
	subscribers []chan interface{}
	access      sync.RWMutex
}

// Channel implements stats.Channel
func (c *Channel) Channel() chan interface{} {
	return c.channel
}

// Subscribers implements stats.Channel
func (c *Channel) Subscribers() []chan interface{} {
	c.access.RLock()
	defer c.access.RUnlock()
	return c.subscribers
}

// Subscribe implements stats.Channel
func (c *Channel) Subscribe() chan interface{} {
	c.access.Lock()
	defer c.access.Unlock()
	ch := make(chan interface{})
	c.subscribers = append(c.subscribers, ch)
	return ch
}

// Unsubscribe implements stats.Channel
func (c *Channel) Unsubscribe(ch chan interface{}) {
	c.access.Lock()
	defer c.access.Unlock()
	for i, s := range c.subscribers {
		if s == ch {
			// Copy to new memory block to prevent modifying original data
			subscribers := make([]chan interface{}, len(c.subscribers)-1)
			copy(subscribers[:i], c.subscribers[:i])
			copy(subscribers[i:], c.subscribers[i+1:])
			c.subscribers = subscribers
			return
		}
	}
}

// Start starts the channel for listening to messsages
func (c *Channel) Start() {
	for message := range c.Channel() {
		subscribers := c.Subscribers() // Store a copy of slice value for concurrency safety
		for _, sub := range subscribers {
			select {
			case sub <- message: // Successfully sent message
			case <-time.After(100 * time.Millisecond):
				c.Unsubscribe(sub) // Remove timeout subscriber
				close(sub)         // Actively close subscriber as notification
			}
		}
	}
}

// Manager is an implementation of stats.Manager.
type Manager struct {
	access   sync.RWMutex
	counters map[string]*Counter
	channels map[string]*Channel
}

func NewManager(ctx context.Context, config *Config) (*Manager, error) {
	m := &Manager{
		counters: make(map[string]*Counter),
		channels: make(map[string]*Channel),
	}

	return m, nil
}

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
	c := &Channel{channel: make(chan interface{})}
	m.channels[name] = c
	go c.Start()
	return c, nil
}

// UnregisterChannel implements stats.Manager.
func (m *Manager) UnregisterChannel(name string) error {
	m.access.Lock()
	defer m.access.Unlock()

	if _, found := m.channels[name]; found {
		newError("remove channel ", name).AtDebug().WriteToLog()
		delete(m.channels, name)
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
	return nil
}

// Close implement common.Closable.
func (m *Manager) Close() error {
	return nil
}

