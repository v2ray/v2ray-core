package dns

import (
	"net"

	"github.com/v2ray/v2ray-core/app"
)

const (
	APP_ID = app.ID(2)
)

// A DnsCache is an internal cache of DNS resolutions.
type Server interface {
	Get(domain string) []net.IP
}

type dnsServerWithContext interface {
	Get(context app.Context, domain string) []net.IP
}

type contextedDnsServer struct {
	context  app.Context
	dnsCache dnsServerWithContext
}

func (this *contextedDnsServer) Get(domain string) []net.IP {
	return this.dnsCache.Get(this.context, domain)
}

func CreateDNSServer(rawConfig interface{}) (Server, error) {
	return nil, nil
}

func init() {
	app.Register(APP_ID, func(context app.Context, obj interface{}) interface{} {
		dcContext := obj.(dnsServerWithContext)
		return &contextedDnsServer{
			context:  context,
			dnsCache: dcContext,
		}
	})
}
