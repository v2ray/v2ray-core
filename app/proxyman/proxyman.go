// Package proxyman defines applications for managing inbound and outbound proxies.
package proxyman

import (
	"context"
)

type key int

const (
	sniffing key = iota
)

func ContextWithSniffingConfig(ctx context.Context, c *SniffingConfig) context.Context {
	return context.WithValue(ctx, sniffing, c)
}

func SniffingConfigFromContext(ctx context.Context) *SniffingConfig {
	if c, ok := ctx.Value(sniffing).(*SniffingConfig); ok {
		return c
	}
	return nil
}
