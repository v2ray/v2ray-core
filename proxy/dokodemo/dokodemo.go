package dokodemo

//go:generate errorgen

import (
	"context"
	"time"

	"v2ray.com/core"
	"v2ray.com/core/common"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/session"
	"v2ray.com/core/common/signal"
	"v2ray.com/core/common/task"
	"v2ray.com/core/transport/internet"
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

type hasHandshakeAddress interface {
	HandshakeAddress() net.Address
}

func (d *DokodemoDoor) Process(ctx context.Context, network net.Network, conn internet.Connection, dispatcher core.Dispatcher) error {
	newError("processing connection from: ", conn.RemoteAddr()).AtDebug().WriteToLog(session.ExportIDToError(ctx))
	dest := net.Destination{
		Network: network,
		Address: d.address,
		Port:    d.port,
	}
	if d.config.FollowRedirect {
		if outbound := session.OutboundFromContext(ctx); outbound != nil && outbound.Target.IsValid() {
			dest = outbound.Target
		} else if handshake, ok := conn.(hasHandshakeAddress); ok {
			addr := handshake.HandshakeAddress()
			if addr != nil {
				dest.Address = addr
			}
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

		reader := buf.NewReader(conn)
		if err := buf.Copy(reader, link.Writer, buf.UpdateActivity(timer)); err != nil {
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
				writer = &buf.SequentialWriter{Writer: conn}
			} else {
				tCtx := internet.ContextWithBindAddress(context.Background(), dest)
				tCtx = internet.ContextWithStreamSettings(tCtx, &internet.MemoryStreamConfig{
					ProtocolName: "udp",
					SocketSettings: &internet.SocketConfig{
						Tproxy: internet.SocketConfig_TProxy,
					},
				})
				tConn, err := internet.DialSystem(tCtx, net.DestinationFromAddr(conn.RemoteAddr()))
				if err != nil {
					return err
				}
				writer = &buf.SequentialWriter{Writer: tConn}
			}
		}

		if err := buf.Copy(link.Reader, writer, buf.UpdateActivity(timer)); err != nil {
			return newError("failed to transport response").Base(err)
		}

		return nil
	}

	if err := task.Run(task.WithContext(ctx),
		task.Parallel(
			task.Single(requestDone, task.OnSuccess(task.Close(link.Writer))),
			responseDone))(); err != nil {
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
