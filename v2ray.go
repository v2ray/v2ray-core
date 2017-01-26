package core

import (
	"context"

	"v2ray.com/core/app"
	"v2ray.com/core/app/dispatcher"
	"v2ray.com/core/app/dns"
	"v2ray.com/core/app/proxyman"
	"v2ray.com/core/common/log"
	v2net "v2ray.com/core/common/net"
)

// Point shell of V2Ray.
type Point struct {
	space app.Space
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

	space := app.NewSpace()
	ctx := app.ContextWithSpace(context.Background(), space)

	vpoint.space = space

	outboundHandlerManager := proxyman.OutboundHandlerManagerFromSpace(space)
	if outboundHandlerManager == nil {
		o, err := app.CreateAppFromConfig(ctx, new(proxyman.OutboundConfig))
		if err != nil {
			return nil, err
		}
		space.AddApplication(o)
		outboundHandlerManager = o.(proxyman.OutboundHandlerManager)
	}

	inboundHandlerManager := proxyman.InboundHandlerManagerFromSpace(space)
	if inboundHandlerManager == nil {
		o, err := app.CreateAppFromConfig(ctx, new(proxyman.InboundConfig))
		if err != nil {
			return nil, err
		}
		space.AddApplication(o)
		inboundHandlerManager = o.(proxyman.InboundHandlerManager)
	}

	for _, appSettings := range pConfig.App {
		settings, err := appSettings.GetInstance()
		if err != nil {
			return nil, err
		}
		application, err := app.CreateAppFromConfig(ctx, settings)
		if err != nil {
			return nil, err
		}
		if err := space.AddApplication(application); err != nil {
			return nil, err
		}
	}

	dnsServer := dns.FromSpace(space)
	if dnsServer == nil {
		dnsConfig := &dns.Config{
			NameServers: []*v2net.Endpoint{{
				Address: v2net.NewIPOrDomain(v2net.LocalHostDomain),
			}},
		}
		d, err := app.CreateAppFromConfig(ctx, dnsConfig)
		if err != nil {
			return nil, err
		}
		space.AddApplication(d)
		dnsServer = d.(dns.Server)
	}

	disp := dispatcher.FromSpace(space)
	if disp == nil {
		d, err := app.CreateAppFromConfig(ctx, new(dispatcher.Config))
		if err != nil {
			return nil, err
		}
		space.AddApplication(d)
		disp = d.(dispatcher.Interface)
	}

	for _, inbound := range pConfig.Inbound {
		if err := inboundHandlerManager.AddHandler(ctx, inbound); err != nil {
			return nil, err
		}
	}

	for _, outbound := range pConfig.Outbound {
		if err := outboundHandlerManager.AddHandler(ctx, outbound); err != nil {
			return nil, err
		}
	}

	if err := vpoint.space.Initialize(); err != nil {
		return nil, err
	}

	return vpoint, nil
}

func (v *Point) Close() {
	ihm := proxyman.InboundHandlerManagerFromSpace(v.space)
	ihm.Close()
}

// Start starts the Point server, and return any error during the process.
// In the case of any errors, the state of the server is unpredicatable.
func (v *Point) Start() error {
	ihm := proxyman.InboundHandlerManagerFromSpace(v.space)
	if err := ihm.Start(); err != nil {
		return err
	}
	log.Warning("V2Ray started.")

	return nil
}
