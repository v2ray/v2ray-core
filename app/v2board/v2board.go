package v2board

//go:generate errorgen

import (
	"context"

	"v2ray.com/core"
	"v2ray.com/core/common"
	"v2ray.com/core/features/inbound"
	"v2ray.com/core/features/stats"
)

type V2Board struct {
	config       V2BoardConfig
	userset      *UserSet
	instance     *core.Instance
	serverConfig *Config
	im           inbound.Manager
	stats        stats.Manager
}

func (v *V2Board) Start() error {
	go v.Loop()

	return nil
}

func (v *V2Board) Close() error {
	return nil
}
func (v *V2Board) Type() interface{} {
	return (*V2Board)(nil)
}

func init() {
	v2board := &V2Board{
		userset: NewUserSet(),
	}
	common.Must(core.RegisterConfigLoader(&core.ConfigFormat{
		Name:      "V2BOARD",
		Extension: []string{"v2board"},
		Loader:    v2board.ConfigLoader,
	}))
	common.Must(common.RegisterConfig((*Config)(nil), func(ctx context.Context, cfg interface{}) (interface{}, error) {
		s := core.MustFromContext(ctx)
		core.RequireFeatures(ctx, func(im inbound.Manager, stat stats.Manager) error {
			v2board.im = im
			v2board.stats = stat
			return nil
		})

		v2board.instance = s
		v2board.serverConfig = cfg.(*Config)
		return v2board, nil
	}))

}
