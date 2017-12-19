package freedom

//go:generate go run $GOPATH/src/v2ray.com/core/common/errors/errorgen/main.go -pkg freedom -path Proxy,Freedom

import (
	"context"

	"v2ray.com/core/app"
	"v2ray.com/core/app/dns"
	"v2ray.com/core/app/policy"
	"v2ray.com/core/common"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/dice"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/retry"
	"v2ray.com/core/common/signal"
	"v2ray.com/core/proxy"
	"v2ray.com/core/transport/internet"
	"v2ray.com/core/transport/ray"
)

// Handler handles Freedom connections.
type Handler struct {
	domainStrategy Config_DomainStrategy
	timeout        uint32
	dns            dns.Server
	destOverride   *DestinationOverride
	policy         policy.Policy
}

// New creates a new Freedom handler.
func New(ctx context.Context, config *Config) (*Handler, error) {
	space := app.SpaceFromContext(ctx)
	if space == nil {
		return nil, newError("no space in context")
	}
	f := &Handler{
		domainStrategy: config.DomainStrategy,
		timeout:        config.Timeout,
		destOverride:   config.DestinationOverride,
	}
	space.On(app.SpaceInitializing, func(interface{}) error {
		if config.DomainStrategy == Config_USE_IP {
			f.dns = dns.FromSpace(space)
			if f.dns == nil {
				return newError("DNS server is not found in the space")
			}
		}
		pm := policy.FromSpace(space)
		if pm == nil {
			return newError("Policy not found in space.")
		}
		f.policy = pm.GetPolicy(config.UserLevel)
		if config.Timeout > 0 && config.UserLevel == 0 {
			f.policy.Timeout.ConnectionIdle.Value = config.Timeout
		}
		return nil
	})
	return f, nil
}

func (h *Handler) resolveIP(ctx context.Context, domain string) net.Address {
	if resolver, ok := proxy.ResolvedIPsFromContext(ctx); ok {
		ips := resolver.Resolve()
		if len(ips) == 0 {
			return nil
		}
		return ips[dice.Roll(len(ips))]
	}

	ips := h.dns.Get(domain)
	if len(ips) == 0 {
		return nil
	}
	return net.IPAddress(ips[dice.Roll(len(ips))])
}

// Process implements proxy.Outbound.
func (h *Handler) Process(ctx context.Context, outboundRay ray.OutboundRay, dialer proxy.Dialer) error {
	destination, _ := proxy.TargetFromContext(ctx)
	if h.destOverride != nil {
		server := h.destOverride.Server
		destination = net.Destination{
			Network: destination.Network,
			Address: server.Address.AsAddress(),
			Port:    net.Port(server.Port),
		}
	}
	newError("opening connection to ", destination).WriteToLog()

	input := outboundRay.OutboundInput()
	output := outboundRay.OutboundOutput()

	if h.domainStrategy == Config_USE_IP && destination.Address.Family().IsDomain() {
		ip := h.resolveIP(ctx, destination.Address.Domain())
		if ip != nil {
			destination = net.Destination{
				Network: destination.Network,
				Address: ip,
				Port:    destination.Port,
			}
			newError("changing destination to ", destination).WriteToLog()
		}
	}

	var conn internet.Connection
	err := retry.ExponentialBackoff(5, 100).On(func() error {
		rawConn, err := dialer.Dial(ctx, destination)
		if err != nil {
			return err
		}
		conn = rawConn
		return nil
	})
	if err != nil {
		return newError("failed to open connection to ", destination).Base(err)
	}
	defer conn.Close()

	ctx, cancel := context.WithCancel(ctx)
	timer := signal.CancelAfterInactivity(ctx, cancel, h.policy.Timeout.ConnectionIdle.Duration())

	requestDone := signal.ExecuteAsync(func() error {
		var writer buf.Writer
		if destination.Network == net.Network_TCP {
			writer = buf.NewWriter(conn)
		} else {
			writer = buf.NewSequentialWriter(conn)
		}
		if err := buf.Copy(input, writer, buf.UpdateActivity(timer)); err != nil {
			return newError("failed to process request").Base(err)
		}
		timer.SetTimeout(h.policy.Timeout.DownlinkOnly.Duration())
		return nil
	})

	responseDone := signal.ExecuteAsync(func() error {
		defer output.Close()

		v2reader := buf.NewReader(conn)
		if err := buf.Copy(v2reader, output, buf.UpdateActivity(timer)); err != nil {
			return newError("failed to process response").Base(err)
		}
		timer.SetTimeout(h.policy.Timeout.UplinkOnly.Duration())
		return nil
	})

	if err := signal.ErrorOrFinish2(ctx, requestDone, responseDone); err != nil {
		input.CloseError()
		output.CloseError()
		return newError("connection ends").Base(err)
	}

	return nil
}

func init() {
	common.Must(common.RegisterConfig((*Config)(nil), func(ctx context.Context, config interface{}) (interface{}, error) {
		return New(ctx, config.(*Config))
	}))
}
