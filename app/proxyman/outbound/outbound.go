package outbound

import (
	"sync"

	"v2ray.com/core/app"
	"v2ray.com/core/app/proxyman"
	"v2ray.com/core/common"
	"v2ray.com/core/proxy"
)

type DefaultOutboundHandlerManager struct {
	sync.RWMutex
	defaultHandler proxy.OutboundHandler
	taggedHandler  map[string]proxy.OutboundHandler
}

func New() *DefaultOutboundHandlerManager {
	return &DefaultOutboundHandlerManager{
		taggedHandler: make(map[string]proxy.OutboundHandler),
	}
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

type OutboundHandlerManagerFactory struct{}

func (v OutboundHandlerManagerFactory) Create(space app.Space, config interface{}) (app.Application, error) {
	return New(), nil
}

func init() {
	common.Must(app.RegisterApplicationFactory((*proxyman.OutboundConfig)(nil), OutboundHandlerManagerFactory{}))
}
