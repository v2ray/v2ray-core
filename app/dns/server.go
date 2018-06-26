package dns

//go:generate go run $GOPATH/src/v2ray.com/core/common/errors/errorgen/main.go -pkg dns -path App,DNS

import (
	"context"
	"sync"
	"time"

	"v2ray.com/core"
	"v2ray.com/core/common"
	"v2ray.com/core/common/net"
)

type Server struct {
	sync.Mutex
	hosts    map[string]net.IP
	servers  []NameServer
	clientIP *Config_ClientIP
}

func New(ctx context.Context, config *Config) (*Server, error) {
	server := &Server{
		servers: make([]NameServer, len(config.NameServers)),
		hosts:   config.GetInternalHosts(),
	}
	if config.ClientIp != nil {
		server.clientIP = config.ClientIp
	}

	v := core.MustFromContext(ctx)
	if err := v.RegisterFeature((*core.DNSClient)(nil), server); err != nil {
		return nil, newError("unable to register DNSClient.").Base(err)
	}

	for idx, destPB := range config.NameServers {
		address := destPB.Address.AsAddress()
		if address.Family().IsDomain() && address.Domain() == "localhost" {
			server.servers[idx] = &LocalNameServer{}
		} else {
			dest := destPB.AsDestination()
			if dest.Network == net.Network_Unknown {
				dest.Network = net.Network_UDP
			}
			if dest.Network == net.Network_UDP {
				server.servers[idx] = NewClassicNameServer(dest, v.Dispatcher(), server.clientIP)
			}
		}
	}
	if len(config.NameServers) == 0 {
		server.servers = append(server.servers, &LocalNameServer{})
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

func (s *Server) LookupIP(domain string) ([]net.IP, error) {
	if ip, found := s.hosts[domain]; found {
		return []net.IP{ip}, nil
	}

	var lastErr error
	for _, server := range s.servers {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
		ips, err := server.QueryIP(ctx, domain)
		cancel()
		if err != nil {
			lastErr = err
		}
		if len(ips) > 0 {
			return ips, nil
		}
	}

	return nil, newError("returning nil for domain ", domain).Base(lastErr)
}

func init() {
	common.Must(common.RegisterConfig((*Config)(nil), func(ctx context.Context, config interface{}) (interface{}, error) {
		return New(ctx, config.(*Config))
	}))
}
