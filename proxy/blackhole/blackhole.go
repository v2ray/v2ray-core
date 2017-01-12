// Package blackhole is an outbound handler that blocks all connections.
package blackhole

import (
	"time"

	"v2ray.com/core/app"
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/proxy"
	"v2ray.com/core/transport/ray"
)

// Handler is an outbound connection that sliently swallow the entire payload.
type Handler struct {
	meta     *proxy.OutboundHandlerMeta
	response ResponseConfig
}

// New creates a new blackhole handler.
func New(space app.Space, config *Config, meta *proxy.OutboundHandlerMeta) (proxy.OutboundHandler, error) {
	response, err := config.GetInternalResponse()
	if err != nil {
		return nil, err
	}
	return &Handler{
		meta:     meta,
		response: response,
	}, nil
}

// Dispatch implements OutboundHandler.Dispatch().
func (v *Handler) Dispatch(destination v2net.Destination, ray ray.OutboundRay) {
	v.response.WriteTo(ray.OutboundOutput())
	// CloseError() will immediately close the connection.
	// Sleep a little here to make sure the response is sent to client.
	time.Sleep(time.Millisecond * 500)
	ray.OutboundInput().CloseError()
	ray.OutboundOutput().CloseError()
}

// Factory is an utility for creating blackhole handlers.
type Factory struct{}

// Create implements OutboundHandlerFactory.Create().
func (v *Factory) Create(space app.Space, config interface{}, meta *proxy.OutboundHandlerMeta) (proxy.OutboundHandler, error) {
	return New(space, config.(*Config), meta)
}
