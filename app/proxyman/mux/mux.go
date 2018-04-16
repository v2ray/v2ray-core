package mux

//go:generate go run $GOPATH/src/v2ray.com/core/common/errors/errorgen/main.go -pkg mux -path App,Proxyman,Mux

import (
	"context"
	"io"
	"sync"
	"time"

	"v2ray.com/core"
	"v2ray.com/core/app/proxyman"
	"v2ray.com/core/common"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/errors"
	"v2ray.com/core/common/log"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/protocol"
	"v2ray.com/core/common/signal"
	"v2ray.com/core/proxy"
	"v2ray.com/core/transport/pipe"
)

const (
	maxTotal = 128
)

type ClientManager struct {
	access  sync.Mutex
	clients []*Client
	proxy   proxy.Outbound
	dialer  proxy.Dialer
	config  *proxyman.MultiplexingConfig
}

func NewClientManager(p proxy.Outbound, d proxy.Dialer, c *proxyman.MultiplexingConfig) *ClientManager {
	return &ClientManager{
		proxy:  p,
		dialer: d,
		config: c,
	}
}

func (m *ClientManager) Dispatch(ctx context.Context, link *core.Link) error {
	m.access.Lock()
	defer m.access.Unlock()

	for _, client := range m.clients {
		if client.Dispatch(ctx, link) {
			return nil
		}
	}

	client, err := NewClient(m.proxy, m.dialer, m)
	if err != nil {
		return newError("failed to create client").Base(err)
	}
	m.clients = append(m.clients, client)
	client.Dispatch(ctx, link)
	return nil
}

func (m *ClientManager) onClientFinish() {
	m.access.Lock()
	defer m.access.Unlock()

	activeClients := make([]*Client, 0, len(m.clients))

	for _, client := range m.clients {
		if !client.Closed() {
			activeClients = append(activeClients, client)
		}
	}
	m.clients = activeClients
}

type Client struct {
	sessionManager *SessionManager
	link           core.Link
	done           *signal.Done
	manager        *ClientManager
	concurrency    uint32
}

var muxCoolAddress = net.DomainAddress("v1.mux.cool")
var muxCoolPort = net.Port(9527)

// NewClient creates a new mux.Client.
func NewClient(p proxy.Outbound, dialer proxy.Dialer, m *ClientManager) (*Client, error) {
	ctx := proxy.ContextWithTarget(context.Background(), net.TCPDestination(muxCoolAddress, muxCoolPort))
	ctx, cancel := context.WithCancel(ctx)
	uplinkReader, upLinkWriter := pipe.New()
	downlinkReader, downlinkWriter := pipe.New()

	c := &Client{
		sessionManager: NewSessionManager(),
		link: core.Link{
			Reader: downlinkReader,
			Writer: upLinkWriter,
		},
		done:        signal.NewDone(),
		manager:     m,
		concurrency: m.config.Concurrency,
	}

	go func() {
		if err := p.Process(ctx, &core.Link{Reader: uplinkReader, Writer: downlinkWriter}, dialer); err != nil {
			errors.New("failed to handler mux client connection").Base(err).WriteToLog()
		}
		c.done.Close()
		cancel()
	}()

	go c.fetchOutput()
	go c.monitor()
	return c, nil
}

// Closed returns true if this Client is closed.
func (m *Client) Closed() bool {
	return m.done.Done()
}

func (m *Client) monitor() {
	defer m.manager.onClientFinish()

	timer := time.NewTicker(time.Second * 16)
	defer timer.Stop()

	for {
		select {
		case <-m.done.Wait():
			m.sessionManager.Close()
			common.Close(m.link.Writer)
			pipe.CloseError(m.link.Reader)
			return
		case <-timer.C:
			size := m.sessionManager.Size()
			if size == 0 && m.sessionManager.CloseIfNoSession() {
				common.Must(m.done.Close())
				return
			}
		}
	}
}

func fetchInput(ctx context.Context, s *Session, output buf.Writer) {
	dest, _ := proxy.TargetFromContext(ctx)
	transferType := protocol.TransferTypeStream
	if dest.Network == net.Network_UDP {
		transferType = protocol.TransferTypePacket
	}
	s.transferType = transferType
	writer := NewWriter(s.ID, dest, output, transferType)
	defer s.Close()

	newError("dispatching request to ", dest).WithContext(ctx).WriteToLog()
	if err := buf.Copy(s.input, writer); err != nil {
		newError("failed to fetch all input").Base(err).WithContext(ctx).WriteToLog()
		writer.hasError = true
	}

	writer.Close()
}

func (m *Client) Dispatch(ctx context.Context, link *core.Link) bool {
	sm := m.sessionManager
	if sm.Size() >= int(m.concurrency) || sm.Count() >= maxTotal {
		return false
	}

	if m.done.Done() {
		return false
	}

	s := sm.Allocate()
	if s == nil {
		return false
	}
	s.input = link.Reader
	s.output = link.Writer
	go fetchInput(ctx, s, m.link.Writer)
	return true
}

func drain(reader *buf.BufferedReader) error {
	return buf.Copy(NewStreamReader(reader), buf.Discard)
}

func (m *Client) handleStatueKeepAlive(meta *FrameMetadata, reader *buf.BufferedReader) error {
	if meta.Option.Has(OptionData) {
		return drain(reader)
	}
	return nil
}

func (m *Client) handleStatusNew(meta *FrameMetadata, reader *buf.BufferedReader) error {
	if meta.Option.Has(OptionData) {
		return drain(reader)
	}
	return nil
}

func (m *Client) handleStatusKeep(meta *FrameMetadata, reader *buf.BufferedReader) error {
	if !meta.Option.Has(OptionData) {
		return nil
	}

	if s, found := m.sessionManager.Get(meta.SessionID); found {
		if err := buf.Copy(s.NewReader(reader), s.output); err != nil {
			drain(reader)
			pipe.CloseError(s.input)
			return s.Close()
		}
		return nil
	}
	return drain(reader)
}

func (m *Client) handleStatusEnd(meta *FrameMetadata, reader *buf.BufferedReader) error {
	if s, found := m.sessionManager.Get(meta.SessionID); found {
		if meta.Option.Has(OptionError) {
			pipe.CloseError(s.input)
			pipe.CloseError(s.output)
		}
		s.Close()
	}
	if meta.Option.Has(OptionData) {
		return drain(reader)
	}
	return nil
}

func (m *Client) fetchOutput() {
	defer m.done.Close()

	reader := buf.NewBufferedReader(m.link.Reader)

	for {
		meta, err := ReadMetadata(reader)
		if err != nil {
			if errors.Cause(err) != io.EOF {
				newError("failed to read metadata").Base(err).WriteToLog()
			}
			break
		}

		switch meta.SessionStatus {
		case SessionStatusKeepAlive:
			err = m.handleStatueKeepAlive(meta, reader)
		case SessionStatusEnd:
			err = m.handleStatusEnd(meta, reader)
		case SessionStatusNew:
			err = m.handleStatusNew(meta, reader)
		case SessionStatusKeep:
			err = m.handleStatusKeep(meta, reader)
		default:
			newError("unknown status: ", meta.SessionStatus).AtError().WriteToLog()
			return
		}

		if err != nil {
			newError("failed to process data").Base(err).WriteToLog()
			return
		}
	}
}

type Server struct {
	dispatcher core.Dispatcher
}

// NewServer creates a new mux.Server.
func NewServer(ctx context.Context) *Server {
	s := &Server{
		dispatcher: core.MustFromContext(ctx).Dispatcher(),
	}
	return s
}

func (s *Server) Dispatch(ctx context.Context, dest net.Destination) (*core.Link, error) {
	if dest.Address != muxCoolAddress {
		return s.dispatcher.Dispatch(ctx, dest)
	}

	uplinkReader, uplinkWriter := pipe.New()
	downlinkReader, downlinkWriter := pipe.New()

	worker := &ServerWorker{
		dispatcher: s.dispatcher,
		link: &core.Link{
			Reader: uplinkReader,
			Writer: downlinkWriter,
		},
		sessionManager: NewSessionManager(),
	}
	go worker.run(ctx)
	return &core.Link{Reader: downlinkReader, Writer: uplinkWriter}, nil
}

func (s *Server) Start() error {
	return nil
}

func (s *Server) Close() error {
	return nil
}

type ServerWorker struct {
	dispatcher     core.Dispatcher
	link           *core.Link
	sessionManager *SessionManager
}

func handle(ctx context.Context, s *Session, output buf.Writer) {
	writer := NewResponseWriter(s.ID, output, s.transferType)
	if err := buf.Copy(s.input, writer); err != nil {
		newError("session ", s.ID, " ends.").Base(err).WithContext(ctx).WriteToLog()
		writer.hasError = true
	}

	writer.Close()
	s.Close()
}

func (w *ServerWorker) handleStatusKeepAlive(meta *FrameMetadata, reader *buf.BufferedReader) error {
	if meta.Option.Has(OptionData) {
		return drain(reader)
	}
	return nil
}

func (w *ServerWorker) handleStatusNew(ctx context.Context, meta *FrameMetadata, reader *buf.BufferedReader) error {
	newError("received request for ", meta.Target).WithContext(ctx).WriteToLog()
	{
		msg := &log.AccessMessage{
			To:     meta.Target,
			Status: log.AccessAccepted,
			Reason: "",
		}
		if src, f := proxy.SourceFromContext(ctx); f {
			msg.From = src
		}
		log.Record(msg)
	}
	link, err := w.dispatcher.Dispatch(ctx, meta.Target)
	if err != nil {
		if meta.Option.Has(OptionData) {
			drain(reader)
		}
		return newError("failed to dispatch request.").Base(err)
	}
	s := &Session{
		input:        link.Reader,
		output:       link.Writer,
		parent:       w.sessionManager,
		ID:           meta.SessionID,
		transferType: protocol.TransferTypeStream,
	}
	if meta.Target.Network == net.Network_UDP {
		s.transferType = protocol.TransferTypePacket
	}
	w.sessionManager.Add(s)
	go handle(ctx, s, w.link.Writer)
	if meta.Option.Has(OptionData) {
		return buf.Copy(s.NewReader(reader), s.output, buf.IgnoreWriterError())
	}
	return nil
}

func (w *ServerWorker) handleStatusKeep(meta *FrameMetadata, reader *buf.BufferedReader) error {
	if !meta.Option.Has(OptionData) {
		return nil
	}
	if s, found := w.sessionManager.Get(meta.SessionID); found {
		if err := buf.Copy(s.NewReader(reader), s.output); err != nil {
			drain(reader)
			pipe.CloseError(s.input)
			return s.Close()
		}
		return nil
	}
	return drain(reader)
}

func (w *ServerWorker) handleStatusEnd(meta *FrameMetadata, reader *buf.BufferedReader) error {
	if s, found := w.sessionManager.Get(meta.SessionID); found {
		if meta.Option.Has(OptionError) {
			pipe.CloseError(s.input)
			pipe.CloseError(s.output)
		}
		s.Close()
	}
	if meta.Option.Has(OptionData) {
		return drain(reader)
	}
	return nil
}

func (w *ServerWorker) handleFrame(ctx context.Context, reader *buf.BufferedReader) error {
	meta, err := ReadMetadata(reader)
	if err != nil {
		return newError("failed to read metadata").Base(err)
	}

	switch meta.SessionStatus {
	case SessionStatusKeepAlive:
		err = w.handleStatusKeepAlive(meta, reader)
	case SessionStatusEnd:
		err = w.handleStatusEnd(meta, reader)
	case SessionStatusNew:
		err = w.handleStatusNew(ctx, meta, reader)
	case SessionStatusKeep:
		err = w.handleStatusKeep(meta, reader)
	default:
		return newError("unknown status: ", meta.SessionStatus).AtError()
	}

	if err != nil {
		return newError("failed to process data").Base(err)
	}
	return nil
}

func (w *ServerWorker) run(ctx context.Context) {
	input := w.link.Reader
	reader := buf.NewBufferedReader(input)

	defer w.sessionManager.Close()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			err := w.handleFrame(ctx, reader)
			if err != nil {
				if errors.Cause(err) != io.EOF {
					newError("unexpected EOF").Base(err).WithContext(ctx).WriteToLog()
					pipe.CloseError(input)
				}
				return
			}
		}
	}
}
