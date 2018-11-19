package dns_test

import (
	"runtime"
	"testing"
	"time"

	"v2ray.com/core"
	"v2ray.com/core/app/dispatcher"
	. "v2ray.com/core/app/dns"
	"v2ray.com/core/app/policy"
	"v2ray.com/core/app/proxyman"
	_ "v2ray.com/core/app/proxyman/outbound"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/serial"
	feature_dns "v2ray.com/core/features/dns"
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

	var clientIP net.IP

	opt := r.IsEdns0()
	if opt != nil {
		for _, o := range opt.Option {
			if o.Option() == dns.EDNS0SUBNET {
				subnet := o.(*dns.EDNS0_SUBNET)
				clientIP = subnet.Address
			}
		}
	}

	for _, q := range r.Question {
		if q.Name == "google.com." && q.Qtype == dns.TypeA {
			if clientIP == nil {
				rr, _ := dns.NewRR("google.com. IN A 8.8.8.8")
				ans.Answer = append(ans.Answer, rr)
			} else {
				rr, _ := dns.NewRR("google.com. IN A 8.8.4.4")
				ans.Answer = append(ans.Answer, rr)
			}
		} else if q.Name == "facebook.com." && q.Qtype == dns.TypeA {
			rr, _ := dns.NewRR("facebook.com. IN A 9.9.9.9")
			ans.Answer = append(ans.Answer, rr)
		}
	}
	w.WriteMsg(ans)
}

func TestUDPServerSubnet(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("doesn't work on Windows due to miekg/dns changes.")
	}
	assert := With(t)

	port := udp.PickPort()

	dnsServer := dns.Server{
		Addr:    "127.0.0.1:" + port.String(),
		Net:     "udp",
		Handler: &staticHandler{},
		UDPSize: 1200,
	}

	go dnsServer.ListenAndServe()
	time.Sleep(time.Second)

	config := &core.Config{
		App: []*serial.TypedMessage{
			serial.ToTypedMessage(&Config{
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
				ClientIp: []byte{7, 8, 9, 10},
			}),
			serial.ToTypedMessage(&dispatcher.Config{}),
			serial.ToTypedMessage(&proxyman.OutboundConfig{}),
			serial.ToTypedMessage(&policy.Config{}),
		},
		Outbound: []*core.OutboundHandlerConfig{
			{
				ProxySettings: serial.ToTypedMessage(&freedom.Config{}),
			},
		},
	}

	v, err := core.New(config)
	assert(err, IsNil)

	client := v.GetFeature(feature_dns.ClientType()).(feature_dns.Client)

	ips, err := client.LookupIP("google.com")
	assert(err, IsNil)
	assert(len(ips), Equals, 1)
	assert([]byte(ips[0]), Equals, []byte{8, 8, 4, 4})
}

func TestUDPServer(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("doesn't work on Windows due to miekg/dns changes.")
	}
	assert := With(t)

	port := udp.PickPort()

	dnsServer := dns.Server{
		Addr:    "127.0.0.1:" + port.String(),
		Net:     "udp",
		Handler: &staticHandler{},
		UDPSize: 1200,
	}

	go dnsServer.ListenAndServe()
	time.Sleep(time.Second)

	config := &core.Config{
		App: []*serial.TypedMessage{
			serial.ToTypedMessage(&Config{
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
			}),
			serial.ToTypedMessage(&dispatcher.Config{}),
			serial.ToTypedMessage(&proxyman.OutboundConfig{}),
			serial.ToTypedMessage(&policy.Config{}),
		},
		Outbound: []*core.OutboundHandlerConfig{
			{
				ProxySettings: serial.ToTypedMessage(&freedom.Config{}),
			},
		},
	}

	v, err := core.New(config)
	assert(err, IsNil)

	client := v.GetFeature(feature_dns.ClientType()).(feature_dns.Client)

	ips, err := client.LookupIP("google.com")
	assert(err, IsNil)
	assert(len(ips), Equals, 1)
	assert([]byte(ips[0]), Equals, []byte{8, 8, 8, 8})

	ips, err = client.LookupIP("facebook.com")
	assert(err, IsNil)
	assert(len(ips), Equals, 1)
	assert([]byte(ips[0]), Equals, []byte{9, 9, 9, 9})

	dnsServer.Shutdown()

	ips, err = client.LookupIP("google.com")
	assert(err, IsNil)
	assert(len(ips), Equals, 1)
	assert([]byte(ips[0]), Equals, []byte{8, 8, 8, 8})
}

func TestPrioritizedDomain(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("doesn't work on Windows due to miekg/dns changes.")
	}
	assert := With(t)

	port := udp.PickPort()

	dnsServer := dns.Server{
		Addr:    "127.0.0.1:" + port.String(),
		Net:     "udp",
		Handler: &staticHandler{},
		UDPSize: 1200,
	}

	go dnsServer.ListenAndServe()
	time.Sleep(time.Second)

	config := &core.Config{
		App: []*serial.TypedMessage{
			serial.ToTypedMessage(&Config{
				NameServers: []*net.Endpoint{
					{
						Network: net.Network_UDP,
						Address: &net.IPOrDomain{
							Address: &net.IPOrDomain_Ip{
								Ip: []byte{127, 0, 0, 1},
							},
						},
						Port: 9999, /* unreachable */
					},
				},
				NameServer: []*NameServer{
					{
						Address: &net.Endpoint{
							Network: net.Network_UDP,
							Address: &net.IPOrDomain{
								Address: &net.IPOrDomain_Ip{
									Ip: []byte{127, 0, 0, 1},
								},
							},
							Port: uint32(port),
						},
						PrioritizedDomain: []*NameServer_PriorityDomain{
							{
								Type:   DomainMatchingType_Full,
								Domain: "google.com",
							},
						},
					},
				},
			}),
			serial.ToTypedMessage(&dispatcher.Config{}),
			serial.ToTypedMessage(&proxyman.OutboundConfig{}),
			serial.ToTypedMessage(&policy.Config{}),
		},
		Outbound: []*core.OutboundHandlerConfig{
			{
				ProxySettings: serial.ToTypedMessage(&freedom.Config{}),
			},
		},
	}

	v, err := core.New(config)
	assert(err, IsNil)

	client := v.GetFeature(feature_dns.ClientType()).(feature_dns.Client)

	startTime := time.Now()
	ips, err := client.LookupIP("google.com")
	assert(err, IsNil)
	assert(len(ips), Equals, 1)
	assert([]byte(ips[0]), Equals, []byte{8, 8, 8, 8})

	endTime := time.Now()
	if startTime.After(endTime.Add(time.Second * 2)) {
		t.Error("DNS query doesn't finish in 2 seconds.")
	}
}
