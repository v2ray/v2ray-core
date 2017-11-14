package dokodemo

//go:generate go run $GOPATH/src/v2ray.com/core/tools/generrorgen/main.go -pkg dokodemo -path Proxy,Dokodemo

import (
	"context"
	"time"

	"v2ray.com/core/app"
	"v2ray.com/core/app/dispatcher"
	"v2ray.com/core/app/log"
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

	timeout := time.Second * time.Duration(d.config.Timeout)
	if timeout == 0 {
		timeout = time.Minute * 5
	}

	ctx, cancel := context.WithCancel(ctx)
	timer := signal.CancelAfterInactivity(ctx, cancel, timeout)

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

		timer.SetTimeout(time.Second * 2)

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
