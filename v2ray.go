package core

import (
	"context"

	"v2ray.com/core/app"
	"v2ray.com/core/app/dispatcher"
	"v2ray.com/core/app/dns"
	"v2ray.com/core/app/log"
	"v2ray.com/core/app/proxyman"
	"v2ray.com/core/common"
	"v2ray.com/core/common/errors"
	"v2ray.com/core/common/net"
)

// Server is an instance of V2Ray. At any time, there must be at most one Server instance running.
type Server interface {
	// Start starts the V2Ray server, and return any error during the process.
	// In the case of any errors, the state of the server is unpredicatable.
	Start() error

	// Close closes the V2Ray server. All inbound and outbound connections will be closed immediately.
	Close()
}

// New creates a new V2Ray server with given config.
func New(config *Config) (Server, error) {
	return newSimpleServer(config)
}

// simpleServer shell of V2Ray.
type simpleServer struct {
	space app.Space
}

// newSimpleServer returns a new Point server based on given configuration.
// The server is not started at this point.
func newSimpleServer(config *Config) (*simpleServer, error) {
	var server = new(simpleServer)

	if err := config.Transport.Apply(); err != nil {
		return nil, err
	}

	space := app.NewSpace()
	ctx := app.ContextWithSpace(context.Background(), space)

	server.space = space

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

	if log.FromSpace(space) == nil {
		l, err := app.CreateAppFromConfig(ctx, &log.Config{
			ErrorLogType:  log.LogType_Console,
			ErrorLogLevel: log.LogLevel_Warning,
			AccessLogType: log.LogType_None,
		})
		if err != nil {
			return nil, errors.Base(err).Message("Core: Failed apply default log settings.")
		}
		space.AddApplication(l)
	}

	outboundHandlerManager := proxyman.OutboundHandlerManagerFromSpace(space)
	if outboundHandlerManager == nil {
		o, err := app.CreateAppFromConfig(ctx, new(proxyman.OutboundConfig))
		if err != nil {
			return nil, err
		}
		if err := space.AddApplication(o); err != nil {
			return nil, errors.Base(err).Message("Core: Failed to add default outbound handler manager.")
		}
		outboundHandlerManager = o.(proxyman.OutboundHandlerManager)
	}

	inboundHandlerManager := proxyman.InboundHandlerManagerFromSpace(space)
	if inboundHandlerManager == nil {
		o, err := app.CreateAppFromConfig(ctx, new(proxyman.InboundConfig))
		if err != nil {
			return nil, err
		}
		if err := space.AddApplication(o); err != nil {
			return nil, errors.Base(err).Message("Core: Failed to add default inbound handler manager.")
		}
		inboundHandlerManager = o.(proxyman.InboundHandlerManager)
	}

	if dns.FromSpace(space) == nil {
		dnsConfig := &dns.Config{
			NameServers: []*net.Endpoint{{
				Address: net.NewIPOrDomain(net.LocalHostDomain),
			}},
		}
		d, err := app.CreateAppFromConfig(ctx, dnsConfig)
		if err != nil {
			return nil, err
		}
		common.Must(space.AddApplication(d))
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

	if err := server.space.Initialize(); err != nil {
		return nil, err
	}

	return server, nil
}

func (s *simpleServer) Close() {
	s.space.Close()
}

func (s *simpleServer) Start() error {
	if err := s.space.Start(); err != nil {
		return err
	}
	log.Warning("V2Ray started.")

	return nil
}
