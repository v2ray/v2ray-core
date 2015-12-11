package app

import (
	"net"
)

// A DnsCache is an internal cache of DNS resolutions.
type DnsCache interface {
	Get(domain string) net.IP
	Add(domain string, ip net.IP)
}
