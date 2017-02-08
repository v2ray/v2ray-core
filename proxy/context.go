package proxy

import (
	"context"

	"v2ray.com/core/common/net"
)

type key int

const (
	sourceKey key = iota
	destinationKey
	originalDestinationKey
	inboundDestinationKey
	inboundTagKey
	outboundTagKey
	resolvedIPsKey
	allowPassiveConnKey
	dispatcherKey
)

func ContextWithSource(ctx context.Context, src net.Destination) context.Context {
	return context.WithValue(ctx, sourceKey, src)
}

func SourceFromContext(ctx context.Context) net.Destination {
	v := ctx.Value(sourceKey)
	if v == nil {
		return net.Destination{}
	}
	return v.(net.Destination)
}

func ContextWithOriginalDestination(ctx context.Context, dest net.Destination) context.Context {
	return context.WithValue(ctx, originalDestinationKey, dest)
}

func OriginalDestinationFromContext(ctx context.Context) net.Destination {
	v := ctx.Value(originalDestinationKey)
	if v == nil {
		return net.Destination{}
	}
	return v.(net.Destination)
}

func ContextWithDestination(ctx context.Context, dest net.Destination) context.Context {
	return context.WithValue(ctx, destinationKey, dest)
}

func DestinationFromContext(ctx context.Context) net.Destination {
	v := ctx.Value(destinationKey)
	if v == nil {
		return net.Destination{}
	}
	return v.(net.Destination)
}

func ContextWithInboundDestination(ctx context.Context, dest net.Destination) context.Context {
	return context.WithValue(ctx, inboundDestinationKey, dest)
}

func InboundDestinationFromContext(ctx context.Context) net.Destination {
	v := ctx.Value(inboundDestinationKey)
	if v == nil {
		return net.Destination{}
	}
	return v.(net.Destination)
}

func ContextWithInboundTag(ctx context.Context, tag string) context.Context {
	return context.WithValue(ctx, inboundTagKey, tag)
}

func InboundTagFromContext(ctx context.Context) string {
	v := ctx.Value(inboundTagKey)
	if v == nil {
		return ""
	}
	return v.(string)
}

func ContextWithOutboundTag(ctx context.Context, tag string) context.Context {
	return context.WithValue(ctx, outboundTagKey, tag)
}

func OutboundTagFromContext(ctx context.Context) string {
	v := ctx.Value(outboundTagKey)
	if v == nil {
		return ""
	}
	return v.(string)
}

func ContextWithResolveIPs(ctx context.Context, ips []net.Address) context.Context {
	return context.WithValue(ctx, resolvedIPsKey, ips)
}

func ResolvedIPsFromContext(ctx context.Context) ([]net.Address, bool) {
	ips, ok := ctx.Value(resolvedIPsKey).([]net.Address)
	return ips, ok
}
