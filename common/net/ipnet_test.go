package net_test

import (
	"net"
	"testing"

	. "github.com/v2ray/v2ray-core/common/net"
	v2testing "github.com/v2ray/v2ray-core/testing"
	"github.com/v2ray/v2ray-core/testing/assert"
)

func parseCIDR(str string) *net.IPNet {
	_, ipNet, err := net.ParseCIDR(str)
	assert.Error(err).IsNil()
	return ipNet
}

func TestIPNet(t *testing.T) {
	v2testing.Current(t)

	ipNet := NewIPNet()
	ipNet.Add(parseCIDR(("0.0.0.0/8")))
	ipNet.Add(parseCIDR(("10.0.0.0/8")))
	ipNet.Add(parseCIDR(("100.64.0.0/10")))
	ipNet.Add(parseCIDR(("127.0.0.0/8")))
	ipNet.Add(parseCIDR(("169.254.0.0/16")))
	ipNet.Add(parseCIDR(("172.16.0.0/12")))
	ipNet.Add(parseCIDR(("192.0.0.0/24")))
	ipNet.Add(parseCIDR(("192.0.2.0/24")))
	ipNet.Add(parseCIDR(("192.168.0.0/16")))
	ipNet.Add(parseCIDR(("198.18.0.0/15")))
	ipNet.Add(parseCIDR(("198.51.100.0/24")))
	ipNet.Add(parseCIDR(("203.0.113.0/24")))
	ipNet.Add(parseCIDR(("8.8.8.8/32")))
	assert.Bool(ipNet.Contains(net.ParseIP("192.168.1.1"))).IsTrue()
	assert.Bool(ipNet.Contains(net.ParseIP("192.0.0.0"))).IsTrue()
	assert.Bool(ipNet.Contains(net.ParseIP("192.0.1.0"))).IsFalse()
	assert.Bool(ipNet.Contains(net.ParseIP("0.1.0.0"))).IsTrue()
	assert.Bool(ipNet.Contains(net.ParseIP("1.0.0.1"))).IsFalse()
	assert.Bool(ipNet.Contains(net.ParseIP("8.8.8.7"))).IsFalse()
	assert.Bool(ipNet.Contains(net.ParseIP("8.8.8.8"))).IsTrue()
}
