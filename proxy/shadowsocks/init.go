package shadowsocks

import (
	"v2ray.com/core/common/serial"
	"v2ray.com/core/proxy"
)

func init() {
	// Must happen after config is initialized
	proxy.MustRegisterOutboundHandlerCreator(serial.GetMessageType(new(ClientConfig)), new(ClientFactory))
	proxy.MustRegisterInboundHandlerCreator(serial.GetMessageType(new(ServerConfig)), new(ServerFactory))
}
