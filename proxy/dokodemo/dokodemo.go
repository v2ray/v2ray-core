package dokodemo

import (
	"sync"

	"v2ray.com/core/app"
	"v2ray.com/core/app/dispatcher"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/errors"
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
			return errors.New("Dokodemo: Dispatcher is not found in the space.")
		}
		d.packetDispatcher = space.GetApp(dispatcher.APP_ID).(dispatcher.PacketDispatcher)
		return nil
	})
	return d
}

func (v *DokodemoDoor) Port() v2net.Port {
	return v.meta.Port
}

func (v *DokodemoDoor) Close() {
	v.accepting = false
	if v.tcpListener != nil {
		v.tcpMutex.Lock()
		v.tcpListener.Close()
		v.tcpListener = nil
		v.tcpMutex.Unlock()
	}
	if v.udpHub != nil {
		v.udpMutex.Lock()
		v.udpHub.Close()
		v.udpHub = nil
		v.udpMutex.Unlock()
	}
}

func (v *DokodemoDoor) Start() error {
	if v.accepting {
		return nil
	}
	v.accepting = true

	if v.config.NetworkList.HasNetwork(v2net.Network_TCP) {
		err := v.ListenTCP()
		if err != nil {
			return err
		}
	}
	if v.config.NetworkList.HasNetwork(v2net.Network_UDP) {
		err := v.ListenUDP()
		if err != nil {
			return err
		}
	}
	return nil
}

func (v *DokodemoDoor) ListenUDP() error {
	v.udpServer = udp.NewUDPServer(v.packetDispatcher)
	udpHub, err := udp.ListenUDP(
		v.meta.Address, v.meta.Port, udp.ListenOption{
			Callback:            v.handleUDPPackets,
			ReceiveOriginalDest: v.config.FollowRedirect,
			Concurrency:         2,
		})
	if err != nil {
		log.Error("Dokodemo failed to listen on ", v.meta.Address, ":", v.meta.Port, ": ", err)
		return err
	}
	v.udpMutex.Lock()
	v.udpHub = udpHub
	v.udpMutex.Unlock()
	return nil
}

func (v *DokodemoDoor) handleUDPPackets(payload *buf.Buffer, session *proxy.SessionInfo) {
	if session.Destination.Network == v2net.Network_Unknown && v.address != nil && v.port > 0 {
		session.Destination = v2net.UDPDestination(v.address, v.port)
	}
	if session.Destination.Network == v2net.Network_Unknown {
		log.Info("Dokodemo: Unknown destination, stop forwarding...")
		return
	}
	session.Inbound = v.meta
	v.udpServer.Dispatch(session, payload, v.handleUDPResponse)
}

func (v *DokodemoDoor) handleUDPResponse(dest v2net.Destination, payload *buf.Buffer) {
	defer payload.Release()
	v.udpMutex.RLock()
	defer v.udpMutex.RUnlock()
	if !v.accepting {
		return
	}
	v.udpHub.WriteTo(payload.Bytes(), dest)
}

func (v *DokodemoDoor) ListenTCP() error {
	tcpListener, err := internet.ListenTCP(v.meta.Address, v.meta.Port, v.HandleTCPConnection, v.meta.StreamSettings)
	if err != nil {
		log.Error("Dokodemo: Failed to listen on ", v.meta.Address, ":", v.meta.Port, ": ", err)
		return err
	}
	v.tcpMutex.Lock()
	v.tcpListener = tcpListener
	v.tcpMutex.Unlock()
	return nil
}

func (v *DokodemoDoor) HandleTCPConnection(conn internet.Connection) {
	defer conn.Close()

	var dest v2net.Destination
	if v.config.FollowRedirect {
		originalDest := GetOriginalDestination(conn)
		if originalDest.Network != v2net.Network_Unknown {
			log.Info("Dokodemo: Following redirect to: ", originalDest)
			dest = originalDest
		}
	}
	if dest.Network == v2net.Network_Unknown && v.address != nil && v.port > v2net.Port(0) {
		dest = v2net.TCPDestination(v.address, v.port)
	}

	if dest.Network == v2net.Network_Unknown {
		log.Info("Dokodemo: Unknown destination, stop forwarding...")
		return
	}
	log.Info("Dokodemo: Handling request to ", dest)

	ray := v.packetDispatcher.DispatchToOutbound(&proxy.SessionInfo{
		Source:      v2net.DestinationFromAddr(conn.RemoteAddr()),
		Destination: dest,
		Inbound:     v.meta,
	})
	defer ray.InboundOutput().Release()

	var wg sync.WaitGroup

	reader := v2net.NewTimeOutReader(v.config.Timeout, conn)
	defer reader.Release()

	wg.Add(1)
	go func() {
		v2reader := v2io.NewAdaptiveReader(reader)
		defer v2reader.Release()

		if err := v2io.PipeUntilEOF(v2reader, ray.InboundInput()); err != nil {
			log.Info("Dokodemo: Failed to transport all TCP request: ", err)
		}
		wg.Done()
		ray.InboundInput().Close()
	}()

	wg.Add(1)
	go func() {
		v2writer := v2io.NewAdaptiveWriter(conn)
		defer v2writer.Release()

		if err := v2io.PipeUntilEOF(ray.InboundOutput(), v2writer); err != nil {
			log.Info("Dokodemo: Failed to transport all TCP response: ", err)
		}
		wg.Done()
	}()

	wg.Wait()
}

type Factory struct{}

func (v *Factory) StreamCapability() v2net.NetworkList {
	return v2net.NetworkList{
		Network: []v2net.Network{v2net.Network_RawTCP},
	}
}

func (v *Factory) Create(space app.Space, rawConfig interface{}, meta *proxy.InboundHandlerMeta) (proxy.InboundHandler, error) {
	return NewDokodemoDoor(rawConfig.(*Config), space, meta), nil
}

func init() {
	registry.MustRegisterInboundHandlerCreator(loader.GetType(new(Config)), new(Factory))
}
