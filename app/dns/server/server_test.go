package server_test

import (
	"context"
	"testing"

	"v2ray.com/core/app"
	"v2ray.com/core/app/dispatcher"
	_ "v2ray.com/core/app/dispatcher/impl"
	. "v2ray.com/core/app/dns"
	_ "v2ray.com/core/app/dns/server"
	"v2ray.com/core/app/proxyman"
	_ "v2ray.com/core/app/proxyman/outbound"
	"v2ray.com/core/common"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/serial"
	"v2ray.com/core/proxy/freedom"
	"v2ray.com/core/testing/servers/udp"
	. "v2ray.com/ext/assert"

	"github.com/miekg/dns"
)

type staticHandler struct {
}

func (*staticHandler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	ans := new(dns.Msg)
	ans.Id = r.Id
	for _, q := range r.Question {
		if q.Name == "google.com." && q.Qtype == dns.TypeA {
			rr, _ := dns.NewRR("google.com. IN A 8.8.8.8")
			ans.Answer = append(ans.Answer, rr)
		} else if q.Name == "facebook.com." && q.Qtype == dns.TypeA {
			rr, _ := dns.NewRR("facebook.com. IN A 9.9.9.9")
			ans.Answer = append(ans.Answer, rr)
		}
	}
	w.WriteMsg(ans)
}

func TestUDPServer(t *testing.T) {
	assert := With(t)

	port := udp.PickPort()

	dnsServer := dns.Server{
		Addr:    "127.0.0.1:" + port.String(),
		Net:     "udp",
		Handler: &staticHandler{},
		UDPSize: 1200,
	}

	go dnsServer.ListenAndServe()

	config := &Config{
		NameServers: []*net.Endpoint{
			{
				Network: net.Network_UDP,
				Address: &net.IPOrDomain{
					Address: &net.IPOrDomain_Ip{
						Ip: []byte{127, 0, 0, 1},
					},
				},
				Port: uint32(port),
			},
		},
	}

	ctx := context.Background()
	space := app.NewSpace()

	ctx = app.ContextWithSpace(ctx, space)
	common.Must(app.AddApplicationToSpace(ctx, config))
	common.Must(app.AddApplicationToSpace(ctx, &dispatcher.Config{}))
	common.Must(app.AddApplicationToSpace(ctx, &proxyman.OutboundConfig{}))

	om := proxyman.OutboundHandlerManagerFromSpace(space)
	om.AddHandler(ctx, &proxyman.OutboundHandlerConfig{
		ProxySettings: serial.ToTypedMessage(&freedom.Config{}),
	})

	common.Must(space.Initialize())
	common.Must(space.Start())

	server := FromSpace(space)
	assert(server, IsNotNil)

	ips := server.Get("google.com")
	assert(len(ips), Equals, 1)
	assert([]byte(ips[0]), Equals, []byte{8, 8, 8, 8})

	ips = server.Get("facebook.com")
	assert(len(ips), Equals, 1)
	assert([]byte(ips[0]), Equals, []byte{9, 9, 9, 9})

	dnsServer.Shutdown()

	ips = server.Get("google.com")
	assert(len(ips), Equals, 1)
	assert([]byte(ips[0]), Equals, []byte{8, 8, 8, 8})
}
