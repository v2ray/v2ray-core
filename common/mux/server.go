package mux

import (
	"context"
	"io"

	"v2ray.com/core"
	"v2ray.com/core/common"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/errors"
	"v2ray.com/core/common/log"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/protocol"
	"v2ray.com/core/common/session"
	"v2ray.com/core/features/routing"
	"v2ray.com/core/transport"
	"v2ray.com/core/transport/pipe"
)

type Server struct {
	dispatcher routing.Dispatcher
}

// NewServer creates a new mux.Server.
func NewServer(ctx context.Context) *Server {
	s := &Server{}
	core.RequireFeatures(ctx, func(d routing.Dispatcher) {
		s.dispatcher = d
	})
	return s
}

// Type implements common.HasType.
func (s *Server) Type() interface{} {
	return s.dispatcher.Type()
}

// Dispatch impliments routing.Dispatcher
func (s *Server) Dispatch(ctx context.Context, dest net.Destination) (*transport.Link, error) {
	if dest.Address != muxCoolAddress {
		return s.dispatcher.Dispatch(ctx, dest)
	}

	opts := pipe.OptionsFromContext(ctx)
	uplinkReader, uplinkWriter := pipe.New(opts...)
	downlinkReader, downlinkWriter := pipe.New(opts...)

	_, err := NewServerWorker(ctx, s.dispatcher, &transport.Link{
		Reader: uplinkReader,
		Writer: downlinkWriter,
	})
	if err != nil {
		return nil, err
	}

	return &transport.Link{Reader: downlinkReader, Writer: uplinkWriter}, nil
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
	link           *transport.Link
	sessionManager *SessionManager
}

func NewServerWorker(ctx context.Context, d routing.Dispatcher, link *transport.Link) (*ServerWorker, error) {
	worker := &ServerWorker{
		dispatcher:     d,
		link:           link,
		sessionManager: NewSessionManager(),
	}
	go worker.run(ctx)
	return worker, nil
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

func (w *ServerWorker) ActiveConnections() uint32 {
	return uint32(w.sessionManager.Size())
}

func (w *ServerWorker) Closed() bool {
	return w.sessionManager.Closed()
}

func (w *ServerWorker) handleStatusKeepAlive(meta *FrameMetadata, reader *buf.BufferedReader) error {
	if meta.Option.Has(OptionData) {
		return buf.Copy(NewStreamReader(reader), buf.Discard)
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
			msg.Email = inbound.User.Email
		}
		log.Record(msg)
	}
	link, err := w.dispatcher.Dispatch(ctx, meta.Target)
	if err != nil {
		if meta.Option.Has(OptionData) {
			buf.Copy(NewStreamReader(reader), buf.Discard)
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
		buf.Copy(rr, buf.Discard)
		common.Interrupt(s.input)
		return s.Close()
	}
	return nil
}

func (w *ServerWorker) handleStatusKeep(meta *FrameMetadata, reader *buf.BufferedReader) error {
	if !meta.Option.Has(OptionData) {
		return nil
	}

	s, found := w.sessionManager.Get(meta.SessionID)
	if !found {
		// Notify remote peer to close this session.
		closingWriter := NewResponseWriter(meta.SessionID, w.link.Writer, protocol.TransferTypeStream)
		closingWriter.Close()

		return buf.Copy(NewStreamReader(reader), buf.Discard)
	}

	rr := s.NewReader(reader)
	err := buf.Copy(rr, s.output)

	if err != nil && buf.IsWriteError(err) {
		newError("failed to write to downstream writer. closing session ", s.ID).Base(err).WriteToLog()

		// Notify remote peer to close this session.
		closingWriter := NewResponseWriter(meta.SessionID, w.link.Writer, protocol.TransferTypeStream)
		closingWriter.Close()

		drainErr := buf.Copy(rr, buf.Discard)
		common.Interrupt(s.input)
		s.Close()
		return drainErr
	}

	return err
}

func (w *ServerWorker) handleStatusEnd(meta *FrameMetadata, reader *buf.BufferedReader) error {
	if s, found := w.sessionManager.Get(meta.SessionID); found {
		if meta.Option.Has(OptionError) {
			common.Interrupt(s.input)
			common.Interrupt(s.output)
		}
		s.Close()
	}
	if meta.Option.Has(OptionData) {
		return buf.Copy(NewStreamReader(reader), buf.Discard)
	}
	return nil
}

func (w *ServerWorker) handleFrame(ctx context.Context, reader *buf.BufferedReader) error {
	var meta FrameMetadata
	err := meta.Unmarshal(reader)
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
					common.Interrupt(input)
				}
				return
			}
		}
	}
}
