package core

import (
	"sync"

	"google.golang.org/grpc"
)

// ServiceRegistryCallback is a callback function for registering services.
type ServiceRegistryCallback func(s *grpc.Server)

// Commander is a feature that accepts commands from external source.
type Commander interface {
	Feature

	// RegisterService registers a service into this Commander.
	RegisterService(ServiceRegistryCallback)
}

type syncCommander struct {
	sync.RWMutex
	Commander
}

func (c *syncCommander) RegisterService(callback ServiceRegistryCallback) {
	c.RLock()
	defer c.RUnlock()

	if c.Commander == nil {
		return
	}

	c.Commander.RegisterService(callback)
}

func (c *syncCommander) Start() error {
	c.RLock()
	defer c.RUnlock()

	if c.Commander == nil {
		return nil
	}

	return c.Commander.Start()
}

func (c *syncCommander) Close() {
	c.RLock()
	defer c.RUnlock()

	if c.Commander == nil {
		return
	}

	c.Commander.Close()
}

func (c *syncCommander) Set(commander Commander) {
	c.Lock()
	defer c.Unlock()

	c.Commander = commander
}
