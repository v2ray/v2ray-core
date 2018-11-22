package policy

import (
	"time"
)

// DefaultManager is the implementation of the Manager.
type DefaultManager struct{}

// Type implements common.HasType.
func (DefaultManager) Type() interface{} {
	return ManagerType()
}

// ForLevel implements Manager.
func (DefaultManager) ForLevel(level uint32) Session {
	p := SessionDefault()
	if level == 1 {
		p.Timeouts.ConnectionIdle = time.Second * 600
	}
	return p
}

// ForSystem implements Manager.
func (DefaultManager) ForSystem() System {
	return System{}
}

// Start implements common.Runnable.
func (DefaultManager) Start() error {
	return nil
}

// Close implements common.Closable.
func (DefaultManager) Close() error {
	return nil
}
