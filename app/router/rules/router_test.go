package rules_test

import (
	"testing"

	"github.com/v2ray/v2ray-core/app"
	"github.com/v2ray/v2ray-core/app/dispatcher"
	dispatchers "github.com/v2ray/v2ray-core/app/dispatcher/impl"
	"github.com/v2ray/v2ray-core/app/dns"
	"github.com/v2ray/v2ray-core/app/proxyman"
	"github.com/v2ray/v2ray-core/app/router"
	. "github.com/v2ray/v2ray-core/app/router/rules"
	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/testing/assert"
)

func TestSimpleRouter(t *testing.T) {
	assert := assert.On(t)

	config := &RouterRuleConfig{
		Rules: []*Rule{
			{
				Tag:       "test",
				Condition: NewNetworkMatcher(v2net.Network("tcp").AsList()),
			},
		},
	}

	space := app.NewSpace()
	space.BindApp(dns.APP_ID, dns.NewCacheServer(space, &dns.Config{}))
	space.BindApp(dispatcher.APP_ID, dispatchers.NewDefaultDispatcher(space))
	space.BindApp(proxyman.APP_ID_OUTBOUND_MANAGER, proxyman.NewDefaultOutboundHandlerManager())
	r := NewRouter(config, space)
	space.BindApp(router.APP_ID, r)
	assert.Error(space.Initialize()).IsNil()

	tag, err := r.TakeDetour(v2net.TCPDestination(v2net.DomainAddress("v2ray.com"), 80))
	assert.Error(err).IsNil()
	assert.String(tag).Equals("test")
}
