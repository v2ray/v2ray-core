package udp

import (
	"net"
	"sync"

	"v2ray.com/core/common/alloc"
	"v2ray.com/core/common/dice"
	"v2ray.com/core/common/log"
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/common/signal"
	"v2ray.com/core/proxy"
	"v2ray.com/core/transport/internet/internal"
)

type UDPPayload struct {
	payload *alloc.Buffer
	session *proxy.SessionInfo
}

type UDPPayloadHandler func(*alloc.Buffer, *proxy.SessionInfo)

type UDPPayloadQueue struct {
	queue    []chan UDPPayload
	callback UDPPayloadHandler
}

func NewUDPPayloadQueue(option ListenOption) *UDPPayloadQueue {
	queue := &UDPPayloadQueue{
		callback: option.Callback,
		queue:    make([]chan UDPPayload, option.Concurrency),
	}
	for i := range queue.queue {
		queue.queue[i] = make(chan UDPPayload, 64)
		go queue.Dequeue(queue.queue[i])
	}
	return queue
}

func (v *UDPPayloadQueue) Enqueue(payload UDPPayload) {
	size := len(v.queue)
	idx := 0
	if size > 1 {
		idx = dice.Roll(size)
	}
	for i := 0; i < size; i++ {
		select {
		case v.queue[idx%size] <- payload:
			return
		default:
			idx++
		}
	}
}

func (v *UDPPayloadQueue) Dequeue(queue <-chan UDPPayload) {
	for payload := range queue {
		v.callback(payload.payload, payload.session)
	}
}

func (v *UDPPayloadQueue) Close() {
	for _, queue := range v.queue {
		close(queue)
	}
}

type ListenOption struct {
	Callback            UDPPayloadHandler
	ReceiveOriginalDest bool
	Concurrency         int
}

type UDPHub struct {
	sync.RWMutex
	conn   *net.UDPConn
	cancel *signal.CancelSignal
	queue  *UDPPayloadQueue
	option ListenOption
}

func ListenUDP(address v2net.Address, port v2net.Port, option ListenOption) (*UDPHub, error) {
	if option.Concurrency < 1 {
		option.Concurrency = 1
	}
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
		queue:  NewUDPPayloadQueue(option),
		option: option,
		cancel: signal.NewCloseSignal(),
	}
	go hub.start()
	return hub, nil
}

func (v *UDPHub) Close() {
	v.Lock()
	defer v.Unlock()

	v.cancel.Cancel()
	v.conn.Close()
	v.cancel.WaitForDone()
	v.queue.Close()
}

func (v *UDPHub) WriteTo(payload []byte, dest v2net.Destination) (int, error) {
	return v.conn.WriteToUDP(payload, &net.UDPAddr{
		IP:   dest.Address.IP(),
		Port: int(dest.Port),
	})
}

func (v *UDPHub) start() {
	v.cancel.WaitThread()
	defer v.cancel.FinishThread()

	oobBytes := make([]byte, 256)
	for v.Running() {
		buffer := alloc.NewSmallBuffer()
		nBytes, noob, _, addr, err := ReadUDPMsg(v.conn, buffer.Bytes(), oobBytes)
		if err != nil {
			log.Info("UDP|Hub: Failed to read UDP msg: ", err)
			buffer.Release()
			continue
		}
		buffer.Slice(0, nBytes)

		session := new(proxy.SessionInfo)
		session.Source = v2net.UDPDestination(v2net.IPAddress(addr.IP), v2net.Port(addr.Port))
		if v.option.ReceiveOriginalDest && noob > 0 {
			session.Destination = RetrieveOriginalDest(oobBytes[:noob])
		}
		v.queue.Enqueue(UDPPayload{
			payload: buffer,
			session: session,
		})
	}
}

func (v *UDPHub) Running() bool {
	return !v.cancel.Cancelled()
}

// Connection return the net.Conn underneath v hub.
// Private: Visible for testing only
func (v *UDPHub) Connection() net.Conn {
	return v.conn
}

func (v *UDPHub) Addr() net.Addr {
	return v.conn.LocalAddr()
}
