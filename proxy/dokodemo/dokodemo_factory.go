package dokodemo

import (
	"github.com/v2ray/v2ray-core/app"
	"github.com/v2ray/v2ray-core/proxy/common/connhandler"
)

type DokodemoDoorFactory struct {
}

func (this DokodemoDoorFactory) Create(space *app.Space, rawConfig interface{}) (connhandler.InboundConnectionHandler, error) {
	config := rawConfig.(Config)
	return NewDokodemoDoor(space, config), nil
}

func init() {
	connhandler.RegisterInboundConnectionHandlerFactory("dokodemo-door", DokodemoDoorFactory{})
}
