// Package proxyman defines applications for managing inbound and outbound proxies.
package proxyman

import (
	"context"
)

type key int

const (
	protocolsKey key = iota
)

func ContextWithProtocolSniffers(ctx context.Context, list []KnownProtocols) context.Context {
	return context.WithValue(ctx, protocolsKey, list)
}

func ProtocolSniffersFromContext(ctx context.Context) []KnownProtocols {
	if list, ok := ctx.Value(protocolsKey).([]KnownProtocols); ok {
		return list
	}
	return nil
}
