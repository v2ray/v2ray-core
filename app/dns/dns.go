package dns

import (
	"net"

	"v2ray.com/core/app"
)

// A DnsCache is an internal cache of DNS resolutions.
type Server interface {
	Get(domain string) []net.IP
}

func FromSpace(space app.Space) Server {
	app := space.GetApplication((*Server)(nil))
	if app == nil {
		return nil
	}
	return app.(Server)
}
