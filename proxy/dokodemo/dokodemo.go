package dokodemo

import (
	"io"
	"net"
	"sync"

	"github.com/v2ray/v2ray-core/app"
	"github.com/v2ray/v2ray-core/common/alloc"
	"github.com/v2ray/v2ray-core/common/log"
	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/proxy"
	"github.com/v2ray/v2ray-core/transport/listener"
)

type DokodemoDoor struct {
	tcpMutex      sync.RWMutex
	udpMutex      sync.RWMutex
	config        *Config
	accepting     bool
	address       v2net.Address
	port          v2net.Port
	space         app.Space
	tcpListener   *listener.TCPListener
	udpConn       *net.UDPConn
	listeningPort v2net.Port
}

func NewDokodemoDoor(space app.Space, config *Config) *DokodemoDoor {
	return &DokodemoDoor{
		config:  config,
		space:   space,
		address: config.Address,
		port:    config.Port,
	}
}

func (this *DokodemoDoor) Port() v2net.Port {
	return this.listeningPort
}

func (this *DokodemoDoor) Close() {
	this.accepting = false
	if this.tcpListener != nil {
		this.tcpMutex.Lock()
		this.tcpListener.Close()
		this.tcpListener = nil
		this.tcpMutex.Unlock()
	}
	if this.udpConn != nil {
		this.udpConn.Close()
		this.udpMutex.Lock()
		this.udpConn = nil
		this.udpMutex.Unlock()
	}
}

func (this *DokodemoDoor) Listen(port v2net.Port) error {
	if this.accepting {
		if this.listeningPort == port {
			return nil
		} else {
			return proxy.ErrorAlreadyListening
		}
	}
	this.listeningPort = port
	this.accepting = true

	if this.config.Network.HasNetwork(v2net.TCPNetwork) {
		err := this.ListenTCP(port)
		if err != nil {
			return err
		}
	}
	if this.config.Network.HasNetwork(v2net.UDPNetwork) {
		err := this.ListenUDP(port)
		if err != nil {
			return err
		}
	}
	return nil
}

func (this *DokodemoDoor) ListenUDP(port v2net.Port) error {
	udpConn, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   []byte{0, 0, 0, 0},
		Port: int(port),
		Zone: "",
	})
	if err != nil {
		log.Error("Dokodemo failed to listen on port ", port, ": ", err)
		return err
	}
	this.udpMutex.Lock()
	this.udpConn = udpConn
	this.udpMutex.Unlock()
	go this.handleUDPPackets()
	return nil
}

func (this *DokodemoDoor) handleUDPPackets() {
	for this.accepting {
		buffer := alloc.NewBuffer()
		this.udpMutex.RLock()
		if !this.accepting {
			this.udpMutex.RUnlock()
			return
		}
		nBytes, addr, err := this.udpConn.ReadFromUDP(buffer.Value)
		this.udpMutex.RUnlock()
		buffer.Slice(0, nBytes)
		if err != nil {
			buffer.Release()
			log.Error("Dokodemo failed to read from UDP: ", err)
			return
		}

		packet := v2net.NewPacket(v2net.UDPDestination(this.address, this.port), buffer, false)
		ray := this.space.PacketDispatcher().DispatchToOutbound(packet)
		close(ray.InboundInput())

		for payload := range ray.InboundOutput() {
			this.udpMutex.RLock()
			if !this.accepting {
				this.udpMutex.RUnlock()
				return
			}
			this.udpConn.WriteToUDP(payload.Value, addr)
			this.udpMutex.RUnlock()
		}
	}
}

func (this *DokodemoDoor) ListenTCP(port v2net.Port) error {
	tcpListener, err := listener.ListenTCP(port, this.HandleTCPConnection)
	if err != nil {
		log.Error("Dokodemo: Failed to listen on port ", port, ": ", err)
		return err
	}
	this.tcpMutex.Lock()
	this.tcpListener = tcpListener
	this.tcpMutex.Unlock()
	return nil
}

func (this *DokodemoDoor) HandleTCPConnection(conn *listener.TCPConn) {
	defer conn.Close()

	packet := v2net.NewPacket(v2net.TCPDestination(this.address, this.port), nil, true)
	ray := this.space.PacketDispatcher().DispatchToOutbound(packet)

	var inputFinish, outputFinish sync.Mutex
	inputFinish.Lock()
	outputFinish.Lock()

	reader := v2net.NewTimeOutReader(this.config.Timeout, conn)
	go dumpInput(reader, ray.InboundInput(), &inputFinish)
	go dumpOutput(conn, ray.InboundOutput(), &outputFinish)

	outputFinish.Lock()
}

func dumpInput(reader io.Reader, input chan<- *alloc.Buffer, finish *sync.Mutex) {
	v2net.ReaderToChan(input, reader)
	finish.Unlock()
	close(input)
}

func dumpOutput(writer io.Writer, output <-chan *alloc.Buffer, finish *sync.Mutex) {
	v2net.ChanToWriter(writer, output)
	finish.Unlock()
}
