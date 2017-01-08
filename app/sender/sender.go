package sender

import (
	"v2ray.com/core/app"
	"v2ray.com/core/common"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/serial"
	"v2ray.com/core/transport/internet"
)

type Sender interface {
	SendTo(net.Destination) (internet.Connection, error)
}

type SenderManager struct {
}

func New(space app.Space, config *Config) (*SenderManager, error) {
	return &SenderManager{}, nil
}

type SenderManagerFactory struct{}

func (SenderManagerFactory) Create(space app.Space, config interface{}) (app.Application, error) {
	return New(space, config.(*Config))
}

func FromSpace(space app.Space) *SenderManager {
	app := space.(app.AppGetter).GetApp(serial.GetMessageType((*Config)(nil)))
	if app == nil {
		return nil
	}
	return app.(*SenderManager)
}

func init() {
	common.Must(app.RegisterApplicationFactory((*Config)(nil), SenderManagerFactory{}))
}
