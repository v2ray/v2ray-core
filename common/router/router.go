package router

import (
	"context"
	"sync"
)

type Router interface {
	Pick(ctx context.Context) (string, bool)
}

type defaultRouter byte

func (defaultRouter) Pick(ctx context.Context) (string, bool) {
	return "", false
}

type syncRouter struct {
	sync.RWMutex
	Router
}

func (r *syncRouter) Pick(ctx context.Context) (string, bool) {
	r.RLock()
	defer r.RUnlock()

	return r.Router.Pick(ctx)
}

func (r *syncRouter) Set(router Router) {
	r.Lock()
	defer r.Unlock()

	r.Router = router
}

var (
	routerInstance = &syncRouter{
		Router: defaultRouter(0),
	}
)

func RegisterRouter(router Router) {
	if router == nil {
		panic("Router is nil.")
	}

	routerInstance.Set(router)
}

func Pick(ctx context.Context) (string, bool) {
	return routerInstance.Router.Pick(ctx)
}
