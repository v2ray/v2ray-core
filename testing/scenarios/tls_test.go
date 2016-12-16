package scenarios

import (
	"v2ray.com/core"
	v2net "v2ray.com/core/common/net"
)

var clientConfig = &core.Config{
	Inbound: []*core.InboundConnectionConfig{
		{
			PortRange: v2net.SinglePortRange(pickPort()),
			ListenOn:  v2net.NewIPOrDomain(v2net.LocalHostIP),
		},
	},
}
