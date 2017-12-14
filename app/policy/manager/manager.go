package manager

import (
	"context"

	"v2ray.com/core/app/policy"
	"v2ray.com/core/common"
)

// Instance is an instance of Policy manager.
type Instance struct {
	levels map[uint32]*policy.Policy
}

// New creates new Policy manager instance.
func New(ctx context.Context, config *policy.Config) (*Instance, error) {
	levels := config.Level
	if levels == nil {
		levels = make(map[uint32]*policy.Policy)
	}
	for _, p := range levels {
		g := global()
		g.OverrideWith(p)
		*p = g
	}
	return &Instance{
		levels: levels,
	}, nil
}

func global() policy.Policy {
	return policy.Policy{
		Timeout: &policy.Policy_Timeout{
			Handshake:      &policy.Second{Value: 4},
			ConnectionIdle: &policy.Second{Value: 300},
			UplinkOnly:     &policy.Second{Value: 5},
			DownlinkOnly:   &policy.Second{Value: 30},
		},
	}
}

// GetPolicy implements policy.Manager.
func (m *Instance) GetPolicy(level uint32) policy.Policy {
	if p, ok := m.levels[level]; ok {
		return *p
	}
	return global()
}

// Start implements app.Application.Start().
func (m *Instance) Start() error {
	return nil
}

// Close implements app.Application.Close().
func (m *Instance) Close() {
}

// Interface implement app.Application.Interface().
func (m *Instance) Interface() interface{} {
	return (*policy.Manager)(nil)
}

func init() {
	common.Must(common.RegisterConfig((*policy.Config)(nil), func(ctx context.Context, config interface{}) (interface{}, error) {
		return New(ctx, config.(*policy.Config))
	}))
}
