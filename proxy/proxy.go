// Package proxy contains all proxies used by V2Ray.

package proxy

import (
	"github.com/v2ray/v2ray-core/app"
	"github.com/v2ray/v2ray-core/proxy/common/connhandler"
	"github.com/v2ray/v2ray-core/proxy/internal"
)

func CreateInboundConnectionHandler(name string, space app.Space, rawConfig []byte) (connhandler.InboundConnectionHandler, error) {
	return internal.CreateInboundConnectionHandler(name, space, rawConfig)
}

func CreateOutboundConnectionHandler(name string, space app.Space, rawConfig []byte) (connhandler.OutboundConnectionHandler, error) {
	return internal.CreateOutboundConnectionHandler(name, space, rawConfig)
}
