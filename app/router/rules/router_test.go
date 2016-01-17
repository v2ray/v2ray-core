package rules_test

import (
	"testing"

	. "github.com/v2ray/v2ray-core/app/router/rules"
	v2net "github.com/v2ray/v2ray-core/common/net"
	v2testing "github.com/v2ray/v2ray-core/testing"
	"github.com/v2ray/v2ray-core/testing/assert"
)

func TestSimpleRouter(t *testing.T) {
	v2testing.Current(t)

	router := NewRouter().AddRule(
		&Rule{
			Tag:       "test",
			Condition: NewNetworkMatcher(v2net.Network("tcp").AsList()),
		})

	tag, err := router.TakeDetour(v2net.TCPDestination(v2net.DomainAddress("v2ray.com"), 80))
	assert.Error(err).IsNil()
	assert.StringLiteral(tag).Equals("test")
}
