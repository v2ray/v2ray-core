package internal

import (
	"github.com/v2ray/v2ray-core/app"
	"github.com/v2ray/v2ray-core/proxy/common/connhandler"
)

type InboundConnectionHandlerCreator func(space app.Space, config interface{}) (connhandler.InboundConnectionHandler, error)
type OutboundConnectionHandlerCreator func(space app.Space, config interface{}) (connhandler.OutboundConnectionHandler, error)
