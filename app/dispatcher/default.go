package dispatcher

//go:generate go run $GOPATH/src/v2ray.com/core/common/errors/errorgen/main.go -pkg impl -path App,Dispatcher,Default

import (
	"context"
	"strings"
	"time"

	"v2ray.com/core"
	"v2ray.com/core/app/proxyman"
	"v2ray.com/core/common"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/protocol"
	"v2ray.com/core/common/session"
	"v2ray.com/core/common/stats"
	"v2ray.com/core/proxy"
	"v2ray.com/core/transport/pipe"
)

var (
	errSniffingTimeout = newError("timeout on sniffing")
)

type cachedReader struct {
	reader *pipe.Reader
	cache  buf.MultiBuffer
}

func (r *cachedReader) Cache(b *buf.Buffer) {
	mb, _ := r.reader.ReadMultiBufferTimeout(time.Millisecond * 100)
	if !mb.IsEmpty() {
		common.Must(r.cache.WriteMultiBuffer(mb))
	}
	common.Must(b.Reset(func(x []byte) (int, error) {
		return r.cache.Copy(x), nil
	}))
}

func (r *cachedReader) ReadMultiBuffer() (buf.MultiBuffer, error) {
	if !r.cache.IsEmpty() {
		mb := r.cache
		r.cache = nil
		return mb, nil
	}

	return r.reader.ReadMultiBuffer()
}

func (r *cachedReader) ReadMultiBufferTimeout(timeout time.Duration) (buf.MultiBuffer, error) {
	if !r.cache.IsEmpty() {
		mb := r.cache
		r.cache = nil
		return mb, nil
	}

	return r.reader.ReadMultiBufferTimeout(timeout)
}

func (r *cachedReader) CloseError() {
	r.cache.Release()
	r.reader.CloseError()
}

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
		return nil, newError("unable to register Dispatcher").Base(err)
	}
	return d, nil
}

// Start implements common.Runnable.
func (*DefaultDispatcher) Start() error {
	return nil
}

// Close implements common.Closable.
func (*DefaultDispatcher) Close() error { return nil }

func (d *DefaultDispatcher) getLink(ctx context.Context) (*core.Link, *core.Link) {
	opt := pipe.OptionsFromContext(ctx)
	uplinkReader, uplinkWriter := pipe.New(opt...)
	downlinkReader, downlinkWriter := pipe.New(opt...)

	inboundLink := &core.Link{
		Reader: downlinkReader,
		Writer: uplinkWriter,
	}

	outboundLink := &core.Link{
		Reader: uplinkReader,
		Writer: downlinkWriter,
	}

	user := protocol.UserFromContext(ctx)
	if user != nil && len(user.Email) > 0 {
		p := d.policy.ForLevel(user.Level)
		if p.Stats.UserUplink {
			name := "user>>>" + user.Email + ">>>traffic>>>uplink"
			if c, _ := core.GetOrRegisterStatCounter(d.stats, name); c != nil {
				inboundLink.Writer = &stats.SizeStatWriter{
					Counter: c,
					Writer:  inboundLink.Writer,
				}
			}
		}
		if p.Stats.UserDownlink {
			name := "user>>>" + user.Email + ">>>traffic>>>downlink"
			if c, _ := core.GetOrRegisterStatCounter(d.stats, name); c != nil {
				outboundLink.Writer = &stats.SizeStatWriter{
					Counter: c,
					Writer:  outboundLink.Writer,
				}
			}
		}
	}

	return inboundLink, outboundLink
}

func shouldOverride(result SniffResult, domainOverride []string) bool {
	for _, p := range domainOverride {
		if strings.HasPrefix(result.Protocol(), p) {
			return true
		}
	}
	return false
}

// Dispatch implements core.Dispatcher.
func (d *DefaultDispatcher) Dispatch(ctx context.Context, destination net.Destination) (*core.Link, error) {
	if !destination.IsValid() {
		panic("Dispatcher: Invalid destination.")
	}
	ctx = proxy.ContextWithTarget(ctx, destination)

	inbound, outbound := d.getLink(ctx)
	sniffingConfig := proxyman.SniffingConfigFromContext(ctx)
	if destination.Network != net.Network_TCP || sniffingConfig == nil || !sniffingConfig.Enabled {
		go d.routedDispatch(ctx, outbound, destination)
	} else {
		go func() {
			cReader := &cachedReader{
				reader: outbound.Reader.(*pipe.Reader),
			}
			outbound.Reader = cReader
			result, err := sniffer(ctx, cReader)
			if err == nil {
				ctx = ContextWithSniffingResult(ctx, result)
			}
			if err == nil && shouldOverride(result, sniffingConfig.DestinationOverride) {
				domain := result.Domain()
				newError("sniffed domain: ", domain).WriteToLog(session.ExportIDToError(ctx))
				destination.Address = net.ParseAddress(domain)
				ctx = proxy.ContextWithTarget(ctx, destination)
			}
			d.routedDispatch(ctx, outbound, destination)
		}()
	}
	return inbound, nil
}

func sniffer(ctx context.Context, cReader *cachedReader) (SniffResult, error) {
	payload := buf.New()
	defer payload.Release()

	sniffer := NewSniffer()
	totalAttempt := 0
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			totalAttempt++
			if totalAttempt > 5 {
				return nil, errSniffingTimeout
			}

			cReader.Cache(payload)
			if !payload.IsEmpty() {
				result, err := sniffer.Sniff(payload.Bytes())
				if err != core.ErrNoClue {
					return result, err
				}
			}
			if payload.IsFull() {
				return nil, errUnknownContent
			}
			time.Sleep(time.Millisecond * 100)
		}
	}
}

func (d *DefaultDispatcher) routedDispatch(ctx context.Context, link *core.Link, destination net.Destination) {
	dispatcher := d.ohm.GetDefaultHandler()
	if d.router != nil {
		if tag, err := d.router.PickRoute(ctx); err == nil {
			if handler := d.ohm.GetHandler(tag); handler != nil {
				newError("taking detour [", tag, "] for [", destination, "]").WriteToLog(session.ExportIDToError(ctx))
				dispatcher = handler
			} else {
				newError("non existing tag: ", tag).AtWarning().WriteToLog(session.ExportIDToError(ctx))
			}
		} else {
			newError("default route for ", destination).WriteToLog(session.ExportIDToError(ctx))
		}
	}
	dispatcher.Dispatch(ctx, link)
}

func init() {
	common.Must(common.RegisterConfig((*Config)(nil), func(ctx context.Context, config interface{}) (interface{}, error) {
		return NewDefaultDispatcher(ctx, config.(*Config))
	}))
}
