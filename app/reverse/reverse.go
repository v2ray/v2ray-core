package reverse

//go:generate errorgen

import (
	"context"

	"v2ray.com/core"
	"v2ray.com/core/common"
	"v2ray.com/core/common/net"
	"v2ray.com/core/features/outbound"
	"v2ray.com/core/features/routing"
)

const (
	internalDomain = "reverse.internal.v2ray.com"
)

func isInternalDomain(dest net.Destination) bool {
	return dest.Address.Family().IsDomain() && dest.Address.Domain() == internalDomain
}

func init() {
	common.Must(common.RegisterConfig((*Config)(nil), func(ctx context.Context, config interface{}) (interface{}, error) {
		r := new(Reverse)
		if err := core.RequireFeatures(ctx, func(d routing.Dispatcher, om outbound.Manager) error {
			return r.Init(config.(*Config), d, om)
		}); err != nil {
			return nil, err
		}
		return r, nil
	}))
}

type Reverse struct {
	bridges []*Bridge
	portals []*Portal
}

func (r *Reverse) Init(config *Config, d routing.Dispatcher, ohm outbound.Manager) error {
	for _, bConfig := range config.BridgeConfig {
		b, err := NewBridge(bConfig, d)
		if err != nil {
			return err
		}
		r.bridges = append(r.bridges, b)
	}

	for _, pConfig := range config.PortalConfig {
		p, err := NewPortal(pConfig, ohm)
		if err != nil {
			return err
		}
		r.portals = append(r.portals, p)
	}

	return nil
}
