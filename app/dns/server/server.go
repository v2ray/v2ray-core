package server

//go:generate go run $GOPATH/src/v2ray.com/core/tools/generrorgen/main.go -pkg server -path App,DNS,Server

import (
	"context"
	"sync"
	"time"

	dnsmsg "github.com/miekg/dns"
	"v2ray.com/core/app"
	"v2ray.com/core/app/dispatcher"
	"v2ray.com/core/app/dns"
	"v2ray.com/core/app/log"
	"v2ray.com/core/common"
	"v2ray.com/core/common/net"
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

func (r *DomainRecord) Inactive() bool {
	now := time.Now()
	return r.Expire.Before(now) || r.LastAccess.Add(time.Minute*5).Before(now)
}

type CacheServer struct {
	sync.Mutex
	hosts   map[string]net.IP
	records map[string]*DomainRecord
	servers []NameServer
}

func NewCacheServer(ctx context.Context, config *dns.Config) (*CacheServer, error) {
	space := app.SpaceFromContext(ctx)
	if space == nil {
		return nil, newError("no space in context")
	}
	server := &CacheServer{
		records: make(map[string]*DomainRecord),
		servers: make([]NameServer, len(config.NameServers)),
		hosts:   config.GetInternalHosts(),
	}
	space.On(app.SpaceInitializing, func(interface{}) error {
		disp := dispatcher.FromSpace(space)
		if disp == nil {
			return newError("dispatcher is not found in the space")
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
					server.servers[idx] = NewUDPNameServer(dest, disp)
				}
			}
		}
		if len(config.NameServers) == 0 {
			server.servers = append(server.servers, &LocalNameServer{})
		}
		return nil
	})
	return server, nil
}

func (*CacheServer) Interface() interface{} {
	return (*dns.Server)(nil)
}

func (*CacheServer) Start() error {
	return nil
}

func (*CacheServer) Close() {}

func (s *CacheServer) GetCached(domain string) []net.IP {
	s.Lock()
	defer s.Unlock()

	if record, found := s.records[domain]; found && !record.Expired() {
		record.LastAccess = time.Now()
		return record.IP
	}
	return nil
}

func (s *CacheServer) tryCleanup() {
	s.Lock()
	defer s.Unlock()

	if len(s.records) > 256 {
		domains := make([]string, 0, 256)
		for d, r := range s.records {
			if r.Expired() {
				domains = append(domains, d)
			}
		}
		for _, d := range domains {
			delete(s.records, d)
		}
	}
}

func (s *CacheServer) Get(domain string) []net.IP {
	if ip, found := s.hosts[domain]; found {
		return []net.IP{ip}
	}

	domain = dnsmsg.Fqdn(domain)
	ips := s.GetCached(domain)
	if ips != nil {
		return ips
	}

	s.tryCleanup()

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
			log.Trace(newError("returning ", len(a.IPs), " IPs for domain ", domain).AtDebug())
			return a.IPs
		case <-time.After(QueryTimeout):
		}
	}

	log.Trace(newError("returning nil for domain ", domain).AtDebug())
	return nil
}

func init() {
	common.Must(common.RegisterConfig((*dns.Config)(nil), func(ctx context.Context, config interface{}) (interface{}, error) {
		return NewCacheServer(ctx, config.(*dns.Config))
	}))
}
