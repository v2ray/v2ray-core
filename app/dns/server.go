package dns

import (
	"net"
	"sync"
	"time"

	"github.com/v2ray/v2ray-core/app"
	"github.com/v2ray/v2ray-core/app/dispatcher"
	"github.com/v2ray/v2ray-core/common/log"

	"github.com/miekg/dns"
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

func NewCacheServer(space app.Space, config *Config) *CacheServer {
	server := &CacheServer{
		records: make(map[string]*DomainRecord),
		servers: make([]NameServer, len(config.NameServers)),
		hosts:   config.Hosts,
	}
	space.InitializeApplication(func() error {
		if !space.HasApp(dispatcher.APP_ID) {
			log.Error("DNS: Dispatcher is not found in the space.")
			return app.ErrMissingApplication
		}

		dispatcher := space.GetApp(dispatcher.APP_ID).(dispatcher.PacketDispatcher)
		for idx, ns := range config.NameServers {
			if ns.Address().Family().IsDomain() && ns.Address().Domain() == "localhost" {
				server.servers[idx] = &LocalNameServer{}
			} else {
				server.servers[idx] = NewUDPNameServer(ns, dispatcher)
			}
		}
		if len(config.NameServers) == 0 {
			server.servers = append(server.servers, &LocalNameServer{})
		}
		return nil
	})
	return server
}

func (this *CacheServer) Release() {

}

//@Private
func (this *CacheServer) GetCached(domain string) []net.IP {
	this.RLock()
	defer this.RUnlock()

	if record, found := this.records[domain]; found && record.A.Expire.After(time.Now()) {
		return record.A.IPs
	}
	return nil
}

func (this *CacheServer) Get(domain string) []net.IP {
	if ip, found := this.hosts[domain]; found {
		return []net.IP{ip}
	}

	domain = dns.Fqdn(domain)
	ips := this.GetCached(domain)
	if ips != nil {
		return ips
	}

	for _, server := range this.servers {
		response := server.QueryA(domain)
		select {
		case a, open := <-response:
			if !open || a == nil {
				continue
			}
			this.Lock()
			this.records[domain] = &DomainRecord{
				A: a,
			}
			this.Unlock()
			log.Debug("DNS: Returning ", len(a.IPs), " IPs for domain ", domain)
			return a.IPs
		case <-time.After(QueryTimeout):
		}
	}

	log.Debug("DNS: Returning nil for domain ", domain)
	return nil
}
