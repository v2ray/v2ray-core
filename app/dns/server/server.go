package server

import (
	"net"
	"sync"
	"time"

	"v2ray.com/core/app"
	"v2ray.com/core/app/dispatcher"
	"v2ray.com/core/app/dns"
	"v2ray.com/core/common/errors"
	"v2ray.com/core/common/log"
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/common/serial"

	dnsmsg "github.com/miekg/dns"
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

func NewCacheServer(space app.Space, config *dns.Config) *CacheServer {
	server := &CacheServer{
		records: make(map[string]*DomainRecord),
		servers: make([]NameServer, len(config.NameServers)),
		hosts:   config.GetInternalHosts(),
	}
	space.InitializeApplication(func() error {
		if !space.HasApp(dispatcher.APP_ID) {
			return errors.New("DNS: Dispatcher is not found in the space.")
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
			log.Debug("DNS: Returning ", len(a.IPs), " IPs for domain ", domain)
			return a.IPs
		case <-time.After(QueryTimeout):
		}
	}

	log.Debug("DNS: Returning nil for domain ", domain)
	return nil
}

type CacheServerFactory struct{}

func (v CacheServerFactory) Create(space app.Space, config interface{}) (app.Application, error) {
	server := NewCacheServer(space, config.(*dns.Config))
	return server, nil
}

func (v CacheServerFactory) AppId() app.ID {
	return dns.APP_ID
}

func init() {
	app.RegisterApplicationFactory(serial.GetMessageType(new(dns.Config)), CacheServerFactory{})
}
