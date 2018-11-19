package dns

import (
	"v2ray.com/core/common/net"
	"v2ray.com/core/features"
)

// Client is a V2Ray feature for querying DNS information.
type Client interface {
	features.Feature

	// LookupIP returns IP address for the given domain. IPs may contain IPv4 and/or IPv6 addresses.
	LookupIP(domain string) ([]net.IP, error)
}

// IPv4Lookup is an optional feature for querying IPv4 addresses only.
type IPv4Lookup interface {
	LookupIPv4(domain string) ([]net.IP, error)
}

// IPv6Lookup is an optional feature for querying IPv6 addresses only.
type IPv6Lookup interface {
	LookupIPv6(domain string) ([]net.IP, error)
}

// ClientType returns the type of Client interface. Can be used for implementing common.HasType.
func ClientType() interface{} {
	return (*Client)(nil)
}
