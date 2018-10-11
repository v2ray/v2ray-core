package core

import (
	"context"
	"sync"

	"v2ray.com/core/common"
	"v2ray.com/core/common/net"
	"v2ray.com/core/features/outbound"
)

// InboundHandler is the interface for handlers that process inbound connections.
type InboundHandler interface {
	common.Runnable
	// The tag of this handler.
	Tag() string

	// Deprecated. Do not use in new code.
	GetRandomInboundProxy() (interface{}, net.Port, int)
}

// InboundHandlerManager is a feature that manages InboundHandlers.
type InboundHandlerManager interface {
	Feature
	// GetHandlers returns an InboundHandler for the given tag.
	GetHandler(ctx context.Context, tag string) (InboundHandler, error)
	// AddHandler adds the given handler into this InboundHandlerManager.
	AddHandler(ctx context.Context, handler InboundHandler) error

	// RemoveHandler removes a handler from InboundHandlerManager.
	RemoveHandler(ctx context.Context, tag string) error
}

type syncInboundHandlerManager struct {
	sync.RWMutex
	InboundHandlerManager
}

func (m *syncInboundHandlerManager) GetHandler(ctx context.Context, tag string) (InboundHandler, error) {
	m.RLock()
	defer m.RUnlock()

	if m.InboundHandlerManager == nil {
		return nil, newError("InboundHandlerManager not set.").AtError()
	}

	return m.InboundHandlerManager.GetHandler(ctx, tag)
}

func (m *syncInboundHandlerManager) AddHandler(ctx context.Context, handler InboundHandler) error {
	m.RLock()
	defer m.RUnlock()

	if m.InboundHandlerManager == nil {
		return newError("InboundHandlerManager not set.").AtError()
	}

	return m.InboundHandlerManager.AddHandler(ctx, handler)
}

func (m *syncInboundHandlerManager) Start() error {
	m.RLock()
	defer m.RUnlock()

	if m.InboundHandlerManager == nil {
		return newError("InboundHandlerManager not set.").AtError()
	}

	return m.InboundHandlerManager.Start()
}

func (m *syncInboundHandlerManager) Close() error {
	m.RLock()
	defer m.RUnlock()

	return common.Close(m.InboundHandlerManager)
}

func (m *syncInboundHandlerManager) Set(manager InboundHandlerManager) {
	if manager == nil {
		return
	}

	m.Lock()
	defer m.Unlock()

	common.Close(m.InboundHandlerManager) // nolint: errcheck
	m.InboundHandlerManager = manager
}

type syncOutboundHandlerManager struct {
	sync.RWMutex
	outbound.HandlerManager
}

func (m *syncOutboundHandlerManager) GetHandler(tag string) outbound.Handler {
	m.RLock()
	defer m.RUnlock()

	if m.HandlerManager == nil {
		return nil
	}

	return m.HandlerManager.GetHandler(tag)
}

func (m *syncOutboundHandlerManager) GetDefaultHandler() outbound.Handler {
	m.RLock()
	defer m.RUnlock()

	if m.HandlerManager == nil {
		return nil
	}

	return m.HandlerManager.GetDefaultHandler()
}

func (m *syncOutboundHandlerManager) AddHandler(ctx context.Context, handler outbound.Handler) error {
	m.RLock()
	defer m.RUnlock()

	if m.HandlerManager == nil {
		return newError("OutboundHandlerManager not set.").AtError()
	}

	return m.HandlerManager.AddHandler(ctx, handler)
}

func (m *syncOutboundHandlerManager) Start() error {
	m.RLock()
	defer m.RUnlock()

	if m.HandlerManager == nil {
		return newError("OutboundHandlerManager not set.").AtError()
	}

	return m.HandlerManager.Start()
}

func (m *syncOutboundHandlerManager) Close() error {
	m.RLock()
	defer m.RUnlock()

	return common.Close(m.HandlerManager)
}

func (m *syncOutboundHandlerManager) Set(manager outbound.HandlerManager) {
	if manager == nil {
		return
	}

	m.Lock()
	defer m.Unlock()

	common.Close(m.HandlerManager) // nolint: errcheck
	m.HandlerManager = manager
}
