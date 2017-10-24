package testing

import (
	"reflect"

	"v2ray.com/core/common/net"
	"v2ray.com/ext/assert"
)

var IsIPv4 = assert.CreateMatcher(reflect.TypeOf(net.Address(nil)), reflect.ValueOf(func(a net.Address) bool {
	return a.Family().IsIPv4()
}), 1, "is IPv4")

var IsIPv6 = assert.CreateMatcher(reflect.TypeOf(net.Address(nil)), reflect.ValueOf(func(a net.Address) bool {
	return a.Family().IsIPv6()
}), 1, "is IPv6")

var IsIP = assert.CreateMatcher(reflect.TypeOf(net.Address(nil)), reflect.ValueOf(func(a net.Address) bool {
	return a.Family().IsIPv4() || a.Family().IsIPv6()
}), 1, "is IP")

var IsTCP = assert.CreateMatcher(reflect.TypeOf(net.Destination{}), reflect.ValueOf(func(a net.Destination) bool {
	return a.Network == net.Network_TCP
}), 1, "is TCP")

var IsUDP = assert.CreateMatcher(reflect.TypeOf(net.Destination{}), reflect.ValueOf(func(a net.Destination) bool {
	return a.Network == net.Network_UDP
}), 1, "is UDP")

var IsDomain = assert.CreateMatcher(reflect.TypeOf(net.Address(nil)), reflect.ValueOf(func(a net.Address) bool {
	return a.Family().IsDomain()
}), 1, "is Domain")

func init() {
	assert.RegisterEqualsMatcher(reflect.TypeOf((*net.Address)(nil)).Elem(), reflect.ValueOf(func(a, b net.Address) bool {
		return a == b
	}))

	assert.RegisterEqualsMatcher(reflect.TypeOf(net.Destination{}), reflect.ValueOf(func(a, b net.Destination) bool {
		return a == b
	}))

	assert.RegisterEqualsMatcher(reflect.TypeOf(net.Port(0)), reflect.ValueOf(func(a, b net.Port) bool {
		return a == b
	}))
}
