package session

import "context"

type sessionKey int

const (
	idSessionKey sessionKey = iota
	inboundSessionKey
	outboundSessionKey
	contentSessionKey
)

// ContextWithID returns a new context with the given ID.
func ContextWithID(ctx context.Context, id ID) context.Context {
	return context.WithValue(ctx, idSessionKey, id)
}

// IDFromContext returns ID in this context, or 0 if not contained.
func IDFromContext(ctx context.Context) ID {
	if id, ok := ctx.Value(idSessionKey).(ID); ok {
		return id
	}
	return 0
}

func ContextWithInbound(ctx context.Context, inbound *Inbound) context.Context {
	return context.WithValue(ctx, inboundSessionKey, inbound)
}

func InboundFromContext(ctx context.Context) *Inbound {
	if inbound, ok := ctx.Value(inboundSessionKey).(*Inbound); ok {
		return inbound
	}
	return nil
}

func ContextWithOutbound(ctx context.Context, outbound *Outbound) context.Context {
	return context.WithValue(ctx, outboundSessionKey, outbound)
}

func OutboundFromContext(ctx context.Context) *Outbound {
	if outbound, ok := ctx.Value(outboundSessionKey).(*Outbound); ok {
		return outbound
	}
	return nil
}

func ContextWithContent(ctx context.Context, content *Content) context.Context {
	return context.WithValue(ctx, contentSessionKey, content)
}

func ContentFromContext(ctx context.Context) *Content {
	if content, ok := ctx.Value(contentSessionKey).(*Content); ok {
		return content
	}
	return nil
}
