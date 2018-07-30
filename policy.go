package core

import (
	"context"
	"sync"
	"time"

	"v2ray.com/core/common"
	"v2ray.com/core/common/platform"
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

// BufferPolicy contains settings for internal buffer.
type BufferPolicy struct {
	// Size of buffer per connection, in bytes. -1 for unlimited buffer.
	PerConnection int32
}

// SystemStatsPolicy contains stat policy settings on system level.
type SystemStatsPolicy struct {
	// Whether or not to enable stat counter for uplink traffic in inbound handlers.
	InboundUplink bool
	// Whether or not to enable stat counter for downlink traffic in inbound handlers.
	InboundDownlink bool
}

// SystemPolicy contains policy settings at system level.
type SystemPolicy struct {
	Stats  SystemStatsPolicy
	Buffer BufferPolicy
}

// Policy is session based settings for controlling V2Ray requests. It contains various settings (or limits) that may differ for different users in the context.
type Policy struct {
	Timeouts TimeoutPolicy // Timeout settings
	Stats    StatsPolicy
	Buffer   BufferPolicy
}

// PolicyManager is a feature that provides Policy for the given user by its id or level.
type PolicyManager interface {
	Feature

	// ForLevel returns the Policy for the given user level.
	ForLevel(level uint32) Policy

	// ForSystem returns the Policy for V2Ray system.
	ForSystem() SystemPolicy
}

var defaultBufferSize int32 = 2 * 1024 * 1024

func init() {
	const key = "v2ray.ray.buffer.size"
	size := platform.EnvFlag{
		Name:    key,
		AltName: platform.NormalizeEnvName(key),
	}.GetValueAsInt(10)
	if size == 0 {
		defaultBufferSize = -1
	} else {
		defaultBufferSize = int32(size) * 1024 * 1024
	}
}

func defaultBufferPolicy() BufferPolicy {
	return BufferPolicy{
		PerConnection: defaultBufferSize,
	}
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
		Buffer: defaultBufferPolicy(),
	}
}

type policyKey int

const (
	bufferPolicyKey policyKey = 0
)

func ContextWithBufferPolicy(ctx context.Context, p BufferPolicy) context.Context {
	return context.WithValue(ctx, bufferPolicyKey, p)
}

func BufferPolicyFromContext(ctx context.Context) BufferPolicy {
	pPolicy := ctx.Value(bufferPolicyKey)
	if pPolicy == nil {
		return defaultBufferPolicy()
	}
	return pPolicy.(BufferPolicy)
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

func (m *syncPolicyManager) ForSystem() SystemPolicy {
	m.RLock()
	defer m.RUnlock()

	if m.PolicyManager == nil {
		return SystemPolicy{}
	}

	return m.PolicyManager.ForSystem()
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

	common.Close(m.PolicyManager) // nolint: errcheck
	m.PolicyManager = manager
}
