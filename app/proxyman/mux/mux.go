package mux

import (
	"context"
	"sync"
	"time"

	"v2ray.com/core/app"
	"v2ray.com/core/app/dispatcher"
	"v2ray.com/core/app/log"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/errors"
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

func (s *session) closeUplink() {
	var allDone bool
	s.Lock()
	s.uplinkClosed = true
	allDone = s.uplinkClosed && s.downlinkClosed
	s.Unlock()
	if allDone {
		s.parent.remove(s.id)
	}
}

func (s *session) closeDownlink() {
	var allDone bool
	s.Lock()
	s.downlinkClosed = true
	allDone = s.uplinkClosed && s.downlinkClosed
	s.Unlock()
	if allDone {
		s.parent.remove(s.id)
	}
}

type ClientManager struct {
	access  sync.Mutex
	clients []*Client
	proxy   proxy.Outbound
	dialer  proxy.Dialer
}

func NewClientManager(p proxy.Outbound, d proxy.Dialer) *ClientManager {
	return &ClientManager{
		proxy:  p,
		dialer: d,
	}
}

func (m *ClientManager) Dispatch(ctx context.Context, outboundRay ray.OutboundRay) error {
	m.access.Lock()
	defer m.access.Unlock()

	for _, client := range m.clients {
		if client.Dispatch(ctx, outboundRay) {
			return nil
		}
	}

	client, err := NewClient(m.proxy, m.dialer, m)
	if err != nil {
		return errors.Base(err).Message("Proxyman|Mux|ClientManager: Failed to create client.")
	}
	m.clients = append(m.clients, client)
	client.Dispatch(ctx, outboundRay)
	return nil
}

func (m *ClientManager) onClientFinish() {
	m.access.Lock()
	defer m.access.Unlock()

	if len(m.clients) < 10 {
		return
	}

	activeClients := make([]*Client, 0, len(m.clients))

	for _, client := range m.clients {
		if !client.Closed() {
			activeClients = append(activeClients, client)
		}
	}
	m.clients = activeClients
}

type Client struct {
	access         sync.RWMutex
	count          uint16
	sessions       map[uint16]*session
	inboundRay     ray.InboundRay
	ctx            context.Context
	cancel         context.CancelFunc
	manager        *ClientManager
	session2Remove chan uint16
}

var muxCoolDestination = net.TCPDestination(net.DomainAddress("v1.mux.cool"), net.Port(9527))

func NewClient(p proxy.Outbound, dialer proxy.Dialer, m *ClientManager) (*Client, error) {
	ctx, cancel := context.WithCancel(context.Background())
	ctx = proxy.ContextWithTarget(ctx, muxCoolDestination)
	pipe := ray.NewRay(ctx)
	go p.Process(ctx, pipe, dialer)
	c := &Client{
		sessions:       make(map[uint16]*session, 256),
		inboundRay:     pipe,
		ctx:            ctx,
		cancel:         cancel,
		manager:        m,
		count:          0,
		session2Remove: make(chan uint16, 16),
	}
	go c.fetchOutput()
	go c.monitor()
	return c, nil
}

func (m *Client) remove(id uint16) {
	select {
	case m.session2Remove <- id:
	default:
		// Probably not gonna happen.
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

func (m *Client) monitor() {
	for {
		select {
		case <-m.ctx.Done():
			m.cleanup()
			return
		case id := <-m.session2Remove:
			m.access.Lock()
			delete(m.sessions, id)
			m.access.Unlock()
		}
	}
}

func (m *Client) cleanup() {
	m.access.Lock()
	defer m.access.Unlock()

	for _, s := range m.sessions {
		s.closeUplink()
		s.closeDownlink()
		s.output.CloseError()
	}
}

func fetchInput(ctx context.Context, s *session, output buf.Writer) {
	dest, _ := proxy.TargetFromContext(ctx)
	writer := &Writer{
		dest:   dest,
		id:     s.id,
		writer: output,
	}
	defer writer.Close()
	defer s.closeUplink()

	log.Info("Proxyman|Mux|Client: Dispatching request to ", dest)
	data, _ := s.input.ReadTimeout(time.Millisecond * 500)
	if data != nil {
		if err := writer.Write(data); err != nil {
			log.Info("Proxyman|Mux|Client: Failed to write first payload: ", err)
			return
		}
	}
	if err := buf.PipeUntilEOF(signal.BackgroundTimer(), s.input, writer); err != nil {
		log.Info("Proxyman|Mux|Client: Failed to fetch all input: ", err)
	}
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
	go fetchInput(ctx, s, m.inboundRay.InboundInput())
	return true
}

func drain(reader *Reader) error {
	for {
		data, more, err := reader.Read()
		if err != nil {
			return err
		}
		data.Release()
		if !more {
			return nil
		}
	}
}

func pipe(reader *Reader, writer buf.Writer) error {
	for {
		data, more, err := reader.Read()
		if err != nil {
			return err
		}
		if err := writer.Write(data); err != nil {
			return err
		}
		if !more {
			return nil
		}
	}
}

func (m *Client) fetchOutput() {
	defer m.cancel()

	reader := NewReader(m.inboundRay.InboundOutput())
	for {
		meta, err := reader.ReadMetadata()
		if err != nil {
			log.Info("Proxyman|Mux|Client: Failed to read metadata: ", err)
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

		if found {
			err = pipe(reader, s.output)
		} else {
			err = drain(reader)
		}

		if err != nil {
			log.Info("Proxyman|Mux|Client: Failed to read data: ", err)
			break
		}
	}
}

type Server struct {
	dispatcher dispatcher.Interface
}

func NewServer(ctx context.Context) *Server {
	s := &Server{}
	space := app.SpaceFromContext(ctx)
	space.OnInitialize(func() error {
		d := dispatcher.FromSpace(space)
		if d == nil {
			return errors.New("Proxyman|Mux: No dispatcher in space.")
		}
		s.dispatcher = d
		return nil
	})
	return s
}

func (s *Server) Dispatch(ctx context.Context, dest net.Destination) (ray.InboundRay, error) {
	if dest != muxCoolDestination {
		return s.dispatcher.Dispatch(ctx, dest)
	}

	ray := ray.NewRay(ctx)
	worker := &ServerWorker{
		dispatcher:  s.dispatcher,
		outboundRay: ray,
		sessions:    make(map[uint16]*session),
	}
	go worker.run(ctx)
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

func handle(ctx context.Context, s *session, output buf.Writer) {
	writer := NewResponseWriter(s.id, output)
	if err := buf.PipeUntilEOF(signal.BackgroundTimer(), s.input, writer); err != nil {
		log.Info("Proxyman|Mux|ServerWorker: Session ", s.id, " ends: ", err)
	}
	writer.Close()
	s.closeDownlink()
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
			log.Info("Proxyman|Mux|Server: Received request for ", meta.Target)
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
			go handle(ctx, s, w.outboundRay.OutboundOutput())
		}

		if !meta.Option.Has(OptionData) {
			continue
		}

		if s != nil {
			err = pipe(reader, s.output)
		} else {
			err = drain(reader)
		}

		if err != nil {
			log.Info("Proxyman|Mux|ServerWorker: Failed to read data: ", err)
			break
		}
	}
}
