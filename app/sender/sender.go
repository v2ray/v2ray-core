package sender

import (
	"context"

	"v2ray.com/core/app"
	"v2ray.com/core/common"
	"v2ray.com/core/common/net"
	"v2ray.com/core/transport/internet"
)

type Sender interface {
	SendTo(net.Destination) (internet.Connection, error)
}

type SenderManager struct {
}

func New(ctx context.Context, config *Config) (*SenderManager, error) {
	return &SenderManager{}, nil
}

func (SenderManager) Interface() interface{} {
	return (*SenderManager)(nil)
}

func FromSpace(space app.Space) *SenderManager {
	app := space.GetApplication((*SenderManager)(nil))
	if app == nil {
		return nil
	}
	return app.(*SenderManager)
}

func init() {
	common.Must(common.RegisterConfig((*Config)(nil), func(ctx context.Context, config interface{}) (interface{}, error) {
		return New(ctx, config.(*Config))
	}))
}
