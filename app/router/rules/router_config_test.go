// +build json

package rules_test

import (
	"testing"

	. "github.com/v2ray/v2ray-core/app/router/rules"
	v2net "github.com/v2ray/v2ray-core/common/net"
	v2testing "github.com/v2ray/v2ray-core/testing"
	"github.com/v2ray/v2ray-core/testing/assert"
)

func TestFieldRule(t *testing.T) {
	v2testing.Current(t)

	rule := ParseRule([]byte(`{
    "type": "field",
    "domain": [
      "ooxx.com",
      "oxox.com"
    ],
    "outboundTag": "direct"
  }`))
	assert.Pointer(rule).IsNotNil()
	assert.Bool(rule.Apply(v2net.TCPDestination(v2net.DomainAddress("www.ooxx.com"), 80))).IsTrue()
}
