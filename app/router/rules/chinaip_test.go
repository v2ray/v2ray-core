package rules

import (
	"net"
	"testing"

	v2net "github.com/v2ray/v2ray-core/common/net"
	v2testing "github.com/v2ray/v2ray-core/testing"
	"github.com/v2ray/v2ray-core/testing/assert"
)

func makeDestination(ip string) v2net.Destination {
	return v2net.TCPDestination(v2net.IPAddress(net.ParseIP(ip)), 80)
}

func TestChinaIP(t *testing.T) {
	v2testing.Current(t)

	rule := NewIPv4Matcher(chinaIPNet)
	assert.Bool(rule.Apply(makeDestination("121.14.1.189"))).IsTrue()    // sina.com.cn
	assert.Bool(rule.Apply(makeDestination("101.226.103.106"))).IsTrue() // qq.com
	assert.Bool(rule.Apply(makeDestination("115.239.210.36"))).IsTrue()  // image.baidu.com
	assert.Bool(rule.Apply(makeDestination("120.135.126.1"))).IsTrue()

	assert.Bool(rule.Apply(makeDestination("8.8.8.8"))).IsFalse()
}
