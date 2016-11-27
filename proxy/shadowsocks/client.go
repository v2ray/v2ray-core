package shadowsocks

import (
	"sync"
	"v2ray.com/core/app"
	"v2ray.com/core/common/alloc"
	v2io "v2ray.com/core/common/io"
	"v2ray.com/core/common/log"
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/common/protocol"
	"v2ray.com/core/common/retry"
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

func (v *Client) Dispatch(destination v2net.Destination, payload *alloc.Buffer, ray ray.OutboundRay) {
	defer payload.Release()
	defer ray.OutboundInput().Release()
	defer ray.OutboundOutput().Close()

	network := destination.Network

	var server *protocol.ServerSpec
	var conn internet.Connection

	err := retry.ExponentialBackoff(5, 100).On(func() error {
		server = v.serverPicker.PickServer()
		dest := server.Destination()
		dest.Network = network
		rawConn, err := internet.Dial(v.meta.Address, dest, v.meta.GetDialerOptions())
		if err != nil {
			return err
		}
		conn = rawConn

		return nil
	})
	if err != nil {
		log.Warning("Shadowsocks|Client: Failed to find an available destination:", err)
		return
	}
	log.Info("Shadowsocks|Client: Tunneling request to ", destination, " via ", server.Destination())

	conn.SetReusable(false)

	request := &protocol.RequestHeader{
		Version: Version,
		Address: destination.Address,
		Port:    destination.Port,
	}
	if destination.Network == v2net.Network_TCP {
		request.Command = protocol.RequestCommandTCP
	} else {
		request.Command = protocol.RequestCommandUDP
	}

	user := server.PickUser()
	rawAccount, err := user.GetTypedAccount()
	if err != nil {
		log.Warning("Shadowsocks|Client: Failed to get a valid user account: ", err)
		return
	}
	account := rawAccount.(*ShadowsocksAccount)
	request.User = user

	if account.OneTimeAuth == Account_Auto || account.OneTimeAuth == Account_Enabled {
		request.Option |= RequestOptionOneTimeAuth
	}

	if request.Command == protocol.RequestCommandTCP {
		bufferedWriter := v2io.NewBufferedWriter(conn)
		defer bufferedWriter.Release()

		bodyWriter, err := WriteTCPRequest(request, bufferedWriter)
		defer bodyWriter.Release()

		if err != nil {
			log.Info("Shadowsock|Client: Failed to write request: ", err)
			return
		}

		if err := bodyWriter.Write(payload); err != nil {
			log.Info("Shadowsocks|Client: Failed to write payload: ", err)
			return
		}

		var responseMutex sync.Mutex
		responseMutex.Lock()

		go func() {
			defer responseMutex.Unlock()

			responseReader, err := ReadTCPResponse(user, conn)
			if err != nil {
				log.Warning("Shadowsocks|Client: Failed to read response: " + err.Error())
				return
			}

			if err := v2io.PipeUntilEOF(responseReader, ray.OutboundOutput()); err != nil {
				log.Info("Shadowsocks|Client: Failed to transport all TCP response: ", err)
			}
		}()

		bufferedWriter.SetCached(false)
		if err := v2io.PipeUntilEOF(ray.OutboundInput(), bodyWriter); err != nil {
			log.Info("Shadowsocks|Client: Failed to trasnport all TCP request: ", err)
		}

		responseMutex.Lock()
	}

	if request.Command == protocol.RequestCommandUDP {
		timedReader := v2net.NewTimeOutReader(16, conn)
		var responseMutex sync.Mutex
		responseMutex.Lock()

		go func() {
			defer responseMutex.Unlock()

			reader := &UDPReader{
				Reader: timedReader,
				User:   user,
			}

			if err := v2io.PipeUntilEOF(reader, ray.OutboundOutput()); err != nil {
				log.Info("Shadowsocks|Client: Failed to transport all UDP response: ", err)
			}
		}()

		writer := &UDPWriter{
			Writer:  conn,
			Request: request,
		}
		if !payload.IsEmpty() {
			if err := writer.Write(payload); err != nil {
				log.Info("Shadowsocks|Client: Failed to write payload: ", err)
				return
			}
		}
		if err := v2io.PipeUntilEOF(ray.OutboundInput(), writer); err != nil {
			log.Info("Shadowsocks|Client: Failed to transport all UDP request: ", err)
		}

		responseMutex.Lock()
	}
}

type ClientFactory struct{}

func (v *ClientFactory) StreamCapability() v2net.NetworkList {
	return v2net.NetworkList{
		Network: []v2net.Network{v2net.Network_TCP, v2net.Network_RawTCP},
	}
}

func (v *ClientFactory) Create(space app.Space, rawConfig interface{}, meta *proxy.OutboundHandlerMeta) (proxy.OutboundHandler, error) {
	return NewClient(rawConfig.(*ClientConfig), space, meta)
}
