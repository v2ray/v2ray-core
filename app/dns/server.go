package dns

//go:generate go run $GOPATH/src/v2ray.com/core/common/errors/errorgen/main.go -pkg dns -path App,DNS

import (
	"context"
	"sync"
	"time"

	"v2ray.com/core"
	"v2ray.com/core/common"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/strmatcher"
)

type Server struct {
	sync.Mutex
	hosts          *StaticHosts
	servers        []NameServerInterface
	clientIP       net.IP
	domainMatcher  strmatcher.IndexMatcher
	domainIndexMap map[uint32]uint32
}

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

	v := core.MustFromContext(ctx)
	if err := v.RegisterFeature((*core.DNSClient)(nil), server); err != nil {
		return nil, newError("unable to register DNSClient.").Base(err)
	}

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
				server.servers = append(server.servers, NewClassicNameServer(dest, v.Dispatcher(), server.clientIP))
			}
		}
		return len(server.servers) - 1
	}

	if len(config.NameServers) > 0 {
		core.PrintDeprecatedFeatureWarning("simple DNS server")
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
					return nil, newError("failed to create proritized domain").Base(err).AtWarning()
				}
				midx := domainMatcher.Add(matcher)
				domainIndexMap[midx] = uint32(idx)
			}
		}

		server.domainMatcher = domainMatcher
		server.domainIndexMap = domainIndexMap
	}

	if len(config.NameServers) == 0 {
		server.servers = append(server.servers, NewLocalNameServer())
	}

	return server, nil
}

// Start implements common.Runnable.
func (s *Server) Start() error {
	return nil
}

// Close implements common.Closable.
func (s *Server) Close() error {
	return nil
}

func (s *Server) queryIPTimeout(server NameServerInterface, domain string) ([]net.IP, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	ips, err := server.QueryIP(ctx, domain)
	cancel()
	return ips, err
}

func (s *Server) LookupIP(domain string) ([]net.IP, error) {
	if ip := s.hosts.LookupIP(domain); len(ip) > 0 {
		return ip, nil
	}

	var lastErr error
	if s.domainMatcher != nil {
		idx := s.domainMatcher.Match(domain)
		if idx > 0 {
			ns := s.servers[idx]
			ips, err := s.queryIPTimeout(ns, domain)
			if len(ips) > 0 {
				return ips, nil
			}
			if err != nil {
				lastErr = err
			}
		}
	}

	for _, server := range s.servers {
		ips, err := s.queryIPTimeout(server, domain)
		if len(ips) > 0 {
			return ips, nil
		}
		if err != nil {
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
