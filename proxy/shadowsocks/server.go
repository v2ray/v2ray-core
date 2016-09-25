// R.I.P Shadowsocks
package shadowsocks

import (
	"crypto/rand"
	"io"
	"sync"

	"v2ray.com/core/app"
	"v2ray.com/core/app/dispatcher"
	"v2ray.com/core/common"
	"v2ray.com/core/common/alloc"
	"v2ray.com/core/common/crypto"
	v2io "v2ray.com/core/common/io"
	"v2ray.com/core/common/log"
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/common/protocol"
	"v2ray.com/core/proxy"
	"v2ray.com/core/proxy/registry"
	"v2ray.com/core/transport/internet"
	"v2ray.com/core/transport/internet/udp"
)

type Server struct {
	packetDispatcher dispatcher.PacketDispatcher
	config           *ServerConfig
	cipher           Cipher
	cipherKey        []byte
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
	account := new(Account)
	if _, err := config.GetUser().GetTypedAccount(account); err != nil {
		return nil, err
	}
	cipher := account.GetCipher()
	s := &Server{
		config:    config,
		meta:      meta,
		cipher:    cipher,
		cipherKey: account.GetCipherKey(cipher.KeySize()),
	}

	space.InitializeApplication(func() error {
		if !space.HasApp(dispatcher.APP_ID) {
			return app.ErrMissingApplication
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
		this.udpServer = udp.NewUDPServer(this.meta, this.packetDispatcher)
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
	defer payload.Release()

	source := session.Source
	ivLen := this.cipher.IVSize()
	iv := payload.Value[:ivLen]
	payload.SliceFrom(ivLen)

	stream, err := this.cipher.NewDecodingStream(this.cipherKey, iv)
	if err != nil {
		log.Error("Shadowsocks: Failed to create decoding stream: ", err)
		return
	}

	reader := crypto.NewCryptionReader(stream, payload)

	request, err := ReadRequest(reader, NewAuthenticator(HeaderKeyGenerator(this.cipherKey, iv)), true)
	if err != nil {
		if err != io.EOF {
			log.Access(source, "", log.AccessRejected, err)
			log.Warning("Shadowsocks: Invalid request from ", source, ": ", err)
		}
		return
	}
	//defer request.Release()

	dest := v2net.UDPDestination(request.Address, request.Port)
	log.Access(source, dest, log.AccessAccepted, "")
	log.Info("Shadowsocks: Tunnelling request to ", dest)

	this.udpServer.Dispatch(&proxy.SessionInfo{Source: source, Destination: dest}, request.DetachUDPPayload(), func(destination v2net.Destination, payload *alloc.Buffer) {
		defer payload.Release()

		response := alloc.NewBuffer().Slice(0, ivLen)
		defer response.Release()

		rand.Read(response.Value)
		respIv := response.Value

		stream, err := this.cipher.NewEncodingStream(this.cipherKey, respIv)
		if err != nil {
			log.Error("Shadowsocks: Failed to create encoding stream: ", err)
			return
		}

		writer := crypto.NewCryptionWriter(stream, response)

		switch request.Address.Family() {
		case v2net.AddressFamilyIPv4:
			writer.Write([]byte{AddrTypeIPv4})
			writer.Write(request.Address.IP())
		case v2net.AddressFamilyIPv6:
			writer.Write([]byte{AddrTypeIPv6})
			writer.Write(request.Address.IP())
		case v2net.AddressFamilyDomain:
			writer.Write([]byte{AddrTypeDomain, byte(len(request.Address.Domain()))})
			writer.Write([]byte(request.Address.Domain()))
		}

		writer.Write(request.Port.Bytes(nil))
		writer.Write(payload.Value)

		if request.OTA {
			respAuth := NewAuthenticator(HeaderKeyGenerator(this.cipherKey, respIv))
			respAuth.Authenticate(response.Value, response.Value[ivLen:])
		}

		this.udpHub.WriteTo(response.Value, source)
	})
}

func (this *Server) handleConnection(conn internet.Connection) {
	defer conn.Close()

	buffer := alloc.NewSmallBuffer()
	defer buffer.Release()

	timedReader := v2net.NewTimeOutReader(16, conn)
	defer timedReader.Release()

	bufferedReader := v2io.NewBufferedReader(timedReader)
	defer bufferedReader.Release()

	ivLen := this.cipher.IVSize()
	_, err := io.ReadFull(bufferedReader, buffer.Value[:ivLen])
	if err != nil {
		if err != io.EOF {
			log.Access(conn.RemoteAddr(), "", log.AccessRejected, err)
			log.Warning("Shadowsocks: Failed to read IV: ", err)
		}
		return
	}

	iv := buffer.Value[:ivLen]

	stream, err := this.cipher.NewDecodingStream(this.cipherKey, iv)
	if err != nil {
		log.Error("Shadowsocks: Failed to create decoding stream: ", err)
		return
	}

	reader := crypto.NewCryptionReader(stream, bufferedReader)

	request, err := ReadRequest(reader, NewAuthenticator(HeaderKeyGenerator(this.cipherKey, iv)), false)
	if err != nil {
		log.Access(conn.RemoteAddr(), "", log.AccessRejected, err)
		log.Warning("Shadowsocks: Invalid request from ", conn.RemoteAddr(), ": ", err)
		return
	}
	defer request.Release()
	bufferedReader.SetCached(false)

	userSettings := this.config.GetUser().GetSettings()
	timedReader.SetTimeOut(userSettings.PayloadReadTimeout)

	dest := v2net.TCPDestination(request.Address, request.Port)
	log.Access(conn.RemoteAddr(), dest, log.AccessAccepted, "")
	log.Info("Shadowsocks: Tunnelling request to ", dest)

	ray := this.packetDispatcher.DispatchToOutbound(this.meta, &proxy.SessionInfo{
		Source:      v2net.DestinationFromAddr(conn.RemoteAddr()),
		Destination: dest,
	})
	defer ray.InboundOutput().Release()

	var writeFinish sync.Mutex
	writeFinish.Lock()
	go func() {
		if payload, err := ray.InboundOutput().Read(); err == nil {
			payload.SliceBack(ivLen)
			rand.Read(payload.Value[:ivLen])

			stream, err := this.cipher.NewEncodingStream(this.cipherKey, payload.Value[:ivLen])
			if err != nil {
				log.Error("Shadowsocks: Failed to create encoding stream: ", err)
				return
			}
			stream.XORKeyStream(payload.Value[ivLen:], payload.Value[ivLen:])

			conn.Write(payload.Value)
			payload.Release()

			writer := crypto.NewCryptionWriter(stream, conn)
			v2writer := v2io.NewAdaptiveWriter(writer)

			v2io.Pipe(ray.InboundOutput(), v2writer)
			writer.Release()
			v2writer.Release()
		}
		writeFinish.Unlock()
	}()

	var payloadReader v2io.Reader
	if request.OTA {
		payloadAuth := NewAuthenticator(ChunkKeyGenerator(iv))
		payloadReader = NewChunkReader(reader, payloadAuth)
	} else {
		payloadReader = v2io.NewAdaptiveReader(reader)
	}

	v2io.Pipe(payloadReader, ray.InboundInput())
	ray.InboundInput().Close()
	payloadReader.Release()

	writeFinish.Lock()
}

type ServerFactory struct{}

func (this *ServerFactory) StreamCapability() internet.StreamConnectionType {
	return internet.StreamConnectionTypeRawTCP
}

func (this *ServerFactory) Create(space app.Space, rawConfig interface{}, meta *proxy.InboundHandlerMeta) (proxy.InboundHandler, error) {
	if !space.HasApp(dispatcher.APP_ID) {
		return nil, common.ErrBadConfiguration
	}
	return NewServer(rawConfig.(*ServerConfig), space, meta)
}

func init() {
	registry.MustRegisterInboundHandlerCreator("shadowsocks", new(ServerFactory))
}
