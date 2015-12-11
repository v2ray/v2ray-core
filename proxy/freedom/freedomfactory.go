package freedom

import (
	"github.com/v2ray/v2ray-core/app"
	"github.com/v2ray/v2ray-core/proxy/common/connhandler"
)

type FreedomFactory struct {
}

func (this FreedomFactory) Create(space app.Space, config interface{}) (connhandler.OutboundConnectionHandler, error) {
	return &FreedomConnection{space: space}, nil
}

func init() {
	connhandler.RegisterOutboundConnectionHandlerFactory("freedom", FreedomFactory{})
}
