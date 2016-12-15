package shadowsocks

import (
	"v2ray.com/core/common/serial"
	"v2ray.com/core/proxy/registry"
)

func init() {
	// Must happen after config is initialized
	registry.MustRegisterOutboundHandlerCreator(serial.GetMessageType(new(ClientConfig)), new(ClientFactory))
	registry.MustRegisterInboundHandlerCreator(serial.GetMessageType(new(ServerConfig)), new(ServerFactory))
}
