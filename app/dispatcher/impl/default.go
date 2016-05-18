package impl

import (
	"github.com/v2ray/v2ray-core/app"
	"github.com/v2ray/v2ray-core/app/dns"
	"github.com/v2ray/v2ray-core/app/proxyman"
	"github.com/v2ray/v2ray-core/app/router"
	"github.com/v2ray/v2ray-core/common/log"
	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/proxy"
	"github.com/v2ray/v2ray-core/transport/ray"
)

type DefaultDispatcher struct {
	ohm    proxyman.OutboundHandlerManager
	router router.Router
	dns    dns.Server
}

func NewDefaultDispatcher(space app.Space) *DefaultDispatcher {
	d := &DefaultDispatcher{}
	space.InitializeApplication(func() error {
		return d.Initialize(space)
	})
	return d
}

// @Private
func (this *DefaultDispatcher) Initialize(space app.Space) error {
	if !space.HasApp(proxyman.APP_ID_OUTBOUND_MANAGER) {
		log.Error("DefaultDispatcher: OutboundHandlerManager is not found in the space.")
		return app.ErrorMissingApplication
	}
	this.ohm = space.GetApp(proxyman.APP_ID_OUTBOUND_MANAGER).(proxyman.OutboundHandlerManager)

	if !space.HasApp(dns.APP_ID) {
		log.Error("DefaultDispatcher: DNS is not found in the space.")
		return app.ErrorMissingApplication
	}
	this.dns = space.GetApp(dns.APP_ID).(dns.Server)

	if space.HasApp(router.APP_ID) {
		this.router = space.GetApp(router.APP_ID).(router.Router)
	}

	return nil
}

func (this *DefaultDispatcher) Release() {

}

func (this *DefaultDispatcher) DispatchToOutbound(destination v2net.Destination) ray.InboundRay {
	direct := ray.NewRay()
	dispatcher := this.ohm.GetDefaultHandler()

	if this.router != nil {
		if tag, err := this.router.TakeDetour(destination); err == nil {
			if handler := this.ohm.GetHandler(tag); handler != nil {
				log.Info("DefaultDispatcher: Taking detour [", tag, "] for [", destination, "].")
				dispatcher = handler
			} else {
				log.Warning("DefaultDispatcher: Nonexisting tag: ", tag)
			}
		} else {
			log.Info("DefaultDispatcher: Default route for ", destination)
		}
	}

	go this.FilterPacketAndDispatch(destination, direct, dispatcher)

	return direct
}

// @Private
func (this *DefaultDispatcher) FilterPacketAndDispatch(destination v2net.Destination, link ray.OutboundRay, dispatcher proxy.OutboundHandler) {
	payload, err := link.OutboundInput().Read()
	if err != nil {
		log.Info("DefaultDispatcher: No payload to dispatch, stopping now.")
		link.OutboundInput().Release()
		link.OutboundOutput().Release()
		return
	}
	dispatcher.Dispatch(destination, payload, link)
}
