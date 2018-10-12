package outbound

import (
	"context"

	"v2ray.com/core"
	"v2ray.com/core/app/proxyman"
	"v2ray.com/core/app/proxyman/mux"
	"v2ray.com/core/common"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/session"
	"v2ray.com/core/common/vio"
	"v2ray.com/core/features/outbound"
	"v2ray.com/core/proxy"
	"v2ray.com/core/transport/internet"
	"v2ray.com/core/transport/pipe"
)

type Handler struct {
	config          *core.OutboundHandlerConfig
	senderSettings  *proxyman.SenderConfig
	streamSettings  *internet.MemoryStreamConfig
	proxy           proxy.Outbound
	outboundManager outbound.Manager
	mux             *mux.ClientManager
}

func NewHandler(ctx context.Context, config *core.OutboundHandlerConfig) (outbound.Handler, error) {
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
			mss, err := internet.ToMemoryStreamConfig(s.StreamSettings)
			if err != nil {
				return nil, newError("failed to parse stream settings").Base(err).AtWarning()
			}
			h.streamSettings = mss
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

// Tag implements outbound.Handler.
func (h *Handler) Tag() string {
	return h.config.Tag
}

// Dispatch implements proxy.Outbound.Dispatch.
func (h *Handler) Dispatch(ctx context.Context, link *vio.Link) {
	if h.mux != nil {
		if err := h.mux.Dispatch(ctx, link); err != nil {
			newError("failed to process mux outbound traffic").Base(err).WriteToLog(session.ExportIDToError(ctx))
			pipe.CloseError(link.Writer)
		}
	} else {
		if err := h.proxy.Process(ctx, link, h); err != nil {
			// Ensure outbound ray is properly closed.
			newError("failed to process outbound traffic").Base(err).WriteToLog(session.ExportIDToError(ctx))
			pipe.CloseError(link.Writer)
		} else {
			common.Must(common.Close(link.Writer))
		}
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
				newError("proxying to ", tag, " for dest ", dest).AtDebug().WriteToLog(session.ExportIDToError(ctx))
				ctx = session.ContextWithOutbound(ctx, &session.Outbound{
					Target: dest,
				})

				opts := pipe.OptionsFromContext(ctx)
				uplinkReader, uplinkWriter := pipe.New(opts...)
				downlinkReader, downlinkWriter := pipe.New(opts...)

				go handler.Dispatch(ctx, &vio.Link{Reader: uplinkReader, Writer: downlinkWriter})
				return net.NewConnection(net.ConnectionInputMulti(uplinkWriter), net.ConnectionOutputMulti(downlinkReader)), nil
			}

			newError("failed to get outbound handler with tag: ", tag).AtWarning().WriteToLog(session.ExportIDToError(ctx))
		}

		if h.senderSettings.Via != nil {
			outbound := session.OutboundFromContext(ctx)
			if outbound == nil {
				outbound = new(session.Outbound)
				ctx = session.ContextWithOutbound(ctx, outbound)
			}
			outbound.Gateway = h.senderSettings.Via.AsAddress()
		}

		ctx = internet.ContextWithStreamSettings(ctx, h.streamSettings)
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
