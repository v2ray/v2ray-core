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
	"v2ray.com/core/common/protocol"
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
	policy core.PolicyManager
	stats  core.StatManager
}

// NewDefaultDispatcher create a new DefaultDispatcher.
func NewDefaultDispatcher(ctx context.Context, config *Config) (*DefaultDispatcher, error) {
	v := core.MustFromContext(ctx)
	d := &DefaultDispatcher{
		ohm:    v.OutboundHandlerManager(),
		router: v.Router(),
		policy: v.PolicyManager(),
		stats:  v.Stats(),
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

func getStatsName(u *protocol.User) string {
	return "user>traffic>" + u.Email
}

// Dispatch implements core.Dispatcher.
func (d *DefaultDispatcher) Dispatch(ctx context.Context, destination net.Destination) (ray.InboundRay, error) {
	if !destination.IsValid() {
		panic("Dispatcher: Invalid destination.")
	}
	ctx = proxy.ContextWithTarget(ctx, destination)

	var rayOptions []ray.Option

	user := protocol.UserFromContext(ctx)
	if user != nil && len(user.Email) > 0 {
		name := getStatsName(user)
		c, err := d.stats.RegisterCounter(name)
		if err != nil {
			c = d.stats.GetCounter(name)
		}
		if c == nil {
			newError("failed to get stats counter ", name).AtWarning().WithContext(ctx).WriteToLog()
		}

		p := d.policy.ForLevel(user.Level)
		if p.Stats.EnablePerUser {
			rayOptions = append(rayOptions, ray.WithStatCounter(c))
		}
	}

	outbound := ray.New(ctx, rayOptions...)
	snifferList := proxyman.ProtocolSniffersFromContext(ctx)
	if destination.Address.Family().IsDomain() || len(snifferList) == 0 {
		go d.routedDispatch(ctx, outbound, destination)
	} else {
		go func() {
			domain, err := sniffer(ctx, snifferList, outbound)
			if err == nil {
				newError("sniffed domain: ", domain).WithContext(ctx).WriteToLog()
				destination.Address = net.ParseAddress(domain)
				ctx = proxy.ContextWithTarget(ctx, destination)
			}
			d.routedDispatch(ctx, outbound, destination)
		}()
	}
	return outbound, nil
}

func sniffer(ctx context.Context, snifferList []proxyman.KnownProtocols, outbound ray.OutboundRay) (string, error) {
	payload := buf.New()
	defer payload.Release()

	sniffer := NewSniffer(snifferList)
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
				newError("taking detour [", tag, "] for [", destination, "]").WithContext(ctx).WriteToLog()
				dispatcher = handler
			} else {
				newError("non existing tag: ", tag).AtWarning().WithContext(ctx).WriteToLog()
			}
		} else {
			newError("default route for ", destination).WithContext(ctx).WriteToLog()
		}
	}
	dispatcher.Dispatch(ctx, outbound)
}

func init() {
	common.Must(common.RegisterConfig((*Config)(nil), func(ctx context.Context, config interface{}) (interface{}, error) {
		return NewDefaultDispatcher(ctx, config.(*Config))
	}))
}
