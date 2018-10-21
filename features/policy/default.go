package policy

import (
	"time"
)

type DefaultManager struct{}

func (DefaultManager) Type() interface{} {
	return ManagerType()
}

func (DefaultManager) ForLevel(level uint32) Session {
	p := SessionDefault()
	if level == 1 {
		p.Timeouts.ConnectionIdle = time.Second * 600
	}
	return p
}

func (DefaultManager) ForSystem() System {
	return System{}
}

func (DefaultManager) Start() error {
	return nil
}

func (DefaultManager) Close() error {
	return nil
}
