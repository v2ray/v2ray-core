// R.I.P Shadowsocks
package shadowsocks

import (
	"crypto/rand"
	"io"
	"sync"

	"github.com/v2ray/v2ray-core/app"
	"github.com/v2ray/v2ray-core/app/dispatcher"
	"github.com/v2ray/v2ray-core/common/alloc"
	"github.com/v2ray/v2ray-core/common/crypto"
	v2io "github.com/v2ray/v2ray-core/common/io"
	"github.com/v2ray/v2ray-core/common/log"
	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/common/protocol"
	"github.com/v2ray/v2ray-core/proxy"
	"github.com/v2ray/v2ray-core/proxy/internal"
	"github.com/v2ray/v2ray-core/transport/hub"
)

type Server struct {
	packetDispatcher dispatcher.PacketDispatcher
	config           *Config
	port             v2net.Port
	address          v2net.Address
	accepting        bool
	tcpHub           *hub.TCPHub
	udpHub           *hub.UDPHub
	udpServer        *hub.UDPServer
}

func NewServer(config *Config, packetDispatcher dispatcher.PacketDispatcher) *Server {
	return &Server{
		config:           config,
		packetDispatcher: packetDispatcher,
	}
}

func (this *Server) Port() v2net.Port {
	return this.port
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

func (this *Server) Listen(address v2net.Address, port v2net.Port) error {
	if this.accepting {
		if this.port == port && this.address.Equals(address) {
			return nil
		} else {
			return proxy.ErrorAlreadyListening
		}
	}

	tcpHub, err := hub.ListenTCP(address, port, this.handleConnection, nil)
	if err != nil {
		log.Error("Shadowsocks: Failed to listen TCP on port ", port, ": ", err)
		return err
	}
	this.tcpHub = tcpHub

	if this.config.UDP {
		this.udpServer = hub.NewUDPServer(this.packetDispatcher)
		udpHub, err := hub.ListenUDP(address, port, this.handlerUDPPayload)
		if err != nil {
			log.Error("Shadowsocks: Failed to listen UDP on port ", port, ": ", err)
			return err
		}
		this.udpHub = udpHub
	}

	this.port = port
	this.address = address
	this.accepting = true

	return nil
}

func (this *Server) handlerUDPPayload(payload *alloc.Buffer, source v2net.Destination) {
	defer payload.Release()

	ivLen := this.config.Cipher.IVSize()
	iv := payload.Value[:ivLen]
	key := this.config.Key
	payload.SliceFrom(ivLen)

	stream, err := this.config.Cipher.NewDecodingStream(key, iv)
	if err != nil {
		log.Error("Shadowsocks: Failed to create decoding stream: ", err)
		return
	}

	reader := crypto.NewCryptionReader(stream, payload)

	request, err := ReadRequest(reader, NewAuthenticator(HeaderKeyGenerator(key, iv)), true)
	if err != nil {
		log.Access(source, "", log.AccessRejected, err)
		log.Warning("Shadowsocks: Invalid request from ", source, ": ", err)
		return
	}
	//defer request.Release()

	dest := v2net.UDPDestination(request.Address, request.Port)
	log.Access(source, dest, log.AccessAccepted, "")
	log.Info("Shadowsocks: Tunnelling request to ", dest)

	this.udpServer.Dispatch(source, dest, request.DetachUDPPayload(), func(destination v2net.Destination, payload *alloc.Buffer) {
		defer payload.Release()

		response := alloc.NewBuffer().Slice(0, ivLen)
		defer response.Release()

		rand.Read(response.Value)
		respIv := response.Value

		stream, err := this.config.Cipher.NewEncodingStream(key, respIv)
		if err != nil {
			log.Error("Shadowsocks: Failed to create encoding stream: ", err)
			return
		}

		writer := crypto.NewCryptionWriter(stream, response)

		switch {
		case request.Address.IsIPv4():
			writer.Write([]byte{AddrTypeIPv4})
			writer.Write(request.Address.IP())
		case request.Address.IsIPv6():
			writer.Write([]byte{AddrTypeIPv6})
			writer.Write(request.Address.IP())
		case request.Address.IsDomain():
			writer.Write([]byte{AddrTypeDomain, byte(len(request.Address.Domain()))})
			writer.Write([]byte(request.Address.Domain()))
		}

		writer.Write(request.Port.Bytes())
		writer.Write(payload.Value)

		if request.OTA {
			respAuth := NewAuthenticator(HeaderKeyGenerator(key, respIv))
			respAuth.Authenticate(response.Value, response.Value[ivLen:])
		}

		this.udpHub.WriteTo(response.Value, source)
	})
}

func (this *Server) handleConnection(conn *hub.Connection) {
	defer conn.Close()

	buffer := alloc.NewSmallBuffer()
	defer buffer.Release()

	timedReader := v2net.NewTimeOutReader(16, conn)
	defer timedReader.Release()

	bufferedReader := v2io.NewBufferedReader(timedReader)
	defer bufferedReader.Release()

	ivLen := this.config.Cipher.IVSize()
	_, err := io.ReadFull(bufferedReader, buffer.Value[:ivLen])
	if err != nil {
		log.Access(conn.RemoteAddr(), "", log.AccessRejected, err)
		log.Error("Shadowsocks: Failed to read IV: ", err)
		return
	}

	iv := buffer.Value[:ivLen]
	key := this.config.Key

	stream, err := this.config.Cipher.NewDecodingStream(key, iv)
	if err != nil {
		log.Error("Shadowsocks: Failed to create decoding stream: ", err)
		return
	}

	reader := crypto.NewCryptionReader(stream, bufferedReader)

	request, err := ReadRequest(reader, NewAuthenticator(HeaderKeyGenerator(key, iv)), false)
	if err != nil {
		log.Access(conn.RemoteAddr(), "", log.AccessRejected, err)
		log.Warning("Shadowsocks: Invalid request from ", conn.RemoteAddr(), ": ", err)
		return
	}
	defer request.Release()
	bufferedReader.SetCached(false)

	userSettings := protocol.GetUserSettings(this.config.Level)
	timedReader.SetTimeOut(userSettings.PayloadReadTimeout)

	dest := v2net.TCPDestination(request.Address, request.Port)
	log.Access(conn.RemoteAddr(), dest, log.AccessAccepted, "")
	log.Info("Shadowsocks: Tunnelling request to ", dest)

	ray := this.packetDispatcher.DispatchToOutbound(dest)
	defer ray.InboundOutput().Release()

	var writeFinish sync.Mutex
	writeFinish.Lock()
	go func() {
		if payload, err := ray.InboundOutput().Read(); err == nil {
			payload.SliceBack(ivLen)
			rand.Read(payload.Value[:ivLen])

			stream, err := this.config.Cipher.NewEncodingStream(key, payload.Value[:ivLen])
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

func init() {
	internal.MustRegisterInboundHandlerCreator("shadowsocks",
		func(space app.Space, rawConfig interface{}) (proxy.InboundHandler, error) {
			if !space.HasApp(dispatcher.APP_ID) {
				return nil, internal.ErrorBadConfiguration
			}
			return NewServer(
				rawConfig.(*Config),
				space.GetApp(dispatcher.APP_ID).(dispatcher.PacketDispatcher)), nil
		})
}
