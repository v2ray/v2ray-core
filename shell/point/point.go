// Package point is a shell of V2Ray to run on various of systems.
// Point server is a full functionality proxying system. It consists of an inbound and an outbound
// connection, as well as any number of inbound and outbound detours. It provides a way internally
// to route network packets.
package point

import (
	"github.com/v2ray/v2ray-core/app"
	"github.com/v2ray/v2ray-core/app/dispatcher"
	dispatchers "github.com/v2ray/v2ray-core/app/dispatcher/impl"
	"github.com/v2ray/v2ray-core/app/dns"
	"github.com/v2ray/v2ray-core/app/proxyman"
	"github.com/v2ray/v2ray-core/app/router"
	"github.com/v2ray/v2ray-core/common/log"
	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/common/retry"
	"github.com/v2ray/v2ray-core/proxy"
	proxyrepo "github.com/v2ray/v2ray-core/proxy/repo"
)

// Point shell of V2Ray.
type Point struct {
	port      v2net.Port
	listen    v2net.Address
	ich       proxy.InboundHandler
	och       proxy.OutboundHandler
	idh       []InboundDetourHandler
	taggedIdh map[string]InboundDetourHandler
	odh       map[string]proxy.OutboundHandler
	router    router.Router
	space     app.Space
}

// NewPoint returns a new Point server based on given configuration.
// The server is not started at this point.
func NewPoint(pConfig *Config) (*Point, error) {
	var vpoint = new(Point)
	vpoint.port = pConfig.InboundConfig.Port
	if vpoint.port == 0 {
		vpoint.port = pConfig.Port // Backward compatibility
	}

	vpoint.listen = pConfig.InboundConfig.ListenOn

	if pConfig.TransportConfig != nil {
		pConfig.TransportConfig.Apply()
	}

	if pConfig.LogConfig != nil {
		logConfig := pConfig.LogConfig
		if len(logConfig.AccessLog) > 0 {
			err := log.InitAccessLogger(logConfig.AccessLog)
			if err != nil {
				return nil, err
			}
		}

		if len(logConfig.ErrorLog) > 0 {
			err := log.InitErrorLogger(logConfig.ErrorLog)
			if err != nil {
				return nil, err
			}
		}

		log.SetLogLevel(logConfig.LogLevel)
	}

	vpoint.space = app.NewSpace()
	vpoint.space.BindApp(proxyman.APP_ID_INBOUND_MANAGER, vpoint)

	outboundHandlerManager := proxyman.NewDefaultOutboundHandlerManager()
	vpoint.space.BindApp(proxyman.APP_ID_OUTBOUND_MANAGER, outboundHandlerManager)

	dnsConfig := pConfig.DNSConfig
	if dnsConfig != nil {
		dnsServer := dns.NewCacheServer(vpoint.space, dnsConfig)
		vpoint.space.BindApp(dns.APP_ID, dnsServer)
	}

	routerConfig := pConfig.RouterConfig
	if routerConfig != nil {
		r, err := router.CreateRouter(routerConfig.Strategy, routerConfig.Settings, vpoint.space)
		if err != nil {
			log.Error("Failed to create router: ", err)
			return nil, ErrorBadConfiguration
		}
		vpoint.space.BindApp(router.APP_ID, r)
		vpoint.router = r
	}

	vpoint.space.BindApp(dispatcher.APP_ID, dispatchers.NewDefaultDispatcher(vpoint.space))

	ichConfig := pConfig.InboundConfig.Settings
	ich, err := proxyrepo.CreateInboundHandler(
		pConfig.InboundConfig.Protocol, vpoint.space, ichConfig, &proxy.InboundHandlerMeta{
			Tag:            "system.inbound",
			Address:        pConfig.InboundConfig.ListenOn,
			Port:           vpoint.port,
			StreamSettings: pConfig.InboundConfig.StreamSettings})
	if err != nil {
		log.Error("Failed to create inbound connection handler: ", err)
		return nil, err
	}
	vpoint.ich = ich

	ochConfig := pConfig.OutboundConfig.Settings
	och, err := proxyrepo.CreateOutboundHandler(
		pConfig.OutboundConfig.Protocol, vpoint.space, ochConfig, &proxy.OutboundHandlerMeta{
			Tag:            "system.outbound",
			Address:        pConfig.OutboundConfig.SendThrough,
			StreamSettings: pConfig.OutboundConfig.StreamSettings,
		})
	if err != nil {
		log.Error("Failed to create outbound connection handler: ", err)
		return nil, err
	}
	vpoint.och = och
	outboundHandlerManager.SetDefaultHandler(och)

	vpoint.taggedIdh = make(map[string]InboundDetourHandler)
	detours := pConfig.InboundDetours
	if len(detours) > 0 {
		vpoint.idh = make([]InboundDetourHandler, len(detours))
		for idx, detourConfig := range detours {
			allocConfig := detourConfig.Allocation
			var detourHandler InboundDetourHandler
			switch allocConfig.Strategy {
			case AllocationStrategyAlways:
				dh, err := NewInboundDetourHandlerAlways(vpoint.space, detourConfig)
				if err != nil {
					log.Error("Point: Failed to create detour handler: ", err)
					return nil, ErrorBadConfiguration
				}
				detourHandler = dh
			case AllocationStrategyRandom:
				dh, err := NewInboundDetourHandlerDynamic(vpoint.space, detourConfig)
				if err != nil {
					log.Error("Point: Failed to create detour handler: ", err)
					return nil, ErrorBadConfiguration
				}
				detourHandler = dh
			default:
				log.Error("Point: Unknown allocation strategy: ", allocConfig.Strategy)
				return nil, ErrorBadConfiguration
			}
			vpoint.idh[idx] = detourHandler
			if len(detourConfig.Tag) > 0 {
				vpoint.taggedIdh[detourConfig.Tag] = detourHandler
			}
		}
	}

	outboundDetours := pConfig.OutboundDetours
	if len(outboundDetours) > 0 {
		vpoint.odh = make(map[string]proxy.OutboundHandler)
		for _, detourConfig := range outboundDetours {
			detourHandler, err := proxyrepo.CreateOutboundHandler(
				detourConfig.Protocol, vpoint.space, detourConfig.Settings, &proxy.OutboundHandlerMeta{
					Tag:            detourConfig.Tag,
					Address:        detourConfig.SendThrough,
					StreamSettings: detourConfig.StreamSettings,
				})
			if err != nil {
				log.Error("Point: Failed to create detour outbound connection handler: ", err)
				return nil, err
			}
			vpoint.odh[detourConfig.Tag] = detourHandler
			outboundHandlerManager.SetHandler(detourConfig.Tag, detourHandler)
		}
	}

	if err := vpoint.space.Initialize(); err != nil {
		return nil, err
	}

	return vpoint, nil
}

func (this *Point) Close() {
	this.ich.Close()
	for _, idh := range this.idh {
		idh.Close()
	}
}

// Start starts the Point server, and return any error during the process.
// In the case of any errors, the state of the server is unpredicatable.
func (this *Point) Start() error {
	if this.port <= 0 {
		log.Error("Point: Invalid port ", this.port)
		return ErrorBadConfiguration
	}

	err := retry.Timed(100 /* times */, 100 /* ms */).On(func() error {
		err := this.ich.Start()
		if err != nil {
			return err
		}
		log.Warning("Point: started on port ", this.port)
		return nil
	})
	if err != nil {
		return err
	}

	for _, detourHandler := range this.idh {
		err := detourHandler.Start()
		if err != nil {
			return err
		}
	}

	return nil
}

func (this *Point) GetHandler(tag string) (proxy.InboundHandler, int) {
	handler, found := this.taggedIdh[tag]
	if !found {
		log.Warning("Point: Unable to find an inbound handler with tag: ", tag)
		return nil, 0
	}
	return handler.GetConnectionHandler()
}

func (this *Point) Release() {

}
