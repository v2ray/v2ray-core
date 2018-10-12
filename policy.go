package core

import (
	"sync"
	"time"

	"v2ray.com/core/common"
	"v2ray.com/core/features/policy"
)

type syncPolicyManager struct {
	sync.RWMutex
	policy.Manager
}

func (*syncPolicyManager) Type() interface{} {
	return policy.ManagerType()
}

func (m *syncPolicyManager) ForLevel(level uint32) policy.Session {
	m.RLock()
	defer m.RUnlock()

	if m.Manager == nil {
		p := policy.SessionDefault()
		if level == 1 {
			p.Timeouts.ConnectionIdle = time.Second * 600
		}
		return p
	}

	return m.Manager.ForLevel(level)
}

func (m *syncPolicyManager) ForSystem() policy.System {
	m.RLock()
	defer m.RUnlock()

	if m.Manager == nil {
		return policy.System{}
	}

	return m.Manager.ForSystem()
}

func (m *syncPolicyManager) Start() error {
	m.RLock()
	defer m.RUnlock()

	if m.Manager == nil {
		return nil
	}

	return m.Manager.Start()
}

func (m *syncPolicyManager) Close() error {
	m.RLock()
	defer m.RUnlock()

	return common.Close(m.Manager)
}

func (m *syncPolicyManager) Set(manager policy.Manager) {
	if manager == nil {
		return
	}

	m.Lock()
	defer m.Unlock()

	common.Close(m.Manager) // nolint: errcheck
	m.Manager = manager
}
