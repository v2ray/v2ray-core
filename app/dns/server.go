// +build !confonly

package dns

//go:generate errorgen

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"v2ray.com/core"
	"v2ray.com/core/app/router"
	"v2ray.com/core/common"
	"v2ray.com/core/common/errors"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/session"
	"v2ray.com/core/common/strmatcher"
	"v2ray.com/core/common/uuid"
	"v2ray.com/core/features"
	"v2ray.com/core/features/dns"
	"v2ray.com/core/features/routing"
)

// Server is a DNS rely server.
type Server struct {
	sync.Mutex
	hosts          *StaticHosts
	clients        []Client
	clientIP       net.IP
	domainMatcher  strmatcher.IndexMatcher
	domainIndexMap map[uint32]uint32
	ipIndexMap     map[uint32]*MultiGeoIPMatcher
	tag            string
}

// MultiGeoIPMatcher for match
type MultiGeoIPMatcher struct {
	matchers []*router.GeoIPMatcher
}

var errExpectedIPNonMatch = errors.New("expected ip not match")

// Match check ip match
func (c *MultiGeoIPMatcher) Match(ip net.IP) bool {
	for _, matcher := range c.matchers {
		if matcher.Match(ip) {
			return true
		}
	}
	return false
}

// HasMatcher check has matcher
func (c *MultiGeoIPMatcher) HasMatcher() bool {
	return len(c.matchers) > 0
}

func generateRandomTag() string {
	id := uuid.New()
	return "v2ray.system." + id.String()
}

// New creates a new DNS server with given configuration.
func New(ctx context.Context, config *Config) (*Server, error) {
	server := &Server{
		clients: make([]Client, 0, len(config.NameServers)+len(config.NameServer)),
		tag:     config.Tag,
	}
	if server.tag == "" {
		server.tag = generateRandomTag()
	}
	if len(config.ClientIp) > 0 {
		if len(config.ClientIp) != 4 && len(config.ClientIp) != 16 {
			return nil, newError("unexpected IP length", len(config.ClientIp))
		}
		server.clientIP = net.IP(config.ClientIp)
	}

	hosts, err := NewStaticHosts(config.StaticHosts, config.Hosts)
	if err != nil {
		return nil, newError("failed to create hosts").Base(err)
	}
	server.hosts = hosts

	addNameServer := func(endpoint *net.Endpoint) int {
		address := endpoint.Address.AsAddress()
		if address.Family().IsDomain() && address.Domain() == "localhost" {
			server.clients = append(server.clients, NewLocalNameServer())
			newError("DNS: localhost inited").AtInfo().WriteToLog()
		} else if address.Family().IsDomain() && strings.HasPrefix(address.Domain(), "DOHL_") {
			dohHost := address.Domain()[5:]
			server.clients = append(server.clients, NewDoHLocalNameServer(dohHost, server.clientIP))
			newError("DNS: DOH - Local inited for https://", dohHost).AtInfo().WriteToLog()
		} else if address.Family().IsDomain() && strings.HasPrefix(address.Domain(), "DOH_") {
			// DOH_ prefix makes net.Address think it's a domain
			// need to process the real address here.
			dohHost := address.Domain()[4:]
			dohAddr := net.ParseAddress(dohHost)
			dohIP := dohHost
			var dests []net.Destination

			if dohAddr.Family().IsDomain() {
				// resolve DOH server in advance
				ips, err := net.LookupIP(dohAddr.Domain())
				if err != nil || len(ips) == 0 {
					return 0
				}
				for _, ip := range ips {
					dohIP := ip.String()
					if len(ip) == net.IPv6len {
						dohIP = fmt.Sprintf("[%s]", dohIP)
					}
					dohdest, _ := net.ParseDestination(fmt.Sprintf("tcp:%s:443", dohIP))
					dests = append(dests, dohdest)
				}
			} else {
				// rfc8484, DOH service only use port 443
				dest, err := net.ParseDestination(fmt.Sprintf("tcp:%s:443", dohIP))
				if err != nil {
					return 0
				}
				dests = []net.Destination{dest}
			}

			// need the core dispatcher, register DOHClient at callback
			idx := len(server.clients)
			server.clients = append(server.clients, nil)
			common.Must(core.RequireFeatures(ctx, func(d routing.Dispatcher) {
				server.clients[idx] = NewDoHNameServer(dests, dohHost, d, server.clientIP)
				newError("DNS: DOH - Remote client inited for https://", dohHost).AtInfo().WriteToLog()
			}))
		} else {
			dest := endpoint.AsDestination()
			if dest.Network == net.Network_Unknown {
				dest.Network = net.Network_UDP
			}
			if dest.Network == net.Network_UDP {
				idx := len(server.clients)
				server.clients = append(server.clients, nil)

				common.Must(core.RequireFeatures(ctx, func(d routing.Dispatcher) {
					server.clients[idx] = NewClassicNameServer(dest, d, server.clientIP)
				}))
			}
			newError("DNS: UDP client inited for ", dest.NetAddr()).AtInfo().WriteToLog()
		}
		return len(server.clients) - 1
	}

	if len(config.NameServers) > 0 {
		features.PrintDeprecatedFeatureWarning("simple DNS server")
		for _, destPB := range config.NameServers {
			addNameServer(destPB)
		}
	}

	if len(config.NameServer) > 0 {
		domainMatcher := &strmatcher.MatcherGroup{}
		domainIndexMap := make(map[uint32]uint32)
		ipIndexMap := make(map[uint32]*MultiGeoIPMatcher)
		var geoIPMatcherContainer router.GeoIPMatcherContainer

		for _, ns := range config.NameServer {
			idx := addNameServer(ns.Address)

			for _, domain := range ns.PrioritizedDomain {
				matcher, err := toStrMatcher(domain.Type, domain.Domain)
				if err != nil {
					return nil, newError("failed to create prioritized domain").Base(err).AtWarning()
				}
				midx := domainMatcher.Add(matcher)
				domainIndexMap[midx] = uint32(idx)
			}

			var matchers []*router.GeoIPMatcher
			for _, geoip := range ns.Geoip {
				matcher, err := geoIPMatcherContainer.Add(geoip)
				if err != nil {
					return nil, newError("failed to create ip matcher").Base(err).AtWarning()
				}
				matchers = append(matchers, matcher)
			}
			matcher := &MultiGeoIPMatcher{matchers: matchers}
			ipIndexMap[uint32(idx)] = matcher
		}

		server.domainMatcher = domainMatcher
		server.domainIndexMap = domainIndexMap
		server.ipIndexMap = ipIndexMap
	}

	if len(server.clients) == 0 {
		server.clients = append(server.clients, NewLocalNameServer())
	}

	return server, nil
}

// Type implements common.HasType.
func (*Server) Type() interface{} {
	return dns.ClientType()
}

// Start implements common.Runnable.
func (s *Server) Start() error {
	return nil
}

// Close implements common.Closable.
func (s *Server) Close() error {
	return nil
}

func (s *Server) IsOwnLink(ctx context.Context) bool {
	inbound := session.InboundFromContext(ctx)
	return inbound != nil && inbound.Tag == s.tag
}

// Match check dns ip match geoip
func (s *Server) Match(idx uint32, client Client, domain string, ips []net.IP) ([]net.IP, error) {
	matcher, exist := s.ipIndexMap[idx]
	if !exist {
		newError("domain ", domain, " server not in ipIndexMap: ", client.Name(), " idx:", idx, " just return").AtDebug().WriteToLog()
		return ips, nil
	}

	if !matcher.HasMatcher() {
		newError("domain ", domain, " server has not valid matcher: ", client.Name(), " idx:", idx, " just return").AtDebug().WriteToLog()
		return ips, nil
	}

	newIps := []net.IP{}
	for _, ip := range ips {
		if matcher.Match(ip) {
			newIps = append(newIps, ip)
			newError("domain ", domain, " ip ", ip, " is match at server ", client.Name(), " idx:", idx).AtDebug().WriteToLog()
		} else {
			newError("domain ", domain, " ip ", ip, " is not match at server ", client.Name(), " idx:", idx).AtDebug().WriteToLog()
		}
	}
	if len(newIps) == 0 {
		return nil, errExpectedIPNonMatch
	}
	return newIps, nil
}

func (s *Server) queryIPTimeout(idx uint32, client Client, domain string, option IPOption) ([]net.IP, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	if len(s.tag) > 0 {
		ctx = session.ContextWithInbound(ctx, &session.Inbound{
			Tag: s.tag,
		})
	}
	ips, err := client.QueryIP(ctx, domain, option)
	cancel()

	if err != nil {
		return ips, err
	}

	ips, err = s.Match(idx, client, domain, ips)
	return ips, err
}

// LookupIP implements dns.Client.
func (s *Server) LookupIP(domain string) ([]net.IP, error) {
	return s.lookupIPInternal(domain, IPOption{
		IPv4Enable: true,
		IPv6Enable: true,
	})
}

// LookupIPv4 implements dns.IPv4Lookup.
func (s *Server) LookupIPv4(domain string) ([]net.IP, error) {
	return s.lookupIPInternal(domain, IPOption{
		IPv4Enable: true,
		IPv6Enable: false,
	})
}

// LookupIPv6 implements dns.IPv6Lookup.
func (s *Server) LookupIPv6(domain string) ([]net.IP, error) {
	return s.lookupIPInternal(domain, IPOption{
		IPv4Enable: false,
		IPv6Enable: true,
	})
}

func (s *Server) lookupStatic(domain string, option IPOption, depth int32) []net.Address {
	ips := s.hosts.LookupIP(domain, option)
	if ips == nil {
		return nil
	}
	if ips[0].Family().IsDomain() && depth < 5 {
		if newIPs := s.lookupStatic(ips[0].Domain(), option, depth+1); newIPs != nil {
			return newIPs
		}
	}
	return ips
}

func toNetIP(ips []net.Address) []net.IP {
	if len(ips) == 0 {
		return nil
	}
	netips := make([]net.IP, 0, len(ips))
	for _, ip := range ips {
		netips = append(netips, ip.IP())
	}
	return netips
}

func (s *Server) lookupIPInternal(domain string, option IPOption) ([]net.IP, error) {
	if domain == "" {
		return nil, newError("empty domain name")
	}

	// normalize the FQDN form query
	if domain[len(domain)-1] == '.' {
		domain = domain[:len(domain)-1]
	}

	// skip domain without any dot
	if strings.Index(domain, ".") == -1 {
		return nil, newError("invalid domain name")
	}

	ips := s.lookupStatic(domain, option, 0)
	if ips != nil && ips[0].Family().IsIP() {
		newError("returning ", len(ips), " IPs for domain ", domain).WriteToLog()
		return toNetIP(ips), nil
	}

	if ips != nil && ips[0].Family().IsDomain() {
		newdomain := ips[0].Domain()
		newError("domain replaced: ", domain, " -> ", newdomain).WriteToLog()
		domain = newdomain
	}

	var lastErr error
	var matchedClient Client
	if s.domainMatcher != nil {
		idx := s.domainMatcher.Match(domain)
		if idx > 0 {
			matchedClient = s.clients[s.domainIndexMap[idx]]
			newError("domain matched, direct lookup ip for domain ", domain, " at ", matchedClient.Name()).WriteToLog()
			ips, err := s.queryIPTimeout(s.domainIndexMap[idx], matchedClient, domain, option)
			if len(ips) > 0 {
				return ips, nil
			}
			if err == dns.ErrEmptyResponse {
				return nil, err
			}
			if err != nil {
				newError("failed to lookup ip for domain ", domain, " at server ", matchedClient.Name()).Base(err).WriteToLog()
				lastErr = err
			}
		}
	}

	for idx, client := range s.clients {
		if client == matchedClient {
			newError("domain ", domain, " at server ", client.Name(), " idx:", idx, " already lookup failed, just ignore").AtDebug().WriteToLog()
			continue
		}

		newError("try to lookup ip for domain ", domain, " at server ", client.Name(), " idx:", idx).AtDebug().WriteToLog()
		ips, err := s.queryIPTimeout(uint32(idx), client, domain, option)
		if len(ips) > 0 {
			newError("lookup ip for domain ", domain, " success: ", ips, " at server ", client.Name(), " idx:", idx).AtDebug().WriteToLog()
			return ips, nil
		}

		if err != nil {
			newError("failed to lookup ip for domain ", domain, " at server ", client.Name()).Base(err).WriteToLog()
			lastErr = err
		}
		if err != context.Canceled && err != context.DeadlineExceeded && err != errExpectedIPNonMatch {
			return nil, err
		}
	}

	return nil, dns.ErrEmptyResponse.Base(lastErr)
}

func init() {
	common.Must(common.RegisterConfig((*Config)(nil), func(ctx context.Context, config interface{}) (interface{}, error) {
		return New(ctx, config.(*Config))
	}))
}
