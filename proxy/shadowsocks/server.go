package shadowsocks

import (
	"context"
	"runtime"
	"time"

	"v2ray.com/core/app"
	"v2ray.com/core/app/dispatcher"
	"v2ray.com/core/app/log"
	"v2ray.com/core/common"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/errors"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/protocol"
	"v2ray.com/core/common/signal"
	"v2ray.com/core/proxy"
	"v2ray.com/core/transport/internet"
	"v2ray.com/core/transport/internet/udp"
)

type Server struct {
	config  *ServerConfig
	user    *protocol.User
	account *ShadowsocksAccount
}

func NewServer(ctx context.Context, config *ServerConfig) (*Server, error) {
	space := app.SpaceFromContext(ctx)
	if space == nil {
		return nil, errors.New("no space in context").Path("Shadowsocks", "Server")
	}
	if config.GetUser() == nil {
		return nil, errors.New("user is not specified").Path("Shadowsocks", "Server")
	}

	rawAccount, err := config.User.GetTypedAccount()
	if err != nil {
		return nil, errors.New("failed to get user account").Base(err).Path("Shadowsocks", "Server")
	}
	account := rawAccount.(*ShadowsocksAccount)

	s := &Server{
		config:  config,
		user:    config.GetUser(),
		account: account,
	}

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
	conn.SetReusable(false)

	switch network {
	case net.Network_TCP:
		return s.handleConnection(ctx, conn, dispatcher)
	case net.Network_UDP:
		return s.handlerUDPPayload(ctx, conn, dispatcher)
	default:
		return errors.New("unknown network: ", network).Path("Shadowsocks", "Server")
	}
}

func (v *Server) handlerUDPPayload(ctx context.Context, conn internet.Connection, dispatcher dispatcher.Interface) error {
	udpServer := udp.NewDispatcher(dispatcher)

	reader := buf.NewReader(conn)
	for {
		payload, err := reader.Read()
		if err != nil {
			break
		}

		request, data, err := DecodeUDPPacket(v.user, payload)
		if err != nil {
			if source, ok := proxy.SourceFromContext(ctx); ok {
				log.Trace(errors.New("dropping invalid UDP packet from: ", source).Base(err).Path("Shadowsocks", "Server"))
				log.Access(source, "", log.AccessRejected, err)
			}
			payload.Release()
			continue
		}

		if request.Option.Has(RequestOptionOneTimeAuth) && v.account.OneTimeAuth == Account_Disabled {
			log.Trace(errors.New("client payload enables OTA but server doesn't allow it").Path("Shadowsocks", "Server"))
			payload.Release()
			continue
		}

		if !request.Option.Has(RequestOptionOneTimeAuth) && v.account.OneTimeAuth == Account_Enabled {
			log.Trace(errors.New("client payload disables OTA but server forces it").Path("Shadowsocks", "Server"))
			payload.Release()
			continue
		}

		dest := request.Destination()
		if source, ok := proxy.SourceFromContext(ctx); ok {
			log.Access(source, dest, log.AccessAccepted, "")
		}
		log.Trace(errors.New("tunnelling request to ", dest).Path("Shadowsocks", "Server"))

		ctx = protocol.ContextWithUser(ctx, request.User)
		udpServer.Dispatch(ctx, dest, data, func(payload *buf.Buffer) {
			defer payload.Release()

			data, err := EncodeUDPPacket(request, payload)
			if err != nil {
				log.Trace(errors.New("failed to encode UDP packet").Base(err).Path("Shadowsocks", "Server").AtWarning())
				return
			}
			defer data.Release()

			conn.Write(data.Bytes())
		})
	}

	return nil
}

func (s *Server) handleConnection(ctx context.Context, conn internet.Connection, dispatcher dispatcher.Interface) error {
	conn.SetReadDeadline(time.Now().Add(time.Second * 8))
	bufferedReader := buf.NewBufferedReader(conn)
	request, bodyReader, err := ReadTCPSession(s.user, bufferedReader)
	if err != nil {
		log.Access(conn.RemoteAddr(), "", log.AccessRejected, err)
		return errors.New("failed to create request from: ", conn.RemoteAddr()).Base(err).Path("Shadowsocks", "Server")
	}
	conn.SetReadDeadline(time.Time{})

	bufferedReader.SetBuffered(false)

	dest := request.Destination()
	log.Access(conn.RemoteAddr(), dest, log.AccessAccepted, "")
	log.Trace(errors.New("tunnelling request to ", dest).Path("Shadowsocks", "Server"))

	ctx = protocol.ContextWithUser(ctx, request.User)

	userSettings := s.user.GetSettings()
	ctx, timer := signal.CancelAfterInactivity(ctx, userSettings.PayloadTimeout)
	ray, err := dispatcher.Dispatch(ctx, dest)
	if err != nil {
		return err
	}

	requestDone := signal.ExecuteAsync(func() error {
		bufferedWriter := buf.NewBufferedWriter(conn)
		responseWriter, err := WriteTCPResponse(request, bufferedWriter)
		if err != nil {
			return errors.New("failed to write response").Base(err).Path("Shadowsocks", "Server")
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

		if err := buf.PipeUntilEOF(timer, ray.InboundOutput(), responseWriter); err != nil {
			return errors.New("failed to transport all TCP response").Base(err).Path("Shadowsocks", "Server")
		}

		return nil
	})

	responseDone := signal.ExecuteAsync(func() error {
		defer ray.InboundInput().Close()

		if err := buf.PipeUntilEOF(timer, bodyReader, ray.InboundInput()); err != nil {
			return errors.New("failed to transport all TCP request").Base(err).Path("Shadowsocks", "Server")
		}
		return nil
	})

	if err := signal.ErrorOrFinish2(ctx, requestDone, responseDone); err != nil {
		ray.InboundInput().CloseError()
		ray.InboundOutput().CloseError()
		return errors.New("connection ends").Base(err).Path("Shadowsocks", "Server")
	}

	runtime.KeepAlive(timer)

	return nil
}

func init() {
	common.Must(common.RegisterConfig((*ServerConfig)(nil), func(ctx context.Context, config interface{}) (interface{}, error) {
		return NewServer(ctx, config.(*ServerConfig))
	}))
}
