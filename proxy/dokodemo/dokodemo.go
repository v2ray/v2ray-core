package dokodemo

//go:generate go run $GOPATH/src/v2ray.com/core/common/errors/errorgen/main.go -pkg dokodemo -path Proxy,Dokodemo

import (
	"context"

	"v2ray.com/core/app"
	"v2ray.com/core/app/dispatcher"
	"v2ray.com/core/app/log"
	"v2ray.com/core/app/policy"
	"v2ray.com/core/common"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/signal"
	"v2ray.com/core/proxy"
	"v2ray.com/core/transport/internet"
	"v2ray.com/core/transport/internet/udp"
)

type DokodemoDoor struct {
	config  *Config
	address net.Address
	port    net.Port
	policy  policy.Policy
}

func New(ctx context.Context, config *Config) (*DokodemoDoor, error) {
	space := app.SpaceFromContext(ctx)
	if space == nil {
		return nil, newError("no space in context")
	}
	if config.NetworkList == nil || config.NetworkList.Size() == 0 {
		return nil, newError("no network specified")
	}
	d := &DokodemoDoor{
		config:  config,
		address: config.GetPredefinedAddress(),
		port:    net.Port(config.Port),
	}
	space.On(app.SpaceInitializing, func(interface{}) error {
		pm := policy.FromSpace(space)
		if pm == nil {
			return newError("Policy not found in space.")
		}
		d.policy = pm.GetPolicy(config.UserLevel)
		if config.Timeout > 0 && config.UserLevel == 0 {
			d.policy.Timeout.ConnectionIdle.Value = config.Timeout
		}
		return nil
	})
	return d, nil
}

func (d *DokodemoDoor) Network() net.NetworkList {
	return *(d.config.NetworkList)
}

func (d *DokodemoDoor) Process(ctx context.Context, network net.Network, conn internet.Connection, dispatcher dispatcher.Interface) error {
	log.Trace(newError("processing connection from: ", conn.RemoteAddr()).AtDebug())
	dest := net.Destination{
		Network: network,
		Address: d.address,
		Port:    d.port,
	}
	if d.config.FollowRedirect {
		if origDest, ok := proxy.OriginalTargetFromContext(ctx); ok {
			dest = origDest
		}
	}
	if !dest.IsValid() || dest.Address == nil {
		return newError("unable to get destination")
	}

	ctx, cancel := context.WithCancel(ctx)
	timer := signal.CancelAfterInactivity(ctx, cancel, d.policy.Timeout.ConnectionIdle.Duration())

	inboundRay, err := dispatcher.Dispatch(ctx, dest)
	if err != nil {
		return newError("failed to dispatch request").Base(err)
	}

	requestDone := signal.ExecuteAsync(func() error {
		defer inboundRay.InboundInput().Close()

		chunkReader := buf.NewReader(conn)

		if err := buf.Copy(chunkReader, inboundRay.InboundInput(), buf.UpdateActivity(timer)); err != nil {
			return newError("failed to transport request").Base(err)
		}

		timer.SetTimeout(d.policy.Timeout.DownlinkOnly.Duration())

		return nil
	})

	responseDone := signal.ExecuteAsync(func() error {
		var writer buf.Writer
		if network == net.Network_TCP {
			writer = buf.NewWriter(conn)
		} else {
			//if we are in TPROXY mode, use linux's udp forging functionality
			if !d.config.FollowRedirect {
				writer = buf.NewSequentialWriter(conn)
			} else {
				srca := net.UDPAddr{IP: dest.Address.IP(), Port: int(dest.Port.Value())}
				origsend, err := udp.TransmitSocket(&srca, conn.RemoteAddr())
				if err != nil {
					return err
				}
				writer = buf.NewSequentialWriter(origsend)
			}
		}

		if err := buf.Copy(inboundRay.InboundOutput(), writer, buf.UpdateActivity(timer)); err != nil {
			return newError("failed to transport response").Base(err)
		}

		timer.SetTimeout(d.policy.Timeout.UplinkOnly.Duration())

		return nil
	})

	if err := signal.ErrorOrFinish2(ctx, requestDone, responseDone); err != nil {
		inboundRay.InboundInput().CloseError()
		inboundRay.InboundOutput().CloseError()
		return newError("connection ends").Base(err)
	}

	return nil
}

func init() {
	common.Must(common.RegisterConfig((*Config)(nil), func(ctx context.Context, config interface{}) (interface{}, error) {
		return New(ctx, config.(*Config))
	}))
}
