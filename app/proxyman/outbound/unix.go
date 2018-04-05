package outbound

import (
	"context"
	"io"

	"v2ray.com/core"
	"v2ray.com/core/app/proxyman"
	"v2ray.com/core/app/proxyman/mux"
	"v2ray.com/core/common"
	"v2ray.com/core/common/errors"
	"v2ray.com/core/common/net"
	"v2ray.com/core/proxy"
	"v2ray.com/core/transport/internet"
	"v2ray.com/core/transport/internet/domainsocket"
	"v2ray.com/core/transport/ray"
)

type UnixHandler struct {
	config          *core.OutboundHandlerConfig
	senderSettings  *proxyman.UnixSenderConfig
	proxy           proxy.Outbound
	outboundManager core.OutboundHandlerManager
	mux             *mux.ClientManager
}

func NewUnixHandler(ctx context.Context, config *core.OutboundHandlerConfig) (core.OutboundHandler, error) {
	v := core.FromContext(ctx)
	if v == nil {
		return nil, newError("V is not in context")
	}
	h := &UnixHandler{
		config:          config,
		outboundManager: v.OutboundHandlerManager(),
	}

	if config.SenderSettings != nil {
		senderSettings, err := config.SenderSettings.GetInstance()
		if err != nil {
			return nil, err
		}
		switch s := senderSettings.(type) {
		case *proxyman.UnixSenderConfig:
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
func (h *UnixHandler) Tag() string {
	return h.config.Tag
}

// Dispatch implements proxy.Outbound.Dispatch.
func (h *UnixHandler) Dispatch(ctx context.Context, outboundRay ray.OutboundRay) {
	if h.mux != nil {
		err := h.mux.Dispatch(ctx, outboundRay)
		if err != nil {
			newError("failed to process outbound traffic").Base(err).WriteToLog()
			outboundRay.OutboundOutput().CloseError()
		}
	} else {
		err := h.proxy.Process(ctx, outboundRay, h)
		// Ensure outbound ray is properly closed.
		if err != nil && errors.Cause(err) != io.EOF {
			newError("failed to process outbound traffic").Base(err).WriteToLog()
			outboundRay.OutboundOutput().CloseError()
		} else {
			outboundRay.OutboundOutput().Close()
		}
		outboundRay.OutboundInput().CloseError()
	}
}

// Dial implements proxy.Dialer.Dial().
func (h *UnixHandler) Dial(ctx context.Context, dest net.Destination) (internet.Connection, error) {
	if h.senderSettings != nil {
		if h.senderSettings.ProxySettings.HasTag() {
			newError("Unix domain socket does not support redirect").AtWarning().WriteToLog()
		}

		if h.senderSettings.StreamSettings != nil {
			newError("Unix domain socket does not support stream setting").AtWarning().WriteToLog()
		}
	}

	return domainsocket.DialDS(ctx, h.senderSettings.GetDomainSockSettings().GetPath())
}

// GetOutbound implements proxy.GetOutbound.
func (h *UnixHandler) GetOutbound() proxy.Outbound {
	return h.proxy
}

// Start implements common.Runnable.
func (h *UnixHandler) Start() error {
	return nil
}

// Close implements common.Runnable.
func (h *UnixHandler) Close() error {
	common.Close(h.mux)
	return nil
}
