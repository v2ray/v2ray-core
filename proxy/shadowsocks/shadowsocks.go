// R.I.P Shadowsocks

package shadowsocks

import (
	"crypto/rand"
	"io"
	"sync"

	"github.com/v2ray/v2ray-core/app"
	"github.com/v2ray/v2ray-core/app/dispatcher"
	"github.com/v2ray/v2ray-core/common/alloc"
	v2io "github.com/v2ray/v2ray-core/common/io"
	"github.com/v2ray/v2ray-core/common/log"
	v2net "github.com/v2ray/v2ray-core/common/net"
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
	this.accepting = true

	tcpHub, err := hub.ListenTCP(port, this.handleConnection)
	if err != nil {
		log.Error("Shadowsocks: Failed to listen TCP on port ", port, ": ", err)
		return err
	}
	this.tcpHub = tcpHub

	if this.config.UDP {
		udpHub, err := hub.ListenUDP(port, this.handlerUDPPayload)
		if err != nil {
			log.Error("Shadowsocks: Failed to listen UDP on port ", port, ": ", err)
		}
		this.udpHub = udpHub
	}

	return nil
}

func (this *Shadowsocks) handlerUDPPayload(payload *alloc.Buffer, source v2net.Destination) {
	defer payload.Release()

	iv := payload.Value[:this.config.Cipher.IVSize()]
	key := this.config.Key
	payload.SliceFrom(this.config.Cipher.IVSize())

	reader, err := this.config.Cipher.NewDecodingStream(key, iv, payload)
	if err != nil {
		log.Error("Shadowsocks: Failed to create decoding stream: ", err)
		return
	}

	request, err := ReadRequest(reader, NewAuthenticator(HeaderKeyGenerator(key, iv)), true)
	if err != nil {
		log.Access(source, serial.StringLiteral(""), log.AccessRejected, serial.StringLiteral(err.Error()))
		log.Warning("Shadowsocks: Invalid request from ", source, ": ", err)
		return
	}

	dest := v2net.UDPDestination(request.Address, request.Port)
	log.Access(source, dest, log.AccessAccepted, serial.StringLiteral(""))
	log.Info("Shadowsocks: Tunnelling request to ", dest)

	packet := v2net.NewPacket(dest, request.UDPPayload, false)
	ray := this.packetDispatcher.DispatchToOutbound(packet)
	close(ray.InboundInput())

	for respChunk := range ray.InboundOutput() {

		response := alloc.NewBuffer().Slice(0, this.config.Cipher.IVSize())
		rand.Read(response.Value)
		respIv := response.Value

		writer, err := this.config.Cipher.NewEncodingStream(key, respIv, response)
		if err != nil {
			log.Error("Shadowsocks: Failed to create encoding stream: ", err)
			return
		}

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
		writer.Write(respChunk.Value)
		respChunk.Release()

		if request.OTA {
			respAuth := NewAuthenticator(HeaderKeyGenerator(key, respIv))
			respAuth.Authenticate(response.Value, response.Value[this.config.Cipher.IVSize():])
		}

		this.udpHub.WriteTo(response.Value, source)
		response.Release()
	}
}

func (this *Shadowsocks) handleConnection(conn *hub.TCPConn) {
	defer conn.Close()

	buffer := alloc.NewSmallBuffer()
	defer buffer.Release()

	_, err := io.ReadFull(conn, buffer.Value[:this.config.Cipher.IVSize()])
	if err != nil {
		log.Access(conn.RemoteAddr(), serial.StringLiteral(""), log.AccessRejected, serial.StringLiteral(err.Error()))
		log.Error("Shadowsocks: Failed to read IV: ", err)
		return
	}

	iv := buffer.Value[:this.config.Cipher.IVSize()]
	key := this.config.Key

	reader, err := this.config.Cipher.NewDecodingStream(key, iv, conn)
	if err != nil {
		log.Error("Shadowsocks: Failed to create decoding stream: ", err)
		return
	}

	request, err := ReadRequest(reader, NewAuthenticator(HeaderKeyGenerator(iv, key)), false)
	if err != nil {
		log.Access(conn.RemoteAddr(), serial.StringLiteral(""), log.AccessRejected, serial.StringLiteral(err.Error()))
		log.Warning("Shadowsocks: Invalid request from ", conn.RemoteAddr(), ": ", err)
		return
	}

	dest := v2net.TCPDestination(request.Address, request.Port)
	log.Access(conn.RemoteAddr(), dest, log.AccessAccepted, serial.StringLiteral(""))
	log.Info("Shadowsocks: Tunnelling request to ", dest)

	packet := v2net.NewPacket(dest, nil, true)
	ray := this.packetDispatcher.DispatchToOutbound(packet)

	var writeFinish sync.Mutex
	writeFinish.Lock()
	go func() {
		if payload, ok := <-ray.InboundOutput(); ok {
			payload.SliceBack(16)
			rand.Read(payload.Value[:16])

			writer, err := this.config.Cipher.NewEncodingStream(key, payload.Value[:16], conn)
			if err != nil {
				log.Error("Shadowsocks: Failed to create encoding stream: ", err)
				return
			}

			writer.Write(payload.Value)
			payload.Release()
			v2io.ChanToWriter(writer, ray.InboundOutput())
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

	v2io.ReaderToChan(ray.InboundInput(), payloadReader)
	close(ray.InboundInput())

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
