// +build !confonly

package dns

//go:generate errorgen

import (
	"context"
	"sync"
	"time"

	"v2ray.com/core"
	"v2ray.com/core/common"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/session"
	"v2ray.com/core/common/strmatcher"
	"v2ray.com/core/features"
	"v2ray.com/core/features/dns"
	"v2ray.com/core/features/routing"
)

// Server is a DNS rely server.
type Server struct {
	sync.Mutex
	hosts          *StaticHosts
	clients        []Client
	clientIP       net.IP
	domainMatcher  strmatcher.IndexMatcher
	domainIndexMap map[uint32]uint32
	tag            string
}

// New creates a new DNS server with given configuration.
func New(ctx context.Context, config *Config) (*Server, error) {
	server := &Server{
		clients: make([]Client, 0, len(config.NameServers)+len(config.NameServer)),
		tag:     config.Tag,
	}
	if len(config.ClientIp) > 0 {
		if len(config.ClientIp) != 4 && len(config.ClientIp) != 16 {
			return nil, newError("unexpected IP length", len(config.ClientIp))
		}
		server.clientIP = net.IP(config.ClientIp)
	}

	hosts, err := NewStaticHosts(config.StaticHosts, config.Hosts)
	if err != nil {
		return nil, newError("failed to create hosts").Base(err)
	}
	server.hosts = hosts

	addNameServer := func(endpoint *net.Endpoint) int {
		address := endpoint.Address.AsAddress()
		if address.Family().IsDomain() && address.Domain() == "localhost" {
			server.clients = append(server.clients, NewLocalNameServer())
		} else {
			dest := endpoint.AsDestination()
			if dest.Network == net.Network_Unknown {
				dest.Network = net.Network_UDP
			}
			if dest.Network == net.Network_UDP {
				idx := len(server.clients)
				server.clients = append(server.clients, nil)

				common.Must(core.RequireFeatures(ctx, func(d routing.Dispatcher) {
					server.clients[idx] = NewClassicNameServer(dest, d, server.clientIP)
				}))
			}
		}
		return len(server.clients) - 1
	}

	if len(config.NameServers) > 0 {
		features.PrintDeprecatedFeatureWarning("simple DNS server")
	}

	for _, destPB := range config.NameServers {
		addNameServer(destPB)
	}

	if len(config.NameServer) > 0 {
		domainMatcher := &strmatcher.MatcherGroup{}
		domainIndexMap := make(map[uint32]uint32)

		for _, ns := range config.NameServer {
			idx := addNameServer(ns.Address)

			for _, domain := range ns.PrioritizedDomain {
				matcher, err := toStrMatcher(domain.Type, domain.Domain)
				if err != nil {
					return nil, newError("failed to create prioritized domain").Base(err).AtWarning()
				}
				midx := domainMatcher.Add(matcher)
				domainIndexMap[midx] = uint32(idx)
			}
		}

		server.domainMatcher = domainMatcher
		server.domainIndexMap = domainIndexMap
	}

	if len(server.clients) == 0 {
		server.clients = append(server.clients, NewLocalNameServer())
	}

	return server, nil
}

// Type implements common.HasType.
func (*Server) Type() interface{} {
	return dns.ClientType()
}

// Start implements common.Runnable.
func (s *Server) Start() error {
	return nil
}

// Close implements common.Closable.
func (s *Server) Close() error {
	return nil
}

func (s *Server) queryIPTimeout(client Client, domain string, option IPOption) ([]net.IP, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	if len(s.tag) > 0 {
		ctx = session.ContextWithInbound(ctx, &session.Inbound{
			Tag: s.tag,
		})
	}
	ips, err := client.QueryIP(ctx, domain, option)
	cancel()
	return ips, err
}

// LookupIP implements dns.Client.
func (s *Server) LookupIP(domain string) ([]net.IP, error) {
	return s.lookupIPInternal(domain, IPOption{
		IPv4Enable: true,
		IPv6Enable: true,
	})
}

// LookupIPv4 implements dns.IPv4Lookup.
func (s *Server) LookupIPv4(domain string) ([]net.IP, error) {
	return s.lookupIPInternal(domain, IPOption{
		IPv4Enable: true,
		IPv6Enable: false,
	})
}

// LookupIPv6 implements dns.IPv6Lookup.
func (s *Server) LookupIPv6(domain string) ([]net.IP, error) {
	return s.lookupIPInternal(domain, IPOption{
		IPv4Enable: false,
		IPv6Enable: true,
	})
}

func (s *Server) lookupStatic(domain string, option IPOption, depth int32) []net.Address {
	ips := s.hosts.LookupIP(domain, option)
	if ips == nil {
		return nil
	}
	if ips[0].Family().IsDomain() && depth < 5 {
		if newIPs := s.lookupStatic(ips[0].Domain(), option, depth+1); newIPs != nil {
			return newIPs
		}
	}
	return ips
}

func toNetIP(ips []net.Address) []net.IP {
	if len(ips) == 0 {
		return nil
	}
	netips := make([]net.IP, 0, len(ips))
	for _, ip := range ips {
		netips = append(netips, ip.IP())
	}
	return netips
}

func (s *Server) lookupIPInternal(domain string, option IPOption) ([]net.IP, error) {
	ips := s.lookupStatic(domain, option, 0)
	if ips != nil && ips[0].Family().IsIP() {
		return toNetIP(ips), nil
	}

	if ips != nil && ips[0].Family().IsDomain() {
		newdomain := ips[0].Domain()
		newError("domain replaced: ", domain, " -> ", newdomain).WriteToLog()
		domain = newdomain
	}

	var lastErr error
	if s.domainMatcher != nil {
		idx := s.domainMatcher.Match(domain)
		if idx > 0 {
			ns := s.clients[s.domainIndexMap[idx]]
			ips, err := s.queryIPTimeout(ns, domain, option)
			if len(ips) > 0 {
				return ips, nil
			}
			if err != nil {
				newError("failed to lookup ip for domain ", domain, " at server ", ns.Name()).Base(err).WriteToLog()
				lastErr = err
			}
		}
	}

	for _, client := range s.clients {
		ips, err := s.queryIPTimeout(client, domain, option)
		if len(ips) > 0 {
			return ips, nil
		}
		if err != nil {
			newError("failed to lookup ip for domain ", domain, " at server ", client.Name()).Base(err).WriteToLog()
			lastErr = err
		}
	}

	return nil, newError("returning nil for domain ", domain).Base(lastErr)
}

func init() {
	common.Must(common.RegisterConfig((*Config)(nil), func(ctx context.Context, config interface{}) (interface{}, error) {
		return New(ctx, config.(*Config))
	}))
}
