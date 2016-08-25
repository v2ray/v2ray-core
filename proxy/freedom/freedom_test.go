package freedom_test

import (
	"net"
	"testing"

	"v2ray.com/core/app"
	"v2ray.com/core/app/dispatcher"
	dispatchers "v2ray.com/core/app/dispatcher/impl"
	"v2ray.com/core/app/dns"
	"v2ray.com/core/app/proxyman"
	"v2ray.com/core/app/router"
	"v2ray.com/core/app/router/rules"
	"v2ray.com/core/common/alloc"
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
			StreamSettings: &internet.StreamSettings{
				Type: internet.StreamConnectionTypeRawTCP,
			},
		})
	space.Initialize()

	traffic := ray.NewRay()
	data2Send := "Data to be sent to remote"
	payload := alloc.NewLocalBuffer(2048).Clear().Append([]byte(data2Send))

	go freedom.Dispatch(v2net.TCPDestination(v2net.LocalHostIP, tcpServer.Port), payload, traffic)
	traffic.InboundInput().Close()

	respPayload, err := traffic.InboundOutput().Read()
	assert.Error(err).IsNil()
	assert.Bytes(respPayload.Value).Equals([]byte("Processed: Data to be sent to remote"))

	tcpServer.Close()
}

func TestUnreachableDestination(t *testing.T) {
	assert := assert.On(t)

	freedom := NewFreedomConnection(
		&Config{},
		app.NewSpace(),
		&proxy.OutboundHandlerMeta{
			Address: v2net.AnyIP,
			StreamSettings: &internet.StreamSettings{
				Type: internet.StreamConnectionTypeRawTCP,
			},
		})
	traffic := ray.NewRay()
	data2Send := "Data to be sent to remote"
	payload := alloc.NewLocalBuffer(2048).Clear().Append([]byte(data2Send))

	err := freedom.Dispatch(v2net.TCPDestination(v2net.IPAddress([]byte{127, 0, 0, 1}), 128), payload, traffic)
	assert.Error(err).IsNotNil()
}

func TestIPResolution(t *testing.T) {
	assert := assert.On(t)

	space := app.NewSpace()
	space.BindApp(proxyman.APP_ID_OUTBOUND_MANAGER, proxyman.NewDefaultOutboundHandlerManager())
	space.BindApp(dispatcher.APP_ID, dispatchers.NewDefaultDispatcher(space))
	r, _ := router.CreateRouter("rules", &rules.RouterRuleConfig{}, space)
	space.BindApp(router.APP_ID, r)
	dnsServer := dns.NewCacheServer(space, &dns.Config{
		Hosts: map[string]net.IP{
			"v2ray.com": net.IP([]byte{127, 0, 0, 1}),
		},
	})
	space.BindApp(dns.APP_ID, dnsServer)

	freedom := NewFreedomConnection(
		&Config{DomainStrategy: Config_USE_IP},
		space,
		&proxy.OutboundHandlerMeta{
			Address: v2net.AnyIP,
			StreamSettings: &internet.StreamSettings{
				Type: internet.StreamConnectionTypeRawTCP,
			},
		})

	space.Initialize()

	ipDest := freedom.ResolveIP(v2net.TCPDestination(v2net.DomainAddress("v2ray.com"), v2net.Port(80)))
	assert.Destination(ipDest).IsTCP()
	assert.Address(ipDest.Address()).Equals(v2net.LocalHostIP)
}
