package point

import (
	"github.com/v2ray/v2ray-core/app"
	"github.com/v2ray/v2ray-core/common/dice"
	"github.com/v2ray/v2ray-core/common/log"
	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/common/retry"
	"github.com/v2ray/v2ray-core/proxy"
	proxyrepo "github.com/v2ray/v2ray-core/proxy/repo"
)

type InboundConnectionHandlerWithPort struct {
	port    v2net.Port
	listen  v2net.Address
	handler proxy.InboundHandler
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
		ich, err := proxyrepo.CreateInboundHandler(config.Protocol, space, ichConfig, config.ListenOn, i)
		if err != nil {
			log.Error("Failed to create inbound connection handler: ", err)
			return nil, err
		}
		handler.ich = append(handler.ich, &InboundConnectionHandlerWithPort{
			port:    i,
			handler: ich,
			listen:  config.ListenOn,
		})
	}
	return handler, nil
}

func (this *InboundDetourHandlerAlways) GetConnectionHandler() (proxy.InboundHandler, int) {
	ich := this.ich[dice.Roll(len(this.ich))]
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
			err := ich.handler.Start()
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
