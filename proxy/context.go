package proxy

import (
	"context"

	"v2ray.com/core/common/net"
)

type key uint32

const (
	resolvedIPsKey key = iota
)

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
