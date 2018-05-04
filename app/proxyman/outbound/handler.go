package outbound

import (
	"context"

	"v2ray.com/core"
	"v2ray.com/core/app/proxyman"
	"v2ray.com/core/app/proxyman/mux"
	"v2ray.com/core/common"
	"v2ray.com/core/common/net"
	"v2ray.com/core/proxy"
	"v2ray.com/core/transport/internet"
	"v2ray.com/core/transport/pipe"
)

type Handler struct {
	config          *core.OutboundHandlerConfig
	senderSettings  *proxyman.SenderConfig
	proxy           proxy.Outbound
	outboundManager core.OutboundHandlerManager
	mux             *mux.ClientManager
}

func NewHandler(ctx context.Context, config *core.OutboundHandlerConfig) (core.OutboundHandler, error) {
	v := core.MustFromContext(ctx)
	h := &Handler{
		config:          config,
		outboundManager: v.OutboundHandlerManager(),
	}

	if config.SenderSettings != nil {
		senderSettings, err := config.SenderSettings.GetInstance()
		if err != nil {
			return nil, err
		}
		switch s := senderSettings.(type) {
		case *proxyman.SenderConfig:
			h.senderSettings = s
		default:
			return nil, newError("settings is not SenderConfig")
		}
	}

	proxyConfig, err := config.ProxySettings.GetInstance()
	if err != nil {
		return nil, err
	}

	rawProxyHandler, err := common.CreateObject(ctx, proxyConfig)
	if err != nil {
		return nil, err
	}

	proxyHandler, ok := rawProxyHandler.(proxy.Outbound)
	if !ok {
		return nil, newError("not an outbound handler")
	}

	if h.senderSettings != nil && h.senderSettings.MultiplexSettings != nil && h.senderSettings.MultiplexSettings.Enabled {
		config := h.senderSettings.MultiplexSettings
		if config.Concurrency < 1 || config.Concurrency > 1024 {
			return nil, newError("invalid mux concurrency: ", config.Concurrency).AtWarning()
		}
		h.mux = mux.NewClientManager(proxyHandler, h, config)
	}

	h.proxy = proxyHandler
	return h, nil
}

// Tag implements core.OutboundHandler.
func (h *Handler) Tag() string {
	return h.config.Tag
}

// Dispatch implements proxy.Outbound.Dispatch.
func (h *Handler) Dispatch(ctx context.Context, link *core.Link) {
	if h.mux != nil {
		if err := h.mux.Dispatch(ctx, link); err != nil {
			newError("failed to process mux outbound traffic").Base(err).WithContext(ctx).WriteToLog()
			pipe.CloseError(link.Writer)
		}
	} else {
		if err := h.proxy.Process(ctx, link, h); err != nil {
			newError("failed to process outbound traffic").Base(err).WithContext(ctx).WriteToLog()
		}
		// read the rest of buffer when any errors
		common.Must(common.Close(link.Writer))
		pipe.CloseError(link.Reader)
	}
}

// Dial implements proxy.Dialer.Dial().
func (h *Handler) Dial(ctx context.Context, dest net.Destination) (internet.Connection, error) {
	if h.senderSettings != nil {
		if h.senderSettings.ProxySettings.HasTag() {
			tag := h.senderSettings.ProxySettings.Tag
			handler := h.outboundManager.GetHandler(tag)
			if handler != nil {
				newError("proxying to ", tag, " for dest ", dest).AtDebug().WithContext(ctx).WriteToLog()
				ctx = proxy.ContextWithTarget(ctx, dest)

				uplinkReader, uplinkWriter := pipe.New()
				downlinkReader, downlinkWriter := pipe.New()

				go handler.Dispatch(ctx, &core.Link{Reader: uplinkReader, Writer: downlinkWriter})
				return net.NewConnection(net.ConnectionInputMulti(uplinkWriter), net.ConnectionOutputMulti(downlinkReader)), nil
			}

			newError("failed to get outbound handler with tag: ", tag).AtWarning().WithContext(ctx).WriteToLog()
		}

		if h.senderSettings.Via != nil {
			ctx = internet.ContextWithDialerSource(ctx, h.senderSettings.Via.AsAddress())
		}

		if h.senderSettings.StreamSettings != nil {
			ctx = internet.ContextWithStreamSettings(ctx, h.senderSettings.StreamSettings)
		}
	}

	return internet.Dial(ctx, dest)
}

// GetOutbound implements proxy.GetOutbound.
func (h *Handler) GetOutbound() proxy.Outbound {
	return h.proxy
}

// Start implements common.Runnable.
func (h *Handler) Start() error {
	return nil
}

// Close implements common.Closable.
func (h *Handler) Close() error {
	common.Close(h.mux)
	return nil
}
