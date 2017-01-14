package outbound

import (
	"context"
	"sync"

	"v2ray.com/core/app/proxyman"
	"v2ray.com/core/common"
	"v2ray.com/core/proxy"
)

type DefaultOutboundHandlerManager struct {
	sync.RWMutex
	defaultHandler proxy.OutboundHandler
	taggedHandler  map[string]proxy.OutboundHandler
}

func New(ctx context.Context, config *proxyman.OutboundConfig) (*DefaultOutboundHandlerManager, error) {
	return &DefaultOutboundHandlerManager{
		taggedHandler: make(map[string]proxy.OutboundHandler),
	}, nil
}

func (DefaultOutboundHandlerManager) Interface() interface{} {
	return (*proxyman.OutboundHandlerManager)(nil)
}

func (v *DefaultOutboundHandlerManager) GetDefaultHandler() proxy.OutboundHandler {
	v.RLock()
	defer v.RUnlock()
	if v.defaultHandler == nil {
		return nil
	}
	return v.defaultHandler
}

func (v *DefaultOutboundHandlerManager) SetDefaultHandler(handler proxy.OutboundHandler) error {
	v.Lock()
	defer v.Unlock()
	v.defaultHandler = handler
	return nil
}

func (v *DefaultOutboundHandlerManager) GetHandler(tag string) proxy.OutboundHandler {
	v.RLock()
	defer v.RUnlock()
	if handler, found := v.taggedHandler[tag]; found {
		return handler
	}
	return nil
}

func (v *DefaultOutboundHandlerManager) SetHandler(tag string, handler proxy.OutboundHandler) error {
	v.Lock()
	defer v.Unlock()

	v.taggedHandler[tag] = handler
	return nil
}

func init() {
	common.Must(common.RegisterConfig((*proxyman.OutboundConfig)(nil), func(ctx context.Context, config interface{}) (interface{}, error) {
		return New(ctx, config.(*proxyman.OutboundConfig))
	}))
}
