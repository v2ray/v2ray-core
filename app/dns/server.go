// +build !confonly

package dns

//go:generate go run v2ray.com/core/common/errors/errorgen

import (
	"context"
	"fmt"
	"log"
	"net/url"
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
	hosts         *StaticHosts
	clientIP      net.IP
	clients       []Client             // clientIdx -> Client
	ipIndexMap    []*MultiGeoIPMatcher // clientIdx -> *MultiGeoIPMatcher
	domainRules   [][]string           // clientIdx -> domainRuleIdx -> DomainRule
	domainMatcher strmatcher.IndexMatcher
	matcherInfos  []DomainMatcherInfo // matcherIdx -> DomainMatcherInfo
	tag           string
}

// DomainMatcherInfo contains information attached to index returned by Server.domainMatcher
type DomainMatcherInfo struct {
	clientIdx     uint16
	domainRuleIdx uint16
}

// MultiGeoIPMatcher for match
type MultiGeoIPMatcher struct {
	matchers []*router.GeoIPMatcher
}

var errExpectedIPNonMatch = errors.New("expectIPs not match")

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
		if len(config.ClientIp) != net.IPv4len && len(config.ClientIp) != net.IPv6len {
			return nil, newError("unexpected IP length", len(config.ClientIp))
		}
		server.clientIP = net.IP(config.ClientIp)
	}

	hosts, err := NewStaticHosts(config.StaticHosts, config.Hosts)
	if err != nil {
		return nil, newError("failed to create hosts").Base(err)
	}
	server.hosts = hosts

	addNameServer := func(ns *NameServer) int {
		endpoint := ns.Address
		address := endpoint.Address.AsAddress()
		if address.Family().IsDomain() && address.Domain() == "localhost" {
			server.clients = append(server.clients, NewLocalNameServer())
			// Priotize local domains with specific TLDs or without any dot to local DNS
			// References:
			// https://www.iana.org/assignments/special-use-domain-names/special-use-domain-names.xhtml
			// https://unix.stackexchange.com/questions/92441/whats-the-difference-between-local-home-and-lan
			localTLDsAndDotlessDomains := []*NameServer_PriorityDomain{
				{Type: DomainMatchingType_Regex, Domain: "^[^.]+$"}, // This will only match domains without any dot
				{Type: DomainMatchingType_Subdomain, Domain: "local"},
				{Type: DomainMatchingType_Subdomain, Domain: "localdomain"},
				{Type: DomainMatchingType_Subdomain, Domain: "localhost"},
				{Type: DomainMatchingType_Subdomain, Domain: "lan"},
				{Type: DomainMatchingType_Subdomain, Domain: "home.arpa"},
				{Type: DomainMatchingType_Subdomain, Domain: "example"},
				{Type: DomainMatchingType_Subdomain, Domain: "invalid"},
				{Type: DomainMatchingType_Subdomain, Domain: "test"},
			}
			ns.PrioritizedDomain = append(ns.PrioritizedDomain, localTLDsAndDotlessDomains...)
		} else if address.Family().IsDomain() && strings.HasPrefix(address.Domain(), "https+local://") {
			// URI schemed string treated as domain
			// DOH Local mode
			u, err := url.Parse(address.Domain())
			if err != nil {
				log.Fatalln(newError("DNS config error").Base(err))
			}
			server.clients = append(server.clients, NewDoHLocalNameServer(u, server.clientIP))
		} else if address.Family().IsDomain() && strings.HasPrefix(address.Domain(), "https://") {
			// DOH Remote mode
			u, err := url.Parse(address.Domain())
			if err != nil {
				log.Fatalln(newError("DNS config error").Base(err))
			}
			idx := len(server.clients)
			server.clients = append(server.clients, nil)

			// need the core dispatcher, register DOHClient at callback
			common.Must(core.RequireFeatures(ctx, func(d routing.Dispatcher) {
				c, err := NewDoHNameServer(u, d, server.clientIP)
				if err != nil {
					log.Fatalln(newError("DNS config error").Base(err))
				}
				server.clients[idx] = c
			}))
		} else {
			// UDP classic DNS mode
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
		}
		server.ipIndexMap = append(server.ipIndexMap, nil)
		return len(server.clients) - 1
	}

	if len(config.NameServers) > 0 {
		features.PrintDeprecatedFeatureWarning("simple DNS server")
		for _, destPB := range config.NameServers {
			addNameServer(&NameServer{Address: destPB})
		}
	}

	if len(config.NameServer) > 0 {
		clientIndices := []int{}
		domainRuleCount := 0
		for _, ns := range config.NameServer {
			idx := addNameServer(ns)
			clientIndices = append(clientIndices, idx)
			domainRuleCount += len(ns.PrioritizedDomain)
		}

		domainRules := make([][]string, len(server.clients))
		domainMatcher := &strmatcher.MatcherGroup{}
		matcherInfos := make([]DomainMatcherInfo, domainRuleCount+1) // matcher index starts from 1
		var geoIPMatcherContainer router.GeoIPMatcherContainer
		for nidx, ns := range config.NameServer {
			idx := clientIndices[nidx]

			// Establish domain rule matcher
			rules := []string{}
			ruleCurr := 0
			ruleIter := 0
			for _, domain := range ns.PrioritizedDomain {
				matcher, err := toStrMatcher(domain.Type, domain.Domain)
				if err != nil {
					return nil, newError("failed to create prioritized domain").Base(err).AtWarning()
				}
				midx := domainMatcher.Add(matcher)
				if midx >= uint32(len(matcherInfos)) { // This rarely happens according to current matcher's implementation
					newError("expanding domain matcher info array to size ", midx, " when adding ", matcher).AtDebug().WriteToLog()
					matcherInfos = append(matcherInfos, make([]DomainMatcherInfo, midx-uint32(len(matcherInfos))+1)...)
				}
				info := &matcherInfos[midx]
				info.clientIdx = uint16(idx)
				if ruleCurr < len(ns.OriginalRules) {
					info.domainRuleIdx = uint16(ruleCurr)
					rule := ns.OriginalRules[ruleCurr]
					if ruleCurr >= len(rules) {
						rules = append(rules, rule.Rule)
					}
					ruleIter++
					if ruleIter >= int(rule.Size) {
						ruleIter = 0
						ruleCurr++
					}
				} else { // No original rule, generate one according to current domain matcher (majorly for compatibility with tests)
					info.domainRuleIdx = uint16(len(rules))
					rules = append(rules, matcher.String())
				}
			}
			domainRules[idx] = rules

			// only add to ipIndexMap if GeoIP is configured
			if len(ns.Geoip) > 0 {
				var matchers []*router.GeoIPMatcher
				for _, geoip := range ns.Geoip {
					matcher, err := geoIPMatcherContainer.Add(geoip)
					if err != nil {
						return nil, newError("failed to create ip matcher").Base(err).AtWarning()
					}
					matchers = append(matchers, matcher)
				}
				matcher := &MultiGeoIPMatcher{matchers: matchers}
				server.ipIndexMap[idx] = matcher
			}
		}
		server.domainRules = domainRules
		server.domainMatcher = domainMatcher
		server.matcherInfos = matcherInfos
	}

	if len(server.clients) == 0 {
		server.clients = append(server.clients, NewLocalNameServer())
		server.ipIndexMap = append(server.ipIndexMap, nil)
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
func (s *Server) Match(idx int, client Client, domain string, ips []net.IP) ([]net.IP, error) {
	var matcher *MultiGeoIPMatcher
	if idx < len(s.ipIndexMap) {
		matcher = s.ipIndexMap[idx]
	}
	if matcher == nil {
		return ips, nil
	}

	if !matcher.HasMatcher() {
		newError("domain ", domain, " server has no valid matcher: ", client.Name(), " idx:", idx).AtDebug().WriteToLog()
		return ips, nil
	}

	newIps := []net.IP{}
	for _, ip := range ips {
		if matcher.Match(ip) {
			newIps = append(newIps, ip)
		}
	}
	if len(newIps) == 0 {
		return nil, errExpectedIPNonMatch
	}
	newError("domain ", domain, " expectIPs ", newIps, " matched at server ", client.Name(), " idx:", idx).AtDebug().WriteToLog()
	return newIps, nil
}

func (s *Server) queryIPTimeout(idx int, client Client, domain string, option IPOption) ([]net.IP, error) {
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
		indices := s.domainMatcher.Match(domain)
		domainRules := []string{}
		matchingDNS := []string{}
		for _, idx := range indices {
			info := s.matcherInfos[idx]
			rule := s.domainRules[info.clientIdx][info.domainRuleIdx]
			domainRules = append(domainRules, fmt.Sprintf("%s(DNS idx:%d)", rule, info.clientIdx))
			matchingDNS = append(matchingDNS, s.clients[info.clientIdx].Name())
		}
		if len(domainRules) > 0 {
			newError("domain ", domain, " matches following rules: ", domainRules).AtDebug().WriteToLog()
		}
		if len(matchingDNS) > 0 {
			newError("domain ", domain, " uses following DNS first: ", matchingDNS).AtDebug().WriteToLog()
		}
		for _, idx := range indices {
			clientIdx := int(s.matcherInfos[idx].clientIdx)
			matchedClient = s.clients[clientIdx]
			ips, err := s.queryIPTimeout(clientIdx, matchedClient, domain, option)
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

		ips, err := s.queryIPTimeout(idx, client, domain, option)
		if len(ips) > 0 {
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

	return nil, newError("returning nil for domain ", domain).Base(lastErr)
}

func init() {
	common.Must(common.RegisterConfig((*Config)(nil), func(ctx context.Context, config interface{}) (interface{}, error) {
		return New(ctx, config.(*Config))
	}))
}
