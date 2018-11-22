package dns

import (
	"context"

	"v2ray.com/core/common/net"
	"v2ray.com/core/features/dns/localdns"
)

type IPOption struct {
	IPv4Enable bool
	IPv6Enable bool
}

type NameServerInterface interface {
	QueryIP(ctx context.Context, domain string, option IPOption) ([]net.IP, error)
}

type localNameServer struct {
	client *localdns.Client
}

func (s *localNameServer) QueryIP(ctx context.Context, domain string, option IPOption) ([]net.IP, error) {
	if option.IPv4Enable && option.IPv6Enable {
		return s.client.LookupIP(domain)
	}

	if option.IPv4Enable {
		return s.client.LookupIPv4(domain)
	}

	if option.IPv6Enable {
		return s.client.LookupIPv6(domain)
	}

	return nil, newError("neither IPv4 nor IPv6 is enabled")
}

func NewLocalNameServer() *localNameServer {
	return &localNameServer{
		client: localdns.New(),
	}
}
