package outbound

//go:generate go run $GOPATH/src/v2ray.com/core/tools/generrorgen/main.go -pkg outbound -path App,Proxyman,Outbound

import (
	"context"
	"sync"

	"v2ray.com/core/app/proxyman"
	"v2ray.com/core/common"
)

// Manager is to manage all outbound handlers.
type Manager struct {
	sync.RWMutex
	defaultHandler *Handler
	taggedHandler  map[string]*Handler
}

// New creates a new Manager.
func New(ctx context.Context, config *proxyman.OutboundConfig) (*Manager, error) {
	return &Manager{
		taggedHandler: make(map[string]*Handler),
	}, nil
}

// Interface implements Application.Interface.
func (*Manager) Interface() interface{} {
	return (*proxyman.OutboundHandlerManager)(nil)
}

// Start implements Application.Start
func (*Manager) Start() error { return nil }

// Close implements Application.Close
func (*Manager) Close() {}

func (m *Manager) GetDefaultHandler() proxyman.OutboundHandler {
	m.RLock()
	defer m.RUnlock()
	if m.defaultHandler == nil {
		return nil
	}
	return m.defaultHandler
}

func (m *Manager) GetHandler(tag string) proxyman.OutboundHandler {
	m.RLock()
	defer m.RUnlock()
	if handler, found := m.taggedHandler[tag]; found {
		return handler
	}
	return nil
}

func (m *Manager) AddHandler(ctx context.Context, config *proxyman.OutboundHandlerConfig) error {
	m.Lock()
	defer m.Unlock()

	handler, err := NewHandler(ctx, config)
	if err != nil {
		return err
	}
	if m.defaultHandler == nil {
		m.defaultHandler = handler
	}

	if len(config.Tag) > 0 {
		m.taggedHandler[config.Tag] = handler
	}

	return nil
}

func init() {
	common.Must(common.RegisterConfig((*proxyman.OutboundConfig)(nil), func(ctx context.Context, config interface{}) (interface{}, error) {
		return New(ctx, config.(*proxyman.OutboundConfig))
	}))
}
