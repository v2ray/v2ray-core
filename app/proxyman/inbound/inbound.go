package inbound

//go:generate go run $GOPATH/src/v2ray.com/core/common/errors/errorgen/main.go -pkg inbound -path App,Proxyman,Inbound

import (
	"context"
	"sync"

	"v2ray.com/core"
	"v2ray.com/core/app/proxyman"
	"v2ray.com/core/common"
)

// Manager is to manage all inbound handlers.
type Manager struct {
	access          sync.RWMutex
	untaggedHandler []core.InboundHandler
	taggedHandlers  map[string]core.InboundHandler
	running         bool
}

// New returns a new Manager for inbound handlers.
func New(ctx context.Context, config *proxyman.InboundConfig) (*Manager, error) {
	m := &Manager{
		taggedHandlers: make(map[string]core.InboundHandler),
	}
	v := core.FromContext(ctx)
	if v == nil {
		return nil, newError("V is not in context")
	}
	if err := v.RegisterFeature((*core.InboundHandlerManager)(nil), m); err != nil {
		return nil, newError("unable to register InboundHandlerManager").Base(err)
	}
	return m, nil
}

// AddHandler implements core.InboundHandlerManager.
func (m *Manager) AddHandler(ctx context.Context, handler core.InboundHandler) error {
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

// GetHandler returns core.InboundHandlerManager.
func (m *Manager) GetHandler(ctx context.Context, tag string) (core.InboundHandler, error) {
	m.access.RLock()
	defer m.access.RUnlock()

	handler, found := m.taggedHandlers[tag]
	if !found {
		return nil, newError("handler not found: ", tag)
	}
	return handler, nil
}

func (m *Manager) RemoveHandler(ctx context.Context, tag string) error {
	if len(tag) == 0 {
		return core.ErrNoClue
	}

	m.access.Lock()
	defer m.access.Unlock()

	if handler, found := m.taggedHandlers[tag]; found {
		handler.Close()
		delete(m.taggedHandlers, tag)
		return nil
	}

	return core.ErrNoClue
}

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

func (m *Manager) Close() error {
	m.access.Lock()
	defer m.access.Unlock()

	m.running = false

	for _, handler := range m.taggedHandlers {
		handler.Close()
	}
	for _, handler := range m.untaggedHandler {
		handler.Close()
	}

	return nil
}

func NewHandler(ctx context.Context, config *core.InboundHandlerConfig) (core.InboundHandler, error) {
	rawReceiverSettings, err := config.ReceiverSettings.GetInstance()
	if err != nil {
		return nil, err
	}
	receiverSettings, ok := rawReceiverSettings.(*proxyman.ReceiverConfig)
	if !ok {
		return nil, newError("not a ReceiverConfig").AtError()
	}
	proxySettings, err := config.ProxySettings.GetInstance()
	if err != nil {
		return nil, err
	}
	tag := config.Tag
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
