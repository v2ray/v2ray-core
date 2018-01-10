package core

import (
	"sync"
	"time"
)

// TimeoutPolicy contains limits for connection timeout.
type TimeoutPolicy struct {
	// Timeout for handshake phase in a connection.
	Handshake time.Duration
	// Timeout for connection being idle, i.e., there is no egress or ingress traffic in this connection.
	ConnectionIdle time.Duration
	// Timeout for an uplink only connection, i.e., the downlink of the connection has ben closed.
	UplinkOnly time.Duration
	// Timeout for an downlink only connection, i.e., the uplink of the connection has ben closed.
	DownlinkOnly time.Duration
}

// OverrideWith overrides the current TimeoutPolicy with another one. All timeouts with zero value will be overridden with the new value.
func (p TimeoutPolicy) OverrideWith(another TimeoutPolicy) TimeoutPolicy {
	if p.Handshake == 0 {
		p.Handshake = another.Handshake
	}
	if p.ConnectionIdle == 0 {
		p.ConnectionIdle = another.ConnectionIdle
	}
	if p.UplinkOnly == 0 {
		p.UplinkOnly = another.UplinkOnly
	}
	if p.DownlinkOnly == 0 {
		p.DownlinkOnly = another.DownlinkOnly
	}
	return p
}

// Policy is session based settings for controlling V2Ray requests. It contains various settings (or limits) that may differ for different users in the context.
type Policy struct {
	Timeouts TimeoutPolicy // Timeout settings
}

// OverrideWith overrides the current Policy with another one. All values with default value will be overridden.
func (p Policy) OverrideWith(another Policy) Policy {
	p.Timeouts.OverrideWith(another.Timeouts)
	return p
}

// PolicyManager is a feature that provides Policy for the given user by its id or level.
type PolicyManager interface {
	Feature

	// ForLevel returns the Policy for the given user level.
	ForLevel(level uint32) Policy
}

// DefaultPolicy returns the Policy when user is not specified.
func DefaultPolicy() Policy {
	return Policy{
		Timeouts: TimeoutPolicy{
			Handshake:      time.Second * 4,
			ConnectionIdle: time.Second * 300,
			UplinkOnly:     time.Second * 5,
			DownlinkOnly:   time.Second * 30,
		},
	}
}

type syncPolicyManager struct {
	sync.RWMutex
	PolicyManager
}

func (m *syncPolicyManager) ForLevel(level uint32) Policy {
	m.RLock()
	defer m.RUnlock()

	if m.PolicyManager == nil {
		p := DefaultPolicy()
		if level == 1 {
			p.Timeouts.ConnectionIdle = time.Second * 600
		}
		return p
	}

	return m.PolicyManager.ForLevel(level)
}

func (m *syncPolicyManager) Start() error {
	m.RLock()
	defer m.RUnlock()

	if m.PolicyManager == nil {
		return nil
	}

	return m.PolicyManager.Start()
}

func (m *syncPolicyManager) Close() {
	m.RLock()
	defer m.RUnlock()

	if m.PolicyManager != nil {
		m.PolicyManager.Close()
	}
}

func (m *syncPolicyManager) Set(manager PolicyManager) {
	if manager == nil {
		return
	}

	m.Lock()
	defer m.Unlock()

	m.PolicyManager = manager
}
