package socks

import (
	"context"
	"io"
	"time"

	"v2ray.com/core"
	"v2ray.com/core/common"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/log"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/protocol"
	"v2ray.com/core/common/signal"
	"v2ray.com/core/proxy"
	"v2ray.com/core/transport/internet"
	"v2ray.com/core/transport/internet/udp"
	"v2ray.com/core/transport/pipe"
)

// Server is a SOCKS 5 proxy server
type Server struct {
	config *ServerConfig
	v      *core.Instance
}

// NewServer creates a new Server object.
func NewServer(ctx context.Context, config *ServerConfig) (*Server, error) {
	s := &Server{
		config: config,
		v:      core.MustFromContext(ctx),
	}
	return s, nil
}

func (s *Server) policy() core.Policy {
	config := s.config
	p := s.v.PolicyManager().ForLevel(config.UserLevel)
	if config.Timeout > 0 && config.UserLevel == 0 {
		p.Timeouts.ConnectionIdle = time.Duration(config.Timeout) * time.Second
	}
	return p
}

// Network implements proxy.Inbound.
func (s *Server) Network() net.NetworkList {
	list := net.NetworkList{
		Network: []net.Network{net.Network_TCP},
	}
	if s.config.UdpEnabled {
		list.Network = append(list.Network, net.Network_UDP)
	}
	return list
}

// Process implements proxy.Inbound.
func (s *Server) Process(ctx context.Context, network net.Network, conn internet.Connection, dispatcher core.Dispatcher) error {
	switch network {
	case net.Network_TCP:
		return s.processTCP(ctx, conn, dispatcher)
	case net.Network_UDP:
		return s.handleUDPPayload(ctx, conn, dispatcher)
	default:
		return newError("unknown network: ", network)
	}
}

func (s *Server) processTCP(ctx context.Context, conn internet.Connection, dispatcher core.Dispatcher) error {
	if err := conn.SetReadDeadline(time.Now().Add(s.policy().Timeouts.Handshake)); err != nil {
		newError("failed to set deadline").Base(err).WithContext(ctx).WriteToLog()
	}

	reader := &buf.BufferedReader{Reader: buf.NewReader(conn)}

	inboundDest, ok := proxy.InboundEntryPointFromContext(ctx)
	if !ok {
		return newError("inbound entry point not specified")
	}
	session := &ServerSession{
		config: s.config,
		port:   inboundDest.Port,
	}

	request, err := session.Handshake(reader, conn)
	if err != nil {
		if source, ok := proxy.SourceFromContext(ctx); ok {
			log.Record(&log.AccessMessage{
				From:   source,
				To:     "",
				Status: log.AccessRejected,
				Reason: err,
			})
		}
		return newError("failed to read request").Base(err)
	}

	if err := conn.SetReadDeadline(time.Time{}); err != nil {
		newError("failed to clear deadline").Base(err).WithContext(ctx).WriteToLog()
	}

	if request.Command == protocol.RequestCommandTCP {
		dest := request.Destination()
		newError("TCP Connect request to ", dest).WithContext(ctx).WriteToLog()
		if source, ok := proxy.SourceFromContext(ctx); ok {
			log.Record(&log.AccessMessage{
				From:   source,
				To:     dest,
				Status: log.AccessAccepted,
				Reason: "",
			})
		}

		return s.transport(ctx, reader, conn, dest, dispatcher)
	}

	if request.Command == protocol.RequestCommandUDP {
		return s.handleUDP(conn)
	}

	return nil
}

func (*Server) handleUDP(c io.Reader) error {
	// The TCP connection closes after this method returns. We need to wait until
	// the client closes it.
	return common.Error2(io.Copy(buf.DiscardBytes, c))
}

func (s *Server) transport(ctx context.Context, reader io.Reader, writer io.Writer, dest net.Destination, dispatcher core.Dispatcher) error {
	ctx, cancel := context.WithCancel(ctx)
	timer := signal.CancelAfterInactivity(ctx, cancel, s.policy().Timeouts.ConnectionIdle)

	plcy := s.policy()
	ctx = core.ContextWithBufferPolicy(ctx, plcy.Buffer)
	link, err := dispatcher.Dispatch(ctx, dest)
	if err != nil {
		return err
	}

	requestDone := func() error {
		defer timer.SetTimeout(plcy.Timeouts.DownlinkOnly)
		defer common.Close(link.Writer)

		v2reader := buf.NewReader(reader)
		if err := buf.Copy(v2reader, link.Writer, buf.UpdateActivity(timer)); err != nil {
			return newError("failed to transport all TCP request").Base(err)
		}

		return nil
	}

	responseDone := func() error {
		defer timer.SetTimeout(plcy.Timeouts.UplinkOnly)

		v2writer := buf.NewWriter(writer)
		if err := buf.Copy(link.Reader, v2writer, buf.UpdateActivity(timer)); err != nil {
			return newError("failed to transport all TCP response").Base(err)
		}

		return nil
	}

	if err := signal.ExecuteParallel(ctx, requestDone, responseDone); err != nil {
		pipe.CloseError(link.Reader)
		pipe.CloseError(link.Writer)
		return newError("connection ends").Base(err)
	}

	return nil
}

func (s *Server) handleUDPPayload(ctx context.Context, conn internet.Connection, dispatcher core.Dispatcher) error {
	udpServer := udp.NewDispatcher(dispatcher)

	if source, ok := proxy.SourceFromContext(ctx); ok {
		newError("client UDP connection from ", source).WithContext(ctx).WriteToLog()
	}

	reader := buf.NewReader(conn)
	for {
		mpayload, err := reader.ReadMultiBuffer()
		if err != nil {
			return err
		}

		for _, payload := range mpayload {
			request, err := DecodeUDPPacket(payload)

			if err != nil {
				newError("failed to parse UDP request").Base(err).WithContext(ctx).WriteToLog()
				payload.Release()
				continue
			}

			if payload.IsEmpty() {
				payload.Release()
				continue
			}

			newError("send packet to ", request.Destination(), " with ", payload.Len(), " bytes").AtDebug().WithContext(ctx).WriteToLog()
			if source, ok := proxy.SourceFromContext(ctx); ok {
				log.Record(&log.AccessMessage{
					From:   source,
					To:     request.Destination(),
					Status: log.AccessAccepted,
					Reason: "",
				})
			}

			udpServer.Dispatch(ctx, request.Destination(), payload, func(payload *buf.Buffer) {
				newError("writing back UDP response with ", payload.Len(), " bytes").AtDebug().WithContext(ctx).WriteToLog()

				udpMessage, err := EncodeUDPPacket(request, payload.Bytes())
				payload.Release()

				defer udpMessage.Release()
				if err != nil {
					newError("failed to write UDP response").AtWarning().Base(err).WithContext(ctx).WriteToLog()
				}

				conn.Write(udpMessage.Bytes())
			})
		}
	}
}

func init() {
	common.Must(common.RegisterConfig((*ServerConfig)(nil), func(ctx context.Context, config interface{}) (interface{}, error) {
		return NewServer(ctx, config.(*ServerConfig))
	}))
}
