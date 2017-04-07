package freedom

import (
	"context"
	"time"

	"runtime"

	"v2ray.com/core/app"
	"v2ray.com/core/app/dns"
	"v2ray.com/core/app/log"
	"v2ray.com/core/common"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/dice"
	"v2ray.com/core/common/errors"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/retry"
	"v2ray.com/core/common/signal"
	"v2ray.com/core/proxy"
	"v2ray.com/core/transport/internet"
	"v2ray.com/core/transport/ray"
)

type Handler struct {
	domainStrategy Config_DomainStrategy
	timeout        uint32
	dns            dns.Server
	destOverride   *DestinationOverride
}

func New(ctx context.Context, config *Config) (*Handler, error) {
	space := app.SpaceFromContext(ctx)
	if space == nil {
		return nil, errors.New("no space in context").Path("Proxy", "Freedom")
	}
	f := &Handler{
		domainStrategy: config.DomainStrategy,
		timeout:        config.Timeout,
		destOverride:   config.DestinationOverride,
	}
	space.OnInitialize(func() error {
		if config.DomainStrategy == Config_USE_IP {
			f.dns = dns.FromSpace(space)
			if f.dns == nil {
				return errors.New("DNS server is not found in the space").Path("Proxy", "Freedom")
			}
		}
		return nil
	})
	return f, nil
}

func (v *Handler) ResolveIP(destination net.Destination) net.Destination {
	if !destination.Address.Family().IsDomain() {
		return destination
	}

	ips := v.dns.Get(destination.Address.Domain())
	if len(ips) == 0 {
		log.Trace(errors.New("DNS returns nil answer. Keep domain as is.").Path("Proxy", "Freedom"))
		return destination
	}

	ip := ips[dice.Roll(len(ips))]
	var newDest net.Destination
	if destination.Network == net.Network_TCP {
		newDest = net.TCPDestination(net.IPAddress(ip), destination.Port)
	} else {
		newDest = net.UDPDestination(net.IPAddress(ip), destination.Port)
	}
	log.Trace(errors.New("changing destination from ", destination, " to ", newDest).Path("Proxy", "Freedom"))
	return newDest
}

func (v *Handler) Process(ctx context.Context, outboundRay ray.OutboundRay, dialer proxy.Dialer) error {
	destination, _ := proxy.TargetFromContext(ctx)
	if v.destOverride != nil {
		server := v.destOverride.Server
		destination = net.Destination{
			Network: destination.Network,
			Address: server.Address.AsAddress(),
			Port:    net.Port(server.Port),
		}
	}
	log.Trace(errors.New("opening connection to ", destination).Path("Proxy", "Freedom"))

	input := outboundRay.OutboundInput()
	output := outboundRay.OutboundOutput()

	var conn internet.Connection
	if v.domainStrategy == Config_USE_IP && destination.Address.Family().IsDomain() {
		destination = v.ResolveIP(destination)
	}

	err := retry.ExponentialBackoff(5, 100).On(func() error {
		rawConn, err := dialer.Dial(ctx, destination)
		if err != nil {
			return err
		}
		conn = rawConn
		return nil
	})
	if err != nil {
		return errors.New("failed to open connection to ", destination).Base(err).Path("Proxy", "Freedom")
	}
	defer conn.Close()

	timeout := time.Second * time.Duration(v.timeout)
	if timeout == 0 {
		timeout = time.Minute * 5
	}
	ctx, timer := signal.CancelAfterInactivity(ctx, timeout)

	requestDone := signal.ExecuteAsync(func() error {
		v2writer := buf.NewWriter(conn)
		if err := buf.PipeUntilEOF(timer, input, v2writer); err != nil {
			return err
		}
		return nil
	})

	responseDone := signal.ExecuteAsync(func() error {
		defer output.Close()

		v2reader := buf.NewReader(conn)
		if err := buf.PipeUntilEOF(timer, v2reader, output); err != nil {
			return err
		}
		return nil
	})

	if err := signal.ErrorOrFinish2(ctx, requestDone, responseDone); err != nil {
		input.CloseError()
		output.CloseError()
		return errors.New("connection ends").Base(err).Path("Proxy", "Freedom")
	}

	runtime.KeepAlive(timer)

	return nil
}

func init() {
	common.Must(common.RegisterConfig((*Config)(nil), func(ctx context.Context, config interface{}) (interface{}, error) {
		return New(ctx, config.(*Config))
	}))
}
