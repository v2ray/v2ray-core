package routing

import (
	"context"

	"v2ray.com/core/common"
	"v2ray.com/core/features"
)

// Router is a feature to choose an outbound tag for the given request.
//
// v2ray:api:stable
type Router interface {
	features.Feature

	// PickRoute returns a tag of an OutboundHandler based on the given context.
	PickRoute(ctx context.Context) (string, error)
}

// RouterType return the type of Router interface. Can be used to implement common.HasType.
//
// v2ray:api:stable
func RouterType() interface{} {
	return (*Router)(nil)
}

// DefaultRouter is an implementation of Router, which always returns ErrNoClue for routing decisions.
type DefaultRouter struct{}

// Type implements common.HasType.
func (DefaultRouter) Type() interface{} {
	return RouterType()
}

// PickRoute implements Router.
func (DefaultRouter) PickRoute(ctx context.Context) (string, error) {
	return "", common.ErrNoClue
}

// Start implements common.Runnable.
func (DefaultRouter) Start() error {
	return nil
}

// Close implements common.Closable.
func (DefaultRouter) Close() error {
	return nil
}
