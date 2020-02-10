package dns

import (
	"strconv"

	"v2ray.com/core/common/net"
)

var domainMapper map[string][]net.Address = make(map[string][]net.Address)
var ipMapper map[net.Address]string = make(map[net.Address]string)

func GetFakeIPForDomain(domain string) []net.Address {
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
