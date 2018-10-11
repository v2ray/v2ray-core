package dns

import (
	"net"

	"v2ray.com/core/features"
)

// Client is a V2Ray feature for querying DNS information.
type Client interface {
	features.Feature
	LookupIP(host string) ([]net.IP, error)
}
