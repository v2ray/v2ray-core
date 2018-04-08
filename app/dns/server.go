package dns

//go:generate go run $GOPATH/src/v2ray.com/core/common/errors/errorgen/main.go -pkg dns -path App,DNS

import (
	"context"
	"sync"
	"time"

	dnsmsg "github.com/miekg/dns"
	"v2ray.com/core"
	"v2ray.com/core/common"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/signal"
)

const (
	QueryTimeout = time.Second * 8
)

type DomainRecord struct {
	IP         []net.IP
	Expire     time.Time
	LastAccess time.Time
}

func (r *DomainRecord) Expired() bool {
	return r.Expire.Before(time.Now())
}

type Server struct {
	sync.Mutex
	hosts   map[string]net.IP
	records map[string]*DomainRecord
	servers []NameServer
	task    *signal.PeriodicTask
}

func New(ctx context.Context, config *Config) (*Server, error) {
	server := &Server{
		records: make(map[string]*DomainRecord),
		servers: make([]NameServer, len(config.NameServers)),
		hosts:   config.GetInternalHosts(),
	}
	server.task = &signal.PeriodicTask{
		Interval: time.Minute * 10,
		Execute: func() error {
			server.cleanup()
			return nil
		},
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
				server.servers[idx] = NewUDPNameServer(dest, v.Dispatcher())
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
	return s.task.Start()
}

// Close implements common.Closable.
func (s *Server) Close() error {
	return s.task.Close()
}

func (s *Server) GetCached(domain string) []net.IP {
	s.Lock()
	defer s.Unlock()

	if record, found := s.records[domain]; found && !record.Expired() {
		record.LastAccess = time.Now()
		return record.IP
	}
	return nil
}

func (s *Server) cleanup() {
	s.Lock()
	defer s.Unlock()

	for d, r := range s.records {
		if r.Expired() {
			delete(s.records, d)
		}
	}

	if len(s.records) == 0 {
		s.records = make(map[string]*DomainRecord)
	}
}

func (s *Server) LookupIP(domain string) ([]net.IP, error) {
	if ip, found := s.hosts[domain]; found {
		return []net.IP{ip}, nil
	}

	domain = dnsmsg.Fqdn(domain)
	ips := s.GetCached(domain)
	if ips != nil {
		return ips, nil
	}

	for _, server := range s.servers {
		response := server.QueryA(domain)
		select {
		case a, open := <-response:
			if !open || a == nil {
				continue
			}
			s.Lock()
			s.records[domain] = &DomainRecord{
				IP:         a.IPs,
				Expire:     a.Expire,
				LastAccess: time.Now(),
			}
			s.Unlock()
			newError("returning ", len(a.IPs), " IPs for domain ", domain).AtDebug().WriteToLog()
			return a.IPs, nil
		case <-time.After(QueryTimeout):
		}
	}

	return nil, newError("returning nil for domain ", domain)
}

func init() {
	common.Must(common.RegisterConfig((*Config)(nil), func(ctx context.Context, config interface{}) (interface{}, error) {
		return New(ctx, config.(*Config))
	}))
}
