package dispatcher

//go:generate errorgen

import (
	"context"
	"strings"
	"sync"
	"time"

	"v2ray.com/core"
	"v2ray.com/core/app/proxyman"
	"v2ray.com/core/common"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/protocol"
	"v2ray.com/core/common/session"
	"v2ray.com/core/features/outbound"
	"v2ray.com/core/features/policy"
	"v2ray.com/core/features/routing"
	"v2ray.com/core/features/stats"
	"v2ray.com/core/transport"
	"v2ray.com/core/transport/pipe"
)

var (
	errSniffingTimeout = newError("timeout on sniffing")
)

type cachedReader struct {
	sync.Mutex
	reader *pipe.Reader
	cache  buf.MultiBuffer
}

func (r *cachedReader) Cache(b *buf.Buffer) {
	mb, _ := r.reader.ReadMultiBufferTimeout(time.Millisecond * 100)
	r.Lock()
	if !mb.IsEmpty() {
		r.cache, _ = buf.MergeMulti(r.cache, mb)
	}
	b.Clear()
	rawBytes := b.Extend(buf.Size)
	n := r.cache.Copy(rawBytes)
	b.Resize(0, int32(n))
	r.Unlock()
}

func (r *cachedReader) readInternal() buf.MultiBuffer {
	r.Lock()
	defer r.Unlock()

	if r.cache != nil && !r.cache.IsEmpty() {
		mb := r.cache
		r.cache = nil
		return mb
	}

	return nil
}

func (r *cachedReader) ReadMultiBuffer() (buf.MultiBuffer, error) {
	mb := r.readInternal()
	if mb != nil {
		return mb, nil
	}

	return r.reader.ReadMultiBuffer()
}

func (r *cachedReader) ReadMultiBufferTimeout(timeout time.Duration) (buf.MultiBuffer, error) {
	mb := r.readInternal()
	if mb != nil {
		return mb, nil
	}

	return r.reader.ReadMultiBufferTimeout(timeout)
}

func (r *cachedReader) CloseError() {
	r.Lock()
	if r.cache != nil {
		r.cache = buf.ReleaseMulti(r.cache)
	}
	r.Unlock()
	r.reader.CloseError()
}

// DefaultDispatcher is a default implementation of Dispatcher.
type DefaultDispatcher struct {
	ohm    outbound.Manager
	router routing.Router
	policy policy.Manager
	stats  stats.Manager
}

func init() {
	common.Must(common.RegisterConfig((*Config)(nil), func(ctx context.Context, config interface{}) (interface{}, error) {
		d := new(DefaultDispatcher)
		if err := core.RequireFeatures(ctx, func(om outbound.Manager, router routing.Router, pm policy.Manager, sm stats.Manager) error {
			return d.Init(config.(*Config), om, router, pm, sm)
		}); err != nil {
			return nil, err
		}
		return d, nil
	}))
}

// Init initializes DefaultDispatcher.
func (d *DefaultDispatcher) Init(config *Config, om outbound.Manager, router routing.Router, pm policy.Manager, sm stats.Manager) error {
	d.ohm = om
	d.router = router
	d.policy = pm
	d.stats = sm
	return nil
}

// Type implements common.HasType.
func (*DefaultDispatcher) Type() interface{} {
	return routing.DispatcherType()
}

// Start implements common.Runnable.
func (*DefaultDispatcher) Start() error {
	return nil
}

// Close implements common.Closable.
func (*DefaultDispatcher) Close() error { return nil }

func (d *DefaultDispatcher) getLink(ctx context.Context) (*transport.Link, *transport.Link) {
	opt := pipe.OptionsFromContext(ctx)
	uplinkReader, uplinkWriter := pipe.New(opt...)
	downlinkReader, downlinkWriter := pipe.New(opt...)

	inboundLink := &transport.Link{
		Reader: downlinkReader,
		Writer: uplinkWriter,
	}

	outboundLink := &transport.Link{
		Reader: uplinkReader,
		Writer: downlinkWriter,
	}

	sessionInbound := session.InboundFromContext(ctx)
	var user *protocol.MemoryUser
	if sessionInbound != nil {
		user = sessionInbound.User
	}

	if user != nil && len(user.Email) > 0 {
		p := d.policy.ForLevel(user.Level)
		if p.Stats.UserUplink {
			name := "user>>>" + user.Email + ">>>traffic>>>uplink"
			if c, _ := stats.GetOrRegisterCounter(d.stats, name); c != nil {
				inboundLink.Writer = &SizeStatWriter{
					Counter: c,
					Writer:  inboundLink.Writer,
				}
			}
		}
		if p.Stats.UserDownlink {
			name := "user>>>" + user.Email + ">>>traffic>>>downlink"
			if c, _ := stats.GetOrRegisterCounter(d.stats, name); c != nil {
				outboundLink.Writer = &SizeStatWriter{
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

// Dispatch implements routing.Dispatcher.
func (d *DefaultDispatcher) Dispatch(ctx context.Context, destination net.Destination) (*transport.Link, error) {
	if !destination.IsValid() {
		panic("Dispatcher: Invalid destination.")
	}
	ob := &session.Outbound{
		Target: destination,
	}
	ctx = session.ContextWithOutbound(ctx, ob)

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
				ob.Target = destination
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
			if totalAttempt > 2 {
				return nil, errSniffingTimeout
			}

			cReader.Cache(payload)
			if !payload.IsEmpty() {
				result, err := sniffer.Sniff(payload.Bytes())
				if err != common.ErrNoClue {
					return result, err
				}
			}
			if payload.IsFull() {
				return nil, errUnknownContent
			}
		}
	}
}

func (d *DefaultDispatcher) routedDispatch(ctx context.Context, link *transport.Link, destination net.Destination) {
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

	if dispatcher == nil {
		newError("default outbound handler not exist").WriteToLog(session.ExportIDToError(ctx))
		common.Close(link.Writer)
		pipe.CloseError(link.Reader)
		return
	}

	dispatcher.Dispatch(ctx, link)
}
