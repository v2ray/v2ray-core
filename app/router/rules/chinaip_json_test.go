// +build json

package rules_test

import (
	"net"
	"testing"

	. "v2ray.com/core/app/router/rules"
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/testing/assert"
)

func makeDestination(ip string) v2net.Destination {
	return v2net.TCPDestination(v2net.IPAddress(net.ParseIP(ip)), 80)
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
	assert.Bool(cond.Apply(makeDestination("121.14.1.189"))).IsTrue()    // sina.com.cn
	assert.Bool(cond.Apply(makeDestination("101.226.103.106"))).IsTrue() // qq.com
	assert.Bool(cond.Apply(makeDestination("115.239.210.36"))).IsTrue()  // image.baidu.com
	assert.Bool(cond.Apply(makeDestination("120.135.126.1"))).IsTrue()

	assert.Bool(cond.Apply(makeDestination("8.8.8.8"))).IsFalse()
}
