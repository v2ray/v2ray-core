package blackhole

import (
	"v2ray.com/core/common/loader"
	"v2ray.com/core/proxy/registry"
)

func init() {
	// Must listed after config.pb.go
	registry.MustRegisterOutboundHandlerCreator(loader.GetType(new(Config)), new(Factory))
}
