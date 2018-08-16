package router_test

import (
	"context"
	"testing"

	"v2ray.com/core"
	"v2ray.com/core/app/dispatcher"
	"v2ray.com/core/app/proxyman"
	_ "v2ray.com/core/app/proxyman/outbound"
	. "v2ray.com/core/app/router"
	"v2ray.com/core/common"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/serial"
	"v2ray.com/core/proxy"
	. "v2ray.com/ext/assert"
)

func TestSimpleRouter(t *testing.T) {
	assert := With(t)

	config := &core.Config{
		App: []*serial.TypedMessage{
			serial.ToTypedMessage(&Config{
				Rule: []*RoutingRule{
					{
						Tag: "test",
						NetworkList: &net.NetworkList{
							Network: []net.Network{net.Network_TCP},
						},
					},
				},
			}),
			serial.ToTypedMessage(&dispatcher.Config{}),
			serial.ToTypedMessage(&proxyman.OutboundConfig{}),
		},
	}

	v, err := core.New(config)
	common.Must(err)

	r := v.Router()

	ctx := proxy.ContextWithTarget(context.Background(), net.TCPDestination(net.DomainAddress("v2ray.com"), 80))
	tag, err := r.PickRoute(ctx)
	assert(err, IsNil)
	assert(tag, Equals, "test")
}
