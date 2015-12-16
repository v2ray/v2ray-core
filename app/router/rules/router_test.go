package rules_test

import (
	"testing"

	. "github.com/v2ray/v2ray-core/app/router/rules"
	testinconfig "github.com/v2ray/v2ray-core/app/router/rules/testing"
	v2net "github.com/v2ray/v2ray-core/common/net"
	v2testing "github.com/v2ray/v2ray-core/testing"
	"github.com/v2ray/v2ray-core/testing/assert"
)

func TestSimpleRouter(t *testing.T) {
	v2testing.Current(t)

	router := NewRouter().AddRule(
		&testinconfig.TestRule{
			TagValue: "test",
			Function: func(dest v2net.Destination) bool {
				return dest.IsTCP()
			},
		})

	tag, err := router.TakeDetour(v2net.TCPDestination(v2net.DomainAddress("v2ray.com"), 80))
	assert.Error(err).IsNil()
	assert.StringLiteral(tag).Equals("test")
}
