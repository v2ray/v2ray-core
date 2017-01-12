package proxy

import (
	"context"

	"v2ray.com/core/common"
	"v2ray.com/core/common/errors"
)

func CreateInboundHandler(ctx context.Context, config interface{}) (InboundHandler, error) {
	handler, err := common.CreateObject(ctx, config)
	if err != nil {
		return nil, err
	}
	switch h := handler.(type) {
	case InboundHandler:
		return h, nil
	default:
		return nil, errors.New("Proxy: Not a InboundHandler.")
	}
}

func CreateOutboundHandler(ctx context.Context, config interface{}) (OutboundHandler, error) {
	handler, err := common.CreateObject(ctx, config)
	if err != nil {
		return nil, err
	}
	switch h := handler.(type) {
	case OutboundHandler:
		return h, nil
	default:
		return nil, errors.New("Proxy: Not a OutboundHandler.")
	}
}
