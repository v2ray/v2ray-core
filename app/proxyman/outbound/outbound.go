package outbound

//go:generate go run $GOPATH/src/v2ray.com/core/common/errors/errorgen/main.go -pkg outbound -path App,Proxyman,Outbound

import (
	"context"
	"sync"

	"v2ray.com/core"
	"v2ray.com/core/app/proxyman"
	"v2ray.com/core/common"
)

// Manager is to manage all outbound handlers.
type Manager struct {
	access           sync.RWMutex
	defaultHandler   core.OutboundHandler
	taggedHandler    map[string]core.OutboundHandler
	untaggedHandlers []core.OutboundHandler
	running          bool
}

// New creates a new Manager.
func New(ctx context.Context, config *proxyman.OutboundConfig) (*Manager, error) {
	m := &Manager{
		taggedHandler: make(map[string]core.OutboundHandler),
	}
	v := core.FromContext(ctx)
	if v == nil {
		return nil, newError("V is not in context")
	}
	if err := v.RegisterFeature((*core.OutboundHandlerManager)(nil), m); err != nil {
		return nil, newError("unable to register OutboundHandlerManager").Base(err)
	}
	return m, nil
}

// Start implements core.Feature
func (m *Manager) Start() error {
	m.access.Lock()
	defer m.access.Unlock()

	m.running = true

	for _, h := range m.taggedHandler {
		if err := h.Start(); err != nil {
			return err
		}
	}

	for _, h := range m.untaggedHandlers {
		if err := h.Start(); err != nil {
			return err
		}
	}

	return nil
}

// Close implements core.Feature
func (m *Manager) Close() error {
	m.access.Lock()
	defer m.access.Unlock()

	m.running = false

	for _, h := range m.taggedHandler {
		h.Close()
	}

	for _, h := range m.untaggedHandlers {
		h.Close()
	}

	return nil
}

// GetDefaultHandler implements core.OutboundHandlerManager.
func (m *Manager) GetDefaultHandler() core.OutboundHandler {
	m.access.RLock()
	defer m.access.RUnlock()

	if m.defaultHandler == nil {
		return nil
	}
	return m.defaultHandler
}

// GetHandler implements core.OutboundHandlerManager.
func (m *Manager) GetHandler(tag string) core.OutboundHandler {
	m.access.RLock()
	defer m.access.RUnlock()
	if handler, found := m.taggedHandler[tag]; found {
		return handler
	}
	return nil
}

// AddHandler implements core.OutboundHandlerManager.
func (m *Manager) AddHandler(ctx context.Context, handler core.OutboundHandler) error {
	m.access.Lock()
	defer m.access.Unlock()

	if m.defaultHandler == nil {
		m.defaultHandler = handler
	}

	tag := handler.Tag()
	if len(tag) > 0 {
		m.taggedHandler[tag] = handler
	} else {
		m.untaggedHandlers = append(m.untaggedHandlers, handler)
	}

	if m.running {
		return handler.Start()
	}

	return nil
}

// RemoveHandler implements core.OutboundHandlerManager.
func (m *Manager) RemoveHandler(ctx context.Context, tag string) error {
	if len(tag) == 0 {
		return core.ErrNoClue
	}
	m.access.Lock()
	defer m.access.Unlock()

	delete(m.taggedHandler, tag)
	if m.defaultHandler.Tag() == tag {
		m.defaultHandler = nil
	}

	return nil
}

func init() {
	common.Must(common.RegisterConfig((*proxyman.OutboundConfig)(nil), func(ctx context.Context, config interface{}) (interface{}, error) {
		return New(ctx, config.(*proxyman.OutboundConfig))
	}))
	common.Must(common.RegisterConfig((*core.OutboundHandlerConfig)(nil), func(ctx context.Context, config interface{}) (interface{}, error) {
		return NewHandler(ctx, config.(*core.OutboundHandlerConfig))
	}))
}
