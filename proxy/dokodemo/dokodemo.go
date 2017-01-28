package dokodemo

import (
	"context"

	"v2ray.com/core/app"
	"v2ray.com/core/app/dispatcher"
	"v2ray.com/core/common"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/errors"
	"v2ray.com/core/common/log"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/signal"
	"v2ray.com/core/proxy"
	"v2ray.com/core/transport/internet"
)

type DokodemoDoor struct {
	config           *Config
	address          net.Address
	port             net.Port
	packetDispatcher dispatcher.Interface
}

func New(ctx context.Context, config *Config) (*DokodemoDoor, error) {
	space := app.SpaceFromContext(ctx)
	if space == nil {
		return nil, errors.New("Dokodemo: No space in context.")
	}
	if config.NetworkList == nil || config.NetworkList.Size() == 0 {
		return nil, errors.New("DokodemoDoor: No network specified.")
	}
	d := &DokodemoDoor{
		config:  config,
		address: config.GetPredefinedAddress(),
		port:    net.Port(config.Port),
	}
	space.OnInitialize(func() error {
		d.packetDispatcher = dispatcher.FromSpace(space)
		if d.packetDispatcher == nil {
			return errors.New("Dokodemo: Dispatcher is not found in the space.")
		}
		return nil
	})
	return d, nil
}

func (d *DokodemoDoor) Network() net.NetworkList {
	return *(d.config.NetworkList)
}

func (d *DokodemoDoor) Process(ctx context.Context, network net.Network, conn internet.Connection) error {
	log.Debug("Dokodemo: processing connection from: ", conn.RemoteAddr())
	conn.SetReusable(false)
	ctx = proxy.ContextWithDestination(ctx, net.Destination{
		Network: network,
		Address: d.address,
		Port:    d.port,
	})
	inboundRay := d.packetDispatcher.DispatchToOutbound(ctx)

	requestDone := signal.ExecuteAsync(func() error {
		defer inboundRay.InboundInput().Close()

		timedReader := net.NewTimeOutReader(d.config.Timeout, conn)
		chunkReader := buf.NewReader(timedReader)

		if err := buf.Pipe(chunkReader, inboundRay.InboundInput()); err != nil {
			log.Info("Dokodemo: Failed to transport request: ", err)
			return err
		}

		return nil
	})

	responseDone := signal.ExecuteAsync(func() error {
		v2writer := buf.NewWriter(conn)

		if err := buf.PipeUntilEOF(inboundRay.InboundOutput(), v2writer); err != nil {
			log.Info("Dokodemo: Failed to transport response: ", err)
			return err
		}
		return nil
	})

	if err := signal.ErrorOrFinish2(ctx, requestDone, responseDone); err != nil {
		inboundRay.InboundInput().CloseError()
		inboundRay.InboundOutput().CloseError()
		log.Info("Dokodemo: Connection ends with ", err)
		return err
	}

	return nil
}

func init() {
	common.Must(common.RegisterConfig((*Config)(nil), func(ctx context.Context, config interface{}) (interface{}, error) {
		return New(ctx, config.(*Config))
	}))
}
