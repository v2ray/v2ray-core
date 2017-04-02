package mux

import (
	"context"
	"sync"
	"time"

	"v2ray.com/core/app/dispatcher"
	"v2ray.com/core/app/log"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/signal"
	"v2ray.com/core/proxy"
	"v2ray.com/core/transport/ray"
)

const (
	maxParallel = 8
	maxTotal    = 128
)

type manager interface {
	remove(id uint16)
}

type session struct {
	sync.Mutex
	input          ray.InputStream
	output         ray.OutputStream
	parent         manager
	id             uint16
	uplinkClosed   bool
	downlinkClosed bool
}

func (s *session) checkAndRemove() {
	s.Lock()
	if s.uplinkClosed && s.downlinkClosed {
		s.parent.remove(s.id)
	}
	s.Unlock()
}

func (s *session) closeUplink() {
	s.Lock()
	s.uplinkClosed = true
	s.Unlock()
	s.checkAndRemove()
}

func (s *session) closeDownlink() {
	s.Lock()
	s.downlinkClosed = true
	s.Unlock()
	s.checkAndRemove()
}

type Client struct {
	access     sync.RWMutex
	count      uint16
	sessions   map[uint16]*session
	inboundRay ray.InboundRay
	ctx        context.Context
	cancel     context.CancelFunc
}

var muxCoolDestination = net.TCPDestination(net.DomainAddress("v1.mux.cool"), net.Port(9527))

func NewClient(p proxy.Outbound, dialer proxy.Dialer) (*Client, error) {
	ctx, cancel := context.WithCancel(context.Background())
	ctx = proxy.ContextWithTarget(ctx, muxCoolDestination)
	pipe := ray.NewRay(ctx)
	err := p.Process(ctx, pipe, dialer)
	if err != nil {
		cancel()
		return nil, err
	}
	return &Client{
		sessions:   make(map[uint16]*session, 256),
		inboundRay: pipe,
		ctx:        ctx,
		cancel:     cancel,
	}, nil
}

func (m *Client) isFullyOccupied() bool {
	m.access.RLock()
	defer m.access.RUnlock()

	return len(m.sessions) >= maxParallel
}

func (m *Client) remove(id uint16) {
	m.access.Lock()
	defer m.access.Unlock()

	delete(m.sessions, id)

	if len(m.sessions) == 0 {
		m.cancel()
		m.inboundRay.InboundInput().Close()
	}
}

func (m *Client) Closed() bool {
	select {
	case <-m.ctx.Done():
		return true
	default:
		return false
	}
}

func (m *Client) fetchInput(ctx context.Context, s *session) {
	dest, _ := proxy.TargetFromContext(ctx)
	writer := &Writer{
		dest:   dest,
		id:     s.id,
		writer: m.inboundRay.InboundInput(),
	}
	_, timer := signal.CancelAfterInactivity(ctx, time.Minute*5)
	buf.PipeUntilEOF(timer, s.input, writer)
	writer.Close()
	s.closeUplink()
}

func (m *Client) Dispatch(ctx context.Context, outboundRay ray.OutboundRay) bool {
	m.access.Lock()
	defer m.access.Unlock()

	if len(m.sessions) >= maxParallel {
		return false
	}

	if m.count >= maxTotal {
		return false
	}

	select {
	case <-m.ctx.Done():
		return false
	default:
	}

	m.count++
	id := m.count
	s := &session{
		input:  outboundRay.OutboundInput(),
		output: outboundRay.OutboundOutput(),
		parent: m,
		id:     id,
	}
	m.sessions[id] = s
	go m.fetchInput(ctx, s)
	return true
}

func (m *Client) fetchOutput() {
	reader := NewReader(m.inboundRay.InboundOutput())
	for {
		meta, err := reader.ReadMetadata()
		if err != nil {
			break
		}
		m.access.RLock()
		s, found := m.sessions[meta.SessionID]
		m.access.RUnlock()
		if found && meta.SessionStatus == SessionStatusEnd {
			s.closeDownlink()
			s.output.Close()
		}
		if !meta.Option.Has(OptionData) {
			continue
		}

		for {
			data, more, err := reader.Read()
			if err != nil {
				break
			}
			if found {
				if err := s.output.Write(data); err != nil {
					break
				}
			}
			if !more {
				break
			}
		}
	}
}

type Server struct {
	dispatcher dispatcher.Interface
}

func (s *Server) Dispatch(ctx context.Context, dest net.Destination) (ray.InboundRay, error) {
	if dest != muxCoolDestination {
		return s.dispatcher.Dispatch(ctx, dest)
	}

	ray := ray.NewRay(ctx)

	return ray, nil
}

type ServerWorker struct {
	dispatcher  dispatcher.Interface
	outboundRay ray.OutboundRay
	sessions    map[uint16]*session
	access      sync.RWMutex
}

func (w *ServerWorker) remove(id uint16) {
	w.access.Lock()
	delete(w.sessions, id)
	w.access.Unlock()
}

func (w *ServerWorker) handle(ctx context.Context, s *session) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			data, err := s.input.Read()
			if err != nil {
				return
			}
			w.outboundRay.OutboundOutput().Write(data)
		}
	}
}

func (w *ServerWorker) run(ctx context.Context) {
	input := w.outboundRay.OutboundInput()
	reader := NewReader(input)
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		meta, err := reader.ReadMetadata()
		if err != nil {
			return
		}

		w.access.RLock()
		s, found := w.sessions[meta.SessionID]
		w.access.RUnlock()

		if found && meta.SessionStatus == SessionStatusEnd {
			s.closeUplink()
			s.output.Close()
		}

		if meta.SessionStatus == SessionStatusNew {
			inboundRay, err := w.dispatcher.Dispatch(ctx, meta.Target)
			if err != nil {
				log.Info("Proxyman|Mux: Failed to dispatch request: ", err)
				continue
			}
			s = &session{
				input:  inboundRay.InboundOutput(),
				output: inboundRay.InboundInput(),
				parent: w,
				id:     meta.SessionID,
			}
			w.access.Lock()
			w.sessions[meta.SessionID] = s
			w.access.Unlock()
			go w.handle(ctx, s)
		}

		if meta.Option.Has(OptionData) {
			for {
				data, more, err := reader.Read()
				if err != nil {
					break
				}
				if s != nil {
					s.output.Write(data)
				}
				if !more {
					break
				}
			}
		}
	}
}
