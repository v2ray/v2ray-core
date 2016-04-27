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
	"github.com/v2ray/v2ray-core/common/serial"
	"github.com/v2ray/v2ray-core/proxy"
	"github.com/v2ray/v2ray-core/proxy/internal"
	"github.com/v2ray/v2ray-core/transport/hub"
)

type Shadowsocks struct {
	packetDispatcher dispatcher.PacketDispatcher
	config           *Config
	port             v2net.Port
	accepting        bool
	tcpHub           *hub.TCPHub
	udpHub           *hub.UDPHub
	udpServer        *hub.UDPServer
}

func NewShadowsocks(config *Config, packetDispatcher dispatcher.PacketDispatcher) *Shadowsocks {
	return &Shadowsocks{
		config:           config,
		packetDispatcher: packetDispatcher,
	}
}

func (this *Shadowsocks) Port() v2net.Port {
	return this.port
}

func (this *Shadowsocks) Close() {
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

func (this *Shadowsocks) Listen(port v2net.Port) error {
	if this.accepting {
		if this.port == port {
			return nil
		} else {
			return proxy.ErrorAlreadyListening
		}
	}

	tcpHub, err := hub.ListenTCP(port, this.handleConnection)
	if err != nil {
		log.Error("Shadowsocks: Failed to listen TCP on port ", port, ": ", err)
		return err
	}
	this.tcpHub = tcpHub

	if this.config.UDP {
		this.udpServer = hub.NewUDPServer(this.packetDispatcher)
		udpHub, err := hub.ListenUDP(port, this.handlerUDPPayload)
		if err != nil {
			log.Error("Shadowsocks: Failed to listen UDP on port ", port, ": ", err)
			return err
		}
		this.udpHub = udpHub
	}

	this.port = port
	this.accepting = true

	return nil
}

func (this *Shadowsocks) handlerUDPPayload(payload *alloc.Buffer, source v2net.Destination) {
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
		log.Access(source, serial.StringLiteral(""), log.AccessRejected, serial.StringLiteral(err.Error()))
		log.Warning("Shadowsocks: Invalid request from ", source, ": ", err)
		return
	}

	dest := v2net.UDPDestination(request.Address, request.Port)
	log.Access(source, dest, log.AccessAccepted, serial.StringLiteral(""))
	log.Info("Shadowsocks: Tunnelling request to ", dest)

	this.udpServer.Dispatch(source, dest, request.UDPPayload, func(destination v2net.Destination, payload *alloc.Buffer) {
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

func (this *Shadowsocks) handleConnection(conn hub.Connection) {
	defer conn.Close()

	buffer := alloc.NewSmallBuffer()
	defer buffer.Release()

	timedReader := v2net.NewTimeOutReader(16, conn)

	ivLen := this.config.Cipher.IVSize()
	_, err := io.ReadFull(timedReader, buffer.Value[:ivLen])
	if err != nil {
		log.Access(conn.RemoteAddr(), serial.StringLiteral(""), log.AccessRejected, serial.StringLiteral(err.Error()))
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

	reader := crypto.NewCryptionReader(stream, timedReader)

	request, err := ReadRequest(reader, NewAuthenticator(HeaderKeyGenerator(key, iv)), false)
	if err != nil {
		log.Access(conn.RemoteAddr(), serial.StringLiteral(""), log.AccessRejected, serial.StringLiteral(err.Error()))
		log.Warning("Shadowsocks: Invalid request from ", conn.RemoteAddr(), ": ", err)
		return
	}

	userSettings := protocol.GetUserSettings(this.config.Level)
	timedReader.SetTimeOut(userSettings.PayloadReadTimeout)

	dest := v2net.TCPDestination(request.Address, request.Port)
	log.Access(conn.RemoteAddr(), dest, log.AccessAccepted, serial.StringLiteral(""))
	log.Info("Shadowsocks: Tunnelling request to ", dest)

	ray := this.packetDispatcher.DispatchToOutbound(dest)

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
			defer writer.Release()

			v2io.Pipe(ray.InboundOutput(), v2writer)
			ray.InboundOutput().Release()
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
			return NewShadowsocks(
				rawConfig.(*Config),
				space.GetApp(dispatcher.APP_ID).(dispatcher.PacketDispatcher)), nil
		})
}
