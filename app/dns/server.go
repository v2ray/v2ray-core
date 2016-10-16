package dns

import (
	"net"
	"sync"
	"time"

	"v2ray.com/core/app"
	"v2ray.com/core/app/dispatcher"
	"v2ray.com/core/common/loader"
	"v2ray.com/core/common/log"
	v2net "v2ray.com/core/common/net"

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
		hosts:   config.GetInternalHosts(),
	}
	space.InitializeApplication(func() error {
		if !space.HasApp(dispatcher.APP_ID) {
			log.Error("DNS: Dispatcher is not found in the space.")
			return app.ErrMissingApplication
		}

		dispatcher := space.GetApp(dispatcher.APP_ID).(dispatcher.PacketDispatcher)
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
					server.servers[idx] = NewUDPNameServer(dest, dispatcher)
				}
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

// Private: Visible for testing.
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

type CacheServerFactory struct{}

func (this CacheServerFactory) Create(space app.Space, config interface{}) (app.Application, error) {
	server := NewCacheServer(space, config.(*Config))
	return server, nil
}

func (this CacheServerFactory) AppId() app.ID {
	return APP_ID
}

func init() {
	app.RegisterApplicationFactory(loader.GetType(new(Config)), CacheServerFactory{})
}
