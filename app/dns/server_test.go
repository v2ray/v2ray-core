package dns_test

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/miekg/dns"

	"v2ray.com/core"
	"v2ray.com/core/app/dispatcher"
	. "v2ray.com/core/app/dns"
	"v2ray.com/core/app/policy"
	"v2ray.com/core/app/proxyman"
	_ "v2ray.com/core/app/proxyman/outbound"
	"v2ray.com/core/app/router"
	"v2ray.com/core/common"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/serial"
	feature_dns "v2ray.com/core/features/dns"
	"v2ray.com/core/proxy/freedom"
	"v2ray.com/core/testing/servers/udp"
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
		} else if q.Name == "ipv6.google.com." && q.Qtype == dns.TypeA {
			rr, err := dns.NewRR("ipv6.google.com. IN A 8.8.8.7")
			common.Must(err)
			ans.Answer = append(ans.Answer, rr)
		} else if q.Name == "ipv6.google.com." && q.Qtype == dns.TypeAAAA {
			rr, err := dns.NewRR("ipv6.google.com. IN AAAA 2001:4860:4860::8888")
			common.Must(err)
			ans.Answer = append(ans.Answer, rr)
		} else if q.Name == "notexist.google.com." && q.Qtype == dns.TypeAAAA {
			ans.MsgHdr.Rcode = dns.RcodeNameError
		}
	}
	w.WriteMsg(ans)
}

func TestUDPServerSubnet(t *testing.T) {
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
	common.Must(err)

	client := v.GetFeature(feature_dns.ClientType()).(feature_dns.Client)

	ips, err := client.LookupIP("google.com")
	if err != nil {
		t.Fatal("unexpected error: ", err)
	}

	if r := cmp.Diff(ips, []net.IP{{8, 8, 4, 4}}); r != "" {
		t.Fatal(r)
	}
}

func TestUDPServer(t *testing.T) {
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
	common.Must(err)

	client := v.GetFeature(feature_dns.ClientType()).(feature_dns.Client)

	{
		ips, err := client.LookupIP("google.com")
		if err != nil {
			t.Fatal("unexpected error: ", err)
		}

		if r := cmp.Diff(ips, []net.IP{{8, 8, 8, 8}}); r != "" {
			t.Fatal(r)
		}
	}

	{
		ips, err := client.LookupIP("facebook.com")
		if err != nil {
			t.Fatal("unexpected error: ", err)
		}

		if r := cmp.Diff(ips, []net.IP{{9, 9, 9, 9}}); r != "" {
			t.Fatal(r)
		}
	}

	{
		_, err := client.LookupIP("notexist.google.com")
		if err == nil {
			t.Fatal("nil error")
		}
		if r := feature_dns.RCodeFromError(err); r != uint16(dns.RcodeNameError) {
			t.Fatal("expected NameError, but got ", r)
		}
	}

	{
		clientv6 := client.(feature_dns.IPv6Lookup)
		ips, err := clientv6.LookupIPv6("ipv4only.google.com")
		if err != feature_dns.ErrEmptyResponse {
			t.Fatal("error: ", err)
		}
		if len(ips) != 0 {
			t.Fatal("ips: ", ips)
		}
	}

	dnsServer.Shutdown()

	{
		ips, err := client.LookupIP("google.com")
		if err != nil {
			t.Fatal("unexpected error: ", err)
		}

		if r := cmp.Diff(ips, []net.IP{{8, 8, 8, 8}}); r != "" {
			t.Fatal(r)
		}
	}
}

func TestPrioritizedDomain(t *testing.T) {
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
	common.Must(err)

	client := v.GetFeature(feature_dns.ClientType()).(feature_dns.Client)

	startTime := time.Now()

	{
		ips, err := client.LookupIP("google.com")
		if err != nil {
			t.Fatal("unexpected error: ", err)
		}

		if r := cmp.Diff(ips, []net.IP{{8, 8, 8, 8}}); r != "" {
			t.Fatal(r)
		}
	}

	endTime := time.Now()
	if startTime.After(endTime.Add(time.Second * 2)) {
		t.Error("DNS query doesn't finish in 2 seconds.")
	}
}

func TestUDPServerIPv6(t *testing.T) {
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
	common.Must(err)

	client := v.GetFeature(feature_dns.ClientType()).(feature_dns.Client)
	client6 := client.(feature_dns.IPv6Lookup)

	{
		ips, err := client6.LookupIPv6("ipv6.google.com")
		if err != nil {
			t.Fatal("unexpected error: ", err)
		}

		if r := cmp.Diff(ips, []net.IP{{32, 1, 72, 96, 72, 96, 0, 0, 0, 0, 0, 0, 0, 0, 136, 136}}); r != "" {
			t.Fatal(r)
		}
	}
}

func TestStaticHostDomain(t *testing.T) {
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
				StaticHosts: []*Config_HostMapping{
					{
						Type:          DomainMatchingType_Full,
						Domain:        "example.com",
						ProxiedDomain: "google.com",
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
	common.Must(err)

	client := v.GetFeature(feature_dns.ClientType()).(feature_dns.Client)

	{
		ips, err := client.LookupIP("example.com")
		if err != nil {
			t.Fatal("unexpected error: ", err)
		}

		if r := cmp.Diff(ips, []net.IP{{8, 8, 8, 8}}); r != "" {
			t.Fatal(r)
		}
	}

	dnsServer.Shutdown()
}

func TestIPMatch(t *testing.T) {
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
				NameServer: []*NameServer{
					// private dns, not match
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
						Geoip: []*router.GeoIP{
							{
								CountryCode: "local",
								Cidr: []*router.CIDR{
									{
										// inner ip, will not match
										Ip:     []byte{192, 168, 11, 1},
										Prefix: 32,
									},
								},
							},
						},
					},
					// second dns, match ip
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
						Geoip: []*router.GeoIP{
							{
								CountryCode: "test",
								Cidr: []*router.CIDR{
									{
										Ip:     []byte{8, 8, 8, 8},
										Prefix: 32,
									},
								},
							},
							{
								CountryCode: "test",
								Cidr: []*router.CIDR{
									{
										Ip:     []byte{8, 8, 8, 4},
										Prefix: 32,
									},
								},
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
	common.Must(err)

	client := v.GetFeature(feature_dns.ClientType()).(feature_dns.Client)

	startTime := time.Now()

	{
		ips, err := client.LookupIP("google.com")
		if err != nil {
			t.Fatal("unexpected error: ", err)
		}

		if r := cmp.Diff(ips, []net.IP{{8, 8, 8, 8}}); r != "" {
			t.Fatal(r)
		}
	}

	endTime := time.Now()
	if startTime.After(endTime.Add(time.Second * 2)) {
		t.Error("DNS query doesn't finish in 2 seconds.")
	}
}
