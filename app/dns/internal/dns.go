package internal

import (
	"net"
	"sync"
	"time"

	"github.com/v2ray/v2ray-core/app"
	"github.com/v2ray/v2ray-core/app/dispatcher"

	"github.com/miekg/dns"
)

const (
	QueryTimeout = time.Second * 2
)

type DomainRecord struct {
	A *ARecord
}

type Server struct {
	sync.RWMutex
	records map[string]*DomainRecord
	servers []NameServer
}

func NewServer(space app.Space, config *Config) *Server {
	server := &Server{
		records: make(map[string]*DomainRecord),
		servers: make([]NameServer, len(config.NameServers)),
	}
	dispatcher := space.GetApp(dispatcher.APP_ID).(dispatcher.PacketDispatcher)
	for idx, ns := range config.NameServers {
		server.servers[idx] = NewUDPNameServer(ns, dispatcher)
	}
	return server
}

//@Private
func (this *Server) GetCached(domain string) []net.IP {
	this.RLock()
	defer this.RUnlock()

	if record, found := this.records[domain]; found && record.A.Expire.After(time.Now()) {
		return record.A.IPs
	}
	return nil
}

func (this *Server) Get(context app.Context, domain string) []net.IP {
	domain = dns.Fqdn(domain)
	ips := this.GetCached(domain)
	if ips != nil {
		return ips
	}

	for _, server := range this.servers {
		response := server.QueryA(domain)
		select {
		case a := <-response:
			this.Lock()
			this.records[domain] = &DomainRecord{
				A: a,
			}
			this.Unlock()
			return a.IPs
		case <-time.Tick(QueryTimeout):
		}
	}

	return nil
}
