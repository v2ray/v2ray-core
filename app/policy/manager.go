package policy

import (
	"context"

	"v2ray.com/core"
	"v2ray.com/core/common"
)

// Instance is an instance of Policy manager.
type Instance struct {
	levels map[uint32]*Policy
}

// New creates new Policy manager instance.
func New(ctx context.Context, config *Config) (*Instance, error) {
	m := &Instance{
		levels: make(map[uint32]*Policy),
	}
	if len(config.Level) > 0 {
		for lv, p := range config.Level {
			pp := defaultPolicy()
			pp.overrideWith(p)
			m.levels[lv] = pp
		}
	}

	v := core.FromContext(ctx)
	if v != nil {
		if err := v.RegisterFeature((*core.PolicyManager)(nil), m); err != nil {
			return nil, newError("unable to register PolicyManager in core").Base(err).AtError()
		}
	}

	return m, nil
}

// ForLevel implements core.PolicyManager.
func (m *Instance) ForLevel(level uint32) core.Policy {
	if p, ok := m.levels[level]; ok {
		return p.ToCorePolicy()
	}
	return core.DefaultPolicy()
}

// Start implements app.Application.Start().
func (m *Instance) Start() error {
	return nil
}

// Close implements app.Application.Close().
func (m *Instance) Close() error {
	return nil
}

func init() {
	common.Must(common.RegisterConfig((*Config)(nil), func(ctx context.Context, config interface{}) (interface{}, error) {
		return New(ctx, config.(*Config))
	}))
}
