package outbound

import (
	"context"

	"v2ray.com/core/common"
	"v2ray.com/core/common/vio"
	"v2ray.com/core/features"
)

// Handler is the interface for handlers that process outbound connections.
type Handler interface {
	common.Runnable
	Tag() string
	Dispatch(ctx context.Context, link *vio.Link)
}

// HandlerManager is a feature that manages outbound.Handlers.
type HandlerManager interface {
	features.Feature
	// GetHandler returns an outbound.Handler for the given tag.
	GetHandler(tag string) Handler
	// GetDefaultHandler returns the default outbound.Handler. It is usually the first outbound.Handler specified in the configuration.
	GetDefaultHandler() Handler
	// AddHandler adds a handler into this outbound.HandlerManager.
	AddHandler(ctx context.Context, handler Handler) error

	// RemoveHandler removes a handler from outbound.HandlerManager.
	RemoveHandler(ctx context.Context, tag string) error
}
