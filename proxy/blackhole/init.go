package blackhole

import (
	"v2ray.com/core/common/serial"
	"v2ray.com/core/proxy"
)

func init() {
	// Must listed after config.pb.go
	proxy.MustRegisterOutboundHandlerCreator(serial.GetMessageType(new(Config)), new(Factory))
}
