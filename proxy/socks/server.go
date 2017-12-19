package socks

import (
	"context"
	"io"
	"time"

	"v2ray.com/core/app"
	"v2ray.com/core/app/dispatcher"
	"v2ray.com/core/app/policy"
	"v2ray.com/core/common"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/log"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/protocol"
	"v2ray.com/core/common/signal"
	"v2ray.com/core/proxy"
	"v2ray.com/core/transport/internet"
	"v2ray.com/core/transport/internet/udp"
)

// Server is a SOCKS 5 proxy server
type Server struct {
	config *ServerConfig
	policy policy.Policy
}

// NewServer creates a new Server object.
func NewServer(ctx context.Context, config *ServerConfig) (*Server, error) {
	space := app.SpaceFromContext(ctx)
	if space == nil {
		return nil, newError("no space in context").AtWarning()
	}
	s := &Server{
		config: config,
	}
	space.On(app.SpaceInitializing, func(interface{}) error {
		pm := policy.FromSpace(space)
		if pm == nil {
			return newError("Policy not found in space.")
		}
		s.policy = pm.GetPolicy(config.UserLevel)
		if config.Timeout > 0 && config.UserLevel == 0 {
			s.policy.Timeout.ConnectionIdle.Value = config.Timeout
		}
		return nil
	})
	return s, nil
}

func (s *Server) Network() net.NetworkList {
	list := net.NetworkList{
		Network: []net.Network{net.Network_TCP},
	}
	if s.config.UdpEnabled {
		list.Network = append(list.Network, net.Network_UDP)
	}
	return list
}

func (s *Server) Process(ctx context.Context, network net.Network, conn internet.Connection, dispatcher dispatcher.Interface) error {
	switch network {
	case net.Network_TCP:
		return s.processTCP(ctx, conn, dispatcher)
	case net.Network_UDP:
		return s.handleUDPPayload(ctx, conn, dispatcher)
	default:
		return newError("unknown network: ", network)
	}
}

func (s *Server) processTCP(ctx context.Context, conn internet.Connection, dispatcher dispatcher.Interface) error {
	conn.SetReadDeadline(time.Now().Add(s.policy.Timeout.Handshake.Duration()))
	reader := buf.NewBufferedReader(buf.NewReader(conn))

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
	conn.SetReadDeadline(time.Time{})

	if request.Command == protocol.RequestCommandTCP {
		dest := request.Destination()
		newError("TCP Connect request to ", dest).WriteToLog()
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

func (*Server) handleUDP(c net.Conn) error {
	// The TCP connection closes after this method returns. We need to wait until
	// the client closes it.
	_, err := io.Copy(buf.DiscardBytes, c)
	return err
}

func (v *Server) transport(ctx context.Context, reader io.Reader, writer io.Writer, dest net.Destination, dispatcher dispatcher.Interface) error {
	ctx, cancel := context.WithCancel(ctx)
	timer := signal.CancelAfterInactivity(ctx, cancel, v.policy.Timeout.ConnectionIdle.Duration())

	ray, err := dispatcher.Dispatch(ctx, dest)
	if err != nil {
		return err
	}

	input := ray.InboundInput()
	output := ray.InboundOutput()

	requestDone := signal.ExecuteAsync(func() error {
		defer input.Close()

		v2reader := buf.NewReader(reader)
		if err := buf.Copy(v2reader, input, buf.UpdateActivity(timer)); err != nil {
			return newError("failed to transport all TCP request").Base(err)
		}
		timer.SetTimeout(v.policy.Timeout.DownlinkOnly.Duration())
		return nil
	})

	responseDone := signal.ExecuteAsync(func() error {
		v2writer := buf.NewWriter(writer)
		if err := buf.Copy(output, v2writer, buf.UpdateActivity(timer)); err != nil {
			return newError("failed to transport all TCP response").Base(err)
		}
		timer.SetTimeout(v.policy.Timeout.UplinkOnly.Duration())
		return nil
	})

	if err := signal.ErrorOrFinish2(ctx, requestDone, responseDone); err != nil {
		input.CloseError()
		output.CloseError()
		return newError("connection ends").Base(err)
	}

	return nil
}

func (v *Server) handleUDPPayload(ctx context.Context, conn internet.Connection, dispatcher dispatcher.Interface) error {
	udpServer := udp.NewDispatcher(dispatcher)

	if source, ok := proxy.SourceFromContext(ctx); ok {
		newError("client UDP connection from ", source).WriteToLog()
	}

	reader := buf.NewReader(conn)
	for {
		mpayload, err := reader.ReadMultiBuffer()
		if err != nil {
			return err
		}

		for _, payload := range mpayload {
			request, data, err := DecodeUDPPacket(payload.Bytes())

			if err != nil {
				newError("failed to parse UDP request").Base(err).WriteToLog()
				continue
			}

			if len(data) == 0 {
				continue
			}

			newError("send packet to ", request.Destination(), " with ", len(data), " bytes").AtDebug().WriteToLog()
			if source, ok := proxy.SourceFromContext(ctx); ok {
				log.Record(&log.AccessMessage{
					From:   source,
					To:     request.Destination,
					Status: log.AccessAccepted,
					Reason: "",
				})
			}

			dataBuf := buf.New()
			dataBuf.Append(data)
			udpServer.Dispatch(ctx, request.Destination(), dataBuf, func(payload *buf.Buffer) {
				defer payload.Release()

				newError("writing back UDP response with ", payload.Len(), " bytes").AtDebug().WriteToLog()

				udpMessage, err := EncodeUDPPacket(request, payload.Bytes())
				defer udpMessage.Release()
				if err != nil {
					newError("failed to write UDP response").AtWarning().Base(err).WriteToLog()
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
