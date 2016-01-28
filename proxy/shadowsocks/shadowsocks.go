// R.I.P Shadowsocks

package shadowsocks

import (
	"sync"

	"github.com/v2ray/v2ray-core/app"
	"github.com/v2ray/v2ray-core/common/alloc"
	"github.com/v2ray/v2ray-core/common/log"
	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/proxy"
	"github.com/v2ray/v2ray-core/proxy/internal"
	"github.com/v2ray/v2ray-core/transport/listener"
)

type Shadowsocks struct {
	space       app.Space
	config      *Config
	port        v2net.Port
	accepting   bool
	tcpListener *listener.TCPListener
}

func (this *Shadowsocks) Port() v2net.Port {
	return this.port
}

func (this *Shadowsocks) Close() {
	this.accepting = false
	this.tcpListener.Close()
	this.tcpListener = nil
}

func (this *Shadowsocks) Listen(port v2net.Port) error {
	if this.accepting {
		if this.port == port {
			return nil
		} else {
			return proxy.ErrorAlreadyListening
		}
	}

	tcpListener, err := listener.ListenTCP(port, this.handleConnection)
	if err != nil {
		log.Error("Shadowsocks: Failed to listen on port ", port, ": ", err)
		return err
	}
	this.tcpListener = tcpListener
	this.accepting = true
	return nil
}

func (this *Shadowsocks) handleConnection(conn *listener.TCPConn) {
	defer conn.Close()

	buffer := alloc.NewSmallBuffer()
	defer buffer.Release()

	_, err := v2net.ReadAllBytes(conn, buffer.Value[:this.config.Cipher.IVSize()])
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

	request, err := ReadRequest(reader)
	if err != nil {
		return
	}

	packet := v2net.NewPacket(v2net.TCPDestination(request.Address, request.Port), nil, true)
	ray := this.space.PacketDispatcher().DispatchToOutbound(packet)

	writer, err := this.config.Cipher.NewEncodingStream(key, iv, conn)
	if err != nil {
		log.Error("Shadowsocks: Failed to create encoding stream: ", err)
		return
	}

	var writeFinish sync.Mutex
	writeFinish.Lock()
	go func() {
		v2net.ChanToWriter(writer, ray.InboundOutput())
		writeFinish.Unlock()
	}()

	v2net.ReaderToChan(ray.InboundInput(), reader)
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
