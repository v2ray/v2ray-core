package outbound

import (
	"sync"
	"v2ray.com/core/app"
	"v2ray.com/core/app/proxyman"
	"v2ray.com/core/common/serial"
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

func (v OutboundHandlerManagerFactory) AppId() app.ID {
	return proxyman.APP_ID_OUTBOUND_MANAGER
}

func init() {
	app.RegisterApplicationFactory(serial.GetMessageType(new(proxyman.OutboundConfig)), OutboundHandlerManagerFactory{})
}
