package testing

import (
	"v2ray.com/core/common/net"
	"v2ray.com/ext/assert"
)

var IsIPv4 = assert.CreateMatcher(func(a net.Address) bool {
	return a.Family().IsIPv4()
}, "is IPv4")

var IsIPv6 = assert.CreateMatcher(func(a net.Address) bool {
	return a.Family().IsIPv6()
}, "is IPv6")

var IsIP = assert.CreateMatcher(func(a net.Address) bool {
	return a.Family().IsIPv4() || a.Family().IsIPv6()
}, "is IP")

var IsTCP = assert.CreateMatcher(func(a net.Destination) bool {
	return a.Network == net.Network_TCP
}, "is TCP")

var IsUDP = assert.CreateMatcher(func(a net.Destination) bool {
	return a.Network == net.Network_UDP
}, "is UDP")

var IsDomain = assert.CreateMatcher(func(a net.Address) bool {
	return a.Family().IsDomain()
}, "is Domain")

func init() {
	assert.RegisterEqualsMatcher(func(a, b net.Address) bool {
		return a == b
	})

	assert.RegisterEqualsMatcher(func(a, b net.Destination) bool {
		return a == b
	})

	assert.RegisterEqualsMatcher(func(a, b net.Port) bool {
		return a == b
	})

	assert.RegisterEqualsMatcher(func(a, b net.IP) bool {
		return a.Equal(b)
	})
}
