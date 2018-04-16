// Package blackhole is an outbound handler that blocks all connections.
package blackhole

//go:generate go run $GOPATH/src/v2ray.com/core/common/errors/errorgen/main.go -pkg blackhole -path Proxy,Blackhole

import (
	"context"
	"time"

	"v2ray.com/core"
	"v2ray.com/core/common"
	"v2ray.com/core/proxy"
	"v2ray.com/core/transport/pipe"
)

// Handler is an outbound connection that silently swallow the entire payload.
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

// Process implements OutboundHandler.Dispatch().
func (h *Handler) Process(ctx context.Context, link *core.Link, dialer proxy.Dialer) error {
	h.response.WriteTo(link.Writer)
	// Sleep a little here to make sure the response is sent to client.
	time.Sleep(time.Second)
	pipe.CloseError(link.Writer)
	return nil
}

func init() {
	common.Must(common.RegisterConfig((*Config)(nil), func(ctx context.Context, config interface{}) (interface{}, error) {
		return New(ctx, config.(*Config))
	}))
}
