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
		} else if q.Name == "api.google.com." && q.Qtype == dns.TypeA {
			rr, _ := dns.NewRR("api.google.com. IN A 8.8.7.7")
			ans.Answer = append(ans.Answer, rr)
		} else if q.Name == "v2.api.google.com." && q.Qtype == dns.TypeA {
			rr, _ := dns.NewRR("v2.api.google.com. IN A 8.8.7.8")
			ans.Answer = append(ans.Answer, rr)
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
		} else if q.Name == "hostname." && q.Qtype == dns.TypeA {
			rr, _ := dns.NewRR("hostname. IN A 127.0.0.1")
			ans.Answer = append(ans.Answer, rr)
		} else if q.Name == "hostname.local." && q.Qtype == dns.TypeA {
			rr, _ := dns.NewRR("hostname.local. IN A 127.0.0.1")
			ans.Answer = append(ans.Answer, rr)
		} else if q.Name == "hostname.localdomain." && q.Qtype == dns.TypeA {
			rr, _ := dns.NewRR("hostname.localdomain. IN A 127.0.0.1")
			ans.Answer = append(ans.Answer, rr)
		} else if q.Name == "localhost." && q.Qtype == dns.TypeA {
			rr, _ := dns.NewRR("localhost. IN A 127.0.0.2")
			ans.Answer = append(ans.Answer, rr)
		} else if q.Name == "localhost-a." && q.Qtype == dns.TypeA {
			rr, _ := dns.NewRR("localhost-a. IN A 127.0.0.3")
			ans.Answer = append(ans.Answer, rr)
		} else if q.Name == "localhost-b." && q.Qtype == dns.TypeA {
			rr, _ := dns.NewRR("localhost-b. IN A 127.0.0.4")
			ans.Answer = append(ans.Answer, rr)
		} else if q.Name == "Mijia\\ Cloud." && q.Qtype == dns.TypeA {
			rr, _ := dns.NewRR("Mijia\\ Cloud. IN A 127.0.0.1")
			ans.Answer = append(ans.Answer, rr)
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

func TestLocalDomain(t *testing.T) {
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
							// Equivalent of dotless:localhost
							{Type: DomainMatchingType_Regex, Domain: "^[^.]*localhost[^.]*$"},
						},
						Geoip: []*router.GeoIP{
							{ // Will match localhost, localhost-a and localhost-b,
								CountryCode: "local",
								Cidr: []*router.CIDR{
									{Ip: []byte{127, 0, 0, 2}, Prefix: 32},
									{Ip: []byte{127, 0, 0, 3}, Prefix: 32},
									{Ip: []byte{127, 0, 0, 4}, Prefix: 32},
								},
							},
						},
					},
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
							// Equivalent of dotless: and domain:local
							{Type: DomainMatchingType_Regex, Domain: "^[^.]*$"},
							{Type: DomainMatchingType_Subdomain, Domain: "local"},
							{Type: DomainMatchingType_Subdomain, Domain: "localdomain"},
						},
					},
				},
				StaticHosts: []*Config_HostMapping{
					{
						Type:   DomainMatchingType_Full,
						Domain: "hostnamestatic",
						Ip:     [][]byte{{127, 0, 0, 53}},
					},
					{
						Type:          DomainMatchingType_Full,
						Domain:        "hostnamealias",
						ProxiedDomain: "hostname.localdomain",
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

	{ // Will match dotless:
		ips, err := client.LookupIP("hostname")
		if err != nil {
			t.Fatal("unexpected error: ", err)
		}

		if r := cmp.Diff(ips, []net.IP{{127, 0, 0, 1}}); r != "" {
			t.Fatal(r)
		}
	}

	{ // Will match domain:local
		ips, err := client.LookupIP("hostname.local")
		if err != nil {
			t.Fatal("unexpected error: ", err)
		}

		if r := cmp.Diff(ips, []net.IP{{127, 0, 0, 1}}); r != "" {
			t.Fatal(r)
		}
	}

	{ // Will match static ip
		ips, err := client.LookupIP("hostnamestatic")
		if err != nil {
			t.Fatal("unexpected error: ", err)
		}

		if r := cmp.Diff(ips, []net.IP{{127, 0, 0, 53}}); r != "" {
			t.Fatal(r)
		}
	}

	{ // Will match domain replacing
		ips, err := client.LookupIP("hostnamealias")
		if err != nil {
			t.Fatal("unexpected error: ", err)
		}

		if r := cmp.Diff(ips, []net.IP{{127, 0, 0, 1}}); r != "" {
			t.Fatal(r)
		}
	}

	{ // Will match dotless:localhost, but not expectIPs: 127.0.0.2, 127.0.0.3, then matches at dotless:
		ips, err := client.LookupIP("localhost")
		if err != nil {
			t.Fatal("unexpected error: ", err)
		}

		if r := cmp.Diff(ips, []net.IP{{127, 0, 0, 2}}); r != "" {
			t.Fatal(r)
		}
	}

	{ // Will match dotless:localhost, and expectIPs: 127.0.0.2, 127.0.0.3
		ips, err := client.LookupIP("localhost-a")
		if err != nil {
			t.Fatal("unexpected error: ", err)
		}

		if r := cmp.Diff(ips, []net.IP{{127, 0, 0, 3}}); r != "" {
			t.Fatal(r)
		}
	}

	{ // Will match dotless:localhost, and expectIPs: 127.0.0.2, 127.0.0.3
		ips, err := client.LookupIP("localhost-b")
		if err != nil {
			t.Fatal("unexpected error: ", err)
		}

		if r := cmp.Diff(ips, []net.IP{{127, 0, 0, 4}}); r != "" {
			t.Fatal(r)
		}
	}

	{ // Will match dotless:
		ips, err := client.LookupIP("Mijia Cloud")
		if err != nil {
			t.Fatal("unexpected error: ", err)
		}

		if r := cmp.Diff(ips, []net.IP{{127, 0, 0, 1}}); r != "" {
			t.Fatal(r)
		}
	}

	endTime := time.Now()
	if startTime.After(endTime.Add(time.Second * 2)) {
		t.Error("DNS query doesn't finish in 2 seconds.")
	}
}

func TestMultiMatchPrioritizedDomain(t *testing.T) {
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
								Type:   DomainMatchingType_Subdomain,
								Domain: "google.com",
							},
						},
						Geoip: []*router.GeoIP{
							{ // Will only match 8.8.8.8 and 8.8.4.4
								Cidr: []*router.CIDR{
									{Ip: []byte{8, 8, 8, 8}, Prefix: 32},
									{Ip: []byte{8, 8, 4, 4}, Prefix: 32},
								},
							},
						},
					},
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
								Type:   DomainMatchingType_Subdomain,
								Domain: "google.com",
							},
						},
						Geoip: []*router.GeoIP{
							{ // Will match 8.8.8.8 and 8.8.8.7, etc
								Cidr: []*router.CIDR{
									{Ip: []byte{8, 8, 8, 7}, Prefix: 24},
								},
							},
						},
					},
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
								Type:   DomainMatchingType_Subdomain,
								Domain: "api.google.com",
							},
						},
						Geoip: []*router.GeoIP{
							{ // Will only match 8.8.7.7 (api.google.com)
								Cidr: []*router.CIDR{
									{Ip: []byte{8, 8, 7, 7}, Prefix: 32},
								},
							},
						},
					},
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
								Domain: "v2.api.google.com",
							},
						},
						Geoip: []*router.GeoIP{
							{ // Will only match 8.8.7.8 (v2.api.google.com)
								Cidr: []*router.CIDR{
									{Ip: []byte{8, 8, 7, 8}, Prefix: 32},
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

	{ // Will match server 1,2 and server 1 returns expected ip
		ips, err := client.LookupIP("google.com")
		if err != nil {
			t.Fatal("unexpected error: ", err)
		}

		if r := cmp.Diff(ips, []net.IP{{8, 8, 8, 8}}); r != "" {
			t.Fatal(r)
		}
	}

	{ // Will match server 1,2 and server 1 returns unexpected ip, then server 2 returns expected one
		clientv4 := client.(feature_dns.IPv4Lookup)
		ips, err := clientv4.LookupIPv4("ipv6.google.com")
		if err != nil {
			t.Fatal("unexpected error: ", err)
		}

		if r := cmp.Diff(ips, []net.IP{{8, 8, 8, 7}}); r != "" {
			t.Fatal(r)
		}
	}

	{ // Will match server 3,1,2 and server 3 returns expected one
		ips, err := client.LookupIP("api.google.com")
		if err != nil {
			t.Fatal("unexpected error: ", err)
		}

		if r := cmp.Diff(ips, []net.IP{{8, 8, 7, 7}}); r != "" {
			t.Fatal(r)
		}
	}

	{ // Will match server 4,3,1,2 and server 4 returns expected one
		ips, err := client.LookupIP("v2.api.google.com")
		if err != nil {
			t.Fatal("unexpected error: ", err)
		}

		if r := cmp.Diff(ips, []net.IP{{8, 8, 7, 8}}); r != "" {
			t.Fatal(r)
		}
	}

	endTime := time.Now()
	if startTime.After(endTime.Add(time.Second * 2)) {
		t.Error("DNS query doesn't finish in 2 seconds.")
	}
}
