package inbound

//go:generate go run $GOPATH/src/v2ray.com/core/tools/generrorgen/main.go -pkg inbound -path App,Proxyman,Inbound

import (
	"context"

	"v2ray.com/core/app/proxyman"
	"v2ray.com/core/common"
)

type DefaultInboundHandlerManager struct {
	handlers       []proxyman.InboundHandler
	taggedHandlers map[string]proxyman.InboundHandler
}

func New(ctx context.Context, config *proxyman.InboundConfig) (*DefaultInboundHandlerManager, error) {
	return &DefaultInboundHandlerManager{
		taggedHandlers: make(map[string]proxyman.InboundHandler),
	}, nil
}

func (m *DefaultInboundHandlerManager) AddHandler(ctx context.Context, config *proxyman.InboundHandlerConfig) error {
	rawReceiverSettings, err := config.ReceiverSettings.GetInstance()
	if err != nil {
		return err
	}
	receiverSettings, ok := rawReceiverSettings.(*proxyman.ReceiverConfig)
	if !ok {
		return newError("not a ReceiverConfig")
	}
	proxySettings, err := config.ProxySettings.GetInstance()
	if err != nil {
		return err
	}
	var handler proxyman.InboundHandler
	tag := config.Tag
	allocStrategy := receiverSettings.AllocationStrategy
	if allocStrategy == nil || allocStrategy.Type == proxyman.AllocationStrategy_Always {
		h, err := NewAlwaysOnInboundHandler(ctx, tag, receiverSettings, proxySettings)
		if err != nil {
			return err
		}
		handler = h
	} else if allocStrategy.Type == proxyman.AllocationStrategy_Random {
		h, err := NewDynamicInboundHandler(ctx, tag, receiverSettings, proxySettings)
		if err != nil {
			return err
		}
		handler = h
	}

	if handler == nil {
		return newError("unknown allocation strategy: ", receiverSettings.AllocationStrategy.Type)
	}

	m.handlers = append(m.handlers, handler)
	if len(tag) > 0 {
		m.taggedHandlers[tag] = handler
	}
	return nil
}

func (m *DefaultInboundHandlerManager) GetHandler(ctx context.Context, tag string) (proxyman.InboundHandler, error) {
	handler, found := m.taggedHandlers[tag]
	if !found {
		return nil, newError("handler not found: ", tag)
	}
	return handler, nil
}

func (m *DefaultInboundHandlerManager) Start() error {
	for _, handler := range m.handlers {
		if err := handler.Start(); err != nil {
			return err
		}
	}
	return nil
}

func (m *DefaultInboundHandlerManager) Close() {
	for _, handler := range m.handlers {
		handler.Close()
	}
}

func (m *DefaultInboundHandlerManager) Interface() interface{} {
	return (*proxyman.InboundHandlerManager)(nil)
}

func init() {
	common.Must(common.RegisterConfig((*proxyman.InboundConfig)(nil), func(ctx context.Context, config interface{}) (interface{}, error) {
		return New(ctx, config.(*proxyman.InboundConfig))
	}))
}
