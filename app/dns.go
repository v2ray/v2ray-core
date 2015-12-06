package app

import (
	"net"
)

type DnsCache interface {
	Get(domain string) net.IP
	Add(domain string, ip net.IP)
}
