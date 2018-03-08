package protocol

import (
	"context"
)

type key int

const (
	userKey key = iota
)

// ContextWithUser returns a context combined with an User.
func ContextWithUser(ctx context.Context, user *User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

// UserFromContext extracts an User from the given context, if any.
func UserFromContext(ctx context.Context) *User {
	v := ctx.Value(userKey)
	if v == nil {
		return nil
	}
	return v.(*User)
}
