package core

import (
	"github.com/v2ray/v2ray-core/common/log"
	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/config"
)

var (
	inboundFactories  = make(map[string]InboundConnectionHandlerFactory)
	outboundFactories = make(map[string]OutboundConnectionHandlerFactory)
)

func RegisterInboundConnectionHandlerFactory(name string, factory InboundConnectionHandlerFactory) error {
	// TODO check name
	inboundFactories[name] = factory
	return nil
}

func RegisterOutboundConnectionHandlerFactory(name string, factory OutboundConnectionHandlerFactory) error {
	// TODO check name
	outboundFactories[name] = factory
	return nil
}

// Point is an single server in V2Ray system.
type Point struct {
	port uint16
	ich  InboundConnectionHandler
	och  OutboundConnectionHandler
}

// NewPoint returns a new Point server based on given configuration.
// The server is not started at this point.
func NewPoint(pConfig config.PointConfig) (*Point, error) {
	var vpoint = new(Point)
	vpoint.port = pConfig.Port()

	ichFactory, ok := inboundFactories[pConfig.InboundConfig().Protocol()]
	if !ok {
		panic(log.Error("Unknown inbound connection handler factory %s", pConfig.InboundConfig().Protocol()))
	}
	ichConfig := pConfig.InboundConfig().Settings(config.TypeInbound)
	ich, err := ichFactory.Create(vpoint, ichConfig)
	if err != nil {
		log.Error("Failed to create inbound connection handler: %v", err)
		return nil, err
	}
	vpoint.ich = ich

	ochFactory, ok := outboundFactories[pConfig.OutboundConfig().Protocol()]
	if !ok {
		panic(log.Error("Unknown outbound connection handler factory %s", pConfig.OutboundConfig().Protocol))
	}
	ochConfig := pConfig.OutboundConfig().Settings(config.TypeOutbound)
	och, err := ochFactory.Create(vpoint, ochConfig)
	if err != nil {
		log.Error("Failed to create outbound connection handler: %v", err)
		return nil, err
	}
	vpoint.och = och

	return vpoint, nil
}

type InboundConnectionHandlerFactory interface {
	Create(vp *Point, config interface{}) (InboundConnectionHandler, error)
}

type InboundConnectionHandler interface {
	Listen(port uint16) error
}

type OutboundConnectionHandlerFactory interface {
	Create(VP *Point, config interface{}) (OutboundConnectionHandler, error)
}

type OutboundConnectionHandler interface {
	Dispatch(firstPacket v2net.Packet, ray OutboundRay) error
}

// Start starts the Point server, and return any error during the process.
// In the case of any errors, the state of the server is unpredicatable.
func (vp *Point) Start() error {
	if vp.port <= 0 {
		return log.Error("Invalid port %d", vp.port)
	}

	err := vp.ich.Listen(vp.port)
	// TODO: handle error
	return err
}

func (p *Point) DispatchToOutbound(packet v2net.Packet) InboundRay {
	ray := NewRay()
	// TODO: handle error
	p.och.Dispatch(packet, ray)
	return ray
}
