package proxy

import (
	"context"

	"v2ray.com/core/common"
)

func CreateInboundHandler(ctx context.Context, config interface{}) (Inbound, error) {
	handler, err := common.CreateObject(ctx, config)
	if err != nil {
		return nil, err
	}
	switch h := handler.(type) {
	case Inbound:
		return h, nil
	default:
		return nil, newError("not a InboundHandler")
	}
}

func CreateOutboundHandler(ctx context.Context, config interface{}) (Outbound, error) {
	handler, err := common.CreateObject(ctx, config)
	if err != nil {
		return nil, err
	}
	switch h := handler.(type) {
	case Outbound:
		return h, nil
	default:
		return nil, newError("not a OutboundHandler")
	}
}
