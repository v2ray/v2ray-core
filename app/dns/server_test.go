package dns_test

import (
	"net"
	"testing"

	"github.com/v2ray/v2ray-core/app"
	"github.com/v2ray/v2ray-core/app/dispatcher"
	. "github.com/v2ray/v2ray-core/app/dns"
	apptesting "github.com/v2ray/v2ray-core/app/testing"
	v2net "github.com/v2ray/v2ray-core/common/net"
	netassert "github.com/v2ray/v2ray-core/common/net/testing/assert"
	"github.com/v2ray/v2ray-core/proxy/freedom"
	v2testing "github.com/v2ray/v2ray-core/testing"
	"github.com/v2ray/v2ray-core/testing/assert"
	"github.com/v2ray/v2ray-core/transport/ray"
)

type TestDispatcher struct {
	freedom *freedom.FreedomConnection
}

func (this *TestDispatcher) DispatchToOutbound(context app.Context, dest v2net.Destination) ray.InboundRay {
	direct := ray.NewRay()

	go func() {
		payload, err := direct.OutboundInput().Read()
		if err != nil {
			direct.OutboundInput().Release()
			direct.OutboundOutput().Release()
			return
		}
		this.freedom.Dispatch(dest, payload, direct)
	}()
	return direct
}

func TestDnsAdd(t *testing.T) {
	v2testing.Current(t)

	d := &TestDispatcher{
		freedom: &freedom.FreedomConnection{},
	}
	spaceController := app.NewController()
	spaceController.Bind(dispatcher.APP_ID, d)
	space := spaceController.ForContext("test")

	domain := "local.v2ray.com"
	server := NewCacheServer(space, &Config{
		NameServers: []v2net.Destination{
			v2net.UDPDestination(v2net.IPAddress([]byte{8, 8, 8, 8}), v2net.Port(53)),
		},
	})
	ips := server.Get(&apptesting.Context{
		CallerTagValue: "a",
	}, domain)
	assert.Int(len(ips)).Equals(1)
	netassert.IP(ips[0].To4()).Equals(net.IP([]byte{127, 0, 0, 1}))
}
