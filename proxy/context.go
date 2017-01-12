package proxy

import (
	"context"
)

type key int

const (
	inboundMetaKey  = key(0)
	outboundMetaKey = key(1)
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
