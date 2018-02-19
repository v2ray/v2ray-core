package freedom

//go:generate go run $GOPATH/src/v2ray.com/core/common/errors/errorgen/main.go -pkg freedom -path Proxy,Freedom

import (
	"context"
	"time"

	"v2ray.com/core"
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
	policyManager core.PolicyManager
	dns           core.DNSClient
	config        Config
}

// New creates a new Freedom handler.
func New(ctx context.Context, config *Config) (*Handler, error) {
	v := core.FromContext(ctx)
	if v == nil {
		return nil, newError("V is not found in context.")
	}

	f := &Handler{
		config:        *config,
		policyManager: v.PolicyManager(),
		dns:           v.DNSClient(),
	}

	return f, nil
}

func (h *Handler) policy() core.Policy {
	p := h.policyManager.ForLevel(h.config.UserLevel)
	if h.config.Timeout > 0 && h.config.UserLevel == 0 {
		p.Timeouts.ConnectionIdle = time.Duration(h.config.Timeout) * time.Second
	}
	return p
}

func (h *Handler) resolveIP(ctx context.Context, domain string) net.Address {
	if resolver, ok := proxy.ResolvedIPsFromContext(ctx); ok {
		ips := resolver.Resolve()
		if len(ips) == 0 {
			return nil
		}
		return ips[dice.Roll(len(ips))]
	}

	ips, err := h.dns.LookupIP(domain)
	if err != nil {
		newError("failed to get IP address for domain ", domain).Base(err).WriteToLog()
	}
	if len(ips) == 0 {
		return nil
	}
	return net.IPAddress(ips[dice.Roll(len(ips))])
}

// Process implements proxy.Outbound.
func (h *Handler) Process(ctx context.Context, outboundRay ray.OutboundRay, dialer proxy.Dialer) error {
	destination, _ := proxy.TargetFromContext(ctx)
	if h.config.DestinationOverride != nil {
		server := h.config.DestinationOverride.Server
		destination = net.Destination{
			Network: destination.Network,
			Address: server.Address.AsAddress(),
			Port:    net.Port(server.Port),
		}
	}
	newError("opening connection to ", destination).WriteToLog()

	input := outboundRay.OutboundInput()
	output := outboundRay.OutboundOutput()

	if h.config.DomainStrategy == Config_USE_IP && destination.Address.Family().IsDomain() {
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
	timer := signal.CancelAfterInactivity(ctx, cancel, h.policy().Timeouts.ConnectionIdle)

	requestDone := signal.ExecuteAsync(func() error {
		defer timer.SetTimeout(h.policy().Timeouts.DownlinkOnly)

		var writer buf.Writer
		if destination.Network == net.Network_TCP {
			writer = buf.NewWriter(conn)
		} else {
			writer = buf.NewSequentialWriter(conn)
		}
		if err := buf.Copy(input, writer, buf.UpdateActivity(timer)); err != nil {
			return newError("failed to process request").Base(err)
		}

		return nil
	})

	responseDone := signal.ExecuteAsync(func() error {
		defer timer.SetTimeout(h.policy().Timeouts.UplinkOnly)

		v2reader := buf.NewReader(conn)
		if err := buf.Copy(v2reader, output, buf.UpdateActivity(timer)); err != nil {
			return newError("failed to process response").Base(err)
		}

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
