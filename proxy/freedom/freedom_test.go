package freedom_test

import (
	"testing"

	"v2ray.com/core/app"
	"v2ray.com/core/app/dispatcher"
	dispatchers "v2ray.com/core/app/dispatcher/impl"
	"v2ray.com/core/app/dns"
	dnsserver "v2ray.com/core/app/dns/server"
	"v2ray.com/core/app/proxyman"
	"v2ray.com/core/app/proxyman/outbound"
	"v2ray.com/core/app/router"
	"v2ray.com/core/common/buf"
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/proxy"
	. "v2ray.com/core/proxy/freedom"
	"v2ray.com/core/testing/assert"
	"v2ray.com/core/testing/servers/tcp"
	"v2ray.com/core/transport/internet"
	"v2ray.com/core/transport/ray"
)

func TestSinglePacket(t *testing.T) {
	assert := assert.On(t)

	tcpServer := &tcp.Server{
		MsgProcessor: func(data []byte) []byte {
			buffer := make([]byte, 0, 2048)
			buffer = append(buffer, []byte("Processed: ")...)
			buffer = append(buffer, data...)
			return buffer
		},
	}
	_, err := tcpServer.Start()
	assert.Error(err).IsNil()

	space := app.NewSpace()
	freedom := NewFreedomConnection(
		&Config{},
		space,
		&proxy.OutboundHandlerMeta{
			Address: v2net.AnyIP,
			StreamSettings: &internet.StreamConfig{
				Network: v2net.Network_RawTCP,
			},
		})
	space.Initialize()

	traffic := ray.NewRay()
	data2Send := "Data to be sent to remote"
	payload := buf.NewLocal(2048)
	payload.Append([]byte(data2Send))

	go freedom.Dispatch(v2net.TCPDestination(v2net.LocalHostIP, tcpServer.Port), payload, traffic)
	traffic.InboundInput().Close()

	respPayload, err := traffic.InboundOutput().Read()
	assert.Error(err).IsNil()
	assert.String(respPayload.String()).Equals("Processed: Data to be sent to remote")

	tcpServer.Close()
}

func TestIPResolution(t *testing.T) {
	assert := assert.On(t)

	space := app.NewSpace()
	space.BindApp(proxyman.APP_ID_OUTBOUND_MANAGER, outbound.New())
	space.BindApp(dispatcher.APP_ID, dispatchers.NewDefaultDispatcher(space))
	r := router.NewRouter(&router.Config{}, space)
	space.BindApp(router.APP_ID, r)
	dnsServer := dnsserver.NewCacheServer(space, &dns.Config{
		Hosts: map[string]*v2net.IPOrDomain{
			"v2ray.com": v2net.NewIPOrDomain(v2net.LocalHostIP),
		},
	})
	space.BindApp(dns.APP_ID, dnsServer)

	freedom := NewFreedomConnection(
		&Config{DomainStrategy: Config_USE_IP},
		space,
		&proxy.OutboundHandlerMeta{
			Address: v2net.AnyIP,
			StreamSettings: &internet.StreamConfig{
				Network: v2net.Network_RawTCP,
			},
		})

	space.Initialize()

	ipDest := freedom.ResolveIP(v2net.TCPDestination(v2net.DomainAddress("v2ray.com"), v2net.Port(80)))
	assert.Destination(ipDest).IsTCP()
	assert.Address(ipDest.Address).Equals(v2net.LocalHostIP)
}
