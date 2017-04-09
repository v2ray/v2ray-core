package server

//go:generate go run $GOPATH/src/v2ray.com/core/tools/generrorgen/main.go -pkg server -path App,DNS,Server

import (
	"context"
	"net"
	"sync"
	"time"

	dnsmsg "github.com/miekg/dns"
	"v2ray.com/core/app"
	"v2ray.com/core/app/dispatcher"
	"v2ray.com/core/app/dns"
	"v2ray.com/core/app/log"
	"v2ray.com/core/common"
	v2net "v2ray.com/core/common/net"
)

const (
	QueryTimeout = time.Second * 8
)

type DomainRecord struct {
	A *ARecord
}

type CacheServer struct {
	sync.RWMutex
	space   app.Space
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
	space.OnInitialize(func() error {
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
				if dest.Network == v2net.Network_Unknown {
					dest.Network = v2net.Network_UDP
				}
				if dest.Network == v2net.Network_UDP {
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

// Private: Visible for testing.
func (v *CacheServer) GetCached(domain string) []net.IP {
	v.RLock()
	defer v.RUnlock()

	if record, found := v.records[domain]; found && record.A.Expire.After(time.Now()) {
		return record.A.IPs
	}
	return nil
}

func (v *CacheServer) Get(domain string) []net.IP {
	if ip, found := v.hosts[domain]; found {
		return []net.IP{ip}
	}

	domain = dnsmsg.Fqdn(domain)
	ips := v.GetCached(domain)
	if ips != nil {
		return ips
	}

	for _, server := range v.servers {
		response := server.QueryA(domain)
		select {
		case a, open := <-response:
			if !open || a == nil {
				continue
			}
			v.Lock()
			v.records[domain] = &DomainRecord{
				A: a,
			}
			v.Unlock()
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
