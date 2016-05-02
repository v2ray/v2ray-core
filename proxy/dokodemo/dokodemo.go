package dokodemo

import (
	"sync"

	"github.com/v2ray/v2ray-core/app/dispatcher"
	"github.com/v2ray/v2ray-core/common/alloc"
	v2io "github.com/v2ray/v2ray-core/common/io"
	"github.com/v2ray/v2ray-core/common/log"
	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/proxy"
	"github.com/v2ray/v2ray-core/transport/hub"
)

type DokodemoDoor struct {
	tcpMutex         sync.RWMutex
	udpMutex         sync.RWMutex
	config           *Config
	accepting        bool
	address          v2net.Address
	port             v2net.Port
	packetDispatcher dispatcher.PacketDispatcher
	tcpListener      *hub.TCPHub
	udpHub           *hub.UDPHub
	udpServer        *hub.UDPServer
	listeningPort    v2net.Port
}

func NewDokodemoDoor(config *Config, packetDispatcher dispatcher.PacketDispatcher) *DokodemoDoor {
	return &DokodemoDoor{
		config:           config,
		packetDispatcher: packetDispatcher,
		address:          config.Address,
		port:             config.Port,
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
	if this.udpHub != nil {
		this.udpMutex.Lock()
		this.udpHub.Close()
		this.udpHub = nil
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
	this.udpServer = hub.NewUDPServer(this.packetDispatcher)
	udpHub, err := hub.ListenUDP(port, this.handleUDPPackets)
	if err != nil {
		log.Error("Dokodemo failed to listen on port ", port, ": ", err)
		return err
	}
	this.udpMutex.Lock()
	this.udpHub = udpHub
	this.udpMutex.Unlock()
	return nil
}

func (this *DokodemoDoor) handleUDPPackets(payload *alloc.Buffer, dest v2net.Destination) {
	this.udpServer.Dispatch(dest, v2net.UDPDestination(this.address, this.port), payload, func(destination v2net.Destination, payload *alloc.Buffer) {
		defer payload.Release()
		this.udpMutex.RLock()
		defer this.udpMutex.RUnlock()
		if !this.accepting {
			return
		}
		this.udpHub.WriteTo(payload.Value, destination)
	})
}

func (this *DokodemoDoor) ListenTCP(port v2net.Port) error {
	tcpListener, err := hub.ListenTCP(port, this.HandleTCPConnection, nil)
	if err != nil {
		log.Error("Dokodemo: Failed to listen on port ", port, ": ", err)
		return err
	}
	this.tcpMutex.Lock()
	this.tcpListener = tcpListener
	this.tcpMutex.Unlock()
	return nil
}

func (this *DokodemoDoor) HandleTCPConnection(conn *hub.Connection) {
	defer conn.Close()

	ray := this.packetDispatcher.DispatchToOutbound(v2net.TCPDestination(this.address, this.port))
	defer ray.InboundOutput().Release()

	var inputFinish, outputFinish sync.Mutex
	inputFinish.Lock()
	outputFinish.Lock()

	reader := v2net.NewTimeOutReader(this.config.Timeout, conn)
	go func() {
		v2reader := v2io.NewAdaptiveReader(reader)
		defer v2reader.Release()

		v2io.Pipe(v2reader, ray.InboundInput())
		inputFinish.Unlock()
		ray.InboundInput().Close()
	}()

	go func() {
		v2writer := v2io.NewAdaptiveWriter(conn)
		defer v2writer.Release()

		v2io.Pipe(ray.InboundOutput(), v2writer)
		outputFinish.Unlock()
	}()

	outputFinish.Lock()
}
