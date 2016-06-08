package hub

import (
	"net"

	"github.com/v2ray/v2ray-core/common/alloc"
	v2net "github.com/v2ray/v2ray-core/common/net"
)

type UDPPayloadHandler func(*alloc.Buffer, v2net.Destination)

type UDPHub struct {
	conn      *net.UDPConn
	callback  UDPPayloadHandler
	accepting bool
}

func ListenUDP(address v2net.Address, port v2net.Port, callback UDPPayloadHandler) (*UDPHub, error) {
	udpConn, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   address.IP(),
		Port: int(port),
	})
	if err != nil {
		return nil, err
	}
	hub := &UDPHub{
		conn:     udpConn,
		callback: callback,
	}
	go hub.start()
	return hub, nil
}

func (this *UDPHub) Close() {
	this.accepting = false
	this.conn.Close()
}

func (this *UDPHub) WriteTo(payload []byte, dest v2net.Destination) (int, error) {
	return this.conn.WriteToUDP(payload, &net.UDPAddr{
		IP:   dest.Address().IP(),
		Port: int(dest.Port()),
	})
}

func (this *UDPHub) start() {
	this.accepting = true
	for this.accepting {
		buffer := alloc.NewBuffer()
		nBytes, addr, err := this.conn.ReadFromUDP(buffer.Value)
		if err != nil {
			buffer.Release()
			continue
		}
		buffer.Slice(0, nBytes)
		dest := v2net.UDPDestination(v2net.IPAddress(addr.IP), v2net.Port(addr.Port))
		go this.callback(buffer, dest)
	}
}
