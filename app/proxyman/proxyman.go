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

func (this *DefaultOutboundHandlerManager) Release() {

}

func (this *DefaultOutboundHandlerManager) GetDefaultHandler() proxy.OutboundHandler {
	this.RLock()
	defer this.RUnlock()
	if this.defaultHandler == nil {
		return nil
	}
	return this.defaultHandler
}

func (this *DefaultOutboundHandlerManager) SetDefaultHandler(handler proxy.OutboundHandler) {
	this.Lock()
	defer this.Unlock()
	this.defaultHandler = handler
}

func (this *DefaultOutboundHandlerManager) GetHandler(tag string) proxy.OutboundHandler {
	this.RLock()
	defer this.RUnlock()
	if handler, found := this.taggedHandler[tag]; found {
		return handler
	}
	return nil
}

func (this *DefaultOutboundHandlerManager) SetHandler(tag string, handler proxy.OutboundHandler) {
	this.Lock()
	defer this.Unlock()

	this.taggedHandler[tag] = handler
}
