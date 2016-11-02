package shadowsocks

import (
	"v2ray.com/core/common/loader"
	"v2ray.com/core/proxy/registry"
)

func init() {
	// Must happen after config is initialized
	registry.MustRegisterOutboundHandlerCreator(loader.GetType(new(ClientConfig)), new(ClientFactory))
	registry.MustRegisterInboundHandlerCreator(loader.GetType(new(ServerConfig)), new(ServerFactory))
}
