package socks

import (
	"v2ray.com/core/common"
	"v2ray.com/core/common/serial"
	"v2ray.com/core/proxy"
)

func init() {
	common.Must(proxy.RegisterOutboundHandlerCreator(serial.GetMessageType((*ClientConfig)(nil)), new(ClientFactory)))
	common.Must(proxy.RegisterInboundHandlerCreator(serial.GetMessageType((*ServerConfig)(nil)), new(ServerFactory)))
}
