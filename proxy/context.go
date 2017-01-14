package proxy

import (
	"context"
)

type key int

const (
	inboundMetaKey key = iota
	outboundMetaKey
	dialerKey
)

func ContextWithInboundMeta(ctx context.Context, meta *InboundHandlerMeta) context.Context {
	return context.WithValue(ctx, inboundMetaKey, meta)
}

func InboundMetaFromContext(ctx context.Context) *InboundHandlerMeta {
	v := ctx.Value(inboundMetaKey)
	if v == nil {
		return nil
	}
	return v.(*InboundHandlerMeta)
}

func ContextWithOutboundMeta(ctx context.Context, meta *OutboundHandlerMeta) context.Context {
	return context.WithValue(ctx, outboundMetaKey, meta)
}

func OutboundMetaFromContext(ctx context.Context) *OutboundHandlerMeta {
	v := ctx.Value(outboundMetaKey)
	if v == nil {
		return nil
	}
	return v.(*OutboundHandlerMeta)
}

func ContextWithDialer(ctx context.Context, dialer Dialer) context.Context {
	return context.WithValue(ctx, dialerKey, dialer)
}

func DialerFromContext(ctx context.Context) Dialer {
	v := ctx.Value(dialerKey)
	if v == nil {
		return nil
	}
	return v.(Dialer)
}
