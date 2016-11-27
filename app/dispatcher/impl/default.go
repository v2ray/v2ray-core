package impl

import (
	"errors"
	"v2ray.com/core/app"
	"v2ray.com/core/app/proxyman"
	"v2ray.com/core/app/router"
	"v2ray.com/core/common/alloc"
	"v2ray.com/core/common/log"
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/proxy"
	"v2ray.com/core/transport/ray"
)

type DefaultDispatcher struct {
	ohm    proxyman.OutboundHandlerManager
	router *router.Router
}

func NewDefaultDispatcher(space app.Space) *DefaultDispatcher {
	d := &DefaultDispatcher{}
	space.InitializeApplication(func() error {
		return d.Initialize(space)
	})
	return d
}

// Private: Used by app.Space only.
func (v *DefaultDispatcher) Initialize(space app.Space) error {
	if !space.HasApp(proxyman.APP_ID_OUTBOUND_MANAGER) {
		return errors.New("DefaultDispatcher: OutboundHandlerManager is not found in the space.")
	}
	v.ohm = space.GetApp(proxyman.APP_ID_OUTBOUND_MANAGER).(proxyman.OutboundHandlerManager)

	if space.HasApp(router.APP_ID) {
		v.router = space.GetApp(router.APP_ID).(*router.Router)
	}

	return nil
}

func (v *DefaultDispatcher) Release() {

}

func (v *DefaultDispatcher) DispatchToOutbound(session *proxy.SessionInfo) ray.InboundRay {
	direct := ray.NewRay()
	dispatcher := v.ohm.GetDefaultHandler()
	destination := session.Destination

	if v.router != nil {
		if tag, err := v.router.TakeDetour(session); err == nil {
			if handler := v.ohm.GetHandler(tag); handler != nil {
				log.Info("DefaultDispatcher: Taking detour [", tag, "] for [", destination, "].")
				dispatcher = handler
			} else {
				log.Warning("DefaultDispatcher: Nonexisting tag: ", tag)
			}
		} else {
			log.Info("DefaultDispatcher: Default route for ", destination)
		}
	}

	if session.Inbound != nil && session.Inbound.AllowPassiveConnection {
		go dispatcher.Dispatch(destination, alloc.NewLocalBuffer(32).Clear(), direct)
	} else {
		go v.FilterPacketAndDispatch(destination, direct, dispatcher)
	}

	return direct
}

// Private: Visible for testing.
func (v *DefaultDispatcher) FilterPacketAndDispatch(destination v2net.Destination, link ray.OutboundRay, dispatcher proxy.OutboundHandler) {
	payload, err := link.OutboundInput().Read()
	if err != nil {
		log.Info("DefaultDispatcher: No payload towards ", destination, ", stopping now.")
		link.OutboundInput().Release()
		link.OutboundOutput().Release()
		return
	}
	dispatcher.Dispatch(destination, payload, link)
}
