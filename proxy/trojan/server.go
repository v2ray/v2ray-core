// +build !confonly

package trojan

import (
	"context"
	"io"
	"time"

	"v2ray.com/core"
	"v2ray.com/core/common"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/errors"
	"v2ray.com/core/common/log"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/protocol"
	udp_proto "v2ray.com/core/common/protocol/udp"
	"v2ray.com/core/common/retry"
	"v2ray.com/core/common/session"
	"v2ray.com/core/common/signal"
	"v2ray.com/core/common/task"
	"v2ray.com/core/features/policy"
	"v2ray.com/core/features/routing"
	"v2ray.com/core/transport/internet"
	"v2ray.com/core/transport/internet/udp"
)

func init() {
	common.Must(common.RegisterConfig((*ServerConfig)(nil), func(ctx context.Context, config interface{}) (interface{}, error) { // nolint: lll
		return NewServer(ctx, config.(*ServerConfig))
	}))
}

// Server is an inbound connection handler that handles messages in trojan protocol.
type Server struct {
	validator     *Validator
	policyManager policy.Manager
	config        *ServerConfig
}

// NewServer creates a new trojan inbound handler.
func NewServer(ctx context.Context, config *ServerConfig) (*Server, error) {
	validator := new(Validator)
	for _, user := range config.Users {
		u, err := user.ToMemoryUser()
		if err != nil {
			return nil, newError("failed to get trojan user").Base(err).AtError()
		}

		if err := validator.Add(u); err != nil {
			return nil, newError("failed to add user").Base(err).AtError()
		}
	}

	v := core.MustFromContext(ctx)
	server := &Server{
		policyManager: v.GetFeature(policy.ManagerType()).(policy.Manager),
		validator:     validator,
		config:        config,
	}

	return server, nil
}

// Network implements proxy.Inbound.Network().
func (s *Server) Network() []net.Network {
	return []net.Network{net.Network_TCP}
}

// Process implements proxy.Inbound.Process().
func (s *Server) Process(ctx context.Context, network net.Network, conn internet.Connection, dispatcher routing.Dispatcher) error { // nolint: funlen,lll
	sessionPolicy := s.policyManager.ForLevel(0)
	if err := conn.SetReadDeadline(time.Now().Add(sessionPolicy.Timeouts.Handshake)); err != nil {
		return newError("unable to set read deadline").Base(err).AtWarning()
	}

	buffer := buf.New()
	defer buffer.Release()

	n, err := buffer.ReadFrom(conn)
	if err != nil {
		return newError("failed to read first request").Base(err)
	}

	bufferedReader := &buf.BufferedReader{
		Reader: buf.NewReader(conn),
		Buffer: buf.MultiBuffer{buffer},
	}

	var user *protocol.MemoryUser
	fallbackEnabled := s.config.Fallback != nil
	shouldFallback := false
	if n < 56 { // nolint: gomnd
		// invalid protocol
		log.Record(&log.AccessMessage{
			From:   conn.RemoteAddr(),
			To:     "",
			Status: log.AccessRejected,
			Reason: newError("not trojan protocol"),
		})

		shouldFallback = true
	} else {
		user = s.validator.Get(hexString(buffer.BytesTo(56))) // nolint: gomnd
		if user == nil {
			// invalid user, let's fallback
			log.Record(&log.AccessMessage{
				From:   conn.RemoteAddr(),
				To:     "",
				Status: log.AccessRejected,
				Reason: newError("not a valid user"),
			})

			shouldFallback = true
		}
	}

	if fallbackEnabled && shouldFallback {
		return s.fallback(ctx, sessionPolicy, bufferedReader, buf.NewWriter(conn))
	} else if shouldFallback {
		return newError("invalid protocol or invalid user")
	}

	clientReader := &ConnReader{Reader: bufferedReader}
	if err := clientReader.ParseHeader(); err != nil {
		log.Record(&log.AccessMessage{
			From:   conn.RemoteAddr(),
			To:     "",
			Status: log.AccessRejected,
			Reason: err,
		})
		return newError("failed to create request from: ", conn.RemoteAddr()).Base(err)
	}

	destination := clientReader.Target
	if err := conn.SetReadDeadline(time.Time{}); err != nil {
		return newError("unable to set read deadline").Base(err).AtWarning()
	}

	inbound := session.InboundFromContext(ctx)
	if inbound == nil {
		panic("no inbound metadata")
	}
	inbound.User = user
	sessionPolicy = s.policyManager.ForLevel(user.Level)

	if destination.Network == net.Network_UDP { // handle udp request
		return s.handleUDPPayload(ctx, &PacketReader{Reader: clientReader}, &PacketWriter{Writer: conn}, dispatcher)
	}

	// handle tcp request

	log.ContextWithAccessMessage(ctx, &log.AccessMessage{
		From:   conn.RemoteAddr(),
		To:     destination,
		Status: log.AccessAccepted,
		Reason: "",
		Email:  user.Email,
	})

	newError("received request for ", destination).WriteToLog(session.ExportIDToError(ctx))
	return s.handleConnection(ctx, sessionPolicy, destination, clientReader, buf.NewWriter(conn), dispatcher)
}

func (s *Server) handleUDPPayload(ctx context.Context, clientReader *PacketReader, clientWriter *PacketWriter, dispatcher routing.Dispatcher) error { // nolint: lll
	udpServer := udp.NewDispatcher(dispatcher, func(ctx context.Context, packet *udp_proto.Packet) {
		common.Must(clientWriter.WriteMultiBufferWithMetadata(buf.MultiBuffer{packet.Payload}, packet.Source))
	})

	inbound := session.InboundFromContext(ctx)
	user := inbound.User

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			p, err := clientReader.ReadMultiBufferWithMetadata()
			if err != nil {
				if errors.Cause(err) != io.EOF {
					return newError("unexpected EOF").Base(err)
				}
				return nil
			}

			log.ContextWithAccessMessage(ctx, &log.AccessMessage{
				From:   inbound.Source,
				To:     p.Target,
				Status: log.AccessAccepted,
				Reason: "",
				Email:  user.Email,
			})
			newError("tunnelling request to ", p.Target).WriteToLog(session.ExportIDToError(ctx))

			for _, b := range p.Buffer {
				udpServer.Dispatch(ctx, p.Target, b)
			}
		}
	}
}

func (s *Server) handleConnection(ctx context.Context, sessionPolicy policy.Session,
	destination net.Destination,
	clientReader buf.Reader,
	clientWriter buf.Writer, dispatcher routing.Dispatcher) error {
	ctx, cancel := context.WithCancel(ctx)
	timer := signal.CancelAfterInactivity(ctx, cancel, sessionPolicy.Timeouts.ConnectionIdle)
	ctx = policy.ContextWithBufferPolicy(ctx, sessionPolicy.Buffer)

	link, err := dispatcher.Dispatch(ctx, destination)
	if err != nil {
		return newError("failed to dispatch request to ", destination).Base(err)
	}

	requestDone := func() error {
		defer timer.SetTimeout(sessionPolicy.Timeouts.DownlinkOnly)

		if err := buf.Copy(clientReader, link.Writer, buf.UpdateActivity(timer)); err != nil {
			return newError("failed to transfer request").Base(err)
		}
		return nil
	}

	responseDone := func() error {
		defer timer.SetTimeout(sessionPolicy.Timeouts.UplinkOnly)

		if err := buf.Copy(link.Reader, clientWriter, buf.UpdateActivity(timer)); err != nil {
			return newError("failed to write response").Base(err)
		}
		return nil
	}

	var requestDonePost = task.OnSuccess(requestDone, task.Close(link.Writer))
	if err := task.Run(ctx, requestDonePost, responseDone); err != nil {
		common.Must(common.Interrupt(link.Reader))
		common.Must(common.Interrupt(link.Writer))
		return newError("connection ends").Base(err)
	}

	return nil
}

func (s *Server) fallback(ctx context.Context, sessionPolicy policy.Session, requestReader buf.Reader, responseWriter buf.Writer) error { // nolint: lll
	ctx, cancel := context.WithCancel(ctx)
	timer := signal.CancelAfterInactivity(ctx, cancel, sessionPolicy.Timeouts.ConnectionIdle)
	ctx = policy.ContextWithBufferPolicy(ctx, sessionPolicy.Buffer)

	var conn net.Conn
	var err error
	fb := s.config.Fallback
	if err := retry.ExponentialBackoff(5, 100).On(func() error { // nolint: gomnd
		var dialer net.Dialer
		conn, err = dialer.DialContext(ctx, fb.Type, fb.Dest)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return newError("failed to dial to " + fb.Dest).Base(err).AtWarning()
	}
	defer conn.Close()

	serverReader := buf.NewReader(conn)
	serverWriter := buf.NewWriter(conn)

	requestDone := func() error {
		defer timer.SetTimeout(sessionPolicy.Timeouts.DownlinkOnly)

		if err := buf.Copy(requestReader, serverWriter, buf.UpdateActivity(timer)); err != nil {
			return newError("failed to fallback request payload").Base(err).AtInfo()
		}
		return nil
	}

	responseDone := func() error {
		defer timer.SetTimeout(sessionPolicy.Timeouts.UplinkOnly)
		if err := buf.Copy(serverReader, responseWriter, buf.UpdateActivity(timer)); err != nil {
			return newError("failed to deliver response payload").Base(err).AtInfo()
		}
		return nil
	}

	if err := task.Run(ctx, task.OnSuccess(requestDone, task.Close(serverWriter)), task.OnSuccess(responseDone, task.Close(responseWriter))); err != nil { // nolint: lll
		common.Must(common.Interrupt(serverReader))
		common.Must(common.Interrupt(serverWriter))
		return newError("fallback ends").Base(err).AtInfo()
	}

	return nil
}
