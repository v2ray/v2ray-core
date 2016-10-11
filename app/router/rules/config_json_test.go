// +build json

package rules_test

import (
	"testing"

	. "v2ray.com/core/app/router/rules"
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/testing/assert"
)

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
	assert.Bool(cond.Apply(v2net.TCPDestination(v2net.DomainAddress("www.ooxx.com"), 80))).IsTrue()
	assert.Bool(cond.Apply(v2net.TCPDestination(v2net.DomainAddress("www.aabb.com"), 80))).IsFalse()
	assert.Bool(cond.Apply(v2net.TCPDestination(v2net.IPAddress([]byte{127, 0, 0, 1}), 80))).IsFalse()
	assert.Bool(cond.Apply(v2net.TCPDestination(v2net.DomainAddress("www.12306.cn"), 80))).IsTrue()
	assert.Bool(cond.Apply(v2net.TCPDestination(v2net.DomainAddress("www.acn.com"), 80))).IsFalse()
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
	assert.Bool(cond.Apply(v2net.TCPDestination(v2net.DomainAddress("www.ooxx.com"), 80))).IsFalse()
	assert.Bool(cond.Apply(v2net.TCPDestination(v2net.IPAddress([]byte{10, 0, 0, 1}), 80))).IsTrue()
	assert.Bool(cond.Apply(v2net.TCPDestination(v2net.IPAddress([]byte{127, 0, 0, 1}), 80))).IsFalse()
	assert.Bool(cond.Apply(v2net.TCPDestination(v2net.IPAddress([]byte{192, 0, 0, 1}), 80))).IsTrue()
}
