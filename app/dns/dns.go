package dns

import (
	"net"

	"v2ray.com/core/app"
)

const (
	APP_ID = app.ID(2)
)

// A DnsCache is an internal cache of DNS resolutions.
type Server interface {
	Get(domain string) []net.IP
}
