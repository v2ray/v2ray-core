package impl

import (
	"time"

	"v2ray.com/core/app"
	"v2ray.com/core/app/dispatcher"
	"v2ray.com/core/app/proxyman"
	"v2ray.com/core/app/router"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/errors"
	"v2ray.com/core/common/log"
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/common/serial"
	"v2ray.com/core/proxy"
	"v2ray.com/core/transport/ray"
)

type DefaultDispatcher struct {
	ohm    proxyman.OutboundHandlerManager
	router *router.Router
}

func NewDefaultDispatcher(space app.Space) *DefaultDispatcher {
	d := &DefaultDispatcher{}
	space.InitializeApplication(func() error {
		return d.Initialize(space)
	})
	return d
}

// Initialize initializes the dispatcher.
// Private: Used by app.Space only.
func (v *DefaultDispatcher) Initialize(space app.Space) error {
	if !space.HasApp(proxyman.APP_ID_OUTBOUND_MANAGER) {
		return errors.New("DefaultDispatcher: OutboundHandlerManager is not found in the space.")
	}
	v.ohm = space.GetApp(proxyman.APP_ID_OUTBOUND_MANAGER).(proxyman.OutboundHandlerManager)

	if space.HasApp(router.APP_ID) {
		v.router = space.GetApp(router.APP_ID).(*router.Router)
	}

	return nil
}

func (v *DefaultDispatcher) DispatchToOutbound(session *proxy.SessionInfo) ray.InboundRay {
	dispatcher := v.ohm.GetDefaultHandler()
	destination := session.Destination

	if v.router != nil {
		if tag, err := v.router.TakeDetour(session); err == nil {
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

	direct := ray.NewRay()
	var waitFunc func() error
	if session.Inbound != nil && session.Inbound.AllowPassiveConnection {
		waitFunc = noOpWait()
	} else {
		wdi := &waitDataInspector{
			hasData: make(chan bool, 1),
		}
		direct.AddInspector(wdi)
		waitFunc = waitForData(wdi)
	}

	go v.waitAndDispatch(waitFunc, destination, direct, dispatcher)

	return direct
}

func (v *DefaultDispatcher) waitAndDispatch(wait func() error, destination v2net.Destination, link ray.OutboundRay, dispatcher proxy.OutboundHandler) {
	if err := wait(); err != nil {
		log.Info("DefaultDispatcher: Failed precondition: ", err)
		link.OutboundInput().ForceClose()
		link.OutboundOutput().Close()
		return
	}

	dispatcher.Dispatch(destination, link)
}

type DefaultDispatcherFactory struct{}

func (v DefaultDispatcherFactory) Create(space app.Space, config interface{}) (app.Application, error) {
	return NewDefaultDispatcher(space), nil
}

func (v DefaultDispatcherFactory) AppId() app.ID {
	return dispatcher.APP_ID
}

func init() {
	app.RegisterApplicationFactory(serial.GetMessageType(new(dispatcher.Config)), DefaultDispatcherFactory{})
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
