package socks

import (
	"context"
	"time"

	"v2ray.com/core"
	"v2ray.com/core/common"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/protocol"
	"v2ray.com/core/common/retry"
	"v2ray.com/core/common/signal"
	"v2ray.com/core/proxy"
	"v2ray.com/core/transport/internet"
	"v2ray.com/core/transport/ray"
)

// Client is a Socks5 client.
type Client struct {
	serverPicker  protocol.ServerPicker
	policyManager core.PolicyManager
}

// NewClient create a new Socks5 client based on the given config.
func NewClient(ctx context.Context, config *ClientConfig) (*Client, error) {
	serverList := protocol.NewServerList()
	for _, rec := range config.Server {
		serverList.AddServer(protocol.NewServerSpecFromPB(*rec))
	}
	if serverList.Size() == 0 {
		return nil, newError("0 target server")
	}

	v := core.MustFromContext(ctx)
	return &Client{
		serverPicker:  protocol.NewRoundRobinServerPicker(serverList),
		policyManager: v.PolicyManager(),
	}, nil
}

// Process implements proxy.Outbound.Process.
func (c *Client) Process(ctx context.Context, ray ray.OutboundRay, dialer proxy.Dialer) error {
	destination, ok := proxy.TargetFromContext(ctx)
	if !ok {
		return newError("target not specified.")
	}

	var server *protocol.ServerSpec
	var conn internet.Connection

	err := retry.ExponentialBackoff(5, 100).On(func() error {
		server = c.serverPicker.PickServer()
		dest := server.Destination()
		rawConn, err := dialer.Dial(ctx, dest)
		if err != nil {
			return err
		}
		conn = rawConn

		return nil
	})

	if err != nil {
		return newError("failed to find an available destination").Base(err)
	}

	defer conn.Close()

	p := c.policyManager.ForLevel(0)

	request := &protocol.RequestHeader{
		Version: socks5Version,
		Command: protocol.RequestCommandTCP,
		Address: destination.Address,
		Port:    destination.Port,
	}
	if destination.Network == net.Network_UDP {
		request.Command = protocol.RequestCommandUDP
	}

	user := server.PickUser()
	if user != nil {
		request.User = user
		p = c.policyManager.ForLevel(user.Level)
	}

	if err := conn.SetDeadline(time.Now().Add(p.Timeouts.Handshake)); err != nil {
		newError("failed to set deadline for handshake").Base(err).WithContext(ctx).WriteToLog()
	}
	udpRequest, err := ClientHandshake(request, conn, conn)
	if err != nil {
		return newError("failed to establish connection to server").AtWarning().Base(err)
	}

	if err := conn.SetDeadline(time.Time{}); err != nil {
		newError("failed to clear deadline after handshake").Base(err).WithContext(ctx).WriteToLog()
	}

	ctx, cancel := context.WithCancel(ctx)
	timer := signal.CancelAfterInactivity(ctx, cancel, p.Timeouts.ConnectionIdle)

	var requestFunc func() error
	var responseFunc func() error
	if request.Command == protocol.RequestCommandTCP {
		requestFunc = func() error {
			defer timer.SetTimeout(p.Timeouts.DownlinkOnly)
			return buf.Copy(ray.OutboundInput(), buf.NewWriter(conn), buf.UpdateActivity(timer))
		}
		responseFunc = func() error {
			defer timer.SetTimeout(p.Timeouts.UplinkOnly)
			return buf.Copy(buf.NewReader(conn), ray.OutboundOutput(), buf.UpdateActivity(timer))
		}
	} else if request.Command == protocol.RequestCommandUDP {
		udpConn, err := dialer.Dial(ctx, udpRequest.Destination())
		if err != nil {
			return newError("failed to create UDP connection").Base(err)
		}
		defer udpConn.Close()
		requestFunc = func() error {
			defer timer.SetTimeout(p.Timeouts.DownlinkOnly)
			return buf.Copy(ray.OutboundInput(), buf.NewSequentialWriter(NewUDPWriter(request, udpConn)), buf.UpdateActivity(timer))
		}
		responseFunc = func() error {
			defer timer.SetTimeout(p.Timeouts.UplinkOnly)
			reader := &UDPReader{reader: udpConn}
			return buf.Copy(reader, ray.OutboundOutput(), buf.UpdateActivity(timer))
		}
	}

	requestDone := signal.ExecuteAsync(requestFunc)
	responseDone := signal.ExecuteAsync(responseFunc)
	if err := signal.ErrorOrFinish2(ctx, requestDone, responseDone); err != nil {
		return newError("connection ends").Base(err)
	}

	return nil
}

func init() {
	common.Must(common.RegisterConfig((*ClientConfig)(nil), func(ctx context.Context, config interface{}) (interface{}, error) {
		return NewClient(ctx, config.(*ClientConfig))
	}))
}
