package mux

//go:generate errorgen

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
	"v2ray.com/core/common/session"
	"v2ray.com/core/common/signal/done"
	"v2ray.com/core/common/vio"
	"v2ray.com/core/features/routing"
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

func (m *ClientManager) Dispatch(ctx context.Context, link *vio.Link) error {
	m.access.Lock()
	defer m.access.Unlock()

	for _, client := range m.clients {
		if client.Dispatch(ctx, link) {
			return nil
		}
	}

	client, err := NewClient(ctx, m.proxy, m.dialer, m)
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
	link           vio.Link
	done           *done.Instance
	manager        *ClientManager
	concurrency    uint32
}

var muxCoolAddress = net.DomainAddress("v1.mux.cool")
var muxCoolPort = net.Port(9527)

// NewClient creates a new mux.Client.
func NewClient(pctx context.Context, p proxy.Outbound, dialer proxy.Dialer, m *ClientManager) (*Client, error) {
	ctx := session.ContextWithOutbound(context.Background(), &session.Outbound{
		Target: net.TCPDestination(muxCoolAddress, muxCoolPort),
	})
	ctx, cancel := context.WithCancel(ctx)

	opts := pipe.OptionsFromContext(pctx)
	uplinkReader, upLinkWriter := pipe.New(opts...)
	downlinkReader, downlinkWriter := pipe.New(opts...)

	c := &Client{
		sessionManager: NewSessionManager(),
		link: vio.Link{
			Reader: downlinkReader,
			Writer: upLinkWriter,
		},
		done:        done.New(),
		manager:     m,
		concurrency: m.config.Concurrency,
	}

	go func() {
		if err := p.Process(ctx, &vio.Link{Reader: uplinkReader, Writer: downlinkWriter}, dialer); err != nil {
			errors.New("failed to handler mux client connection").Base(err).WriteToLog()
		}
		common.Must(c.done.Close())
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
			common.Close(m.link.Writer)    // nolint: errcheck
			pipe.CloseError(m.link.Reader) // nolint: errcheck
			return
		case <-timer.C:
			size := m.sessionManager.Size()
			if size == 0 && m.sessionManager.CloseIfNoSession() {
				common.Must(m.done.Close())
			}
		}
	}
}

func writeFirstPayload(reader buf.Reader, writer *Writer) error {
	err := buf.CopyOnceTimeout(reader, writer, time.Millisecond*100)
	if err == buf.ErrNotTimeoutReader || err == buf.ErrReadTimeout {
		return writer.WriteMultiBuffer(buf.MultiBuffer{})
	}

	if err != nil {
		return err
	}

	return nil
}

func fetchInput(ctx context.Context, s *Session, output buf.Writer) {
	dest := session.OutboundFromContext(ctx).Target
	transferType := protocol.TransferTypeStream
	if dest.Network == net.Network_UDP {
		transferType = protocol.TransferTypePacket
	}
	s.transferType = transferType
	writer := NewWriter(s.ID, dest, output, transferType)
	defer s.Close()      // nolint: errcheck
	defer writer.Close() // nolint: errcheck

	newError("dispatching request to ", dest).WriteToLog(session.ExportIDToError(ctx))
	if err := writeFirstPayload(s.input, writer); err != nil {
		newError("failed to write first payload").Base(err).WriteToLog(session.ExportIDToError(ctx))
		writer.hasError = true
		pipe.CloseError(s.input)
		return
	}

	if err := buf.Copy(s.input, writer); err != nil {
		newError("failed to fetch all input").Base(err).WriteToLog(session.ExportIDToError(ctx))
		writer.hasError = true
		pipe.CloseError(s.input)
		return
	}
}

func (m *Client) Dispatch(ctx context.Context, link *vio.Link) bool {
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

func drain(reader buf.Reader) error {
	return buf.Copy(reader, buf.Discard)
}

func (m *Client) handleStatueKeepAlive(meta *FrameMetadata, reader *buf.BufferedReader) error {
	if meta.Option.Has(OptionData) {
		return drain(NewStreamReader(reader))
	}
	return nil
}

func (m *Client) handleStatusNew(meta *FrameMetadata, reader *buf.BufferedReader) error {
	if meta.Option.Has(OptionData) {
		return drain(NewStreamReader(reader))
	}
	return nil
}

func (m *Client) handleStatusKeep(meta *FrameMetadata, reader *buf.BufferedReader) error {
	if !meta.Option.Has(OptionData) {
		return nil
	}

	if s, found := m.sessionManager.Get(meta.SessionID); found {
		rr := s.NewReader(reader)
		if err := buf.Copy(rr, s.output); err != nil {
			drain(rr)
			pipe.CloseError(s.input)
			return s.Close()
		}
		return nil
	}
	return drain(NewStreamReader(reader))
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
		return drain(NewStreamReader(reader))
	}
	return nil
}

func (m *Client) fetchOutput() {
	defer func() {
		common.Must(m.done.Close())
	}()

	reader := &buf.BufferedReader{Reader: m.link.Reader}

	var meta FrameMetadata
	for {
		err := meta.ReadFrom(reader)
		if err != nil {
			if errors.Cause(err) != io.EOF {
				newError("failed to read metadata").Base(err).WriteToLog()
			}
			break
		}

		switch meta.SessionStatus {
		case SessionStatusKeepAlive:
			err = m.handleStatueKeepAlive(&meta, reader)
		case SessionStatusEnd:
			err = m.handleStatusEnd(&meta, reader)
		case SessionStatusNew:
			err = m.handleStatusNew(&meta, reader)
		case SessionStatusKeep:
			err = m.handleStatusKeep(&meta, reader)
		default:
			status := meta.SessionStatus
			newError("unknown status: ", status).AtError().WriteToLog()
			return
		}

		if err != nil {
			newError("failed to process data").Base(err).WriteToLog()
			return
		}
	}
}

type Server struct {
	dispatcher routing.Dispatcher
}

// NewServer creates a new mux.Server.
func NewServer(ctx context.Context) *Server {
	s := &Server{
		dispatcher: core.MustFromContext(ctx).Dispatcher(),
	}
	return s
}

func (s *Server) Type() interface{} {
	return s.dispatcher.Type()
}

func (s *Server) Dispatch(ctx context.Context, dest net.Destination) (*vio.Link, error) {
	if dest.Address != muxCoolAddress {
		return s.dispatcher.Dispatch(ctx, dest)
	}

	opts := pipe.OptionsFromContext(ctx)
	uplinkReader, uplinkWriter := pipe.New(opts...)
	downlinkReader, downlinkWriter := pipe.New(opts...)

	worker := &ServerWorker{
		dispatcher: s.dispatcher,
		link: &vio.Link{
			Reader: uplinkReader,
			Writer: downlinkWriter,
		},
		sessionManager: NewSessionManager(),
	}
	go worker.run(ctx)
	return &vio.Link{Reader: downlinkReader, Writer: uplinkWriter}, nil
}

// Start implements common.Runnable.
func (s *Server) Start() error {
	return nil
}

// Close implements common.Closable.
func (s *Server) Close() error {
	return nil
}

type ServerWorker struct {
	dispatcher     routing.Dispatcher
	link           *vio.Link
	sessionManager *SessionManager
}

func handle(ctx context.Context, s *Session, output buf.Writer) {
	writer := NewResponseWriter(s.ID, output, s.transferType)
	if err := buf.Copy(s.input, writer); err != nil {
		newError("session ", s.ID, " ends.").Base(err).WriteToLog(session.ExportIDToError(ctx))
		writer.hasError = true
	}

	writer.Close()
	s.Close()
}

func (w *ServerWorker) handleStatusKeepAlive(meta *FrameMetadata, reader *buf.BufferedReader) error {
	if meta.Option.Has(OptionData) {
		return drain(NewStreamReader(reader))
	}
	return nil
}

func (w *ServerWorker) handleStatusNew(ctx context.Context, meta *FrameMetadata, reader *buf.BufferedReader) error {
	newError("received request for ", meta.Target).WriteToLog(session.ExportIDToError(ctx))
	{
		msg := &log.AccessMessage{
			To:     meta.Target,
			Status: log.AccessAccepted,
			Reason: "",
		}
		if inbound := session.InboundFromContext(ctx); inbound != nil && inbound.Source.IsValid() {
			msg.From = inbound.Source
		}
		log.Record(msg)
	}
	link, err := w.dispatcher.Dispatch(ctx, meta.Target)
	if err != nil {
		if meta.Option.Has(OptionData) {
			drain(NewStreamReader(reader))
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
	if !meta.Option.Has(OptionData) {
		return nil
	}

	rr := s.NewReader(reader)
	if err := buf.Copy(rr, s.output); err != nil {
		drain(rr)
		pipe.CloseError(s.input)
		return s.Close()
	}
	return nil
}

func (w *ServerWorker) handleStatusKeep(meta *FrameMetadata, reader *buf.BufferedReader) error {
	if !meta.Option.Has(OptionData) {
		return nil
	}
	if s, found := w.sessionManager.Get(meta.SessionID); found {
		rr := s.NewReader(reader)
		if err := buf.Copy(rr, s.output); err != nil {
			drain(rr)
			pipe.CloseError(s.input)
			return s.Close()
		}
		return nil
	}
	return drain(NewStreamReader(reader))
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
		return drain(NewStreamReader(reader))
	}
	return nil
}

func (w *ServerWorker) handleFrame(ctx context.Context, reader *buf.BufferedReader) error {
	var meta FrameMetadata
	err := meta.ReadFrom(reader)
	if err != nil {
		return newError("failed to read metadata").Base(err)
	}

	switch meta.SessionStatus {
	case SessionStatusKeepAlive:
		err = w.handleStatusKeepAlive(&meta, reader)
	case SessionStatusEnd:
		err = w.handleStatusEnd(&meta, reader)
	case SessionStatusNew:
		err = w.handleStatusNew(ctx, &meta, reader)
	case SessionStatusKeep:
		err = w.handleStatusKeep(&meta, reader)
	default:
		status := meta.SessionStatus
		return newError("unknown status: ", status).AtError()
	}

	if err != nil {
		return newError("failed to process data").Base(err)
	}
	return nil
}

func (w *ServerWorker) run(ctx context.Context) {
	input := w.link.Reader
	reader := &buf.BufferedReader{Reader: input}

	defer w.sessionManager.Close() // nolint: errcheck

	for {
		select {
		case <-ctx.Done():
			return
		default:
			err := w.handleFrame(ctx, reader)
			if err != nil {
				if errors.Cause(err) != io.EOF {
					newError("unexpected EOF").Base(err).WriteToLog(session.ExportIDToError(ctx))
					pipe.CloseError(input)
				}
				return
			}
		}
	}
}
