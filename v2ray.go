package core

import (
	"v2ray.com/core/app"
	"v2ray.com/core/app/dispatcher"
	dispatchers "v2ray.com/core/app/dispatcher/impl"
	"v2ray.com/core/app/dns"
	"v2ray.com/core/app/proxyman"
	"v2ray.com/core/app/router"
	"v2ray.com/core/common"
	"v2ray.com/core/common/log"
	"v2ray.com/core/proxy"
	proxyregistry "v2ray.com/core/proxy/registry"
)

// Point shell of V2Ray.
type Point struct {
	inboundHandlers       []InboundDetourHandler
	taggedInboundHandlers map[string]InboundDetourHandler

	outboundHandlers       []proxy.OutboundHandler
	taggedOutboundHandlers map[string]proxy.OutboundHandler

	router *router.Router
	space  app.Space
}

// NewPoint returns a new Point server based on given configuration.
// The server is not started at this point.
func NewPoint(pConfig *Config) (*Point, error) {
	var vpoint = new(Point)

	if err := pConfig.Transport.Apply(); err != nil {
		return nil, err
	}

	if err := pConfig.Log.Apply(); err != nil {
		return nil, err
	}

	vpoint.space = app.NewSpace()
	vpoint.space.BindApp(proxyman.APP_ID_INBOUND_MANAGER, vpoint)

	outboundHandlerManager := proxyman.NewDefaultOutboundHandlerManager()
	vpoint.space.BindApp(proxyman.APP_ID_OUTBOUND_MANAGER, outboundHandlerManager)

	dnsConfig := pConfig.Dns
	if dnsConfig != nil {
		dnsServer := dns.NewCacheServer(vpoint.space, dnsConfig)
		vpoint.space.BindApp(dns.APP_ID, dnsServer)
	}

	routerConfig := pConfig.Router
	if routerConfig != nil {
		r := router.NewRouter(routerConfig, vpoint.space)
		vpoint.space.BindApp(router.APP_ID, r)
		vpoint.router = r
	}

	vpoint.space.BindApp(dispatcher.APP_ID, dispatchers.NewDefaultDispatcher(vpoint.space))

	vpoint.inboundHandlers = make([]InboundDetourHandler, 8)
	vpoint.taggedInboundHandlers = make(map[string]InboundDetourHandler)
	for _, inbound := range pConfig.Inbound {
		allocConfig := inbound.GetAllocationStrategyValue()
		var inboundHandler InboundDetourHandler
		switch allocConfig.Type {
		case AllocationStrategy_Always:
			dh, err := NewInboundDetourHandlerAlways(vpoint.space, inbound)
			if err != nil {
				log.Error("Point: Failed to create detour handler: ", err)
				return nil, common.ErrBadConfiguration
			}
			inboundHandler = dh
		case AllocationStrategy_Random:
			dh, err := NewInboundDetourHandlerDynamic(vpoint.space, inbound)
			if err != nil {
				log.Error("Point: Failed to create detour handler: ", err)
				return nil, common.ErrBadConfiguration
			}
			inboundHandler = dh
		default:
			log.Error("Point: Unknown allocation strategy: ", allocConfig.Type)
			return nil, common.ErrBadConfiguration
		}
		vpoint.inboundHandlers = append(vpoint.inboundHandlers, inboundHandler)
		if len(inbound.Tag) > 0 {
			vpoint.taggedInboundHandlers[inbound.Tag] = inboundHandler
		}
	}

	vpoint.outboundHandlers = make([]proxy.OutboundHandler, 8)
	vpoint.taggedOutboundHandlers = make(map[string]proxy.OutboundHandler)
	for idx, outbound := range pConfig.Outbound {
		outboundHandler, err := proxyregistry.CreateOutboundHandler(
			outbound.Protocol, vpoint.space, outbound.Settings, &proxy.OutboundHandlerMeta{
				Tag:            outbound.Tag,
				Address:        outbound.SendThrough.AsAddress(),
				StreamSettings: outbound.StreamSettings,
			})
		if err != nil {
			log.Error("Point: Failed to create detour outbound connection handler: ", err)
			return nil, err
		}
		if idx == 0 {
			outboundHandlerManager.SetDefaultHandler(outboundHandler)
		}
		if len(outbound.Tag) > 0 {
			outboundHandlerManager.SetHandler(outbound.Tag, outboundHandler)
		}
	}

	if err := vpoint.space.Initialize(); err != nil {
		return nil, err
	}

	return vpoint, nil
}

func (this *Point) Close() {
	for _, inbound := range this.inboundHandlers {
		inbound.Close()
	}
}

// Start starts the Point server, and return any error during the process.
// In the case of any errors, the state of the server is unpredicatable.
func (this *Point) Start() error {
	for _, inbound := range this.inboundHandlers {
		err := inbound.Start()
		if err != nil {
			return err
		}
	}
	log.Warning("V2Ray started.")

	return nil
}

func (this *Point) GetHandler(tag string) (proxy.InboundHandler, int) {
	handler, found := this.taggedInboundHandlers[tag]
	if !found {
		log.Warning("Point: Unable to find an inbound handler with tag: ", tag)
		return nil, 0
	}
	return handler.GetConnectionHandler()
}

func (this *Point) Release() {

}
