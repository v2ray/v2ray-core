package json

import (
	"testing"

	v2net "github.com/v2ray/v2ray-core/common/net"
	v2nettesting "github.com/v2ray/v2ray-core/common/net/testing"
	"github.com/v2ray/v2ray-core/testing/unit"
)

func TestDomainMatching(t *testing.T) {
	assert := unit.Assert(t)

	rule := &FieldRule{
		Domain: "v2ray.com",
	}
	dest := v2net.NewTCPDestination(v2net.DomainAddress("www.v2ray.com", 80))
	assert.Bool(rule.Apply(dest)).IsTrue()
}

func TestPortMatching(t *testing.T) {
	assert := unit.Assert(t)

	rule := &FieldRule{
		Port: &v2nettesting.PortRange{
			FromValue: 0,
			ToValue:   100,
		},
	}
	dest := v2net.NewTCPDestination(v2net.DomainAddress("www.v2ray.com", 80))
	assert.Bool(rule.Apply(dest)).IsTrue()
}
