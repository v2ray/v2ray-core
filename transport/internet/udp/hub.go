package udp

import (
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/net"
)

// Payload represents a single UDP payload.
type Payload struct {
	payload      *buf.Buffer
	source       net.Destination
	originalDest net.Destination
}

// PayloadHandler is function to handle Payload.
type PayloadHandler func(payload *buf.Buffer, source net.Destination, originalDest net.Destination)

type HubOption func(h *Hub)

func HubCapacity(cap int) HubOption {
	return func(h *Hub) {
		h.capacity = cap
	}
}

func HubReceiveOriginalDestination(r bool) HubOption {
	return func(h *Hub) {
		h.recvOrigDest = r
	}
}

type Hub struct {
	conn         *net.UDPConn
	callback     PayloadHandler
	capacity     int
	recvOrigDest bool
}

func ListenUDP(address net.Address, port net.Port, callback PayloadHandler, options ...HubOption) (*Hub, error) {
	udpConn, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   address.IP(),
		Port: int(port),
	})
	if err != nil {
		return nil, err
	}
	newError("listening UDP on ", address, ":", port).WriteToLog()
	hub := &Hub{
		conn:         udpConn,
		capacity:     256,
		callback:     callback,
		recvOrigDest: false,
	}
	for _, opt := range options {
		opt(hub)
	}

	if hub.recvOrigDest {
		rawConn, err := udpConn.SyscallConn()
		if err != nil {
			return nil, newError("failed to get fd").Base(err)
		}
		err = rawConn.Control(func(fd uintptr) {
			if err := SetOriginalDestOptions(int(fd)); err != nil {
				newError("failed to set socket options").Base(err).WriteToLog()
			}
		})
		if err != nil {
			return nil, newError("failed to control socket").Base(err)
		}
	}

	c := make(chan *Payload, hub.capacity)

	go hub.start(c)
	go hub.process(c)
	return hub, nil
}

// Close implements net.Listener.
func (h *Hub) Close() error {
	h.conn.Close()
	return nil
}

func (h *Hub) WriteTo(payload []byte, dest net.Destination) (int, error) {
	return h.conn.WriteToUDP(payload, &net.UDPAddr{
		IP:   dest.Address.IP(),
		Port: int(dest.Port),
	})
}

func (h *Hub) process(c <-chan *Payload) {
	for p := range c {
		h.callback(p.payload, p.source, p.originalDest)
	}
}

func (h *Hub) start(c chan<- *Payload) {
	defer close(c)

	oobBytes := make([]byte, 256)

	for {
		buffer := buf.New()
		var noob int
		var addr *net.UDPAddr
		err := buffer.AppendSupplier(func(b []byte) (int, error) {
			n, nb, _, a, e := ReadUDPMsg(h.conn, b, oobBytes)
			noob = nb
			addr = a
			return n, e
		})

		if err != nil {
			newError("failed to read UDP msg").Base(err).WriteToLog()
			buffer.Release()
			break
		}

		payload := &Payload{
			payload: buffer,
		}
		payload.source = net.UDPDestination(net.IPAddress(addr.IP), net.Port(addr.Port))
		if h.recvOrigDest && noob > 0 {
			payload.originalDest = RetrieveOriginalDest(oobBytes[:noob])
			if payload.originalDest.IsValid() {
				newError("UDP original destination: ", payload.originalDest).AtDebug().WriteToLog()
			} else {
				newError("failed to read UDP original destination").WriteToLog()
			}
		}

		select {
		case c <- payload:
		default:
		}

	}
}

// Addr implements net.Listener.
func (h *Hub) Addr() net.Addr {
	return h.conn.LocalAddr()
}
