package routing

import (
	"v2ray.com/core/common"
	"v2ray.com/core/features"
)

// Router is a feature to choose an outbound tag for the given request.
//
// v2ray:api:stable
type Router interface {
	features.Feature

	// PickRoute returns a route decision based on the given routing context.
	PickRoute(ctx Context) (Route, error)
}

// Route is the routing result of Router feature.
//
// v2ray:api:stable
type Route interface {
	// A Route is also a routing context.
	Context

	// GetOutboundGroupTags returns the detoured outbound group tags in sequence before a final outbound is chosen.
	GetOutboundGroupTags() []string

	// GetOutboundTag returns the tag of the outbound the connection was dispatched to.
	GetOutboundTag() string
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
func (DefaultRouter) PickRoute(ctx Context) (Route, error) {
	return nil, common.ErrNoClue
}

// Start implements common.Runnable.
func (DefaultRouter) Start() error {
	return nil
}

// Close implements common.Closable.
func (DefaultRouter) Close() error {
	return nil
}
