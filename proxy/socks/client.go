package socks

import (
	"context"
	"runtime"
	"time"

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
	serverPicker protocol.ServerPicker
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

	return &Client{
		serverPicker: protocol.NewRoundRobinServerPicker(serverList),
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
	}

	udpRequest, err := ClientHandshake(request, conn, conn)
	if err != nil {
		return newError("failed to establish connection to server").AtWarning().Base(err)
	}

	ctx, timer := signal.CancelAfterInactivity(ctx, time.Minute*2)

	var requestFunc func() error
	var responseFunc func() error
	if request.Command == protocol.RequestCommandTCP {
		requestFunc = func() error {
			return buf.Copy(timer, ray.OutboundInput(), buf.NewWriter(conn))
		}
		responseFunc = func() error {
			defer ray.OutboundOutput().Close()
			return buf.Copy(timer, buf.NewReader(conn), ray.OutboundOutput())
		}
	} else if request.Command == protocol.RequestCommandUDP {
		udpConn, err := dialer.Dial(ctx, udpRequest.Destination())
		if err != nil {
			return newError("failed to create UDP connection").Base(err)
		}
		defer udpConn.Close()
		requestFunc = func() error {
			return buf.Copy(timer, ray.OutboundInput(), buf.NewSequentialWriter(NewUDPWriter(request, udpConn)))
		}
		responseFunc = func() error {
			defer ray.OutboundOutput().Close()
			reader := &UDPReader{reader: udpConn}
			return buf.Copy(timer, reader, ray.OutboundOutput())
		}
	}

	requestDone := signal.ExecuteAsync(requestFunc)
	responseDone := signal.ExecuteAsync(responseFunc)
	if err := signal.ErrorOrFinish2(ctx, requestDone, responseDone); err != nil {
		return newError("connection ends").Base(err)
	}

	runtime.KeepAlive(timer)

	return nil
}

func init() {
	common.Must(common.RegisterConfig((*ClientConfig)(nil), func(ctx context.Context, config interface{}) (interface{}, error) {
		return NewClient(ctx, config.(*ClientConfig))
	}))
}
