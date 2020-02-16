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
	ips, err := net.LookupIP(host)
	if err != nil {
		return nil, err
	}
	parsedIPs := make([]net.IP, 0, len(ips))
	for _, ip := range ips {
		parsed := net.IPAddress(ip)
		if parsed != nil {
			parsedIPs = append(parsedIPs, parsed.IP())
		}
	}
	if len(parsedIPs) == 0 {
		return nil, dns.ErrEmptyResponse
	}
	return parsedIPs, nil
}

// LookupIPv4 implements IPv4Lookup.
func (c *Client) LookupIPv4(host string) ([]net.IP, error) {
	ips, err := c.LookupIP(host)
	if err != nil {
		return nil, err
	}
	ipv4 := make([]net.IP, 0, len(ips))
	for _, ip := range ips {
		if len(ip) == net.IPv4len {
			ipv4 = append(ipv4, ip)
		}
	}
	if len(ipv4) == 0 {
		return nil, dns.ErrEmptyResponse
	}
	return ipv4, nil
}

// LookupIPv6 implements IPv6Lookup.
func (c *Client) LookupIPv6(host string) ([]net.IP, error) {
	ips, err := c.LookupIP(host)
	if err != nil {
		return nil, err
	}
	ipv6 := make([]net.IP, 0, len(ips))
	for _, ip := range ips {
		if len(ip) == net.IPv6len {
			ipv6 = append(ipv6, ip)
		}
	}
	if len(ipv6) == 0 {
		return nil, dns.ErrEmptyResponse
	}
	return ipv6, nil
}

// LookupRealIP implements Client.
func (c *Client) LookupRealIP(host string) ([]net.IP, error) {
	return c.LookupIP(host)
}

// LookupRealIPv4 implements IPv4Lookup.
func (c *Client) LookupRealIPv4(host string) ([]net.IP, error) {
	return c.LookupIPv4(host)
}

// LookupRealIPv6 implements IPv6Lookup.
func (c *Client) LookupRealIPv6(host string) ([]net.IP, error) {
	return c.LookupIPv6(host)
}

// New create a new dns.Client that queries localhost for DNS.
func New() *Client {
	return &Client{}
}
