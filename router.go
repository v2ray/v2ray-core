package core

import (
	"context"
	"sync"

	"v2ray.com/core/common"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/vio"
	"v2ray.com/core/features/routing"
)

type syncDispatcher struct {
	sync.RWMutex
	routing.Dispatcher
}

func (d *syncDispatcher) Dispatch(ctx context.Context, dest net.Destination) (*vio.Link, error) {
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

func (d *syncDispatcher) Set(disp routing.Dispatcher) {
	if disp == nil {
		return
	}

	d.Lock()
	defer d.Unlock()

	common.Close(d.Dispatcher) // nolint: errorcheck
	d.Dispatcher = disp
}

type syncRouter struct {
	sync.RWMutex
	routing.Router
}

func (r *syncRouter) PickRoute(ctx context.Context) (string, error) {
	r.RLock()
	defer r.RUnlock()

	if r.Router == nil {
		return "", common.ErrNoClue
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

func (r *syncRouter) Set(router routing.Router) {
	if router == nil {
		return
	}

	r.Lock()
	defer r.Unlock()

	common.Close(r.Router) // nolint: errcheck
	r.Router = router
}
