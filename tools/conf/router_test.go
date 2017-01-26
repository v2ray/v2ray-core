package conf_test

import (
	"net"
	"testing"

	"context"

	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/proxy"
	"v2ray.com/core/testing/assert"
	. "v2ray.com/core/tools/conf"
)

func makeDestination(ip string) v2net.Destination {
	return v2net.TCPDestination(v2net.IPAddress(net.ParseIP(ip)), 80)
}

func makeDomainDestination(domain string) v2net.Destination {
	return v2net.TCPDestination(v2net.DomainAddress(domain), 80)
}

func TestChinaIPJson(t *testing.T) {
	assert := assert.On(t)

	rule := ParseRule([]byte(`{
    "type": "chinaip",
    "outboundTag": "x"
  }`))
	assert.String(rule.Tag).Equals("x")
	cond, err := rule.BuildCondition()
	assert.Error(err).IsNil()
	assert.Bool(cond.Apply(proxy.ContextWithDestination(context.Background(), v2net.TCPDestination(v2net.ParseAddress("121.14.1.189"), 80)))).IsTrue()    // sina.com.cn
	assert.Bool(cond.Apply(proxy.ContextWithDestination(context.Background(), v2net.TCPDestination(v2net.ParseAddress("101.226.103.106"), 80)))).IsTrue() // qq.com
	assert.Bool(cond.Apply(proxy.ContextWithDestination(context.Background(), v2net.TCPDestination(v2net.ParseAddress("115.239.210.36"), 80)))).IsTrue()  // image.baidu.com
	assert.Bool(cond.Apply(proxy.ContextWithDestination(context.Background(), v2net.TCPDestination(v2net.ParseAddress("120.135.126.1"), 80)))).IsTrue()

	assert.Bool(cond.Apply(proxy.ContextWithDestination(context.Background(), v2net.TCPDestination(v2net.ParseAddress("8.8.8.8"), 80)))).IsFalse()
}

func TestChinaSitesJson(t *testing.T) {
	assert := assert.On(t)

	rule := ParseRule([]byte(`{
    "type": "chinasites",
    "outboundTag": "y"
  }`))
	assert.String(rule.Tag).Equals("y")
	cond, err := rule.BuildCondition()
	assert.Error(err).IsNil()
	assert.Bool(cond.Apply(proxy.ContextWithDestination(context.Background(), v2net.TCPDestination(v2net.ParseAddress("v.qq.com"), 80)))).IsTrue()
	assert.Bool(cond.Apply(proxy.ContextWithDestination(context.Background(), v2net.TCPDestination(v2net.ParseAddress("www.163.com"), 80)))).IsTrue()
	assert.Bool(cond.Apply(proxy.ContextWithDestination(context.Background(), v2net.TCPDestination(v2net.ParseAddress("ngacn.cc"), 80)))).IsTrue()
	assert.Bool(cond.Apply(proxy.ContextWithDestination(context.Background(), v2net.TCPDestination(v2net.ParseAddress("12306.cn"), 80)))).IsTrue()

	assert.Bool(cond.Apply(proxy.ContextWithDestination(context.Background(), v2net.TCPDestination(v2net.ParseAddress("v2ray.com"), 80)))).IsFalse()
}

func TestDomainRule(t *testing.T) {
	assert := assert.On(t)

	rule := ParseRule([]byte(`{
    "type": "field",
    "domain": [
      "ooxx.com",
      "oxox.com",
      "regexp:\\.cn$"
    ],
    "network": "tcp",
    "outboundTag": "direct"
  }`))
	assert.Pointer(rule).IsNotNil()
	cond, err := rule.BuildCondition()
	assert.Error(err).IsNil()
	assert.Bool(cond.Apply(proxy.ContextWithDestination(context.Background(), v2net.TCPDestination(v2net.ParseAddress("www.ooxx.com"), 80)))).IsTrue()
	assert.Bool(cond.Apply(proxy.ContextWithDestination(context.Background(), v2net.TCPDestination(v2net.ParseAddress("www.aabb.com"), 80)))).IsFalse()
	assert.Bool(cond.Apply(proxy.ContextWithDestination(context.Background(), v2net.TCPDestination(v2net.IPAddress([]byte{127, 0, 0, 1}), 80)))).IsFalse()
	assert.Bool(cond.Apply(proxy.ContextWithDestination(context.Background(), v2net.TCPDestination(v2net.ParseAddress("www.12306.cn"), 80)))).IsTrue()
	assert.Bool(cond.Apply(proxy.ContextWithDestination(context.Background(), v2net.TCPDestination(v2net.ParseAddress("www.acn.com"), 80)))).IsFalse()
}

func TestIPRule(t *testing.T) {
	assert := assert.On(t)

	rule := ParseRule([]byte(`{
    "type": "field",
    "ip": [
      "10.0.0.0/8",
      "192.0.0.0/24"
    ],
    "network": "tcp",
    "outboundTag": "direct"
  }`))
	assert.Pointer(rule).IsNotNil()
	cond, err := rule.BuildCondition()
	assert.Error(err).IsNil()
	assert.Bool(cond.Apply(proxy.ContextWithDestination(context.Background(), v2net.TCPDestination(v2net.DomainAddress("www.ooxx.com"), 80)))).IsFalse()
	assert.Bool(cond.Apply(proxy.ContextWithDestination(context.Background(), v2net.TCPDestination(v2net.IPAddress([]byte{10, 0, 0, 1}), 80)))).IsTrue()
	assert.Bool(cond.Apply(proxy.ContextWithDestination(context.Background(), v2net.TCPDestination(v2net.IPAddress([]byte{127, 0, 0, 1}), 80)))).IsFalse()
	assert.Bool(cond.Apply(proxy.ContextWithDestination(context.Background(), v2net.TCPDestination(v2net.IPAddress([]byte{192, 0, 0, 1}), 80)))).IsTrue()
}

func TestSourceIPRule(t *testing.T) {
	assert := assert.On(t)

	rule := ParseRule([]byte(`{
    "type": "field",
    "source": [
      "10.0.0.0/8",
      "192.0.0.0/24"
    ],
    "outboundTag": "direct"
  }`))
	assert.Pointer(rule).IsNotNil()
	cond, err := rule.BuildCondition()
	assert.Error(err).IsNil()
	assert.Bool(cond.Apply(proxy.ContextWithSource(context.Background(), v2net.TCPDestination(v2net.DomainAddress("www.ooxx.com"), 80)))).IsFalse()
	assert.Bool(cond.Apply(proxy.ContextWithSource(context.Background(), v2net.TCPDestination(v2net.IPAddress([]byte{10, 0, 0, 1}), 80)))).IsTrue()
	assert.Bool(cond.Apply(proxy.ContextWithSource(context.Background(), v2net.TCPDestination(v2net.IPAddress([]byte{127, 0, 0, 1}), 80)))).IsFalse()
	assert.Bool(cond.Apply(proxy.ContextWithSource(context.Background(), v2net.TCPDestination(v2net.IPAddress([]byte{192, 0, 0, 1}), 80)))).IsTrue()
}
