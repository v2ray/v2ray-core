// Package blackhole is an outbound handler that blocks all connections.
package blackhole

import (
	"context"
	"errors"
	"time"

	"v2ray.com/core/common"
	"v2ray.com/core/proxy"
	"v2ray.com/core/transport/ray"
)

var (
	errConnectionBlocked = errors.New("Blackhole: connection blocked.")
)

// Handler is an outbound connection that sliently swallow the entire payload.
type Handler struct {
	response ResponseConfig
}

// New creates a new blackhole handler.
func New(ctx context.Context, config *Config) (*Handler, error) {
	response, err := config.GetInternalResponse()
	if err != nil {
		return nil, err
	}
	return &Handler{
		response: response,
	}, nil
}

// Dispatch implements OutboundHandler.Dispatch().
func (v *Handler) Process(ctx context.Context, outboundRay ray.OutboundRay, dialer proxy.Dialer) error {
	v.response.WriteTo(outboundRay.OutboundOutput())
	// Sleep a little here to make sure the response is sent to client.
	time.Sleep(time.Second)
	return errConnectionBlocked
}

func init() {
	common.Must(common.RegisterConfig((*Config)(nil), func(ctx context.Context, config interface{}) (interface{}, error) {
		return New(ctx, config.(*Config))
	}))
}
