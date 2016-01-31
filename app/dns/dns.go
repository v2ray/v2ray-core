package dns

import (
	"net"

	"github.com/v2ray/v2ray-core/app"
)

const (
	APP_ID = app.ID(2)
)

// A DnsCache is an internal cache of DNS resolutions.
type DnsCache interface {
	Get(domain string) net.IP
	Add(domain string, ip net.IP)
}

type dnsCacheWithContext interface {
	Get(context app.Context, domain string) net.IP
	Add(contaxt app.Context, domain string, ip net.IP)
}

type contextedDnsCache struct {
	context  app.Context
	dnsCache dnsCacheWithContext
}

func (this *contextedDnsCache) Get(domain string) net.IP {
	return this.dnsCache.Get(this.context, domain)
}

func (this *contextedDnsCache) Add(domain string, ip net.IP) {
	this.dnsCache.Add(this.context, domain, ip)
}

func init() {
	app.RegisterApp(APP_ID, func(context app.Context, obj interface{}) interface{} {
		dcContext := obj.(dnsCacheWithContext)
		return &contextedDnsCache{
			context:  context,
			dnsCache: dcContext,
		}
	})
}
