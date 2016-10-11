// +build json

package rules_test

import (
	"testing"

	. "v2ray.com/core/app/router/rules"
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/testing/assert"
)

func makeDomainDestination(domain string) v2net.Destination {
	return v2net.TCPDestination(v2net.DomainAddress(domain), 80)
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
	assert.Bool(cond.Apply(makeDomainDestination("v.qq.com"))).IsTrue()
	assert.Bool(cond.Apply(makeDomainDestination("www.163.com"))).IsTrue()
	assert.Bool(cond.Apply(makeDomainDestination("ngacn.cc"))).IsTrue()
	assert.Bool(cond.Apply(makeDomainDestination("12306.cn"))).IsTrue()

	assert.Bool(cond.Apply(makeDomainDestination("v2ray.com"))).IsFalse()
}
