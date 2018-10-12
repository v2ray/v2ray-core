package routing

import (
	"context"

	"v2ray.com/core/features"
)

// Router is a feature to choose an outbound tag for the given request.
type Router interface {
	features.Feature

	// PickRoute returns a tag of an OutboundHandler based on the given context.
	PickRoute(ctx context.Context) (string, error)
}
