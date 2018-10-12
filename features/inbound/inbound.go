package inbound

import (
	"context"

	"v2ray.com/core/common"
	"v2ray.com/core/common/net"
	"v2ray.com/core/features"
)

// Handler is the interface for handlers that process inbound connections.
type Handler interface {
	common.Runnable
	// The tag of this handler.
	Tag() string

	// Deprecated. Do not use in new code.
	GetRandomInboundProxy() (interface{}, net.Port, int)
}

// Manager is a feature that manages InboundHandlers.
type Manager interface {
	features.Feature
	// GetHandlers returns an InboundHandler for the given tag.
	GetHandler(ctx context.Context, tag string) (Handler, error)
	// AddHandler adds the given handler into this Manager.
	AddHandler(ctx context.Context, handler Handler) error

	// RemoveHandler removes a handler from Manager.
	RemoveHandler(ctx context.Context, tag string) error
}

func ManagerType() interface{} {
	return (*Manager)(nil)
}
