package policy

import (
	"context"
	"runtime"
	"time"

	"v2ray.com/core/common/platform"
	"v2ray.com/core/features"
)

// Timeout contains limits for connection timeout.
type Timeout struct {
	// Timeout for handshake phase in a connection.
	Handshake time.Duration
	// Timeout for connection being idle, i.e., there is no egress or ingress traffic in this connection.
	ConnectionIdle time.Duration
	// Timeout for an uplink only connection, i.e., the downlink of the connection has been closed.
	UplinkOnly time.Duration
	// Timeout for an downlink only connection, i.e., the uplink of the connection has been closed.
	DownlinkOnly time.Duration
}

// Stats contains settings for stats counters.
type Stats struct {
	// Whether or not to enable stat counter for user uplink traffic.
	UserUplink bool
	// Whether or not to enable stat counter for user downlink traffic.
	UserDownlink bool
}

// Buffer contains settings for internal buffer.
type Buffer struct {
	// Size of buffer per connection, in bytes. -1 for unlimited buffer.
	PerConnection int32
}

// SystemStats contains stat policy settings on system level.
type SystemStats struct {
	// Whether or not to enable stat counter for uplink traffic in inbound handlers.
	InboundUplink bool
	// Whether or not to enable stat counter for downlink traffic in inbound handlers.
	InboundDownlink bool
}

// System contains policy settings at system level.
type System struct {
	Stats  SystemStats
	Buffer Buffer
}

// Session is session based settings for controlling V2Ray requests. It contains various settings (or limits) that may differ for different users in the context.
type Session struct {
	Timeouts Timeout // Timeout settings
	Stats    Stats
	Buffer   Buffer
}

// Manager is a feature that provides Policy for the given user by its id or level.
//
// v2ray:api:stable
type Manager interface {
	features.Feature

	// ForLevel returns the Session policy for the given user level.
	ForLevel(level uint32) Session

	// ForSystem returns the System policy for V2Ray system.
	ForSystem() System
}

// ManagerType returns the type of Manager interface. Can be used to implement common.HasType.
//
// v2ray:api:stable
func ManagerType() interface{} {
	return (*Manager)(nil)
}

var defaultBufferSize int32

func init() {
	const key = "v2ray.ray.buffer.size"
	const defaultValue = -17
	size := platform.EnvFlag{
		Name:    key,
		AltName: platform.NormalizeEnvName(key),
	}.GetValueAsInt(defaultValue)

	switch size {
	case 0:
		defaultBufferSize = -1 // For pipe to use unlimited size
	case defaultValue: // Env flag not defined. Use default values per CPU-arch.
		switch runtime.GOARCH {
		case "arm", "mips", "mipsle":
			defaultBufferSize = 0
		case "arm64", "mips64", "mips64le":
			defaultBufferSize = 4 * 1024 // 4k cache for low-end devices
		default:
			defaultBufferSize = 512 * 1024
		}
	default:
		defaultBufferSize = int32(size) * 1024 * 1024
	}
}

func defaultBufferPolicy() Buffer {
	return Buffer{
		PerConnection: defaultBufferSize,
	}
}

// SessionDefault returns the Policy when user is not specified.
func SessionDefault() Session {
	return Session{
		Timeouts: Timeout{
			Handshake:      time.Second * 4,
			ConnectionIdle: time.Second * 300,
			UplinkOnly:     time.Second * 1,
			DownlinkOnly:   time.Second * 1,
		},
		Stats: Stats{
			UserUplink:   false,
			UserDownlink: false,
		},
		Buffer: defaultBufferPolicy(),
	}
}

type policyKey int32

const (
	bufferPolicyKey policyKey = 0
)

func ContextWithBufferPolicy(ctx context.Context, p Buffer) context.Context {
	return context.WithValue(ctx, bufferPolicyKey, p)
}

func BufferPolicyFromContext(ctx context.Context) Buffer {
	pPolicy := ctx.Value(bufferPolicyKey)
	if pPolicy == nil {
		return defaultBufferPolicy()
	}
	return pPolicy.(Buffer)
}
