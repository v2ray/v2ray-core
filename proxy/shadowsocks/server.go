package shadowsocks

import (
	"v2ray.com/core/app"
	"v2ray.com/core/app/dispatcher"
	"v2ray.com/core/common"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/bufio"
	"v2ray.com/core/common/errors"
	"v2ray.com/core/common/log"
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/common/protocol"
	"v2ray.com/core/common/signal"
	"v2ray.com/core/proxy"
	"v2ray.com/core/transport/internet"
	"v2ray.com/core/transport/internet/udp"
)

type Server struct {
	packetDispatcher dispatcher.PacketDispatcher
	config           *ServerConfig
	user             *protocol.User
	account          *ShadowsocksAccount
	meta             *proxy.InboundHandlerMeta
	accepting        bool
	tcpHub           *internet.TCPHub
	udpHub           *udp.Hub
	udpServer        *udp.Server
}

func NewServer(config *ServerConfig, space app.Space, meta *proxy.InboundHandlerMeta) (*Server, error) {
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
		meta:    meta,
		user:    config.GetUser(),
		account: account,
	}

	space.InitializeApplication(func() error {
		if !space.HasApp(dispatcher.APP_ID) {
			return errors.New("Shadowsocks|Server: Dispatcher is not found in space.")
		}
		s.packetDispatcher = space.GetApp(dispatcher.APP_ID).(dispatcher.PacketDispatcher)
		return nil
	})

	return s, nil
}

func (v *Server) Port() v2net.Port {
	return v.meta.Port
}

func (v *Server) Close() {
	v.accepting = false
	// TODO: synchronization
	if v.tcpHub != nil {
		v.tcpHub.Close()
		v.tcpHub = nil
	}

	if v.udpHub != nil {
		v.udpHub.Close()
		v.udpHub = nil
	}
}

func (v *Server) Start() error {
	if v.accepting {
		return nil
	}

	tcpHub, err := internet.ListenTCP(v.meta.Address, v.meta.Port, v.handleConnection, v.meta.StreamSettings)
	if err != nil {
		log.Error("Shadowsocks: Failed to listen TCP on ", v.meta.Address, ":", v.meta.Port, ": ", err)
		return err
	}
	v.tcpHub = tcpHub

	if v.config.UdpEnabled {
		v.udpServer = udp.NewServer(v.packetDispatcher)
		udpHub, err := udp.ListenUDP(v.meta.Address, v.meta.Port, udp.ListenOption{Callback: v.handlerUDPPayload})
		if err != nil {
			log.Error("Shadowsocks: Failed to listen UDP on ", v.meta.Address, ":", v.meta.Port, ": ", err)
			return err
		}
		v.udpHub = udpHub
	}

	v.accepting = true

	return nil
}

func (v *Server) handlerUDPPayload(payload *buf.Buffer, session *proxy.SessionInfo) {
	source := session.Source
	request, data, err := DecodeUDPPacket(v.user, payload)
	if err != nil {
		log.Info("Shadowsocks|Server: Skipping invalid UDP packet from: ", source, ": ", err)
		log.Access(source, "", log.AccessRejected, err)
		payload.Release()
		return
	}

	if request.Option.Has(RequestOptionOneTimeAuth) && v.account.OneTimeAuth == Account_Disabled {
		log.Info("Shadowsocks|Server: Client payload enables OTA but server doesn't allow it.")
		payload.Release()
		return
	}

	if !request.Option.Has(RequestOptionOneTimeAuth) && v.account.OneTimeAuth == Account_Enabled {
		log.Info("Shadowsocks|Server: Client payload disables OTA but server forces it.")
		payload.Release()
		return
	}

	dest := request.Destination()
	log.Access(source, dest, log.AccessAccepted, "")
	log.Info("Shadowsocks|Server: Tunnelling request to ", dest)

	v.udpServer.Dispatch(&proxy.SessionInfo{Source: source, Destination: dest, User: request.User, Inbound: v.meta}, data, func(destination v2net.Destination, payload *buf.Buffer) {
		defer payload.Release()

		data, err := EncodeUDPPacket(request, payload)
		if err != nil {
			log.Warning("Shadowsocks|Server: Failed to encode UDP packet: ", err)
			return
		}
		defer data.Release()

		v.udpHub.WriteTo(data.Bytes(), source)
	})
}

func (v *Server) handleConnection(conn internet.Connection) {
	defer conn.Close()
	conn.SetReusable(false)

	timedReader := v2net.NewTimeOutReader(16, conn)
	bufferedReader := bufio.NewReader(timedReader)
	request, bodyReader, err := ReadTCPSession(v.user, bufferedReader)
	if err != nil {
		log.Access(conn.RemoteAddr(), "", log.AccessRejected, err)
		log.Info("Shadowsocks|Server: Failed to create request from: ", conn.RemoteAddr(), ": ", err)
		return
	}

	bufferedReader.SetBuffered(false)

	userSettings := v.user.GetSettings()
	timedReader.SetTimeOut(userSettings.PayloadReadTimeout)

	dest := request.Destination()
	log.Access(conn.RemoteAddr(), dest, log.AccessAccepted, "")
	log.Info("Shadowsocks|Server: Tunnelling request to ", dest)

	ray := v.packetDispatcher.DispatchToOutbound(&proxy.SessionInfo{
		Source:      v2net.DestinationFromAddr(conn.RemoteAddr()),
		Destination: dest,
		User:        request.User,
		Inbound:     v.meta,
	})
	defer ray.InboundOutput().ForceClose()
	defer ray.InboundInput().Close()

	requestDone := signal.ExecuteAsync(func() error {
		defer ray.InboundOutput().ForceClose()

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

	if err := signal.ErrorOrFinish2(requestDone, responseDone); err != nil {
		log.Info("Shadowsocks|Server: Connection ends with ", err)
	}
}

type ServerFactory struct{}

func (v *ServerFactory) StreamCapability() v2net.NetworkList {
	return v2net.NetworkList{
		Network: []v2net.Network{v2net.Network_TCP},
	}
}

func (v *ServerFactory) Create(space app.Space, rawConfig interface{}, meta *proxy.InboundHandlerMeta) (proxy.InboundHandler, error) {
	if !space.HasApp(dispatcher.APP_ID) {
		return nil, common.ErrBadConfiguration
	}
	return NewServer(rawConfig.(*ServerConfig), space, meta)
}
