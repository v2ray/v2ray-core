// R.I.P Shadowsocks

package shadowsocks

import (
	"crypto/rand"
	"io"
	"sync"

	"github.com/v2ray/v2ray-core/app"
	"github.com/v2ray/v2ray-core/common/alloc"
	v2io "github.com/v2ray/v2ray-core/common/io"
	"github.com/v2ray/v2ray-core/common/log"
	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/proxy"
	"github.com/v2ray/v2ray-core/proxy/internal"
	"github.com/v2ray/v2ray-core/transport/hub"
)

type Shadowsocks struct {
	space     app.Space
	config    *Config
	port      v2net.Port
	accepting bool
	tcpHub    *hub.TCPHub
	udpHub    *hub.UDPHub
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
		return
	}

	dest := v2net.UDPDestination(request.Address, request.Port)
	log.Info("Shadowsocks: Tunnelling request to ", dest)
	packet := v2net.NewPacket(dest, request.UDPPayload, false)
	ray := this.space.PacketDispatcher().DispatchToOutbound(packet)
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
		return
	}

	packet := v2net.NewPacket(v2net.TCPDestination(request.Address, request.Port), nil, true)
	ray := this.space.PacketDispatcher().DispatchToOutbound(packet)

	var writeFinish sync.Mutex
	writeFinish.Lock()
	go func() {
		firstChunk := alloc.NewBuffer().Slice(0, this.config.Cipher.IVSize())
		defer firstChunk.Release()

		writer, err := this.config.Cipher.NewEncodingStream(key, firstChunk.Value, conn)
		if err != nil {
			log.Error("Shadowsocks: Failed to create encoding stream: ", err)
			return
		}

		if payload, ok := <-ray.InboundOutput(); ok {
			firstChunk.Append(payload.Value)
			payload.Release()

			writer.Write(firstChunk.Value)
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
			config := rawConfig.(*Config)
			return &Shadowsocks{
				space:  space,
				config: config,
			}, nil
		})
}
