package app

import (
	"net"
)

type DnsCache interface {
	Get(domain string) net.IP
	Add(domain string, ip net.IP)
}

type DnsCacheWithContext interface {
	Get(context Context, domain string) net.IP
	Add(contaxt Context, domain string, ip net.IP)
}

type contextedDnsCache struct {
	context  Context
	dnsCache DnsCacheWithContext
}

func (this *contextedDnsCache) Get(domain string) net.IP {
	return this.dnsCache.Get(this.context, domain)
}

func (this *contextedDnsCache) Add(domain string, ip net.IP) {
	this.dnsCache.Add(this.context, domain, ip)
}
