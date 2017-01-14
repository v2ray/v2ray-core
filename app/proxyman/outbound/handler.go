package outbound

import (
	"context"

	"v2ray.com/core/app/proxyman"
	"v2ray.com/core/common/errors"
	"v2ray.com/core/common/net"
	"v2ray.com/core/proxy"
	"v2ray.com/core/transport/internet"
)

type Handler struct {
	config           *proxyman.OutboundHandlerConfig
	streamSettings   *proxyman.StreamSenderConfig
	datagramSettings *proxyman.DatagramSenderConfig
	proxy            proxy.OutboundHandler
}

func NewHandler(ctx context.Context, config *proxyman.OutboundHandlerConfig) (*Handler, error) {
	h := &Handler{
		config: config,
	}
	for _, rawSettings := range config.SenderSettings {
		settings, err := rawSettings.GetInstance()
		if err != nil {
			return nil, err
		}
		switch ts := settings.(type) {
		case *proxyman.StreamSenderConfig:
			h.streamSettings = ts
		case *proxyman.DatagramSenderConfig:
			h.datagramSettings = ts
		default:
			return nil, errors.New("Proxyman|DefaultOutboundHandler: Unknown sender settings: ", rawSettings.Type)
		}
	}
	proxyHandler, err := config.GetProxyHandler(proxy.ContextWithDialer(ctx, h))
	if err != nil {
		return nil, err
	}

	h.proxy = proxyHandler
	return h, nil
}

func (h *Handler) Dial(ctx context.Context, destination net.Destination) (internet.Connection, error) {
	switch destination.Network {
	case net.Network_TCP:
		return h.dialStream(ctx, destination)
	case net.Network_UDP:
		return h.dialDatagram(ctx, destination)
	default:
		panic("Proxyman|DefaultOutboundHandler: unexpected network.")
	}
}

func (h *Handler) dialStream(ctx context.Context, destination net.Destination) (internet.Connection, error) {
	var src net.Address
	if h.streamSettings != nil {
		src = h.streamSettings.Via.AsAddress()
	}
	var options internet.DialerOptions
	if h.streamSettings != nil {
		options.Proxy = h.streamSettings.ProxySettings
		options.Stream = h.streamSettings.StreamSettings
	}
	return internet.Dial(src, destination, options)
}

func (h *Handler) dialDatagram(ctx context.Context, destination net.Destination) (internet.Connection, error) {
	var src net.Address
	if h.datagramSettings != nil {
		src = h.datagramSettings.Via.AsAddress()
	}
	var options internet.DialerOptions
	if h.datagramSettings != nil {
		options.Proxy = h.datagramSettings.ProxySettings
	}
	return internet.Dial(src, destination, options)
}
