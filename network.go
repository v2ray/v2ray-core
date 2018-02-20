package core

import (
	"context"
	"sync"

	"v2ray.com/core/common"
	"v2ray.com/core/common/net"
	"v2ray.com/core/transport/ray"
)

// InboundHandler is the interface for handlers that process inbound connections.
type InboundHandler interface {
	common.Runnable
	// The tag of this handler.
	Tag() string

	// Deprecated. Do not use in new code.
	GetRandomInboundProxy() (interface{}, net.Port, int)
}

// OutboundHandler is the interface for handlers that process outbound connections.
type OutboundHandler interface {
	common.Runnable
	Tag() string
	Dispatch(ctx context.Context, outboundRay ray.OutboundRay)
}

// InboundHandlerManager is a feature that managers InboundHandlers.
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

	m.Close()

	m.Lock()
	defer m.Unlock()

	m.InboundHandlerManager = manager
}

// OutboundHandlerManager is a feature that manages OutboundHandlers.
type OutboundHandlerManager interface {
	Feature
	// GetHandler returns an OutboundHandler will given tag.
	GetHandler(tag string) OutboundHandler
	// GetDefaultHandler returns the default OutboundHandler. It is usually the first OutboundHandler specified in the configuration.
	GetDefaultHandler() OutboundHandler
	// AddHandler adds a handler into this OutboundHandlerManager.
	AddHandler(ctx context.Context, handler OutboundHandler) error

	// RemoveHandler removes a handler from OutboundHandlerManager.
	RemoveHandler(ctx context.Context, tag string) error
}

type syncOutboundHandlerManager struct {
	sync.RWMutex
	OutboundHandlerManager
}

func (m *syncOutboundHandlerManager) GetHandler(tag string) OutboundHandler {
	m.RLock()
	defer m.RUnlock()

	if m.OutboundHandlerManager == nil {
		return nil
	}

	return m.OutboundHandlerManager.GetHandler(tag)
}

func (m *syncOutboundHandlerManager) GetDefaultHandler() OutboundHandler {
	m.RLock()
	defer m.RUnlock()

	if m.OutboundHandlerManager == nil {
		return nil
	}

	return m.OutboundHandlerManager.GetDefaultHandler()
}

func (m *syncOutboundHandlerManager) AddHandler(ctx context.Context, handler OutboundHandler) error {
	m.RLock()
	defer m.RUnlock()

	if m.OutboundHandlerManager == nil {
		return newError("OutboundHandlerManager not set.").AtError()
	}

	return m.OutboundHandlerManager.AddHandler(ctx, handler)
}

func (m *syncOutboundHandlerManager) Start() error {
	m.RLock()
	defer m.RUnlock()

	if m.OutboundHandlerManager == nil {
		return newError("OutboundHandlerManager not set.").AtError()
	}

	return m.OutboundHandlerManager.Start()
}

func (m *syncOutboundHandlerManager) Close() error {
	m.RLock()
	defer m.RUnlock()

	return common.Close(m.OutboundHandlerManager)
}

func (m *syncOutboundHandlerManager) Set(manager OutboundHandlerManager) {
	if manager == nil {
		return
	}

	m.Close()
	m.Lock()
	defer m.Unlock()

	m.OutboundHandlerManager = manager
}
