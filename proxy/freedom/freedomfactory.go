package freedom

import (
	"github.com/v2ray/v2ray-core/proxy/common/connhandler"
)

type FreedomFactory struct {
}

func (this FreedomFactory) Create(config interface{}) (connhandler.OutboundConnectionHandler, error) {
	return NewFreedomConnection(), nil
}

func init() {
	connhandler.RegisterOutboundConnectionHandlerFactory("freedom", FreedomFactory{})
}
