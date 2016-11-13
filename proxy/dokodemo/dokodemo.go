package dokodemo

import (
	"sync"

	"v2ray.com/core/app"
	"v2ray.com/core/app/dispatcher"
	"v2ray.com/core/common/alloc"
	v2io "v2ray.com/core/common/io"
	"v2ray.com/core/common/loader"
	"v2ray.com/core/common/log"
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/proxy"
	"v2ray.com/core/proxy/registry"
	"v2ray.com/core/transport/internet"
	"v2ray.com/core/transport/internet/udp"
)

type DokodemoDoor struct {
	tcpMutex         sync.RWMutex
	udpMutex         sync.RWMutex
	config           *Config
	accepting        bool
	address          v2net.Address
	port             v2net.Port
	packetDispatcher dispatcher.PacketDispatcher
	tcpListener      *internet.TCPHub
	udpHub           *udp.UDPHub
	udpServer        *udp.UDPServer
	meta             *proxy.InboundHandlerMeta
}

func NewDokodemoDoor(config *Config, space app.Space, meta *proxy.InboundHandlerMeta) *DokodemoDoor {
	d := &DokodemoDoor{
		config:  config,
		address: config.GetPredefinedAddress(),
		port:    v2net.Port(config.Port),
		meta:    meta,
	}
	space.InitializeApplication(func() error {
		if !space.HasApp(dispatcher.APP_ID) {
			log.Error("Dokodemo: Dispatcher is not found in the space.")
			return app.ErrMissingApplication
		}
		d.packetDispatcher = space.GetApp(dispatcher.APP_ID).(dispatcher.PacketDispatcher)
		return nil
	})
	return d
}

func (this *DokodemoDoor) Port() v2net.Port {
	return this.meta.Port
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

func (this *DokodemoDoor) Start() error {
	if this.accepting {
		return nil
	}
	this.accepting = true

	if this.config.NetworkList.HasNetwork(v2net.Network_TCP) {
		err := this.ListenTCP()
		if err != nil {
			return err
		}
	}
	if this.config.NetworkList.HasNetwork(v2net.Network_UDP) {
		err := this.ListenUDP()
		if err != nil {
			return err
		}
	}
	return nil
}

func (this *DokodemoDoor) ListenUDP() error {
	this.udpServer = udp.NewUDPServer(this.packetDispatcher)
	udpHub, err := udp.ListenUDP(
		this.meta.Address, this.meta.Port, udp.ListenOption{
			Callback:            this.handleUDPPackets,
			ReceiveOriginalDest: this.config.FollowRedirect,
		})
	if err != nil {
		log.Error("Dokodemo failed to listen on ", this.meta.Address, ":", this.meta.Port, ": ", err)
		return err
	}
	this.udpMutex.Lock()
	this.udpHub = udpHub
	this.udpMutex.Unlock()
	return nil
}

func (this *DokodemoDoor) handleUDPPackets(payload *alloc.Buffer, session *proxy.SessionInfo) {
	if session.Destination.Network == v2net.Network_Unknown && this.address != nil && this.port > 0 {
		session.Destination = v2net.UDPDestination(this.address, this.port)
	}
	if session.Destination.Network == v2net.Network_Unknown {
		log.Info("Dokodemo: Unknown destination, stop forwarding...")
		return
	}
	session.Inbound = this.meta
	this.udpServer.Dispatch(session, payload, this.handleUDPResponse)
}

func (this *DokodemoDoor) handleUDPResponse(dest v2net.Destination, payload *alloc.Buffer) {
	defer payload.Release()
	this.udpMutex.RLock()
	defer this.udpMutex.RUnlock()
	if !this.accepting {
		return
	}
	this.udpHub.WriteTo(payload.Value, dest)
}

func (this *DokodemoDoor) ListenTCP() error {
	tcpListener, err := internet.ListenTCP(this.meta.Address, this.meta.Port, this.HandleTCPConnection, this.meta.StreamSettings)
	if err != nil {
		log.Error("Dokodemo: Failed to listen on ", this.meta.Address, ":", this.meta.Port, ": ", err)
		return err
	}
	this.tcpMutex.Lock()
	this.tcpListener = tcpListener
	this.tcpMutex.Unlock()
	return nil
}

func (this *DokodemoDoor) HandleTCPConnection(conn internet.Connection) {
	defer conn.Close()

	var dest v2net.Destination
	if this.config.FollowRedirect {
		originalDest := GetOriginalDestination(conn)
		if originalDest.Network != v2net.Network_Unknown {
			log.Info("Dokodemo: Following redirect to: ", originalDest)
			dest = originalDest
		}
	}
	if dest.Network == v2net.Network_Unknown && this.address != nil && this.port > v2net.Port(0) {
		dest = v2net.TCPDestination(this.address, this.port)
	}

	if dest.Network == v2net.Network_Unknown {
		log.Info("Dokodemo: Unknown destination, stop forwarding...")
		return
	}
	log.Info("Dokodemo: Handling request to ", dest)

	ray := this.packetDispatcher.DispatchToOutbound(&proxy.SessionInfo{
		Source:      v2net.DestinationFromAddr(conn.RemoteAddr()),
		Destination: dest,
		Inbound:     this.meta,
	})
	defer ray.InboundOutput().Release()

	var wg sync.WaitGroup

	reader := v2net.NewTimeOutReader(this.config.Timeout, conn)
	defer reader.Release()

	wg.Add(1)
	go func() {
		v2reader := v2io.NewAdaptiveReader(reader)
		defer v2reader.Release()

		v2io.Pipe(v2reader, ray.InboundInput())
		wg.Done()
		ray.InboundInput().Close()
	}()

	wg.Add(1)
	go func() {
		v2writer := v2io.NewAdaptiveWriter(conn)
		defer v2writer.Release()

		v2io.Pipe(ray.InboundOutput(), v2writer)
		wg.Done()
	}()

	wg.Wait()
}

type Factory struct{}

func (this *Factory) StreamCapability() v2net.NetworkList {
	return v2net.NetworkList{
		Network: []v2net.Network{v2net.Network_RawTCP},
	}
}

func (this *Factory) Create(space app.Space, rawConfig interface{}, meta *proxy.InboundHandlerMeta) (proxy.InboundHandler, error) {
	return NewDokodemoDoor(rawConfig.(*Config), space, meta), nil
}

func init() {
	registry.MustRegisterInboundHandlerCreator(loader.GetType(new(Config)), new(Factory))
}
