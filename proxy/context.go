package proxy

import (
	"context"

	"v2ray.com/core/common/net"
)

type key int

const (
	sourceKey key = iota
	targetKey
	originalTargetKey
	inboundEntryPointKey
	inboundTagKey
	resolvedIPsKey
)

// ContextWithSource creates a new context with given source.
func ContextWithSource(ctx context.Context, src net.Destination) context.Context {
	return context.WithValue(ctx, sourceKey, src)
}

// SourceFromContext retrieves source from the given context.
func SourceFromContext(ctx context.Context) (net.Destination, bool) {
	v, ok := ctx.Value(sourceKey).(net.Destination)
	return v, ok
}

func ContextWithOriginalTarget(ctx context.Context, dest net.Destination) context.Context {
	return context.WithValue(ctx, originalTargetKey, dest)
}

func OriginalTargetFromContext(ctx context.Context) (net.Destination, bool) {
	v, ok := ctx.Value(originalTargetKey).(net.Destination)
	return v, ok
}

func ContextWithTarget(ctx context.Context, dest net.Destination) context.Context {
	return context.WithValue(ctx, targetKey, dest)
}

func TargetFromContext(ctx context.Context) (net.Destination, bool) {
	v, ok := ctx.Value(targetKey).(net.Destination)
	return v, ok
}

func ContextWithInboundEntryPoint(ctx context.Context, dest net.Destination) context.Context {
	return context.WithValue(ctx, inboundEntryPointKey, dest)
}

func InboundEntryPointFromContext(ctx context.Context) (net.Destination, bool) {
	v, ok := ctx.Value(inboundEntryPointKey).(net.Destination)
	return v, ok
}

func ContextWithInboundTag(ctx context.Context, tag string) context.Context {
	return context.WithValue(ctx, inboundTagKey, tag)
}

func InboundTagFromContext(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(inboundTagKey).(string)
	return v, ok
}

type IPResolver interface {
	Resolve() []net.Address
}

func ContextWithResolveIPs(ctx context.Context, f IPResolver) context.Context {
	return context.WithValue(ctx, resolvedIPsKey, f)
}

func ResolvedIPsFromContext(ctx context.Context) (IPResolver, bool) {
	ips, ok := ctx.Value(resolvedIPsKey).(IPResolver)
	return ips, ok
}
