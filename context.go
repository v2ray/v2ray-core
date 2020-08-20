// +build !confonly

package core

import (
	"context"
)

// V2rayKey is the key type of Instance in Context, exported for test.
type V2rayKey int

const v2rayKey V2rayKey = 1

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
