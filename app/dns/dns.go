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
