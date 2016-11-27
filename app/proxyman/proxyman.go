package proxyman

import (
	"sync"

	"v2ray.com/core/app"
	"v2ray.com/core/proxy"
)

const (
	APP_ID_INBOUND_MANAGER  = app.ID(4)
	APP_ID_OUTBOUND_MANAGER = app.ID(6)
)

type InboundHandlerManager interface {
	GetHandler(tag string) (proxy.InboundHandler, int)
}

type OutboundHandlerManager interface {
	GetHandler(tag string) proxy.OutboundHandler
	GetDefaultHandler() proxy.OutboundHandler
}

type DefaultOutboundHandlerManager struct {
	sync.RWMutex
	defaultHandler proxy.OutboundHandler
	taggedHandler  map[string]proxy.OutboundHandler
}

func NewDefaultOutboundHandlerManager() *DefaultOutboundHandlerManager {
	return &DefaultOutboundHandlerManager{
		taggedHandler: make(map[string]proxy.OutboundHandler),
	}
}

func (v *DefaultOutboundHandlerManager) Release() {

}

func (v *DefaultOutboundHandlerManager) GetDefaultHandler() proxy.OutboundHandler {
	v.RLock()
	defer v.RUnlock()
	if v.defaultHandler == nil {
		return nil
	}
	return v.defaultHandler
}

func (v *DefaultOutboundHandlerManager) SetDefaultHandler(handler proxy.OutboundHandler) {
	v.Lock()
	defer v.Unlock()
	v.defaultHandler = handler
}

func (v *DefaultOutboundHandlerManager) GetHandler(tag string) proxy.OutboundHandler {
	v.RLock()
	defer v.RUnlock()
	if handler, found := v.taggedHandler[tag]; found {
		return handler
	}
	return nil
}

func (v *DefaultOutboundHandlerManager) SetHandler(tag string, handler proxy.OutboundHandler) {
	v.Lock()
	defer v.Unlock()

	v.taggedHandler[tag] = handler
}
