package outbound

//go:generate go run $GOPATH/src/v2ray.com/core/tools/generrorgen/main.go -pkg outbound -path App,Proxyman,Outbound

import (
	"context"
	"sync"

	"v2ray.com/core/app/proxyman"
	"v2ray.com/core/common"
)

type DefaultOutboundHandlerManager struct {
	sync.RWMutex
	defaultHandler *Handler
	taggedHandler  map[string]*Handler
}

func New(ctx context.Context, config *proxyman.OutboundConfig) (*DefaultOutboundHandlerManager, error) {
	return &DefaultOutboundHandlerManager{
		taggedHandler: make(map[string]*Handler),
	}, nil
}

func (*DefaultOutboundHandlerManager) Interface() interface{} {
	return (*proxyman.OutboundHandlerManager)(nil)
}

func (*DefaultOutboundHandlerManager) Start() error { return nil }

func (*DefaultOutboundHandlerManager) Close() {}

func (v *DefaultOutboundHandlerManager) GetDefaultHandler() proxyman.OutboundHandler {
	v.RLock()
	defer v.RUnlock()
	if v.defaultHandler == nil {
		return nil
	}
	return v.defaultHandler
}

func (v *DefaultOutboundHandlerManager) GetHandler(tag string) proxyman.OutboundHandler {
	v.RLock()
	defer v.RUnlock()
	if handler, found := v.taggedHandler[tag]; found {
		return handler
	}
	return nil
}

func (v *DefaultOutboundHandlerManager) AddHandler(ctx context.Context, config *proxyman.OutboundHandlerConfig) error {
	v.Lock()
	defer v.Unlock()

	handler, err := NewHandler(ctx, config)
	if err != nil {
		return err
	}
	if v.defaultHandler == nil {
		v.defaultHandler = handler
	}

	if len(config.Tag) > 0 {
		v.taggedHandler[config.Tag] = handler
	}

	return nil
}

func init() {
	common.Must(common.RegisterConfig((*proxyman.OutboundConfig)(nil), func(ctx context.Context, config interface{}) (interface{}, error) {
		return New(ctx, config.(*proxyman.OutboundConfig))
	}))
}
