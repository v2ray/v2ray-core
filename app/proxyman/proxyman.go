package proxyman

import (
	"github.com/v2ray/v2ray-core/app"
	"github.com/v2ray/v2ray-core/proxy"
)

const (
	APP_ID_INBOUND_MANAGER = app.ID(4)
)

type InboundHandlerManager interface {
	GetHandler(tag string) (proxy.InboundHandler, int)
}

type inboundHandlerManagerWithContext interface {
	GetHandler(context app.Context, tag string) (proxy.InboundHandler, int)
}

type inboundHandlerManagerWithContextImpl struct {
	context app.Context
	manager inboundHandlerManagerWithContext
}

func (this *inboundHandlerManagerWithContextImpl) GetHandler(tag string) (proxy.InboundHandler, int) {
	return this.manager.GetHandler(this.context, tag)
}

func init() {
	app.Register(APP_ID_INBOUND_MANAGER, func(context app.Context, obj interface{}) interface{} {
		manager := obj.(inboundHandlerManagerWithContext)
		return &inboundHandlerManagerWithContextImpl{
			context: context,
			manager: manager,
		}
	})
}
