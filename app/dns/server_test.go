package dns_test

import (
	"net"
	"testing"

	"v2ray.com/core/app"
	"v2ray.com/core/app/dispatcher"
	dispatchers "v2ray.com/core/app/dispatcher/impl"
	. "v2ray.com/core/app/dns"
	"v2ray.com/core/app/proxyman"
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/proxy"
	"v2ray.com/core/proxy/freedom"
	"v2ray.com/core/testing/assert"
	"v2ray.com/core/transport/internet"
)

func TestDnsAdd(t *testing.T) {
	assert := assert.On(t)

	space := app.NewSpace()

	outboundHandlerManager := proxyman.NewDefaultOutboundHandlerManager()
	outboundHandlerManager.SetDefaultHandler(
		freedom.NewFreedomConnection(
			&freedom.Config{},
			space,
			&proxy.OutboundHandlerMeta{
				Address: v2net.AnyIP,
				StreamSettings: &internet.StreamConfig{
					Network: v2net.Network_RawTCP,
				},
			}))
	space.BindApp(proxyman.APP_ID_OUTBOUND_MANAGER, outboundHandlerManager)
	space.BindApp(dispatcher.APP_ID, dispatchers.NewDefaultDispatcher(space))

	domain := "local.v2ray.com"
	server := NewCacheServer(space, &Config{
		NameServers: []*v2net.DestinationPB{{
			Network: v2net.Network_UDP,
			Address: &v2net.AddressPB{
				Address: &v2net.AddressPB_Ip{
					Ip: []byte{8, 8, 8, 8},
				},
			},
			Port: 53,
		}},
	})
	space.BindApp(APP_ID, server)
	space.Initialize()

	ips := server.Get(domain)
	assert.Int(len(ips)).Equals(1)
	assert.IP(ips[0].To4()).Equals(net.IP([]byte{127, 0, 0, 1}))
}
