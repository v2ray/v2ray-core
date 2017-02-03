package dokodemo

import (
	"context"
	"runtime"
	"time"

	"v2ray.com/core/app"
	"v2ray.com/core/app/dispatcher"
	"v2ray.com/core/app/log"
	"v2ray.com/core/common"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/errors"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/signal"
	"v2ray.com/core/proxy"
	"v2ray.com/core/transport/internet"
)

type DokodemoDoor struct {
	config  *Config
	address net.Address
	port    net.Port
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
	return d, nil
}

func (d *DokodemoDoor) Network() net.NetworkList {
	return *(d.config.NetworkList)
}

func (d *DokodemoDoor) Process(ctx context.Context, network net.Network, conn internet.Connection, dispatcher dispatcher.Interface) error {
	log.Debug("Dokodemo: processing connection from: ", conn.RemoteAddr())
	conn.SetReusable(false)
	dest := net.Destination{
		Network: network,
		Address: d.address,
		Port:    d.port,
	}
	if d.config.FollowRedirect {
		if origDest := proxy.OriginalDestinationFromContext(ctx); origDest.IsValid() {
			dest = origDest
		}
	}
	if !dest.IsValid() || dest.Address == nil {
		log.Info("Dokodemo: Invalid destination. Discarding...")
		return errors.New("Dokodemo: Unable to get destination.")
	}
	ctx, cancel := context.WithCancel(ctx)
	timeout := time.Second * time.Duration(d.config.Timeout)
	if timeout == 0 {
		timeout = time.Minute * 2
	}
	timer := signal.CancelAfterInactivity(ctx, cancel, timeout)

	inboundRay, err := dispatcher.Dispatch(ctx, dest)
	if err != nil {
		return err
	}

	requestDone := signal.ExecuteAsync(func() error {
		defer inboundRay.InboundInput().Close()

		chunkReader := buf.NewReader(conn)

		if err := buf.PipeUntilEOF(timer, chunkReader, inboundRay.InboundInput()); err != nil {
			log.Info("Dokodemo: Failed to transport request: ", err)
			return err
		}

		return nil
	})

	responseDone := signal.ExecuteAsync(func() error {
		v2writer := buf.NewWriter(conn)

		if err := buf.PipeUntilEOF(timer, inboundRay.InboundOutput(), v2writer); err != nil {
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

	runtime.KeepAlive(timer)

	return nil
}

func init() {
	common.Must(common.RegisterConfig((*Config)(nil), func(ctx context.Context, config interface{}) (interface{}, error) {
		return New(ctx, config.(*Config))
	}))
}
