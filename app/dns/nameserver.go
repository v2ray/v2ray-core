package dns

import (
	"context"
	"time"

	"v2ray.com/core/common/net"
)

var (
	multiQuestionDNS = map[net.Address]bool{
		net.IPAddress([]byte{8, 8, 8, 8}): true,
		net.IPAddress([]byte{8, 8, 4, 4}): true,
		net.IPAddress([]byte{9, 9, 9, 9}): true,
	}
)

type ARecord struct {
	IPs    []net.IP
	Expire time.Time
}

type NameServer interface {
	QueryIP(ctx context.Context, domain string) ([]net.IP, error)
}

type LocalNameServer struct {
}

func (*LocalNameServer) QueryIP(ctx context.Context, domain string) ([]net.IP, error) {
	return net.LookupIP(domain)
}
