// +build json

package rules_test

import (
	"testing"

	. "github.com/v2ray/v2ray-core/app/router/rules"
	v2testing "github.com/v2ray/v2ray-core/testing"
	"github.com/v2ray/v2ray-core/testing/assert"
)

func TestChinaSitesJson(t *testing.T) {
	v2testing.Current(t)

	rule := ParseRule([]byte(`{
    "type": "chinasites",
    "outboundTag": "y"
  }`))
	assert.StringLiteral(rule.Tag).Equals("y")
	assert.Bool(rule.Apply(makeDomainDestination("v.qq.com"))).IsTrue()
	assert.Bool(rule.Apply(makeDomainDestination("www.163.com"))).IsTrue()
	assert.Bool(rule.Apply(makeDomainDestination("ngacn.cc"))).IsTrue()
	assert.Bool(rule.Apply(makeDomainDestination("12306.cn"))).IsTrue()

	assert.Bool(rule.Apply(makeDomainDestination("v2ray.com"))).IsFalse()
}
