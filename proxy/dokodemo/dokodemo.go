package dokodemo

import (
	"io"
	"net"
	"sync"

	"github.com/v2ray/v2ray-core/app"
	"github.com/v2ray/v2ray-core/common/alloc"
	"github.com/v2ray/v2ray-core/common/log"
	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/common/retry"
	"github.com/v2ray/v2ray-core/proxy/dokodemo/config/json"
)

type DokodemoDoor struct {
	config    *json.DokodemoConfig
	accepting bool
	address   v2net.Address
	space     *app.Space
}

func NewDokodemoDoor(space *app.Space, config *json.DokodemoConfig) *DokodemoDoor {
	d := &DokodemoDoor{
		config: config,
		space:  space,
	}
	ip := net.ParseIP(config.Host)
	if ip != nil {
		d.address = v2net.IPAddress(ip, config.Port)
	} else {
		d.address = v2net.DomainAddress(config.Host, config.Port)
	}
	return d
}

func (this *DokodemoDoor) Listen(port v2net.Port) error {
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
		log.Error("Dokodemo failed to listen on port %d: %v", port, err)
		return err
	}
	go this.handleUDPPackets(udpConn)
	return nil
}

func (this *DokodemoDoor) handleUDPPackets(udpConn *net.UDPConn) {
	defer udpConn.Close()
	for this.accepting {
		buffer := alloc.NewBuffer()
		nBytes, addr, err := udpConn.ReadFromUDP(buffer.Value)
		buffer.Slice(0, nBytes)
		if err != nil {
			buffer.Release()
			log.Error("Dokodemo failed to read from UDP: %v", err)
			return
		}

		packet := v2net.NewPacket(v2net.NewUDPDestination(this.address), buffer, false)
		ray := this.space.PacketDispatcher().DispatchToOutbound(packet)
		close(ray.InboundInput())

		for payload := range ray.InboundOutput() {
			udpConn.WriteToUDP(payload.Value, addr)
		}
	}
}

func (this *DokodemoDoor) ListenTCP(port v2net.Port) error {
	tcpListener, err := net.ListenTCP("tcp", &net.TCPAddr{
		IP:   []byte{0, 0, 0, 0},
		Port: int(port),
		Zone: "",
	})
	if err != nil {
		log.Error("Dokodemo failed to listen on port %d: %v", port, err)
		return err
	}
	go this.AcceptTCPConnections(tcpListener)
	return nil
}

func (this *DokodemoDoor) AcceptTCPConnections(tcpListener *net.TCPListener) {
	for this.accepting {
		retry.Timed(100, 100).On(func() error {
			connection, err := tcpListener.AcceptTCP()
			if err != nil {
				log.Error("Dokodemo failed to accept new connections: %v", err)
				return err
			}
			go this.HandleTCPConnection(connection)
			return nil
		})
	}
}

func (this *DokodemoDoor) HandleTCPConnection(conn *net.TCPConn) {
	defer conn.Close()

	packet := v2net.NewPacket(v2net.NewTCPDestination(this.address), nil, true)
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
