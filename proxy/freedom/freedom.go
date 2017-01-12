package freedom

import (
	"io"

	"v2ray.com/core/app"
	"v2ray.com/core/app/dns"
	"v2ray.com/core/common"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/dice"
	"v2ray.com/core/common/errors"
	"v2ray.com/core/common/log"
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/common/retry"
	"v2ray.com/core/common/serial"
	"v2ray.com/core/common/signal"
	"v2ray.com/core/proxy"
	"v2ray.com/core/transport/internet"
	"v2ray.com/core/transport/ray"
)

type Handler struct {
	domainStrategy Config_DomainStrategy
	timeout        uint32
	dns            dns.Server
	meta           *proxy.OutboundHandlerMeta
}

func New(config *Config, space app.Space, meta *proxy.OutboundHandlerMeta) *Handler {
	f := &Handler{
		domainStrategy: config.DomainStrategy,
		timeout:        config.Timeout,
		meta:           meta,
	}
	space.OnInitialize(func() error {
		if config.DomainStrategy == Config_USE_IP {
			f.dns = dns.FromSpace(space)
			if f.dns == nil {
				return errors.New("Freedom: DNS server is not found in the space.")
			}
		}
		return nil
	})
	return f
}

// Private: Visible for testing.
func (v *Handler) ResolveIP(destination v2net.Destination) v2net.Destination {
	if !destination.Address.Family().IsDomain() {
		return destination
	}

	ips := v.dns.Get(destination.Address.Domain())
	if len(ips) == 0 {
		log.Info("Freedom: DNS returns nil answer. Keep domain as is.")
		return destination
	}

	ip := ips[dice.Roll(len(ips))]
	var newDest v2net.Destination
	if destination.Network == v2net.Network_TCP {
		newDest = v2net.TCPDestination(v2net.IPAddress(ip), destination.Port)
	} else {
		newDest = v2net.UDPDestination(v2net.IPAddress(ip), destination.Port)
	}
	log.Info("Freedom: Changing destination from ", destination, " to ", newDest)
	return newDest
}

func (v *Handler) Dispatch(destination v2net.Destination, ray ray.OutboundRay) {
	log.Info("Freedom: Opening connection to ", destination)

	input := ray.OutboundInput()
	output := ray.OutboundOutput()

	var conn internet.Connection
	if v.domainStrategy == Config_USE_IP && destination.Address.Family().IsDomain() {
		destination = v.ResolveIP(destination)
	}
	err := retry.ExponentialBackoff(5, 100).On(func() error {
		rawConn, err := internet.Dial(v.meta.Address, destination, v.meta.GetDialerOptions())
		if err != nil {
			return err
		}
		conn = rawConn
		return nil
	})
	if err != nil {
		log.Warning("Freedom: Failed to open connection to ", destination, ": ", err)
		return
	}
	defer conn.Close()

	conn.SetReusable(false)

	requestDone := signal.ExecuteAsync(func() error {
		v2writer := buf.NewWriter(conn)
		if err := buf.PipeUntilEOF(input, v2writer); err != nil {
			return err
		}
		return nil
	})

	var reader io.Reader = conn

	timeout := v.timeout
	if destination.Network == v2net.Network_UDP {
		timeout = 16
	}
	if timeout > 0 {
		reader = v2net.NewTimeOutReader(timeout /* seconds */, conn)
	}

	responseDone := signal.ExecuteAsync(func() error {
		defer output.Close()

		v2reader := buf.NewReader(reader)
		if err := buf.PipeUntilEOF(v2reader, output); err != nil {
			return err
		}
		return nil
	})

	if err := signal.ErrorOrFinish2(requestDone, responseDone); err != nil {
		log.Info("Freedom: Connection ending with ", err)
		input.CloseError()
		output.CloseError()
	}
}

type Factory struct{}

func (v *Factory) Create(space app.Space, config interface{}, meta *proxy.OutboundHandlerMeta) (proxy.OutboundHandler, error) {
	return New(config.(*Config), space, meta), nil
}

func init() {
	common.Must(proxy.RegisterOutboundHandlerCreator(serial.GetMessageType(new(Config)), new(Factory)))
}
