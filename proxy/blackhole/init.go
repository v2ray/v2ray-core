package blackhole

import (
	"v2ray.com/core/common/serial"
	"v2ray.com/core/proxy/registry"
)

func init() {
	// Must listed after config.pb.go
	registry.MustRegisterOutboundHandlerCreator(serial.GetMessageType(new(Config)), new(Factory))
}
