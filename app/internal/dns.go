package internal

import (
	"net"

	"github.com/v2ray/v2ray-core/app"
)

type DnsCacheWithContext interface {
	Get(context app.Context, domain string) net.IP
	Add(contaxt app.Context, domain string, ip net.IP)
}

type contextedDnsCache struct {
	context  app.Context
	dnsCache DnsCacheWithContext
}

func (this *contextedDnsCache) Get(domain string) net.IP {
	return this.dnsCache.Get(this.context, domain)
}

func (this *contextedDnsCache) Add(domain string, ip net.IP) {
	this.dnsCache.Add(this.context, domain, ip)
}
