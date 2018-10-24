package mux

import (
	"context"
	"io"
	"sync"
	"time"

	"v2ray.com/core/common"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/errors"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/protocol"
	"v2ray.com/core/common/session"
	"v2ray.com/core/common/signal/done"
	"v2ray.com/core/common/vio"
	"v2ray.com/core/proxy"
	"v2ray.com/core/transport/internet"
	"v2ray.com/core/transport/pipe"
)

type ClientManager struct {
	access      sync.Mutex
	clients     []*Client
	proxy       proxy.Outbound
	dialer      internet.Dialer
	concurrency uint32
}

func NewClientManager(p proxy.Outbound, d internet.Dialer, c uint32) *ClientManager {
	return &ClientManager{
		proxy:       p,
		dialer:      d,
		concurrency: c,
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
func NewClient(pctx context.Context, p proxy.Outbound, dialer internet.Dialer, m *ClientManager) (*Client, error) {
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
		concurrency: m.concurrency,
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

func (m *Client) handleStatueKeepAlive(meta *FrameMetadata, reader *buf.BufferedReader) error {
	if meta.Option.Has(OptionData) {
		return buf.Copy(NewStreamReader(reader), buf.Discard)
	}
	return nil
}

func (m *Client) handleStatusNew(meta *FrameMetadata, reader *buf.BufferedReader) error {
	if meta.Option.Has(OptionData) {
		return buf.Copy(NewStreamReader(reader), buf.Discard)
	}
	return nil
}

func (m *Client) handleStatusKeep(meta *FrameMetadata, reader *buf.BufferedReader) error {
	if !meta.Option.Has(OptionData) {
		return nil
	}

	s, found := m.sessionManager.Get(meta.SessionID)
	if !found {
		return buf.Copy(NewStreamReader(reader), buf.Discard)
	}

	rr := s.NewReader(reader)
	err := buf.Copy(rr, s.output)
	if err != nil && buf.IsWriteError(err) {
		newError("failed to write to downstream. closing session ", s.ID).Base(err).WriteToLog()

		drainErr := buf.Copy(rr, buf.Discard)
		pipe.CloseError(s.input)
		s.Close()
		return drainErr
	}

	return err
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
		return buf.Copy(NewStreamReader(reader), buf.Discard)
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
		err := meta.Unmarshal(reader)
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
