package dispatcher

//go:generate go run $GOPATH/src/v2ray.com/core/common/errors/errorgen/main.go -pkg impl -path App,Dispatcher,Default

import (
	"context"
	"time"

	"v2ray.com/core"
	"v2ray.com/core/app/proxyman"
	"v2ray.com/core/common"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/net"
	"v2ray.com/core/proxy"
	"v2ray.com/core/transport/ray"
)

var (
	errSniffingTimeout = newError("timeout on sniffing")
)

// DefaultDispatcher is a default implementation of Dispatcher.
type DefaultDispatcher struct {
	ohm    core.OutboundHandlerManager
	router core.Router
}

// NewDefaultDispatcher create a new DefaultDispatcher.
func NewDefaultDispatcher(ctx context.Context, config *Config) (*DefaultDispatcher, error) {
	v := core.FromContext(ctx)
	if v == nil {
		return nil, newError("V is not in context.")
	}

	d := &DefaultDispatcher{
		ohm:    v.OutboundHandlerManager(),
		router: v.Router(),
	}

	if err := v.RegisterFeature((*core.Dispatcher)(nil), d); err != nil {
		return nil, newError("unable to register Dispatcher")
	}
	return d, nil
}

// Start implements app.Application.
func (*DefaultDispatcher) Start() error {
	return nil
}

// Close implements app.Application.
func (*DefaultDispatcher) Close() error { return nil }

// Dispatch implements core.Dispatcher.
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
			domain, err := snifer(ctx, sniferList, outbound)
			if err == nil {
				newError("sniffed domain: ", domain).WriteToLog()
				destination.Address = net.ParseAddress(domain)
				ctx = proxy.ContextWithTarget(ctx, destination)
			}
			d.routedDispatch(ctx, outbound, destination)
		}()
	}
	return outbound, nil
}

func snifer(ctx context.Context, sniferList []proxyman.KnownProtocols, outbound ray.OutboundRay) (string, error) {
	payload := buf.New()
	defer payload.Release()

	sniffer := NewSniffer(sniferList)
	totalAttempt := 0
	for {
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		default:
			totalAttempt++
			if totalAttempt > 5 {
				return "", errSniffingTimeout
			}
			outbound.OutboundInput().Peek(payload)
			if !payload.IsEmpty() {
				domain, err := sniffer.Sniff(payload.Bytes())
				if err != ErrMoreData {
					return domain, err
				}
			}
			if payload.IsFull() {
				return "", ErrInvalidData
			}
			time.Sleep(time.Millisecond * 100)
		}
	}
}

func (d *DefaultDispatcher) routedDispatch(ctx context.Context, outbound ray.OutboundRay, destination net.Destination) {
	dispatcher := d.ohm.GetDefaultHandler()
	if d.router != nil {
		if tag, err := d.router.PickRoute(ctx); err == nil {
			if handler := d.ohm.GetHandler(tag); handler != nil {
				newError("taking detour [", tag, "] for [", destination, "]").WriteToLog()
				dispatcher = handler
			} else {
				newError("nonexisting tag: ", tag).AtWarning().WriteToLog()
			}
		} else {
			newError("default route for ", destination).WriteToLog()
		}
	}
	dispatcher.Dispatch(ctx, outbound)
}

func init() {
	common.Must(common.RegisterConfig((*Config)(nil), func(ctx context.Context, config interface{}) (interface{}, error) {
		return NewDefaultDispatcher(ctx, config.(*Config))
	}))
}
