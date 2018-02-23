package shadowsocks

import (
	"context"
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
)

type Server struct {
	config  *ServerConfig
	user    *protocol.User
	account *MemoryAccount
	v       *core.Instance
}

// NewServer create a new Shadowsocks server.
func NewServer(ctx context.Context, config *ServerConfig) (*Server, error) {
	if config.GetUser() == nil {
		return nil, newError("user is not specified")
	}

	rawAccount, err := config.User.GetTypedAccount()
	if err != nil {
		return nil, newError("failed to get user account").Base(err)
	}
	account := rawAccount.(*MemoryAccount)

	s := &Server{
		config:  config,
		user:    config.GetUser(),
		account: account,
		v:       core.MustFromContext(ctx),
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

func (s *Server) Process(ctx context.Context, network net.Network, conn internet.Connection, dispatcher core.Dispatcher) error {
	switch network {
	case net.Network_TCP:
		return s.handleConnection(ctx, conn, dispatcher)
	case net.Network_UDP:
		return s.handlerUDPPayload(ctx, conn, dispatcher)
	default:
		return newError("unknown network: ", network)
	}
}

func (s *Server) handlerUDPPayload(ctx context.Context, conn internet.Connection, dispatcher core.Dispatcher) error {
	udpServer := udp.NewDispatcher(dispatcher)

	reader := buf.NewReader(conn)
	for {
		mpayload, err := reader.ReadMultiBuffer()
		if err != nil {
			break
		}

		for _, payload := range mpayload {
			request, data, err := DecodeUDPPacket(s.user, payload)
			if err != nil {
				if source, ok := proxy.SourceFromContext(ctx); ok {
					newError("dropping invalid UDP packet from: ", source).Base(err).WithContext(ctx).WriteToLog()
					log.Record(&log.AccessMessage{
						From:   source,
						To:     "",
						Status: log.AccessRejected,
						Reason: err,
					})
				}
				payload.Release()
				continue
			}

			if request.Option.Has(RequestOptionOneTimeAuth) && s.account.OneTimeAuth == Account_Disabled {
				newError("client payload enables OTA but server doesn't allow it").WithContext(ctx).WriteToLog()
				payload.Release()
				continue
			}

			if !request.Option.Has(RequestOptionOneTimeAuth) && s.account.OneTimeAuth == Account_Enabled {
				newError("client payload disables OTA but server forces it").WithContext(ctx).WriteToLog()
				payload.Release()
				continue
			}

			dest := request.Destination()
			if source, ok := proxy.SourceFromContext(ctx); ok {
				log.Record(&log.AccessMessage{
					From:   source,
					To:     dest,
					Status: log.AccessAccepted,
					Reason: "",
				})
			}
			newError("tunnelling request to ", dest).WithContext(ctx).WriteToLog()

			ctx = protocol.ContextWithUser(ctx, request.User)
			udpServer.Dispatch(ctx, dest, data, func(payload *buf.Buffer) {
				defer payload.Release()

				data, err := EncodeUDPPacket(request, payload.Bytes())
				if err != nil {
					newError("failed to encode UDP packet").Base(err).AtWarning().WithContext(ctx).WriteToLog()
					return
				}
				defer data.Release()

				conn.Write(data.Bytes())
			})
		}
	}

	return nil
}

func (s *Server) handleConnection(ctx context.Context, conn internet.Connection, dispatcher core.Dispatcher) error {
	sessionPolicy := s.v.PolicyManager().ForLevel(s.user.Level)
	conn.SetReadDeadline(time.Now().Add(sessionPolicy.Timeouts.Handshake))
	bufferedReader := buf.NewBufferedReader(buf.NewReader(conn))
	request, bodyReader, err := ReadTCPSession(s.user, bufferedReader)
	if err != nil {
		log.Record(&log.AccessMessage{
			From:   conn.RemoteAddr(),
			To:     "",
			Status: log.AccessRejected,
			Reason: err,
		})
		return newError("failed to create request from: ", conn.RemoteAddr()).Base(err)
	}
	conn.SetReadDeadline(time.Time{})

	bufferedReader.SetBuffered(false)

	dest := request.Destination()
	log.Record(&log.AccessMessage{
		From:   conn.RemoteAddr(),
		To:     dest,
		Status: log.AccessAccepted,
		Reason: "",
	})
	newError("tunnelling request to ", dest).WithContext(ctx).WriteToLog()

	ctx = protocol.ContextWithUser(ctx, request.User)

	ctx, cancel := context.WithCancel(ctx)
	timer := signal.CancelAfterInactivity(ctx, cancel, sessionPolicy.Timeouts.ConnectionIdle)
	ray, err := dispatcher.Dispatch(ctx, dest)
	if err != nil {
		return err
	}

	responseDone := signal.ExecuteAsync(func() error {
		defer timer.SetTimeout(sessionPolicy.Timeouts.UplinkOnly)

		bufferedWriter := buf.NewBufferedWriter(buf.NewWriter(conn))
		responseWriter, err := WriteTCPResponse(request, bufferedWriter)
		if err != nil {
			return newError("failed to write response").Base(err)
		}

		payload, err := ray.InboundOutput().ReadMultiBuffer()
		if err != nil {
			return err
		}
		if err := responseWriter.WriteMultiBuffer(payload); err != nil {
			return err
		}
		payload.Release()

		if err := bufferedWriter.SetBuffered(false); err != nil {
			return err
		}

		if err := buf.Copy(ray.InboundOutput(), responseWriter, buf.UpdateActivity(timer)); err != nil {
			return newError("failed to transport all TCP response").Base(err)
		}

		return nil
	})

	requestDone := signal.ExecuteAsync(func() error {
		defer timer.SetTimeout(sessionPolicy.Timeouts.DownlinkOnly)
		defer ray.InboundInput().Close()

		if err := buf.Copy(bodyReader, ray.InboundInput(), buf.UpdateActivity(timer)); err != nil {
			return newError("failed to transport all TCP request").Base(err)
		}

		return nil
	})

	if err := signal.ErrorOrFinish2(ctx, requestDone, responseDone); err != nil {
		ray.InboundInput().CloseError()
		ray.InboundOutput().CloseError()
		return newError("connection ends").Base(err)
	}

	return nil
}

func init() {
	common.Must(common.RegisterConfig((*ServerConfig)(nil), func(ctx context.Context, config interface{}) (interface{}, error) {
		return NewServer(ctx, config.(*ServerConfig))
	}))
}
