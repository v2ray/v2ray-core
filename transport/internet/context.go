package internet

import (
	"context"

	"v2ray.com/core/common/net"
)

type key int

const (
	streamSettingsKey key = iota
	dialerSrcKey
)

func ContextWithStreamSettings(ctx context.Context, streamSettings *MemoryStreamConfig) context.Context {
	return context.WithValue(ctx, streamSettingsKey, streamSettings)
}

func StreamSettingsFromContext(ctx context.Context) *MemoryStreamConfig {
	ss := ctx.Value(streamSettingsKey)
	if ss == nil {
		return nil
	}
	return ss.(*MemoryStreamConfig)
}

func ContextWithDialerSource(ctx context.Context, addr net.Address) context.Context {
	return context.WithValue(ctx, dialerSrcKey, addr)
}

func DialerSourceFromContext(ctx context.Context) net.Address {
	if addr, ok := ctx.Value(dialerSrcKey).(net.Address); ok {
		return addr
	}
	return net.AnyIP
}
