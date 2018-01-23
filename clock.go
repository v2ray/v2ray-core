package core

import (
	"sync"
	"time"
)

// Clock is a V2Ray feature that returns current time.
type Clock interface {
	Feature

	// Now returns current time.
	Now() time.Time
}

type syncClock struct {
	sync.RWMutex
	Clock
}

func (c *syncClock) Now() time.Time {
	c.RLock()
	defer c.RUnlock()

	if c.Clock == nil {
		return time.Now()
	}

	return c.Clock.Now()
}

func (c *syncClock) Start() error {
	c.RLock()
	defer c.RUnlock()

	if c.Clock == nil {
		return nil
	}

	return c.Clock.Start()
}

func (c *syncClock) Close() {
	c.RLock()
	defer c.RUnlock()

	if c.Clock == nil {
		return
	}

	c.Clock.Close()
}

func (c *syncClock) Set(clock Clock) {
	c.Lock()
	defer c.Unlock()

	c.Clock = clock
}
