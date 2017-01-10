package socks

import (
	"v2ray.com/core/app"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/log"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/protocol"
	"v2ray.com/core/common/retry"
	"v2ray.com/core/common/signal"
	"v2ray.com/core/proxy"
	"v2ray.com/core/transport/internet"
	"v2ray.com/core/transport/ray"
)

type Client struct {
	serverPicker protocol.ServerPicker
	meta         *proxy.OutboundHandlerMeta
}

func NewClient(config *ClientConfig, space app.Space, meta *proxy.OutboundHandlerMeta) (*Client, error) {
	serverList := protocol.NewServerList()
	for _, rec := range config.Server {
		serverList.AddServer(protocol.NewServerSpecFromPB(*rec))
	}
	client := &Client{
		serverPicker: protocol.NewRoundRobinServerPicker(serverList),
		meta:         meta,
	}

	return client, nil
}

func (c *Client) Dispatch(destination net.Destination, ray ray.OutboundRay) {
	var server *protocol.ServerSpec
	var conn internet.Connection

	err := retry.ExponentialBackoff(5, 100).On(func() error {
		server = c.serverPicker.PickServer()
		dest := server.Destination()
		rawConn, err := internet.Dial(c.meta.Address, dest, c.meta.GetDialerOptions())
		if err != nil {
			return err
		}
		conn = rawConn

		return nil
	})

	if err != nil {
		log.Warning("Socks|Client: Failed to find an available destination.")
		return
	}

	defer conn.Close()
	conn.SetReusable(false)

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
		log.Warning("Socks|Client: Failed to establish connection to server: ", err)
		return
	}

	var requestFunc func() error
	var responseFunc func() error
	if request.Command == protocol.RequestCommandTCP {
		requestFunc = func() error {
			return buf.PipeUntilEOF(ray.OutboundInput(), buf.NewWriter(conn))
		}
		responseFunc = func() error {
			defer ray.OutboundOutput().Close()
			return buf.PipeUntilEOF(buf.NewReader(conn), ray.OutboundOutput())
		}
	} else if request.Command == protocol.RequestCommandUDP {
		udpConn, err := internet.Dial(c.meta.Address, udpRequest.Destination(), c.meta.GetDialerOptions())
		if err != nil {
			log.Info("Socks|Client: Failed to create UDP connection: ", err)
			return
		}
		defer udpConn.Close()
		requestFunc = func() error {
			return buf.PipeUntilEOF(ray.OutboundInput(), &UDPWriter{request: request, writer: udpConn})
		}
		responseFunc = func() error {
			defer ray.OutboundOutput().Close()
			reader := &UDPReader{reader: net.NewTimeOutReader(16, udpConn)}
			return buf.PipeUntilEOF(reader, ray.OutboundOutput())
		}
	}

	requestDone := signal.ExecuteAsync(requestFunc)
	responseDone := signal.ExecuteAsync(responseFunc)
	if err := signal.ErrorOrFinish2(requestDone, responseDone); err != nil {
		log.Info("Socks|Client: Connection ends with ", err)
		ray.OutboundInput().CloseError()
		ray.OutboundOutput().CloseError()
	}
}

type ClientFactory struct{}

func (ClientFactory) StreamCapability() net.NetworkList {
	return net.NetworkList{
		Network: []net.Network{net.Network_TCP},
	}
}

func (ClientFactory) Create(space app.Space, rawConfig interface{}, meta *proxy.OutboundHandlerMeta) (proxy.OutboundHandler, error) {
	return NewClient(rawConfig.(*ClientConfig), space, meta)
}
