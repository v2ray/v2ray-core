package shadowsocks

import (
	"v2ray.com/core/common/alloc"
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
		log.Error("Shadowsocks|Client: Failed to find an available destination:", err)
		return err
	}
	log.Info("Shadowsocks|Client: Tunneling request to ", destination, " via ", server.Destination())

	return nil
}
