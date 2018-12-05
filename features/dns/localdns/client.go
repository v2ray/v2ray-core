package localdns

import (
	"v2ray.com/core/common/net"
	"v2ray.com/core/features/dns"
)

// Client is an implementation of dns.Client, which queries localhost for DNS.
type Client struct{}

// Type implements common.HasType.
func (*Client) Type() interface{} {
	return dns.ClientType()
}

// Start implements common.Runnable.
func (*Client) Start() error { return nil }

// Close implements common.Closable.
func (*Client) Close() error { return nil }

// LookupIP implements Client.
func (*Client) LookupIP(host string) ([]net.IP, error) {
	return net.LookupIP(host)
}

// LookupIPv4 implements IPv4Lookup.
func (c *Client) LookupIPv4(host string) ([]net.IP, error) {
	ips, err := c.LookupIP(host)
	if err != nil {
		return nil, err
	}
	var ipv4 []net.IP
	for _, ip := range ips {
		parsed := net.IPAddress(ip)
		if parsed.Family().IsIPv4() {
			ipv4 = append(ipv4, parsed.IP())
		}
	}
	return ipv4, nil
}

// LookupIPv6 implements IPv6Lookup.
func (c *Client) LookupIPv6(host string) ([]net.IP, error) {
	ips, err := c.LookupIP(host)
	if err != nil {
		return nil, err
	}
	var ipv6 []net.IP
	for _, ip := range ips {
		parsed := net.IPAddress(ip)
		if parsed.Family().IsIPv6() {
			ipv6 = append(ipv6, parsed.IP())
		}
	}
	return ipv6, nil
}

// New create a new dns.Client that queries localhost for DNS.
func New() *Client {
	return &Client{}
}
