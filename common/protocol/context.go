package protocol

import (
	"context"
)

type key int

const (
	userKey key = iota
)

func ContextWithUser(ctx context.Context, user *User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

func UserFromContext(ctx context.Context) *User {
	v := ctx.Value(userKey)
	if v == nil {
		return nil
	}
	return v.(*User)
}
