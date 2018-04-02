package core

import (
	"context"
)

type key int

const v2rayKey key = 1

// FromContext returns an Instance from the given context, or nil if the context doesn't contain one.
func FromContext(ctx context.Context) *Instance {
	if s, ok := ctx.Value(v2rayKey).(*Instance); ok {
		return s
	}
	return nil
}

// MustFromContext returns an Instance from the given context, or panics if not present.
func MustFromContext(ctx context.Context) *Instance {
	v := FromContext(ctx)
	if v == nil {
		panic("V is not in context.")
	}
	return v
}
