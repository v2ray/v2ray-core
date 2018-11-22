package dns

//go:generate errorgen

import (
	"context"
	"sync"
	"time"

	"v2ray.com/core"
	"v2ray.com/core/common"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/strmatcher"
	"v2ray.com/core/features"
	"v2ray.com/core/features/dns"
	"v2ray.com/core/features/routing"
)

// Server is a DNS rely server.
type Server struct {
	sync.Mutex
	hosts          *StaticHosts
	servers        []NameServerInterface
	clientIP       net.IP
	domainMatcher  strmatcher.IndexMatcher
	domainIndexMap map[uint32]uint32
}

// New creates a new DNS server with given configuration.
func New(ctx context.Context, config *Config) (*Server, error) {
	server := &Server{
		servers: make([]NameServerInterface, 0, len(config.NameServers)+len(config.NameServer)),
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
			server.servers = append(server.servers, NewLocalNameServer())
		} else {
			dest := endpoint.AsDestination()
			if dest.Network == net.Network_Unknown {
				dest.Network = net.Network_UDP
			}
			if dest.Network == net.Network_UDP {
				idx := len(server.servers)
				server.servers = append(server.servers, nil)

				common.Must(core.RequireFeatures(ctx, func(d routing.Dispatcher) {
					server.servers[idx] = NewClassicNameServer(dest, d, server.clientIP)
				}))
			}
		}
		return len(server.servers) - 1
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

	if len(server.servers) == 0 {
		server.servers = append(server.servers, NewLocalNameServer())
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

func (s *Server) queryIPTimeout(server NameServerInterface, domain string, option IPOption) ([]net.IP, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	ips, err := server.QueryIP(ctx, domain, option)
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

func (s *Server) lookupIPInternal(domain string, option IPOption) ([]net.IP, error) {
	if ip := s.hosts.LookupIP(domain, option); len(ip) > 0 {
		return ip, nil
	}

	var lastErr error
	if s.domainMatcher != nil {
		idx := s.domainMatcher.Match(domain)
		if idx > 0 {
			ns := s.servers[s.domainIndexMap[idx]]
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

	for _, server := range s.servers {
		ips, err := s.queryIPTimeout(server, domain, option)
		if len(ips) > 0 {
			return ips, nil
		}
		if err != nil {
			newError("failed to lookup ip for domain ", domain, " at server ", server.Name()).Base(err).WriteToLog()
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
