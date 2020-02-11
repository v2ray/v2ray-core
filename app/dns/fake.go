package dns

import (
	"strconv"

	"v2ray.com/core/common/net"
	"v2ray.com/core/common/strmatcher"
)

var domainMapper map[string][]net.Address = make(map[string][]net.Address)
var ipMapper map[net.Address]string = make(map[net.Address]string)
var matcher *strmatcher.OrMatcher = new(strmatcher.OrMatcher)

func InitFakeIPServer(patterns []string, externalRules map[string][]string) {
	matcher.New()
	for _, pattern := range patterns {
		matcher.ParsePattern(pattern, externalRules)
	}
}
func GetFakeIPForDomain(domain string) []net.Address {
	if !matcher.Match(domain) {
		return nil
	}

	if domainMapper[domain] == nil {
		addressCounter := len(ipMapper) + 1
		as := "224."
		as += strconv.Itoa((0xff0000&addressCounter)>>16) + "."
		as += strconv.Itoa((0xff00&addressCounter)>>8) + "."
		as += strconv.Itoa(0xff & addressCounter)
		ip := net.ParseAddress(as)
		ipMapper[ip] = domain
		domainMapper[domain] = []net.Address{ip}
	}
	return domainMapper[domain]
}
func GetDomainForFakeIP(ip net.Address) string {
	if len(ipMapper) == 0 {
		return ""
	}
	return ipMapper[ip]
}
