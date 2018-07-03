package protocol

import (
	"context"
)

type key int

const (
	userKey key = iota
	requestKey
)

// ContextWithUser returns a context combined with a User.
func ContextWithUser(ctx context.Context, user *User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

// UserFromContext extracts a User from the given context, if any.
func UserFromContext(ctx context.Context) *User {
	v := ctx.Value(userKey)
	if v == nil {
		return nil
	}
	return v.(*User)
}

func ContextWithRequestHeader(ctx context.Context, request *RequestHeader) context.Context {
	return context.WithValue(ctx, requestKey, request)
}

func RequestHeaderFromContext(ctx context.Context) *RequestHeader {
	request := ctx.Value(requestKey)
	if request == nil {
		return nil
	}
	return request.(*RequestHeader)
}
