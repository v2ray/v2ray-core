package freedom_test

import (
	"net"
	"testing"

	"github.com/v2ray/v2ray-core/app"
	"github.com/v2ray/v2ray-core/app/dispatcher"
	dispatchers "github.com/v2ray/v2ray-core/app/dispatcher/impl"
	"github.com/v2ray/v2ray-core/app/dns"
	"github.com/v2ray/v2ray-core/app/proxyman"
	"github.com/v2ray/v2ray-core/app/router"
	"github.com/v2ray/v2ray-core/app/router/rules"
	"github.com/v2ray/v2ray-core/common/alloc"
	v2net "github.com/v2ray/v2ray-core/common/net"
	v2nettesting "github.com/v2ray/v2ray-core/common/net/testing"
	netassert "github.com/v2ray/v2ray-core/common/net/testing/assert"
	. "github.com/v2ray/v2ray-core/proxy/freedom"
	v2testing "github.com/v2ray/v2ray-core/testing"
	"github.com/v2ray/v2ray-core/testing/assert"
	"github.com/v2ray/v2ray-core/testing/servers/tcp"
	"github.com/v2ray/v2ray-core/transport/ray"
)

func TestSinglePacket(t *testing.T) {
	v2testing.Current(t)
	port := v2nettesting.PickPort()

	tcpServer := &tcp.Server{
		Port: port,
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
	freedom := NewFreedomConnection(&Config{}, space)
	space.Initialize()

	traffic := ray.NewRay()
	data2Send := "Data to be sent to remote"
	payload := alloc.NewSmallBuffer().Clear().Append([]byte(data2Send))

	go freedom.Dispatch(v2net.TCPDestination(v2net.LocalHostIP, port), payload, traffic)
	traffic.InboundInput().Close()

	respPayload, err := traffic.InboundOutput().Read()
	assert.Error(err).IsNil()
	assert.Bytes(respPayload.Value).Equals([]byte("Processed: Data to be sent to remote"))

	tcpServer.Close()
}

func TestUnreachableDestination(t *testing.T) {
	v2testing.Current(t)

	freedom := NewFreedomConnection(&Config{}, app.NewSpace())
	traffic := ray.NewRay()
	data2Send := "Data to be sent to remote"
	payload := alloc.NewSmallBuffer().Clear().Append([]byte(data2Send))

	err := freedom.Dispatch(v2net.TCPDestination(v2net.IPAddress([]byte{127, 0, 0, 1}), 128), payload, traffic)
	assert.Error(err).IsNotNil()
}

func TestIPResolution(t *testing.T) {
	v2testing.Current(t)

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

	freedom := NewFreedomConnection(&Config{DomainStrategy: DomainStrategyUseIP}, space)

	space.Initialize()

	ipDest := freedom.ResolveIP(v2net.TCPDestination(v2net.DomainAddress("v2ray.com"), v2net.Port(80)))
	netassert.Destination(ipDest).IsTCP()
	netassert.Address(ipDest.Address()).Equals(v2net.LocalHostIP)
}
