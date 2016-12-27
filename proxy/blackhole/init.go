package blackhole

import (
	"v2ray.com/core/common"
	"v2ray.com/core/common/serial"
	"v2ray.com/core/proxy"
)

func init() {
	// Must listed after config.pb.go
	common.Must(proxy.RegisterOutboundHandlerCreator(serial.GetMessageType(new(Config)), new(Factory)))
}
