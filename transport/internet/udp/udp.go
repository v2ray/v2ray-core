package udp

import (
	"net"
	"sync"

	"github.com/v2ray/v2ray-core/common/alloc"
	"github.com/v2ray/v2ray-core/common/log"
	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/proxy"
	"github.com/v2ray/v2ray-core/transport/internet/internal"
)

type UDPPayloadHandler func(*alloc.Buffer, *proxy.SessionInfo)

type UDPHub struct {
	sync.RWMutex
	conn      *net.UDPConn
	option    ListenOption
	accepting bool
}

type ListenOption struct {
	Callback            UDPPayloadHandler
	ReceiveOriginalDest bool
}

func ListenUDP(address v2net.Address, port v2net.Port, option ListenOption) (*UDPHub, error) {
	udpConn, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   address.IP(),
		Port: int(port),
	})
	if err != nil {
		return nil, err
	}
	if option.ReceiveOriginalDest {
		fd, err := internal.GetSysFd(udpConn)
		if err != nil {
			log.Warning("UDP|Listener: Failed to get fd: ", err)
			return nil, err
		}
		err = SetOriginalDestOptions(fd)
		if err != nil {
			log.Warning("UDP|Listener: Failed to set socket options: ", err)
			return nil, err
		}
	}
	hub := &UDPHub{
		conn:   udpConn,
		option: option,
	}
	go hub.start()
	return hub, nil
}

func (this *UDPHub) Close() {
	this.Lock()
	defer this.Unlock()

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
	this.Lock()
	this.accepting = true
	this.Unlock()

	oobBytes := make([]byte, 256)
	for this.Running() {
		buffer := alloc.NewBuffer()
		nBytes, noob, _, addr, err := this.conn.ReadMsgUDP(buffer.Value, oobBytes)
		if err != nil {
			buffer.Release()
			continue
		}
		buffer.Slice(0, nBytes)

		session := new(proxy.SessionInfo)
		session.Source = v2net.UDPDestination(v2net.IPAddress(addr.IP), v2net.Port(addr.Port))
		if this.option.ReceiveOriginalDest && noob > 0 {
			session.Destination = RetrieveOriginalDest(oobBytes[:noob])
		}
		go this.option.Callback(buffer, session)
	}
}

func (this *UDPHub) Running() bool {
	this.RLock()
	defer this.RUnlock()

	return this.accepting
}
