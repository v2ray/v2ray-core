package shadowsocks

import (
	"context"

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

type Server struct {
	packetDispatcher dispatcher.Interface
	config           *ServerConfig
	user             *protocol.User
	account          *ShadowsocksAccount
	udpServer        *udp.Dispatcher
}

func NewServer(ctx context.Context, config *ServerConfig) (*Server, error) {
	space := app.SpaceFromContext(ctx)
	if space == nil {
		return nil, errors.New("Shadowsocks|Server: No space in context.")
	}
	if config.GetUser() == nil {
		return nil, protocol.ErrUserMissing
	}

	rawAccount, err := config.User.GetTypedAccount()
	if err != nil {
		return nil, errors.Base(err).Message("Shadowsocks|Server: Failed to get user account.")
	}
	account := rawAccount.(*ShadowsocksAccount)

	s := &Server{
		config:  config,
		user:    config.GetUser(),
		account: account,
	}

	space.OnInitialize(func() error {
		s.packetDispatcher = dispatcher.FromSpace(space)
		if s.packetDispatcher == nil {
			return errors.New("Shadowsocks|Server: Dispatcher is not found in space.")
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

func (s *Server) Process(ctx context.Context, network net.Network, conn internet.Connection) error {
	conn.SetReusable(false)

	switch network {
	case net.Network_TCP:
		return s.handleConnection(ctx, conn)
	case net.Network_UDP:
		return s.handlerUDPPayload(ctx, conn)
	default:
		return errors.New("Shadowsocks|Server: Unknown network: ", network)
	}
}

func (v *Server) handlerUDPPayload(ctx context.Context, conn internet.Connection) error {
	source := proxy.SourceFromContext(ctx)

	reader := buf.NewReader(conn)
	for {
		payload, err := reader.Read()
		if err != nil {
			break
		}

		request, data, err := DecodeUDPPacket(v.user, payload)
		if err != nil {
			log.Info("Shadowsocks|Server: Skipping invalid UDP packet from: ", source, ": ", err)
			log.Access(source, "", log.AccessRejected, err)
			payload.Release()
			continue
		}

		if request.Option.Has(RequestOptionOneTimeAuth) && v.account.OneTimeAuth == Account_Disabled {
			log.Info("Shadowsocks|Server: Client payload enables OTA but server doesn't allow it.")
			payload.Release()
			continue
		}

		if !request.Option.Has(RequestOptionOneTimeAuth) && v.account.OneTimeAuth == Account_Enabled {
			log.Info("Shadowsocks|Server: Client payload disables OTA but server forces it.")
			payload.Release()
			continue
		}

		dest := request.Destination()
		log.Access(source, dest, log.AccessAccepted, "")
		log.Info("Shadowsocks|Server: Tunnelling request to ", dest)

		ctx = protocol.ContextWithUser(ctx, request.User)
		v.udpServer.Dispatch(ctx, dest, data, func(payload *buf.Buffer) {
			defer payload.Release()

			data, err := EncodeUDPPacket(request, payload)
			if err != nil {
				log.Warning("Shadowsocks|Server: Failed to encode UDP packet: ", err)
				return
			}
			defer data.Release()

			conn.Write(data.Bytes())
		})
	}

	return nil
}

func (s *Server) handleConnection(ctx context.Context, conn internet.Connection) error {
	timedReader := net.NewTimeOutReader(16, conn)
	bufferedReader := bufio.NewReader(timedReader)
	request, bodyReader, err := ReadTCPSession(s.user, bufferedReader)
	if err != nil {
		log.Access(conn.RemoteAddr(), "", log.AccessRejected, err)
		log.Info("Shadowsocks|Server: Failed to create request from: ", conn.RemoteAddr(), ": ", err)
		return err
	}

	bufferedReader.SetBuffered(false)

	userSettings := s.user.GetSettings()
	timedReader.SetTimeOut(userSettings.PayloadReadTimeout)

	dest := request.Destination()
	log.Access(conn.RemoteAddr(), dest, log.AccessAccepted, "")
	log.Info("Shadowsocks|Server: Tunnelling request to ", dest)

	ctx = proxy.ContextWithDestination(ctx, dest)
	ctx = protocol.ContextWithUser(ctx, request.User)
	ray := s.packetDispatcher.DispatchToOutbound(ctx)

	requestDone := signal.ExecuteAsync(func() error {
		bufferedWriter := bufio.NewWriter(conn)
		responseWriter, err := WriteTCPResponse(request, bufferedWriter)
		if err != nil {
			log.Warning("Shadowsocks|Server: Failed to write response: ", err)
			return err
		}

		payload, err := ray.InboundOutput().Read()
		if err != nil {
			return err
		}
		if err := responseWriter.Write(payload); err != nil {
			return err
		}
		payload.Release()

		if err := bufferedWriter.SetBuffered(false); err != nil {
			return err
		}

		if err := buf.PipeUntilEOF(ray.InboundOutput(), responseWriter); err != nil {
			log.Info("Shadowsocks|Server: Failed to transport all TCP response: ", err)
			return err
		}

		return nil
	})

	responseDone := signal.ExecuteAsync(func() error {
		defer ray.InboundInput().Close()

		if err := buf.PipeUntilEOF(bodyReader, ray.InboundInput()); err != nil {
			log.Info("Shadowsocks|Server: Failed to transport all TCP request: ", err)
			return err
		}
		return nil
	})

	if err := signal.ErrorOrFinish2(ctx, requestDone, responseDone); err != nil {
		log.Info("Shadowsocks|Server: Connection ends with ", err)
		ray.InboundInput().CloseError()
		ray.InboundOutput().CloseError()
		return err
	}

	return nil
}

func init() {
	common.Must(common.RegisterConfig((*ServerConfig)(nil), func(ctx context.Context, config interface{}) (interface{}, error) {
		return NewServer(ctx, config.(*ServerConfig))
	}))
}
