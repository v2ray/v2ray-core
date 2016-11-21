// R.I.P Shadowsocks
package shadowsocks

import (
	"sync"

	"errors"
	"v2ray.com/core/app"
	"v2ray.com/core/app/dispatcher"
	"v2ray.com/core/common"
	"v2ray.com/core/common/alloc"
	v2io "v2ray.com/core/common/io"
	"v2ray.com/core/common/log"
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/common/protocol"
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
	udpHub           *udp.UDPHub
	udpServer        *udp.UDPServer
}

func NewServer(config *ServerConfig, space app.Space, meta *proxy.InboundHandlerMeta) (*Server, error) {
	if config.GetUser() == nil {
		return nil, protocol.ErrUserMissing
	}

	rawAccount, err := config.User.GetTypedAccount()
	if err != nil {
		return nil, errors.New("Shadowsocks|Server: Failed to get user account: " + err.Error())
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

func (this *Server) Port() v2net.Port {
	return this.meta.Port
}

func (this *Server) Close() {
	this.accepting = false
	// TODO: synchronization
	if this.tcpHub != nil {
		this.tcpHub.Close()
		this.tcpHub = nil
	}

	if this.udpHub != nil {
		this.udpHub.Close()
		this.udpHub = nil
	}
}

func (this *Server) Start() error {
	if this.accepting {
		return nil
	}

	tcpHub, err := internet.ListenTCP(this.meta.Address, this.meta.Port, this.handleConnection, this.meta.StreamSettings)
	if err != nil {
		log.Error("Shadowsocks: Failed to listen TCP on ", this.meta.Address, ":", this.meta.Port, ": ", err)
		return err
	}
	this.tcpHub = tcpHub

	if this.config.UdpEnabled {
		this.udpServer = udp.NewUDPServer(this.packetDispatcher)
		udpHub, err := udp.ListenUDP(this.meta.Address, this.meta.Port, udp.ListenOption{Callback: this.handlerUDPPayload})
		if err != nil {
			log.Error("Shadowsocks: Failed to listen UDP on ", this.meta.Address, ":", this.meta.Port, ": ", err)
			return err
		}
		this.udpHub = udpHub
	}

	this.accepting = true

	return nil
}

func (this *Server) handlerUDPPayload(payload *alloc.Buffer, session *proxy.SessionInfo) {
	source := session.Source
	request, data, err := DecodeUDPPacket(this.user, payload)
	if err != nil {
		log.Info("Shadowsocks|Server: Skipping invalid UDP packet from: ", source, ": ", err)
		log.Access(source, "", log.AccessRejected, err)
		payload.Release()
		return
	}

	if request.Option.Has(RequestOptionOneTimeAuth) && this.account.OneTimeAuth == Account_Disabled {
		log.Info("Shadowsocks|Server: Client payload enables OTA but server doesn't allow it.")
		payload.Release()
		return
	}

	if !request.Option.Has(RequestOptionOneTimeAuth) && this.account.OneTimeAuth == Account_Enabled {
		log.Info("Shadowsocks|Server: Client payload disables OTA but server forces it.")
		payload.Release()
		return
	}

	dest := request.Destination()
	log.Access(source, dest, log.AccessAccepted, "")
	log.Info("Shadowsocks|Server: Tunnelling request to ", dest)

	this.udpServer.Dispatch(&proxy.SessionInfo{Source: source, Destination: dest, User: request.User, Inbound: this.meta}, data, func(destination v2net.Destination, payload *alloc.Buffer) {
		defer payload.Release()

		data, err := EncodeUDPPacket(request, payload)
		if err != nil {
			log.Warning("Shadowsocks|Server: Failed to encode UDP packet: ", err)
			return
		}
		defer data.Release()

		this.udpHub.WriteTo(data.Value, source)
	})
}

func (this *Server) handleConnection(conn internet.Connection) {
	defer conn.Close()
	conn.SetReusable(false)

	timedReader := v2net.NewTimeOutReader(16, conn)
	defer timedReader.Release()

	bufferedReader := v2io.NewBufferedReader(timedReader)
	defer bufferedReader.Release()

	request, bodyReader, err := ReadTCPSession(this.user, bufferedReader)
	if err != nil {
		log.Access(conn.RemoteAddr(), "", log.AccessRejected, err)
		log.Info("Shadowsocks|Server: Failed to create request from: ", conn.RemoteAddr(), ": ", err)
		return
	}
	defer bodyReader.Release()

	bufferedReader.SetCached(false)

	userSettings := this.user.GetSettings()
	timedReader.SetTimeOut(userSettings.PayloadReadTimeout)

	dest := request.Destination()
	log.Access(conn.RemoteAddr(), dest, log.AccessAccepted, "")
	log.Info("Shadowsocks|Server: Tunnelling request to ", dest)

	ray := this.packetDispatcher.DispatchToOutbound(&proxy.SessionInfo{
		Source:      v2net.DestinationFromAddr(conn.RemoteAddr()),
		Destination: dest,
		User:        request.User,
		Inbound:     this.meta,
	})
	defer ray.InboundOutput().Release()

	var writeFinish sync.Mutex
	writeFinish.Lock()
	go func() {
		defer writeFinish.Unlock()

		bufferedWriter := v2io.NewBufferedWriter(conn)
		defer bufferedWriter.Release()

		responseWriter, err := WriteTCPResponse(request, bufferedWriter)
		if err != nil {
			log.Warning("Shadowsocks|Server: Failed to write response: ", err)
			return
		}
		defer responseWriter.Release()

		if payload, err := ray.InboundOutput().Read(); err == nil {
			responseWriter.Write(payload)
			bufferedWriter.SetCached(false)

			if err := v2io.PipeUntilEOF(ray.InboundOutput(), responseWriter); err != nil {
				log.Info("Shadowsocks|Server: Failed to transport all TCP response: ", err)
			}
		}
	}()

	if err := v2io.PipeUntilEOF(bodyReader, ray.InboundInput()); err != nil {
		log.Info("Shadowsocks|Server: Failed to transport all TCP request: ", err)
	}
	ray.InboundInput().Close()

	writeFinish.Lock()
}

type ServerFactory struct{}

func (this *ServerFactory) StreamCapability() v2net.NetworkList {
	return v2net.NetworkList{
		Network: []v2net.Network{v2net.Network_TCP, v2net.Network_RawTCP},
	}
}

func (this *ServerFactory) Create(space app.Space, rawConfig interface{}, meta *proxy.InboundHandlerMeta) (proxy.InboundHandler, error) {
	if !space.HasApp(dispatcher.APP_ID) {
		return nil, common.ErrBadConfiguration
	}
	return NewServer(rawConfig.(*ServerConfig), space, meta)
}
