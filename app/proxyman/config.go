package proxyman

import (
	"context"

	"v2ray.com/core/common/errors"
	"v2ray.com/core/proxy"
)

func (s *AllocationStrategy) GetConcurrencyValue() uint32 {
	if s == nil || s.Concurrency == nil {
		return 3
	}
	return s.Concurrency.Value
}

func (s *AllocationStrategy) GetRefreshValue() uint32 {
	if s == nil || s.Refresh == nil {
		return 5
	}
	return s.Refresh.Value
}

func (c *OutboundHandlerConfig) GetProxyHandler(ctx context.Context) (proxy.Outbound, error) {
	if c == nil {
		return nil, errors.New("OutboundHandlerConfig is nil").Path("App", "Proxyman", "Outbound", "OutboundHandlerConfig")
	}
	config, err := c.ProxySettings.GetInstance()
	if err != nil {
		return nil, err
	}
	return proxy.CreateOutboundHandler(ctx, config)
}
