package dns

import (
	"strconv"
	"strings"

	"v2ray.com/core/common/net"
	"v2ray.com/core/common/strmatcher"
)

var domainMapper = make(map[string][]net.Address)
var ipMapper = make(map[uint32]string)
var matcher = strmatcher.NewOrMatcher()
var ip net.IP

// Begin of ipMapper index
var begin uint32
var upperLimit uint32

func getIPSum(i net.IP) uint32 {
	return (uint32(i[0]) << 24) | (uint32(i[1]) << 16) | (uint32(i[2]) << 8) | uint32(i[3])
}

// InitFakeIPServer initializes matcher for domain name checking
func InitFakeIPServer(fake *Config_Fake, externalRules map[string][]string) error {
	if fake != nil {
		if fake.FakeRules == nil {
			return newError("no rules for fake ip").AtWarning()
		}
		for _, pattern := range fake.FakeRules {
			if err := matcher.ParsePattern(pattern, externalRules); err != nil {
				newError("failed to parse pattern: ", pattern).Base(err).AtWarning().WriteToLog()
			}
		}
		nd := strings.Split(fake.FakeNet, "/")
		mask, err := strconv.Atoi(nd[1])
		if err != nil {
			return newError("failed to parse fakeNet: ", fake.FakeNet).Base(err).AtWarning()
		}
		upperLimit = (1 << (32 - mask)) - 1
		ip = ((net.ParseAddress(nd[0])).IP())
		ip[0] &= ^byte(upperLimit >> 24)
		ip[1] &= ^byte(upperLimit >> 16)
		ip[2] &= ^byte(upperLimit >> 8)
		ip[3] &= ^byte(upperLimit)
		begin = getIPSum(ip)
	}
	return nil
}

// GetFakeIPForDomain checks if we should response with a fake IP for a domain name
func GetFakeIPForDomain(domain string) []net.Address {
	if !matcher.Match(domain) {
		return nil
	}

	if domainMapper[domain] == nil {
		if uint32(len(ipMapper)) >= upperLimit {
			return nil
		}
		var add byte = 1
		tmp := ip[3] / 255
		ip[3] -= tmp * 255
		ip[3] += add
		add = tmp
		tmp = ip[2] / 255
		ip[2] -= tmp * 255
		ip[2] += add
		add = tmp
		tmp = ip[1] / 255
		ip[1] -= tmp * 255
		ip[1] += add
		add = tmp
		tmp = ip[0] / 255
		ip[0] -= tmp * 255
		ip[0] += add
		ipMapper[begin|uint32(len(ipMapper)+1)] = domain
		domainMapper[domain] = []net.Address{net.IPAddress(ip)}
	}
	return domainMapper[domain]
}

// GetDomainForFakeIP checks if a IP is a fake IP and return its corresponding domain name
func GetDomainForFakeIP(ip net.Address) string {
	if len(ipMapper) == 0 {
		return ""
	}
	return ipMapper[begin|getIPSum(ip.IP())]
}

// ResetFakeIPServer is for testing only
func ResetFakeIPServer() {
	domainMapper = make(map[string][]net.Address)
	ipMapper = make(map[uint32]string)
	matcher = new(strmatcher.OrMatcher)
}
