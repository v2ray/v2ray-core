package rules

import (
	"testing"

	"github.com/v2ray/v2ray-core/app/router/rules/config"
	testinconfig "github.com/v2ray/v2ray-core/app/router/rules/config/testing"
	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/testing/unit"
)

func TestSimpleRouter(t *testing.T) {
	assert := unit.Assert(t)

	router := &Router{
		rules: []config.Rule{
			&testinconfig.TestRule{
				TagValue: "test",
				Function: func(dest v2net.Destination) bool {
					return dest.IsTCP()
				},
			},
		},
	}

	tag, err := router.TakeDetour(v2net.NewTCPDestination(v2net.DomainAddress("v2ray.com", 80)))
	assert.Error(err).IsNil()
	assert.String(tag).Equals("test")
}
