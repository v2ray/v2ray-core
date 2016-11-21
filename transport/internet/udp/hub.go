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

func (this *UDPPayloadQueue) Enqueue(payload UDPPayload) {
	size := len(this.queue)
	idx := 0
	if size > 1 {
		idx = dice.Roll(size)
	}
	for i := 0; i < size; i++ {
		select {
		case this.queue[idx%size] <- payload:
			return
		default:
			idx++
		}
	}
}

func (this *UDPPayloadQueue) Dequeue(queue <-chan UDPPayload) {
	for payload := range queue {
		this.callback(payload.payload, payload.session)
	}
}

func (this *UDPPayloadQueue) Close() {
	for _, queue := range this.queue {
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

func (this *UDPHub) Close() {
	this.Lock()
	defer this.Unlock()

	this.cancel.Cancel()
	this.conn.Close()
	this.cancel.WaitForDone()
	this.queue.Close()
}

func (this *UDPHub) WriteTo(payload []byte, dest v2net.Destination) (int, error) {
	return this.conn.WriteToUDP(payload, &net.UDPAddr{
		IP:   dest.Address.IP(),
		Port: int(dest.Port),
	})
}

func (this *UDPHub) start() {
	this.cancel.WaitThread()
	defer this.cancel.FinishThread()

	oobBytes := make([]byte, 256)
	for this.Running() {
		buffer := alloc.NewSmallBuffer()
		nBytes, noob, _, addr, err := ReadUDPMsg(this.conn, buffer.Value, oobBytes)
		if err != nil {
			log.Info("UDP|Hub: Failed to read UDP msg: ", err)
			buffer.Release()
			continue
		}
		buffer.Slice(0, nBytes)

		session := new(proxy.SessionInfo)
		session.Source = v2net.UDPDestination(v2net.IPAddress(addr.IP), v2net.Port(addr.Port))
		if this.option.ReceiveOriginalDest && noob > 0 {
			session.Destination = RetrieveOriginalDest(oobBytes[:noob])
		}
		this.queue.Enqueue(UDPPayload{
			payload: buffer,
			session: session,
		})
	}
}

func (this *UDPHub) Running() bool {
	return !this.cancel.Cancelled()
}

// Connection return the net.Conn underneath this hub.
// Private: Visible for testing only
func (this *UDPHub) Connection() net.Conn {
	return this.conn
}

func (this *UDPHub) Addr() net.Addr {
	return this.conn.LocalAddr()
}
