package impl

//go:generate go run $GOPATH/src/v2ray.com/core/tools/generrorgen/main.go -pkg impl -path App,Dispatcher,Default

import (
	"context"
	"time"

	"v2ray.com/core/app"
	"v2ray.com/core/app/dispatcher"
	"v2ray.com/core/app/log"
	"v2ray.com/core/app/proxyman"
	"v2ray.com/core/app/router"
	"v2ray.com/core/common"
	"v2ray.com/core/common/net"
	"v2ray.com/core/proxy"
	"v2ray.com/core/transport/ray"
)

var (
	errSniffingTimeout = newError("timeout on sniffing")
)

type DefaultDispatcher struct {
	ohm    proxyman.OutboundHandlerManager
	router *router.Router
}

func NewDefaultDispatcher(ctx context.Context, config *dispatcher.Config) (*DefaultDispatcher, error) {
	space := app.SpaceFromContext(ctx)
	if space == nil {
		return nil, newError("no space in context")
	}
	d := &DefaultDispatcher{}
	space.OnInitialize(func() error {
		d.ohm = proxyman.OutboundHandlerManagerFromSpace(space)
		if d.ohm == nil {
			return newError("OutboundHandlerManager is not found in the space")
		}
		d.router = router.FromSpace(space)
		return nil
	})
	return d, nil
}

func (DefaultDispatcher) Start() error {
	return nil
}

func (DefaultDispatcher) Close() {}

func (DefaultDispatcher) Interface() interface{} {
	return (*dispatcher.Interface)(nil)
}

type domainOrError struct {
	domain string
	err    error
}

func (d *DefaultDispatcher) Dispatch(ctx context.Context, destination net.Destination) (ray.InboundRay, error) {
	if !destination.IsValid() {
		panic("Dispatcher: Invalid destination.")
	}
	ctx = proxy.ContextWithTarget(ctx, destination)

	outbound := ray.NewRay(ctx)
	sniferList := proxyman.ProtocoSniffersFromContext(ctx)
	if destination.Address.Family().IsDomain() || len(sniferList) == 0 {
		go d.routedDispatch(ctx, outbound, destination)
	} else {
		go func() {
			done := make(chan domainOrError)
			go snifer(ctx, sniferList, outbound, done)
			de := <-done
			if de.err != nil {
				log.Trace(newError("failed to snif").Base(de.err))
				return
			}
			log.Trace(newError("sniffed domain: ", de.domain))
			destination.Address = net.ParseAddress(de.domain)
			ctx = proxy.ContextWithTarget(ctx, destination)
			d.routedDispatch(ctx, outbound, destination)
		}()
	}
	return outbound, nil
}

func snifer(ctx context.Context, sniferList []proxyman.KnownProtocols, outbound ray.OutboundRay, done chan<- domainOrError) {
	payload := make([]byte, 2048)
	totalAttempt := 0
	for {
		select {
		case <-ctx.Done():
			done <- domainOrError{
				domain: "",
				err:    ctx.Err(),
			}
			return
		case <-time.After(time.Millisecond * 100):
			totalAttempt++
			if totalAttempt > 5 {
				done <- domainOrError{
					domain: "",
					err:    errSniffingTimeout,
				}
				return
			}
			mb := outbound.OutboundInput().Peek()
			if mb.IsEmpty() {
				continue
			}
			nBytes, _ := mb.Read(payload)
			for _, protocol := range sniferList {
				var f func([]byte) (string, error)
				switch protocol {
				case proxyman.KnownProtocols_HTTP:
					f = SniffHTTP
				case proxyman.KnownProtocols_TLS:
					f = SniffTLS
				default:
					panic("Unsupported protocol")
				}

				domain, err := f(payload[:nBytes])
				if err != ErrMoreData {
					done <- domainOrError{
						domain: domain,
						err:    err,
					}
					return
				}
			}
			if nBytes == 2048 {
				done <- domainOrError{
					domain: "",
					err:    ErrInvalidData,
				}
				return
			}
		}
	}
}

func (d *DefaultDispatcher) routedDispatch(ctx context.Context, outbound ray.OutboundRay, destination net.Destination) {
	dispatcher := d.ohm.GetDefaultHandler()
	if d.router != nil {
		if tag, err := d.router.TakeDetour(ctx); err == nil {
			if handler := d.ohm.GetHandler(tag); handler != nil {
				log.Trace(newError("taking detour [", tag, "] for [", destination, "]"))
				dispatcher = handler
			} else {
				log.Trace(newError("nonexisting tag: ", tag).AtWarning())
			}
		} else {
			log.Trace(newError("default route for ", destination))
		}
	}
	dispatcher.Dispatch(ctx, outbound)
}

func init() {
	common.Must(common.RegisterConfig((*dispatcher.Config)(nil), func(ctx context.Context, config interface{}) (interface{}, error) {
		return NewDefaultDispatcher(ctx, config.(*dispatcher.Config))
	}))
}
