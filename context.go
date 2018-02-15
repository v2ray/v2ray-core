package core

import (
	"context"
)

type key int

const v2rayKey key = 1

// FromContext returns a Instance from the given context, or nil if the context doesn't contain one.
func FromContext(ctx context.Context) *Instance {
	if s, ok := ctx.Value(v2rayKey).(*Instance); ok {
		return s
	}
	return nil
}
