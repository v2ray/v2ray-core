package impl

import (
	"context"
	"time"

	"v2ray.com/core/app"
	"v2ray.com/core/app/dispatcher"
	"v2ray.com/core/app/proxyman"
	"v2ray.com/core/app/router"
	"v2ray.com/core/common"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/errors"
	"v2ray.com/core/common/log"
	"v2ray.com/core/proxy"
	"v2ray.com/core/transport/ray"
)

type DefaultDispatcher struct {
	ohm    proxyman.OutboundHandlerManager
	router *router.Router
}

func NewDefaultDispatcher(ctx context.Context, config *dispatcher.Config) (*DefaultDispatcher, error) {
	space := app.SpaceFromContext(ctx)
	if space == nil {
		return nil, errors.New("DefaultDispatcher: No space in context.")
	}
	d := &DefaultDispatcher{}
	space.OnInitialize(func() error {
		d.ohm = proxyman.OutboundHandlerManagerFromSpace(space)
		if d.ohm == nil {
			return errors.New("DefaultDispatcher: OutboundHandlerManager is not found in the space.")
		}
		d.router = router.FromSpace(space)
		return nil
	})
	return d, nil
}

func (DefaultDispatcher) Interface() interface{} {
	return (*dispatcher.Interface)(nil)
}

func (v *DefaultDispatcher) DispatchToOutbound(ctx context.Context) ray.InboundRay {
	dispatcher := v.ohm.GetDefaultHandler()
	destination := proxy.DestinationFromContext(ctx)
	if !destination.IsValid() {
		panic("Dispatcher: Invalid destination.")
	}

	if v.router != nil {
		if tag, err := v.router.TakeDetour(ctx); err == nil {
			if handler := v.ohm.GetHandler(tag); handler != nil {
				log.Info("DefaultDispatcher: Taking detour [", tag, "] for [", destination, "].")
				dispatcher = handler
			} else {
				log.Warning("DefaultDispatcher: Nonexisting tag: ", tag)
			}
		} else {
			log.Info("DefaultDispatcher: Default route for ", destination)
		}
	}

	direct := ray.NewRay(ctx)
	var waitFunc func() error
	if allowPassiveConnection, ok := proxy.AllowPassiveConnectionFromContext(ctx); ok && allowPassiveConnection {
		waitFunc = noOpWait()
	} else {
		wdi := &waitDataInspector{
			hasData: make(chan bool, 1),
		}
		direct.AddInspector(wdi)
		waitFunc = waitForData(wdi)
	}

	go v.waitAndDispatch(ctx, waitFunc, direct, dispatcher)

	return direct
}

func (v *DefaultDispatcher) waitAndDispatch(ctx context.Context, wait func() error, link ray.OutboundRay, dispatcher proxyman.OutboundHandler) {
	if err := wait(); err != nil {
		log.Info("DefaultDispatcher: Failed precondition: ", err)
		link.OutboundInput().CloseError()
		link.OutboundOutput().CloseError()
		return
	}

	dispatcher.Dispatch(ctx, link)
}

func init() {
	common.Must(common.RegisterConfig((*dispatcher.Config)(nil), func(ctx context.Context, config interface{}) (interface{}, error) {
		return NewDefaultDispatcher(ctx, config.(*dispatcher.Config))
	}))
}

type waitDataInspector struct {
	hasData chan bool
}

func (wdi *waitDataInspector) Input(*buf.Buffer) {
	select {
	case wdi.hasData <- true:
	default:
	}
}

func (wdi *waitDataInspector) WaitForData() bool {
	select {
	case <-wdi.hasData:
		return true
	case <-time.After(time.Minute):
		return false
	}
}

func waitForData(wdi *waitDataInspector) func() error {
	return func() error {
		if wdi.WaitForData() {
			return nil
		}
		return errors.New("DefaultDispatcher: No data.")
	}
}

func noOpWait() func() error {
	return func() error {
		return nil
	}
}
