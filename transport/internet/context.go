package internet

import (
	"context"

	"v2ray.com/core/common/net"
)

type key int

const (
	streamSettingsKey key = iota
	dialerSrcKey
	transportSettingsKey
	securitySettingsKey
)

func ContextWithStreamSettings(ctx context.Context, streamSettings *StreamConfig) context.Context {
	return context.WithValue(ctx, streamSettingsKey, streamSettings)
}

func StreamSettingsFromContext(ctx context.Context) (*StreamConfig, bool) {
	ss, ok := ctx.Value(streamSettingsKey).(*StreamConfig)
	return ss, ok
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

func ContextWithTransportSettings(ctx context.Context, transportSettings interface{}) context.Context {
	return context.WithValue(ctx, transportSettingsKey, transportSettings)
}

func TransportSettingsFromContext(ctx context.Context) interface{} {
	return ctx.Value(transportSettingsKey)
}

func ContextWithSecuritySettings(ctx context.Context, securitySettings interface{}) context.Context {
	return context.WithValue(ctx, securitySettingsKey, securitySettings)
}

func SecuritySettingsFromContext(ctx context.Context) interface{} {
	return ctx.Value(securitySettingsKey)
}
