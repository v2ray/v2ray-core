package core

import (
	"context"
	"sync"

	"v2ray.com/core/common"
	"v2ray.com/core/common/errors"
	"v2ray.com/core/common/net"
	"v2ray.com/core/transport/ray"
)

// Dispatcher is a feature that dispatches inbound requests to outbound handlers based on rules.
// Dispatcher is required to be registered in a V2Ray instance to make V2Ray function properly.
type Dispatcher interface {
	Feature

	// Dispatch returns a Ray for transporting data for the given request.
	Dispatch(ctx context.Context, dest net.Destination) (ray.InboundRay, error)
}

type syncDispatcher struct {
	sync.RWMutex
	Dispatcher
}

func (d *syncDispatcher) Dispatch(ctx context.Context, dest net.Destination) (ray.InboundRay, error) {
	d.RLock()
	defer d.RUnlock()

	if d.Dispatcher == nil {
		return nil, newError("Dispatcher not set.").AtError()
	}

	return d.Dispatcher.Dispatch(ctx, dest)
}

func (d *syncDispatcher) Start() error {
	d.RLock()
	defer d.RUnlock()

	if d.Dispatcher == nil {
		return newError("Dispatcher not set.").AtError()
	}

	return d.Dispatcher.Start()
}

func (d *syncDispatcher) Close() error {
	d.RLock()
	defer d.RUnlock()

	return common.Close(d.Dispatcher)
}

func (d *syncDispatcher) Set(disp Dispatcher) {
	if disp == nil {
		return
	}

	d.Lock()
	defer d.Unlock()

	common.Close(d.Dispatcher)
	d.Dispatcher = disp
}

var (
	// ErrNoClue is for the situation that existing information is not enough to make a decision. For example, Router may return this error when there is no suitable route.
	ErrNoClue = errors.New("not enough information for making a decision")
)

// Router is a feature to choose a outbound tag for the given request.
type Router interface {
	Feature

	// PickRoute returns a tag of an OutboundHandler based on the given context.
	PickRoute(ctx context.Context) (string, error)
}

type syncRouter struct {
	sync.RWMutex
	Router
}

func (r *syncRouter) PickRoute(ctx context.Context) (string, error) {
	r.RLock()
	defer r.RUnlock()

	if r.Router == nil {
		return "", ErrNoClue
	}

	return r.Router.PickRoute(ctx)
}

func (r *syncRouter) Start() error {
	r.RLock()
	defer r.RUnlock()

	if r.Router == nil {
		return nil
	}

	return r.Router.Start()
}

func (r *syncRouter) Close() error {
	r.RLock()
	defer r.RUnlock()

	return common.Close(r.Router)
}

func (r *syncRouter) Set(router Router) {
	if router == nil {
		return
	}

	r.Lock()
	defer r.Unlock()

	common.Close(r.Router)
	r.Router = router
}
