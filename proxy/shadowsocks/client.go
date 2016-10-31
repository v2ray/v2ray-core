package shadowsocks

import (
	"errors"

	"sync"
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

func (this *Client) Dispatch(destination v2net.Destination, payload *alloc.Buffer, ray ray.OutboundRay) error {
	defer payload.Release()
	defer ray.OutboundInput().Release()
	defer ray.OutboundOutput().Close()

	network := destination.Network

	var server *protocol.ServerSpec
	var conn internet.Connection

	err := retry.Timed(5, 100).On(func() error {
		server = this.serverPicker.PickServer()
		dest := server.Destination()
		dest.Network = network
		rawConn, err := internet.Dial(this.meta.Address, dest, this.meta.StreamSettings)
		if err != nil {
			return err
		}
		conn = rawConn

		return nil
	})
	if err != nil {
		return errors.New("Shadowsocks|Client: Failed to find an available destination:" + err.Error())
	}
	log.Info("Shadowsocks|Client: Tunneling request to ", destination, " via ", server.Destination())

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
		return errors.New("Shadowsocks|Client: Failed to get a valid user account: " + err.Error())
	}
	account := rawAccount.(*ShadowsocksAccount)
	request.User = user

	if account.OneTimeAuth {
		request.Option |= RequestOptionOneTimeAuth
	}

	if request.Command == protocol.RequestCommandTCP {
		bufferedWriter := v2io.NewBufferedWriter(conn)
		defer bufferedWriter.Release()

		bodyWriter, err := WriteTCPRequest(request, bufferedWriter)
		defer bodyWriter.Release()

		if err != nil {
			return errors.New("Shadowsock|Client: Failed to write request: " + err.Error())
		}

		if err := bodyWriter.Write(payload); err != nil {
			return errors.New("Shadowsocks|Client: Failed to write payload: " + err.Error())
		}

		bufferedWriter.SetCached(false)
		v2io.Pipe(ray.OutboundInput(), bodyWriter)

		var responseMutex sync.Mutex
		responseMutex.Lock()

		go func() {
			defer responseMutex.Unlock()

			responseReader, err := ReadTCPResponse(user, conn)
			if err != nil {
				log.Warning("Shadowsocks|Client: Failed to read response: " + err.Error())
				return
			}

			v2io.Pipe(responseReader, ray.OutboundOutput())
		}()

		responseMutex.Lock()
	}

	if request.Command == protocol.RequestCommandUDP {
		writer := &UDPWriter{
			Writer:  conn,
			Request: request,
		}
		if err := writer.Write(payload); err != nil {
			return errors.New("Shadowsocks|Client: Failed to write payload: " + err.Error())
		}
		v2io.Pipe(ray.OutboundInput(), writer)

		timedReader := v2net.NewTimeOutReader(16, conn)
		var responseMutex sync.Mutex
		responseMutex.Lock()

		go func() {
			defer responseMutex.Unlock()

			reader := &UDPReader{
				Reader: timedReader,
				User:   user,
			}

			v2io.Pipe(reader, ray.OutboundOutput())
		}()

		responseMutex.Lock()
	}

	return nil
}
