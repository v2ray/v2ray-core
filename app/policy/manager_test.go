package policy_test

import (
	"context"
	"testing"
	"time"

	"v2ray.com/core"
	. "v2ray.com/core/app/policy"
	. "v2ray.com/ext/assert"
)

func TestPolicy(t *testing.T) {
	assert := With(t)

	manager, err := New(context.Background(), &Config{
		Level: map[uint32]*Policy{
			0: {
				Timeout: &Policy_Timeout{
					Handshake: &Second{
						Value: 2,
					},
				},
			},
		},
	})
	assert(err, IsNil)

	pDefault := core.DefaultPolicy()

	p0 := manager.ForLevel(0)
	assert(p0.Timeouts.Handshake, Equals, 2*time.Second)
	assert(p0.Timeouts.ConnectionIdle, Equals, pDefault.Timeouts.ConnectionIdle)

	p1 := manager.ForLevel(1)
	assert(p1.Timeouts.Handshake, Equals, pDefault.Timeouts.Handshake)
}
