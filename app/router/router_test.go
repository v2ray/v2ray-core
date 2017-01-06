package router_test

import (
	"testing"

	"v2ray.com/core/app"
	"v2ray.com/core/app/dispatcher"
	_ "v2ray.com/core/app/dispatcher/impl"
	"v2ray.com/core/app/dns"
	_ "v2ray.com/core/app/dns/server"
	"v2ray.com/core/app/proxyman"
	_ "v2ray.com/core/app/proxyman/outbound"
	. "v2ray.com/core/app/router"
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/proxy"
	"v2ray.com/core/testing/assert"
)

func TestSimpleRouter(t *testing.T) {
	assert := assert.On(t)

	config := &Config{
		Rule: []*RoutingRule{
			{
				Tag: "test",
				NetworkList: &v2net.NetworkList{
					Network: []v2net.Network{v2net.Network_TCP},
				},
			},
		},
	}

	space := app.NewSpace()
	assert.Error(space.AddApp(new(dns.Config))).IsNil()
	assert.Error(space.AddApp(new(dispatcher.Config))).IsNil()
	assert.Error(space.AddApp(new(proxyman.OutboundConfig))).IsNil()
	assert.Error(space.AddApp(config)).IsNil()
	assert.Error(space.Initialize()).IsNil()

	r := FromSpace(space)

	tag, err := r.TakeDetour(&proxy.SessionInfo{Destination: v2net.TCPDestination(v2net.DomainAddress("v2ray.com"), 80)})
	assert.Error(err).IsNil()
	assert.String(tag).Equals("test")
}
