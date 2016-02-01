// +build json

package rules_test

import (
	"testing"

	. "github.com/v2ray/v2ray-core/app/router/rules"
	v2testing "github.com/v2ray/v2ray-core/testing"
	"github.com/v2ray/v2ray-core/testing/assert"
)

func TestChinaIPJson(t *testing.T) {
	v2testing.Current(t)

	rule := ParseRule([]byte(`{
    "type": "chinaip",
    "outboundTag": "x"
  }`))
	assert.StringLiteral(rule.Tag).Equals("x")
	assert.Bool(rule.Apply(makeDestination("121.14.1.189"))).IsTrue()    // sina.com.cn
	assert.Bool(rule.Apply(makeDestination("101.226.103.106"))).IsTrue() // qq.com
	assert.Bool(rule.Apply(makeDestination("115.239.210.36"))).IsTrue()  // image.baidu.com
	assert.Bool(rule.Apply(makeDestination("120.135.126.1"))).IsTrue()

	assert.Bool(rule.Apply(makeDestination("8.8.8.8"))).IsFalse()
}
