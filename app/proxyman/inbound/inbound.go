package inbound

//go:generate go run v2ray.com/core/common/errors/errorgen

import (
	"context"
	"sync"

	"v2ray.com/core"
	"v2ray.com/core/app/proxyman"
	"v2ray.com/core/common"
	"v2ray.com/core/common/serial"
	"v2ray.com/core/common/session"
	"v2ray.com/core/features/inbound"
)

// Manager is to manage all inbound handlers.
type Manager struct {
	access          sync.RWMutex
	untaggedHandler []inbound.Handler
	taggedHandlers  map[string]inbound.Handler
	running         bool
}

// New returns a new Manager for inbound handlers.
func New(ctx context.Context, config *proxyman.InboundConfig) (*Manager, error) {
	m := &Manager{
		taggedHandlers: make(map[string]inbound.Handler),
	}
	return m, nil
}

// Type implements common.HasType.
func (*Manager) Type() interface{} {
	return inbound.ManagerType()
}

// AddHandler implements inbound.Manager.
func (m *Manager) AddHandler(ctx context.Context, handler inbound.Handler) error {
	m.access.Lock()
	defer m.access.Unlock()

	tag := handler.Tag()
	if len(tag) > 0 {
		m.taggedHandlers[tag] = handler
	} else {
		m.untaggedHandler = append(m.untaggedHandler, handler)
	}

	if m.running {
		return handler.Start()
	}

	return nil
}

// GetHandler implements inbound.Manager.
func (m *Manager) GetHandler(ctx context.Context, tag string) (inbound.Handler, error) {
	m.access.RLock()
	defer m.access.RUnlock()

	handler, found := m.taggedHandlers[tag]
	if !found {
		return nil, newError("handler not found: ", tag)
	}
	return handler, nil
}

// RemoveHandler implements inbound.Manager.
func (m *Manager) RemoveHandler(ctx context.Context, tag string) error {
	if tag == "" {
		return common.ErrNoClue
	}

	m.access.Lock()
	defer m.access.Unlock()

	if handler, found := m.taggedHandlers[tag]; found {
		if err := handler.Close(); err != nil {
			newError("failed to close handler ", tag).Base(err).AtWarning().WriteToLog(session.ExportIDToError(ctx))
		}
		delete(m.taggedHandlers, tag)
		return nil
	}

	return common.ErrNoClue
}

// Start implements common.Runnable.
func (m *Manager) Start() error {
	m.access.Lock()
	defer m.access.Unlock()

	m.running = true

	for _, handler := range m.taggedHandlers {
		if err := handler.Start(); err != nil {
			return err
		}
	}

	for _, handler := range m.untaggedHandler {
		if err := handler.Start(); err != nil {
			return err
		}
	}
	return nil
}

// Close implements common.Closable.
func (m *Manager) Close() error {
	m.access.Lock()
	defer m.access.Unlock()

	m.running = false

	var errors []interface{}
	for _, handler := range m.taggedHandlers {
		if err := handler.Close(); err != nil {
			errors = append(errors, err)
		}
	}
	for _, handler := range m.untaggedHandler {
		if err := handler.Close(); err != nil {
			errors = append(errors, err)
		}
	}

	if len(errors) > 0 {
		return newError("failed to close all handlers").Base(newError(serial.Concat(errors...)))
	}

	return nil
}

// NewHandler creates a new inbound.Handler based on the given config.
func NewHandler(ctx context.Context, config *core.InboundHandlerConfig) (inbound.Handler, error) {
	rawReceiverSettings, err := config.ReceiverSettings.GetInstance()
	if err != nil {
		return nil, err
	}
	proxySettings, err := config.ProxySettings.GetInstance()
	if err != nil {
		return nil, err
	}
	tag := config.Tag

	receiverSettings, ok := rawReceiverSettings.(*proxyman.ReceiverConfig)
	if !ok {
		return nil, newError("not a ReceiverConfig").AtError()
	}

	streamSettings := receiverSettings.StreamSettings
	if streamSettings != nil && streamSettings.SocketSettings != nil {
		ctx = session.ContextWithSockopt(ctx, &session.Sockopt{
			Mark: streamSettings.SocketSettings.Mark,
		})
	}

	allocStrategy := receiverSettings.AllocationStrategy
	if allocStrategy == nil || allocStrategy.Type == proxyman.AllocationStrategy_Always {
		return NewAlwaysOnInboundHandler(ctx, tag, receiverSettings, proxySettings)
	}

	if allocStrategy.Type == proxyman.AllocationStrategy_Random {
		return NewDynamicInboundHandler(ctx, tag, receiverSettings, proxySettings)
	}
	return nil, newError("unknown allocation strategy: ", receiverSettings.AllocationStrategy.Type).AtError()
}

func init() {
	common.Must(common.RegisterConfig((*proxyman.InboundConfig)(nil), func(ctx context.Context, config interface{}) (interface{}, error) {
		return New(ctx, config.(*proxyman.InboundConfig))
	}))
	common.Must(common.RegisterConfig((*core.InboundHandlerConfig)(nil), func(ctx context.Context, config interface{}) (interface{}, error) {
		return NewHandler(ctx, config.(*core.InboundHandlerConfig))
	}))
}
