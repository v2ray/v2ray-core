package socks

import (
	"io"
	"sync"
	"time"

	"v2ray.com/core/app"
	"v2ray.com/core/app/dispatcher"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/bufio"
	"v2ray.com/core/common/errors"
	"v2ray.com/core/common/log"
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/common/protocol"
	"v2ray.com/core/common/signal"
	"v2ray.com/core/proxy"
	"v2ray.com/core/transport/internet"
	"v2ray.com/core/transport/internet/udp"
)

// Server is a SOCKS 5 proxy server
type Server struct {
	tcpMutex         sync.RWMutex
	udpMutex         sync.RWMutex
	accepting        bool
	packetDispatcher dispatcher.PacketDispatcher
	config           *ServerConfig
	tcpListener      *internet.TCPHub
	udpHub           *udp.Hub
	udpAddress       v2net.Destination
	udpServer        *udp.Server
	meta             *proxy.InboundHandlerMeta
}

// NewServer creates a new Server object.
func NewServer(config *ServerConfig, space app.Space, meta *proxy.InboundHandlerMeta) *Server {
	s := &Server{
		config: config,
		meta:   meta,
	}
	space.OnInitialize(func() error {
		s.packetDispatcher = dispatcher.FromSpace(space)
		if s.packetDispatcher == nil {
			return errors.New("Socks|Server: Dispatcher is not found in the space.")
		}
		return nil
	})
	return s
}

// Port implements InboundHandler.Port().
func (v *Server) Port() v2net.Port {
	return v.meta.Port
}

// Close implements InboundHandler.Close().
func (v *Server) Close() {
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

// Start implements InboundHandler.Start().
func (v *Server) Start() error {
	if v.accepting {
		return nil
	}

	listener, err := internet.ListenTCP(
		v.meta.Address,
		v.meta.Port,
		v.handleConnection,
		v.meta.StreamSettings)
	if err != nil {
		log.Error("Socks: failed to listen on ", v.meta.Address, ":", v.meta.Port, ": ", err)
		return err
	}
	v.accepting = true
	v.tcpMutex.Lock()
	v.tcpListener = listener
	v.tcpMutex.Unlock()
	if v.config.UdpEnabled {
		v.listenUDP()
	}
	return nil
}

func (v *Server) handleConnection(connection internet.Connection) {
	defer connection.Close()

	connection.SetReusable(false)

	timedReader := v2net.NewTimeOutReader(v.config.Timeout, connection)
	reader := bufio.NewReader(timedReader)

	session := &ServerSession{
		config: v.config,
		meta:   v.meta,
	}

	clientAddr := v2net.DestinationFromAddr(connection.RemoteAddr())

	request, err := session.Handshake(reader, connection)
	if err != nil {
		log.Access(clientAddr, "", log.AccessRejected, err)
		log.Info("Socks|Server: Failed to read request: ", err)
		return
	}

	if request.Command == protocol.RequestCommandTCP {
		dest := request.Destination()
		session := &proxy.SessionInfo{
			Source:      clientAddr,
			Destination: dest,
			Inbound:     v.meta,
		}
		log.Info("Socks|Server: TCP Connect request to ", dest)
		log.Access(clientAddr, dest, log.AccessAccepted, "")

		v.transport(reader, connection, session)
		return
	}

	if request.Command == protocol.RequestCommandUDP {
		v.handleUDP()
		return
	}
}

func (v *Server) handleUDP() error {
	// The TCP connection closes after v method returns. We need to wait until
	// the client closes it.
	// TODO: get notified from UDP part
	<-time.After(5 * time.Minute)

	return nil
}

func (v *Server) transport(reader io.Reader, writer io.Writer, session *proxy.SessionInfo) {
	ray := v.packetDispatcher.DispatchToOutbound(session)
	input := ray.InboundInput()
	output := ray.InboundOutput()

	requestDone := signal.ExecuteAsync(func() error {
		defer input.Close()

		v2reader := buf.NewReader(reader)
		if err := buf.PipeUntilEOF(v2reader, input); err != nil {
			log.Info("Socks|Server: Failed to transport all TCP request: ", err)
			return err
		}
		return nil
	})

	responseDone := signal.ExecuteAsync(func() error {
		defer output.ForceClose()

		v2writer := buf.NewWriter(writer)
		if err := buf.PipeUntilEOF(output, v2writer); err != nil {
			log.Info("Socks|Server: Failed to transport all TCP response: ", err)
			return err
		}
		return nil

	})

	if err := signal.ErrorOrFinish2(requestDone, responseDone); err != nil {
		log.Info("Socks|Server: Connection ends with ", err)
	}
}

type ServerFactory struct{}

func (v *ServerFactory) StreamCapability() v2net.NetworkList {
	return v2net.NetworkList{
		Network: []v2net.Network{v2net.Network_TCP},
	}
}

func (v *ServerFactory) Create(space app.Space, rawConfig interface{}, meta *proxy.InboundHandlerMeta) (proxy.InboundHandler, error) {
	return NewServer(rawConfig.(*ServerConfig), space, meta), nil
}
