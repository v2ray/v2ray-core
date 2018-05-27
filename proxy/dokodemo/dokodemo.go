package dokodemo

//go:generate go run $GOPATH/src/v2ray.com/core/common/errors/errorgen/main.go -pkg dokodemo -path Proxy,Dokodemo

import (
	"context"
	"time"

	"v2ray.com/core"
	"v2ray.com/core/common"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/functions"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/signal"
	"v2ray.com/core/proxy"
	"v2ray.com/core/transport/internet"
	"v2ray.com/core/transport/internet/udp"
	"v2ray.com/core/transport/pipe"
)

type DokodemoDoor struct {
	policyManager core.PolicyManager
	config        *Config
	address       net.Address
	port          net.Port
}

func New(ctx context.Context, config *Config) (*DokodemoDoor, error) {
	if config.NetworkList == nil || config.NetworkList.Size() == 0 {
		return nil, newError("no network specified")
	}
	v := core.MustFromContext(ctx)
	d := &DokodemoDoor{
		config:        config,
		address:       config.GetPredefinedAddress(),
		port:          net.Port(config.Port),
		policyManager: v.PolicyManager(),
	}

	return d, nil
}

func (d *DokodemoDoor) Network() net.NetworkList {
	return *(d.config.NetworkList)
}

func (d *DokodemoDoor) policy() core.Policy {
	config := d.config
	p := d.policyManager.ForLevel(config.UserLevel)
	if config.Timeout > 0 && config.UserLevel == 0 {
		p.Timeouts.ConnectionIdle = time.Duration(config.Timeout) * time.Second
	}
	return p
}

func (d *DokodemoDoor) Process(ctx context.Context, network net.Network, conn internet.Connection, dispatcher core.Dispatcher) error {
	newError("processing connection from: ", conn.RemoteAddr()).AtDebug().WithContext(ctx).WriteToLog()
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

	plcy := d.policy()
	ctx, cancel := context.WithCancel(ctx)
	timer := signal.CancelAfterInactivity(ctx, cancel, plcy.Timeouts.ConnectionIdle)

	ctx = core.ContextWithBufferPolicy(ctx, plcy.Buffer)
	link, err := dispatcher.Dispatch(ctx, dest)
	if err != nil {
		return newError("failed to dispatch request").Base(err)
	}

	requestDone := func() error {
		defer timer.SetTimeout(plcy.Timeouts.DownlinkOnly)

		chunkReader := buf.NewReader(conn)

		if err := buf.Copy(chunkReader, link.Writer, buf.UpdateActivity(timer)); err != nil {
			return newError("failed to transport request").Base(err)
		}

		return nil
	}

	responseDone := func() error {
		defer timer.SetTimeout(plcy.Timeouts.UplinkOnly)

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

		if err := buf.Copy(link.Reader, writer, buf.UpdateActivity(timer)); err != nil {
			return newError("failed to transport response").Base(err)
		}

		return nil
	}

	if err := signal.ExecuteParallel(ctx, functions.OnSuccess(requestDone, functions.Close(link.Writer)), responseDone); err != nil {
		pipe.CloseError(link.Reader)
		pipe.CloseError(link.Writer)
		return newError("connection ends").Base(err)
	}

	return nil
}

func init() {
	common.Must(common.RegisterConfig((*Config)(nil), func(ctx context.Context, config interface{}) (interface{}, error) {
		return New(ctx, config.(*Config))
	}))
}
