package core

import (
	"context"
	"sync"

	"v2ray.com/core/common"
	"v2ray.com/core/features/inbound"
	"v2ray.com/core/features/outbound"
)

type syncInboundHandlerManager struct {
	sync.RWMutex
	inbound.Manager
}

func (m *syncInboundHandlerManager) GetHandler(ctx context.Context, tag string) (inbound.Handler, error) {
	m.RLock()
	defer m.RUnlock()

	if m.Manager == nil {
		return nil, newError("inbound.Manager not set.").AtError()
	}

	return m.Manager.GetHandler(ctx, tag)
}

func (m *syncInboundHandlerManager) AddHandler(ctx context.Context, handler inbound.Handler) error {
	m.RLock()
	defer m.RUnlock()

	if m.Manager == nil {
		return newError("inbound.Manager not set.").AtError()
	}

	return m.Manager.AddHandler(ctx, handler)
}

func (m *syncInboundHandlerManager) Start() error {
	m.RLock()
	defer m.RUnlock()

	if m.Manager == nil {
		return newError("inbound.Manager not set.").AtError()
	}

	return m.Manager.Start()
}

func (m *syncInboundHandlerManager) Close() error {
	m.RLock()
	defer m.RUnlock()

	return common.Close(m.Manager)
}

func (m *syncInboundHandlerManager) Set(manager inbound.Manager) {
	if manager == nil {
		return
	}

	m.Lock()
	defer m.Unlock()

	common.Close(m.Manager) // nolint: errcheck
	m.Manager = manager
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
