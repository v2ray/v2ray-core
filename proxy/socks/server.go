package socks

import (
	"context"
	"io"
	"runtime"
	"time"

	"v2ray.com/core/app"
	"v2ray.com/core/app/dispatcher"
	"v2ray.com/core/common"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/bufio"
	"v2ray.com/core/common/errors"
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
	packetDispatcher dispatcher.Interface
	config           *ServerConfig
	udpServer        *udp.Dispatcher
}

// NewServer creates a new Server object.
func NewServer(ctx context.Context, config *ServerConfig) (*Server, error) {
	space := app.SpaceFromContext(ctx)
	if space == nil {
		return nil, errors.New("Socks|Server: No space in context.")
	}
	s := &Server{
		config: config,
	}
	space.OnInitialize(func() error {
		s.packetDispatcher = dispatcher.FromSpace(space)
		if s.packetDispatcher == nil {
			return errors.New("Socks|Server: Dispatcher is not found in the space.")
		}
		s.udpServer = udp.NewDispatcher(s.packetDispatcher)
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

func (s *Server) Process(ctx context.Context, network net.Network, conn internet.Connection) error {
	conn.SetReusable(false)

	switch network {
	case net.Network_TCP:
		return s.processTCP(ctx, conn)
	case net.Network_UDP:
		return s.handleUDPPayload(ctx, conn)
	default:
		return errors.New("Socks|Server: Unknown network: ", network)
	}
}

func (s *Server) processTCP(ctx context.Context, conn internet.Connection) error {
	conn.SetReadDeadline(time.Now().Add(time.Second * 8))
	reader := bufio.NewReader(conn)

	inboundDest := proxy.InboundDestinationFromContext(ctx)
	session := &ServerSession{
		config: s.config,
		port:   inboundDest.Port,
	}

	source := proxy.SourceFromContext(ctx)
	request, err := session.Handshake(reader, conn)
	if err != nil {
		log.Access(source, "", log.AccessRejected, err)
		log.Info("Socks|Server: Failed to read request: ", err)
		return err
	}
	conn.SetReadDeadline(time.Time{})

	if request.Command == protocol.RequestCommandTCP {
		dest := request.Destination()
		log.Info("Socks|Server: TCP Connect request to ", dest)
		log.Access(source, dest, log.AccessAccepted, "")

		ctx = proxy.ContextWithDestination(ctx, dest)
		return s.transport(ctx, reader, conn)
	}

	if request.Command == protocol.RequestCommandUDP {
		return s.handleUDP()
	}

	return nil
}

func (*Server) handleUDP() error {
	// The TCP connection closes after v method returns. We need to wait until
	// the client closes it.
	// TODO: get notified from UDP part
	<-time.After(5 * time.Minute)

	return nil
}

func (v *Server) transport(ctx context.Context, reader io.Reader, writer io.Writer) error {
	ctx, cancel := context.WithCancel(ctx)
	timeout := time.Second * time.Duration(v.config.Timeout)
	if timeout == 0 {
		timeout = time.Minute * 2
	}
	timer := signal.CancelAfterInactivity(ctx, cancel, timeout)

	ray := v.packetDispatcher.DispatchToOutbound(ctx)
	input := ray.InboundInput()
	output := ray.InboundOutput()

	requestDone := signal.ExecuteAsync(func() error {
		defer input.Close()

		v2reader := buf.NewReader(reader)
		if err := buf.PipeUntilEOF(timer, v2reader, input); err != nil {
			log.Info("Socks|Server: Failed to transport all TCP request: ", err)
			return err
		}
		return nil
	})

	responseDone := signal.ExecuteAsync(func() error {
		v2writer := buf.NewWriter(writer)
		if err := buf.PipeUntilEOF(timer, output, v2writer); err != nil {
			log.Info("Socks|Server: Failed to transport all TCP response: ", err)
			return err
		}
		return nil
	})

	if err := signal.ErrorOrFinish2(ctx, requestDone, responseDone); err != nil {
		log.Info("Socks|Server: Connection ends with ", err)
		input.CloseError()
		output.CloseError()
		return err
	}

	runtime.KeepAlive(timer)

	return nil
}

func (v *Server) handleUDPPayload(ctx context.Context, conn internet.Connection) error {
	source := proxy.SourceFromContext(ctx)
	log.Info("Socks|Server: Client UDP connection from ", source)

	reader := buf.NewReader(conn)
	for {
		payload, err := reader.Read()
		if err != nil {
			return err
		}
		request, data, err := DecodeUDPPacket(payload.Bytes())

		if err != nil {
			log.Info("Socks|Server: Failed to parse UDP request: ", err)
			continue
		}

		if len(data) == 0 {
			continue
		}

		log.Info("Socks: Send packet to ", request.Destination(), " with ", len(data), " bytes")
		log.Access(source, request.Destination, log.AccessAccepted, "")

		dataBuf := buf.NewSmall()
		dataBuf.Append(data)
		v.udpServer.Dispatch(ctx, request.Destination(), dataBuf, func(payload *buf.Buffer) {
			defer payload.Release()

			log.Info("Socks|Server: Writing back UDP response with ", payload.Len(), " bytes")

			udpMessage := EncodeUDPPacket(request, payload.Bytes())
			defer udpMessage.Release()

			conn.Write(udpMessage.Bytes())
		})
	}
}

func init() {
	common.Must(common.RegisterConfig((*ServerConfig)(nil), func(ctx context.Context, config interface{}) (interface{}, error) {
		return NewServer(ctx, config.(*ServerConfig))
	}))
}
