package core

import (
	"sync"
	"time"

	"v2ray.com/core/common"
)

// TimeoutPolicy contains limits for connection timeout.
type TimeoutPolicy struct {
	// Timeout for handshake phase in a connection.
	Handshake time.Duration
	// Timeout for connection being idle, i.e., there is no egress or ingress traffic in this connection.
	ConnectionIdle time.Duration
	// Timeout for an uplink only connection, i.e., the downlink of the connection has been closed.
	UplinkOnly time.Duration
	// Timeout for an downlink only connection, i.e., the uplink of the connection has been closed.
	DownlinkOnly time.Duration
}

// StatsPolicy contains settings for stats counters.
type StatsPolicy struct {
	// Whether or not to enable stat counter for user uplink traffic.
	UserUplink bool
	// Whether or not to enable stat counter for user downlink traffic.
	UserDownlink bool
}

type SystemStatsPolicy struct {
	// Whether or not to enable stat counter for uplink traffic in inbound handlers.
	InboundUplink bool
	// Whether or not to enable stat counter for downlink traffic in inbound handlers.
	InboundDownlink bool
}

type SystemPolicy struct {
	Stats SystemStatsPolicy
}

// Policy is session based settings for controlling V2Ray requests. It contains various settings (or limits) that may differ for different users in the context.
type Policy struct {
	Timeouts TimeoutPolicy // Timeout settings
	Stats    StatsPolicy
}

// PolicyManager is a feature that provides Policy for the given user by its id or level.
type PolicyManager interface {
	Feature

	// ForLevel returns the Policy for the given user level.
	ForLevel(level uint32) Policy

	// ForSystem returns the Policy for V2Ray system.
	ForSystem() SystemPolicy
}

// DefaultPolicy returns the Policy when user is not specified.
func DefaultPolicy() Policy {
	return Policy{
		Timeouts: TimeoutPolicy{
			Handshake:      time.Second * 4,
			ConnectionIdle: time.Second * 300,
			UplinkOnly:     time.Second * 2,
			DownlinkOnly:   time.Second * 5,
		},
		Stats: StatsPolicy{
			UserUplink:   false,
			UserDownlink: false,
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

func (m *syncPolicyManager) Close() error {
	m.RLock()
	defer m.RUnlock()

	return common.Close(m.PolicyManager)
}

func (m *syncPolicyManager) Set(manager PolicyManager) {
	if manager == nil {
		return
	}

	m.Lock()
	defer m.Unlock()

	common.Close(m.PolicyManager)
	m.PolicyManager = manager
}
