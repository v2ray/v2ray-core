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

// NewServer create a new Shadowsocks server.
func NewServer(ctx context.Context, config *ServerConfig) (*Server, error) {
	space := app.SpaceFromContext(ctx)
	if space == nil {
		return nil, newError("no space in context")
	}
	if config.GetUser() == nil {
		return nil, newError("user is not specified")
	}

	rawAccount, err := config.User.GetTypedAccount()
	if err != nil {
		return nil, newError("failed to get user account").Base(err)
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
	switch network {
	case net.Network_TCP:
		return s.handleConnection(ctx, conn, dispatcher)
	case net.Network_UDP:
		return s.handlerUDPPayload(ctx, conn, dispatcher)
	default:
		return newError("unknown network: ", network)
	}
}

func (v *Server) handlerUDPPayload(ctx context.Context, conn internet.Connection, dispatcher dispatcher.Interface) error {
	udpServer := udp.NewDispatcher(dispatcher)

	reader := buf.NewReader(conn)
	for {
		mpayload, err := reader.Read()
		if err != nil {
			break
		}

		for _, payload := range mpayload {
			request, data, err := DecodeUDPPacket(v.user, payload)
			if err != nil {
				if source, ok := proxy.SourceFromContext(ctx); ok {
					log.Trace(newError("dropping invalid UDP packet from: ", source).Base(err))
					log.Access(source, "", log.AccessRejected, err)
				}
				payload.Release()
				continue
			}

			if request.Option.Has(RequestOptionOneTimeAuth) && v.account.OneTimeAuth == Account_Disabled {
				log.Trace(newError("client payload enables OTA but server doesn't allow it"))
				payload.Release()
				continue
			}

			if !request.Option.Has(RequestOptionOneTimeAuth) && v.account.OneTimeAuth == Account_Enabled {
				log.Trace(newError("client payload disables OTA but server forces it"))
				payload.Release()
				continue
			}

			dest := request.Destination()
			if source, ok := proxy.SourceFromContext(ctx); ok {
				log.Access(source, dest, log.AccessAccepted, "")
			}
			log.Trace(newError("tunnelling request to ", dest))

			ctx = protocol.ContextWithUser(ctx, request.User)
			udpServer.Dispatch(ctx, dest, data, func(payload *buf.Buffer) {
				defer payload.Release()

				data, err := EncodeUDPPacket(request, payload)
				if err != nil {
					log.Trace(newError("failed to encode UDP packet").Base(err).AtWarning())
					return
				}
				defer data.Release()

				conn.Write(data.Bytes())
			})
		}
	}

	return nil
}

func (s *Server) handleConnection(ctx context.Context, conn internet.Connection, dispatcher dispatcher.Interface) error {
	conn.SetReadDeadline(time.Now().Add(time.Second * 8))
	bufferedReader := buf.NewBufferedReader(conn)
	request, bodyReader, err := ReadTCPSession(s.user, bufferedReader)
	if err != nil {
		log.Access(conn.RemoteAddr(), "", log.AccessRejected, err)
		return newError("failed to create request from: ", conn.RemoteAddr()).Base(err)
	}
	conn.SetReadDeadline(time.Time{})

	bufferedReader.SetBuffered(false)

	dest := request.Destination()
	log.Access(conn.RemoteAddr(), dest, log.AccessAccepted, "")
	log.Trace(newError("tunnelling request to ", dest))

	ctx = protocol.ContextWithUser(ctx, request.User)

	userSettings := s.user.GetSettings()
	ctx, timer := signal.CancelAfterInactivity(ctx, userSettings.PayloadTimeout)
	ray, err := dispatcher.Dispatch(ctx, dest)
	if err != nil {
		return err
	}

	responseDone := signal.ExecuteAsync(func() error {
		bufferedWriter := buf.NewBufferedWriter(conn)
		responseWriter, err := WriteTCPResponse(request, bufferedWriter)
		if err != nil {
			return newError("failed to write response").Base(err)
		}

		mergeReader := buf.NewMergingReader(ray.InboundOutput())
		payload, err := mergeReader.Read()
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

		if err := buf.PipeUntilEOF(timer, mergeReader, responseWriter); err != nil {
			return newError("failed to transport all TCP response").Base(err)
		}

		return nil
	})

	requestDone := signal.ExecuteAsync(func() error {
		defer ray.InboundInput().Close()

		if err := buf.PipeUntilEOF(timer, bodyReader, ray.InboundInput()); err != nil {
			return newError("failed to transport all TCP request").Base(err)
		}
		return nil
	})

	if err := signal.ErrorOrFinish2(ctx, requestDone, responseDone); err != nil {
		ray.InboundInput().CloseError()
		ray.InboundOutput().CloseError()
		return newError("connection ends").Base(err)
	}

	runtime.KeepAlive(timer)

	return nil
}

func init() {
	common.Must(common.RegisterConfig((*ServerConfig)(nil), func(ctx context.Context, config interface{}) (interface{}, error) {
		return NewServer(ctx, config.(*ServerConfig))
	}))
}
