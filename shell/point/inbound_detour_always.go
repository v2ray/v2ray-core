package point

import (
	"math/rand"

	"github.com/v2ray/v2ray-core/app"
	"github.com/v2ray/v2ray-core/common/log"
	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/common/retry"
	"github.com/v2ray/v2ray-core/proxy"
	proxyrepo "github.com/v2ray/v2ray-core/proxy/repo"
)

type InboundConnectionHandlerWithPort struct {
	port    v2net.Port
	handler proxy.InboundConnectionHandler
}

// Handler for inbound detour connections.
type InboundDetourHandlerAlways struct {
	space  app.Space
	config *InboundDetourConfig
	ich    []*InboundConnectionHandlerWithPort
}

func NewInboundDetourHandlerAlways(space app.Space, config *InboundDetourConfig) (*InboundDetourHandlerAlways, error) {
	handler := &InboundDetourHandlerAlways{
		space:  space,
		config: config,
	}
	ports := config.PortRange
	handler.ich = make([]*InboundConnectionHandlerWithPort, 0, ports.To-ports.From+1)
	for i := ports.From; i <= ports.To; i++ {
		ichConfig := config.Settings
		ich, err := proxyrepo.CreateInboundConnectionHandler(config.Protocol, space, ichConfig)
		if err != nil {
			log.Error("Failed to create inbound connection handler: ", err)
			return nil, err
		}
		handler.ich = append(handler.ich, &InboundConnectionHandlerWithPort{
			port:    i,
			handler: ich,
		})
	}
	return handler, nil
}

func (this *InboundDetourHandlerAlways) GetConnectionHandler() (proxy.InboundConnectionHandler, int) {
	idx := rand.Intn(len(this.ich))
	ich := this.ich[idx]
	return ich.handler, this.config.Allocation.Refresh
}

func (this *InboundDetourHandlerAlways) Close() {
	for _, ich := range this.ich {
		ich.handler.Close()
	}
}

// Starts the inbound connection handler.
func (this *InboundDetourHandlerAlways) Start() error {
	for _, ich := range this.ich {
		err := retry.Timed(100 /* times */, 100 /* ms */).On(func() error {
			err := ich.handler.Listen(ich.port)
			if err != nil {
				log.Error("Failed to start inbound detour on port ", ich.port, ": ", err)
				return err
			}
			return nil
		})
		if err != nil {
			return err
		}
	}
	return nil
}
