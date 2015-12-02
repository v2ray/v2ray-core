package point

import (
	"github.com/v2ray/v2ray-core/app/router"
	"github.com/v2ray/v2ray-core/common/log"
	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/common/retry"
	"github.com/v2ray/v2ray-core/proxy/common/connhandler"
	"github.com/v2ray/v2ray-core/shell/point/config"
	"github.com/v2ray/v2ray-core/transport/ray"
)

// Point is an single server in V2Ray system.
type Point struct {
	port   v2net.Port
	ich    connhandler.InboundConnectionHandler
	och    connhandler.OutboundConnectionHandler
	idh    []*InboundDetourHandler
	odh    map[string]connhandler.OutboundConnectionHandler
	router router.Router
}

// NewPoint returns a new Point server based on given configuration.
// The server is not started at this point.
func NewPoint(pConfig config.PointConfig) (*Point, error) {
	var vpoint = new(Point)
	vpoint.port = pConfig.Port()

	ichFactory := connhandler.GetInboundConnectionHandlerFactory(pConfig.InboundConfig().Protocol())
	if ichFactory == nil {
		log.Error("Unknown inbound connection handler factory %s", pConfig.InboundConfig().Protocol())
		return nil, config.BadConfiguration
	}
	ichConfig := pConfig.InboundConfig().Settings()
	ich, err := ichFactory.Create(vpoint, ichConfig)
	if err != nil {
		log.Error("Failed to create inbound connection handler: %v", err)
		return nil, err
	}
	vpoint.ich = ich

	ochFactory := connhandler.GetOutboundConnectionHandlerFactory(pConfig.OutboundConfig().Protocol())
	if ochFactory == nil {
		log.Error("Unknown outbound connection handler factory %s", pConfig.OutboundConfig().Protocol())
		return nil, config.BadConfiguration
	}
	ochConfig := pConfig.OutboundConfig().Settings()
	och, err := ochFactory.Create(ochConfig)
	if err != nil {
		log.Error("Failed to create outbound connection handler: %v", err)
		return nil, err
	}
	vpoint.och = och

	detours := pConfig.InboundDetours()
	if len(detours) > 0 {
		vpoint.idh = make([]*InboundDetourHandler, len(detours))
		for idx, detourConfig := range detours {
			detourHandler := &InboundDetourHandler{
				point:  vpoint,
				config: detourConfig,
			}
			err := detourHandler.Initialize()
			if err != nil {
				return nil, err
			}
			vpoint.idh[idx] = detourHandler
		}
	}

	outboundDetours := pConfig.OutboundDetours()
	if len(outboundDetours) > 0 {
		vpoint.odh = make(map[string]connhandler.OutboundConnectionHandler)
		for _, detourConfig := range outboundDetours {
			detourFactory := connhandler.GetOutboundConnectionHandlerFactory(detourConfig.Protocol())
			if detourFactory == nil {
				log.Error("Unknown detour outbound connection handler factory %s", detourConfig.Protocol())
				return nil, config.BadConfiguration
			}
			detourHandler, err := detourFactory.Create(detourConfig.Settings())
			if err != nil {
				log.Error("Failed to create detour outbound connection handler: %v", err)
				return nil, err
			}
			vpoint.odh[detourConfig.Tag()] = detourHandler
		}
	}

	routerConfig := pConfig.RouterConfig()
	if routerConfig != nil {
		r, err := router.CreateRouter(routerConfig.Strategy(), routerConfig.Settings())
		if err != nil {
			log.Error("Failed to create router: %v", err)
			return nil, config.BadConfiguration
		}
		vpoint.router = r
	}

	return vpoint, nil
}

// Start starts the Point server, and return any error during the process.
// In the case of any errors, the state of the server is unpredicatable.
func (this *Point) Start() error {
	if this.port <= 0 {
		log.Error("Invalid port %d", this.port)
		return config.BadConfiguration
	}

	err := retry.Timed(100 /* times */, 100 /* ms */).On(func() error {
		err := this.ich.Listen(this.port)
		if err != nil {
			return err
		}
		log.Warning("Point server started on port %d", this.port)
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

func (this *Point) DispatchToOutbound(packet v2net.Packet) ray.InboundRay {
	direct := ray.NewRay()
	dest := packet.Destination()

	if this.router != nil {
		tag, err := this.router.TakeDetour(dest)
		if err == nil {
			handler, found := this.odh[tag]
			if found {
				go handler.Dispatch(packet, direct)
				return direct
			}
		}
	}

	go this.och.Dispatch(packet, direct)
	return direct
}
