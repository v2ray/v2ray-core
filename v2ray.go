package core

import (
	"context"

	"v2ray.com/core/app"
	"v2ray.com/core/app/dispatcher"
	"v2ray.com/core/app/dns"
	"v2ray.com/core/app/log"
	"v2ray.com/core/app/proxyman"
	"v2ray.com/core/common/net"
)

// Point shell of V2Ray.
type Point struct {
	space app.Space
}

// NewPoint returns a new Point server based on given configuration.
// The server is not started at this point.
func NewPoint(config *Config) (*Point, error) {
	var pt = new(Point)

	if err := config.Transport.Apply(); err != nil {
		return nil, err
	}

	space := app.NewSpace()
	ctx := app.ContextWithSpace(context.Background(), space)

	pt.space = space

	for _, appSettings := range config.App {
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

	logger := log.FromSpace(space)
	if logger == nil {
		l, err := app.CreateAppFromConfig(ctx, &log.Config{
			ErrorLogType:  log.LogType_Console,
			ErrorLogLevel: log.LogLevel_Warning,
			AccessLogType: log.LogType_None,
		})
		if err != nil {
			return nil, err
		}
		space.AddApplication(l)
	}

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

	dnsServer := dns.FromSpace(space)
	if dnsServer == nil {
		dnsConfig := &dns.Config{
			NameServers: []*net.Endpoint{{
				Address: net.NewIPOrDomain(net.LocalHostDomain),
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

	for _, inbound := range config.Inbound {
		if err := inboundHandlerManager.AddHandler(ctx, inbound); err != nil {
			return nil, err
		}
	}

	for _, outbound := range config.Outbound {
		if err := outboundHandlerManager.AddHandler(ctx, outbound); err != nil {
			return nil, err
		}
	}

	if err := pt.space.Initialize(); err != nil {
		return nil, err
	}

	return pt, nil
}

func (v *Point) Close() {
	v.space.Close()
}

// Start starts the Point server, and return any error during the process.
// In the case of any errors, the state of the server is unpredicatable.
func (v *Point) Start() error {
	if err := v.space.Start(); err != nil {
		return err
	}
	log.Warning("V2Ray started.")

	return nil
}
