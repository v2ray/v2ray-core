package shadowsocks

import (
	"v2ray.com/core/common"
	"v2ray.com/core/common/serial"
	"v2ray.com/core/proxy"
)

func init() {
	// Must happen after config is initialized
	common.Must(proxy.RegisterOutboundHandlerCreator(serial.GetMessageType(new(ClientConfig)), new(ClientFactory)))
	common.Must(proxy.RegisterInboundHandlerCreator(serial.GetMessageType(new(ServerConfig)), new(ServerFactory)))
}
