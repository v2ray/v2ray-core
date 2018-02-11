package core

import (
	"sync"
)

// Commander is a feature that accepts commands from external source.
type Commander interface {
	Feature
}

type syncCommander struct {
	sync.RWMutex
	Commander
}

func (c *syncCommander) Start() error {
	c.RLock()
	defer c.RUnlock()

	if c.Commander == nil {
		return nil
	}

	return c.Commander.Start()
}

func (c *syncCommander) Close() error {
	c.RLock()
	defer c.RUnlock()

	if c.Commander == nil {
		return nil
	}

	return c.Commander.Close()
}

func (c *syncCommander) Set(commander Commander) {
	c.Lock()
	defer c.Unlock()

	c.Commander = commander
}
