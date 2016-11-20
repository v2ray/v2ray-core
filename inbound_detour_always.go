package core

import (
	"v2ray.com/core/app"
	"v2ray.com/core/common/dice"
	"v2ray.com/core/common/log"
	"v2ray.com/core/common/retry"
	"v2ray.com/core/proxy"
	proxyregistry "v2ray.com/core/proxy/registry"
)

// Handler for inbound detour connections.
type InboundDetourHandlerAlways struct {
	space  app.Space
	config *InboundConnectionConfig
	ich    []proxy.InboundHandler
}

func NewInboundDetourHandlerAlways(space app.Space, config *InboundConnectionConfig) (*InboundDetourHandlerAlways, error) {
	handler := &InboundDetourHandlerAlways{
		space:  space,
		config: config,
	}
	ports := config.PortRange
	handler.ich = make([]proxy.InboundHandler, 0, ports.To-ports.From+1)
	for i := ports.FromPort(); i <= ports.ToPort(); i++ {
		ichConfig, err := config.GetTypedSettings()
		if err != nil {
			return nil, err
		}
		ich, err := proxyregistry.CreateInboundHandler(config.Settings.Type, space, ichConfig, &proxy.InboundHandlerMeta{
			Address:                config.GetListenOnValue(),
			Port:                   i,
			Tag:                    config.Tag,
			StreamSettings:         config.StreamSettings,
			AllowPassiveConnection: config.AllowPassiveConnection,
		})
		if err != nil {
			log.Error("Failed to create inbound connection handler: ", err)
			return nil, err
		}
		handler.ich = append(handler.ich, ich)
	}
	return handler, nil
}

func (this *InboundDetourHandlerAlways) GetConnectionHandler() (proxy.InboundHandler, int) {
	ich := this.ich[dice.Roll(len(this.ich))]
	return ich, int(this.config.GetAllocationStrategyValue().Refresh.GetValue())
}

func (this *InboundDetourHandlerAlways) Close() {
	for _, ich := range this.ich {
		ich.Close()
	}
}

// Starts the inbound connection handler.
func (this *InboundDetourHandlerAlways) Start() error {
	for _, ich := range this.ich {
		err := retry.ExponentialBackoff(10 /* times */, 200 /* ms */).On(func() error {
			err := ich.Start()
			if err != nil {
				log.Error("Failed to start inbound detour:", err)
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
