package kcp

import (
	"errors"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/v2ray/v2ray-core/common/alloc"
	"github.com/v2ray/v2ray-core/common/log"
)

var (
	errTimeout          = errors.New("i/o timeout")
	errBrokenPipe       = errors.New("broken pipe")
	errClosedListener   = errors.New("Listener closed.")
	errClosedConnection = errors.New("Connection closed.")
)

type State int32

const (
	StateActive       State = 0
	StateReadyToClose State = 1
	StatePeerClosed   State = 2
	StateTerminating  State = 3
	StateTerminated   State = 4
)

const (
	headerSize uint32 = 2
)

func nowMillisec() int64 {
	now := time.Now()
	return now.Unix()*1000 + int64(now.Nanosecond()/1000000)
}

type RountTripInfo struct {
	sync.RWMutex
	variation uint32
	srtt      uint32
	rto       uint32
	minRtt    uint32
}

func (this *RountTripInfo) Update(rtt uint32) {
	if rtt > 0x7FFFFFFF {
		return
	}
	this.Lock()
	defer this.Unlock()

	// https://tools.ietf.org/html/rfc6298
	if this.srtt == 0 {
		this.srtt = rtt
		this.variation = rtt / 2
	} else {
		delta := rtt - this.srtt
		if this.srtt > rtt {
			delta = this.srtt - rtt
		}
		this.variation = (3*this.variation + delta) / 4
		this.srtt = (7*this.srtt + rtt) / 8
		if this.srtt < this.minRtt {
			this.srtt = this.minRtt
		}
	}
	var rto uint32
	if this.minRtt < 4*this.variation {
		rto = this.srtt + 4*this.variation
	} else {
		rto = this.srtt + this.variation
	}

	if rto > 10000 {
		rto = 10000
	}
	this.rto = rto * 3 / 2
}

func (this *RountTripInfo) Timeout() uint32 {
	this.RLock()
	defer this.RUnlock()

	return this.rto
}

func (this *RountTripInfo) SmoothedTime() uint32 {
	this.RLock()
	defer this.RUnlock()

	return this.srtt
}

// Connection is a KCP connection over UDP.
type Connection struct {
	block         Authenticator
	local, remote net.Addr
	rd            time.Time
	wd            time.Time // write deadline
	writer        io.WriteCloser
	since         int64
	dataInputCond *sync.Cond

	conv             uint16
	state            State
	stateBeginTime   uint32
	lastIncomingTime uint32
	lastPayloadTime  uint32
	sendingUpdated   bool
	lastPingTime     uint32

	mss       uint32
	roundTrip *RountTripInfo
	interval  uint32

	receivingWorker *ReceivingWorker
	sendingWorker   *SendingWorker

	fastresend        uint32
	congestionControl bool
	output            *BufferedSegmentWriter
}

// NewConnection create a new KCP connection between local and remote.
func NewConnection(conv uint16, writerCloser io.WriteCloser, local *net.UDPAddr, remote *net.UDPAddr, block Authenticator) *Connection {
	log.Info("KCP|Connection: creating connection ", conv)

	conn := new(Connection)
	conn.local = local
	conn.remote = remote
	conn.block = block
	conn.writer = writerCloser
	conn.since = nowMillisec()
	conn.dataInputCond = sync.NewCond(new(sync.Mutex))

	authWriter := &AuthenticationWriter{
		Authenticator: block,
		Writer:        writerCloser,
	}
	conn.conv = conv
	conn.output = NewSegmentWriter(authWriter)

	conn.mss = authWriter.Mtu() - DataSegmentOverhead
	conn.roundTrip = &RountTripInfo{
		rto:    100,
		minRtt: effectiveConfig.Tti,
	}
	conn.interval = effectiveConfig.Tti
	conn.receivingWorker = NewReceivingWorker(conn)
	conn.fastresend = 2
	conn.congestionControl = effectiveConfig.Congestion
	conn.sendingWorker = NewSendingWorker(conn)

	go conn.updateTask()

	return conn
}

func (this *Connection) Elapsed() uint32 {
	return uint32(nowMillisec() - this.since)
}

// Read implements the Conn Read method.
func (this *Connection) Read(b []byte) (int, error) {
	if this == nil {
		return 0, io.EOF
	}

	for {
		if this.State() == StateReadyToClose || this.State() == StateTerminating || this.State() == StateTerminated {
			return 0, io.EOF
		}
		nBytes := this.receivingWorker.Read(b)
		if nBytes > 0 {
			return nBytes, nil
		}
		var timer *time.Timer
		if !this.rd.IsZero() && this.rd.Before(time.Now()) {
			timer = time.AfterFunc(this.rd.Sub(time.Now()), this.dataInputCond.Signal)
		}
		this.dataInputCond.L.Lock()
		this.dataInputCond.Wait()
		this.dataInputCond.L.Unlock()
		if timer != nil {
			timer.Stop()
		}
		if !this.rd.IsZero() && this.rd.Before(time.Now()) {
			return 0, errTimeout
		}
	}
}

// Write implements the Conn Write method.
func (this *Connection) Write(b []byte) (int, error) {
	if this == nil || this.State() != StateActive {
		return 0, io.ErrClosedPipe
	}
	totalWritten := 0

	for {
		if this == nil || this.State() != StateActive {
			return totalWritten, io.ErrClosedPipe
		}

		nBytes := this.sendingWorker.Push(b[totalWritten:])
		if nBytes > 0 {
			totalWritten += nBytes
			if totalWritten == len(b) {
				return totalWritten, nil
			}
		}

		if !this.wd.IsZero() && this.wd.Before(time.Now()) {
			return totalWritten, errTimeout
		}

		// Sending windows is 1024 for the moment. This amount is not gonna sent in 1 sec.
		time.Sleep(time.Second)
	}
}

func (this *Connection) SetState(state State) {
	current := this.Elapsed()
	atomic.StoreInt32((*int32)(&this.state), int32(state))
	atomic.StoreUint32(&this.stateBeginTime, current)
	log.Info("KCP|Connection: Entering state ", state, " at ", current)

	switch state {
	case StateReadyToClose:
		this.receivingWorker.CloseRead()
	case StatePeerClosed:
		this.sendingWorker.CloseWrite()
	case StateTerminating:
		this.receivingWorker.CloseRead()
		this.sendingWorker.CloseWrite()
	case StateTerminated:
		this.receivingWorker.CloseRead()
		this.sendingWorker.CloseWrite()
	}
}

// Close closes the connection.
func (this *Connection) Close() error {
	if this == nil {
		return errClosedConnection
	}

	this.dataInputCond.Broadcast()

	state := this.State()
	if state == StateReadyToClose ||
		state == StateTerminating ||
		state == StateTerminated {
		return errClosedConnection
	}
	log.Info("KCP|Connection: Closing connection to ", this.remote)

	if state == StateActive {
		this.SetState(StateReadyToClose)
	}
	if state == StatePeerClosed {
		this.SetState(StateTerminating)
	}

	return nil
}

// LocalAddr returns the local network address. The Addr returned is shared by all invocations of LocalAddr, so do not modify it.
func (this *Connection) LocalAddr() net.Addr {
	if this == nil {
		return nil
	}
	return this.local
}

// RemoteAddr returns the remote network address. The Addr returned is shared by all invocations of RemoteAddr, so do not modify it.
func (this *Connection) RemoteAddr() net.Addr {
	if this == nil {
		return nil
	}
	return this.remote
}

// SetDeadline sets the deadline associated with the listener. A zero time value disables the deadline.
func (this *Connection) SetDeadline(t time.Time) error {
	if err := this.SetReadDeadline(t); err != nil {
		return err
	}
	if err := this.SetWriteDeadline(t); err != nil {
		return err
	}
	return nil
}

// SetReadDeadline implements the Conn SetReadDeadline method.
func (this *Connection) SetReadDeadline(t time.Time) error {
	if this == nil || this.State() != StateActive {
		return errClosedConnection
	}
	this.rd = t
	return nil
}

// SetWriteDeadline implements the Conn SetWriteDeadline method.
func (this *Connection) SetWriteDeadline(t time.Time) error {
	if this == nil || this.State() != StateActive {
		return errClosedConnection
	}
	this.wd = t
	return nil
}

// kcp update, input loop
func (this *Connection) updateTask() {
	for this.State() != StateTerminated {
		this.flush()

		interval := time.Duration(effectiveConfig.Tti) * time.Millisecond
		if this.State() == StateTerminating {
			interval = time.Second
		}
		time.Sleep(interval)
	}
	this.Terminate()
}

func (this *Connection) FetchInputFrom(conn net.Conn) {
	go func() {
		for {
			payload := alloc.NewBuffer()
			nBytes, err := conn.Read(payload.Value)
			if err != nil {
				payload.Release()
				return
			}
			payload.Slice(0, nBytes)
			if this.block.Open(payload) {
				this.Input(payload.Value)
			} else {
				log.Info("KCP|Connection: Invalid response from ", conn.RemoteAddr())
			}
			payload.Release()
		}
	}()
}

func (this *Connection) Reusable() bool {
	return false
}

func (this *Connection) SetReusable(b bool) {}

func (this *Connection) Terminate() {
	if this == nil || this.writer == nil {
		return
	}
	log.Info("KCP|Connection: Terminating connection to ", this.RemoteAddr())

	this.writer.Close()
}

func (this *Connection) HandleOption(opt SegmentOption) {
	if (opt & SegmentOptionClose) == SegmentOptionClose {
		this.OnPeerClosed()
	}
}

func (this *Connection) OnPeerClosed() {
	state := this.State()
	if state == StateReadyToClose {
		this.SetState(StateTerminating)
	}
	if state == StateActive {
		this.SetState(StatePeerClosed)
	}
}

// Input when you received a low level packet (eg. UDP packet), call it
func (this *Connection) Input(data []byte) int {
	current := this.Elapsed()
	atomic.StoreUint32(&this.lastIncomingTime, current)

	var seg Segment
	for {
		seg, data = ReadSegment(data)
		if seg == nil {
			break
		}

		switch seg := seg.(type) {
		case *DataSegment:
			this.HandleOption(seg.Opt)
			this.receivingWorker.ProcessSegment(seg)
			atomic.StoreUint32(&this.lastPayloadTime, current)
			this.dataInputCond.Signal()
		case *AckSegment:
			this.HandleOption(seg.Opt)
			this.sendingWorker.ProcessSegment(current, seg)
			atomic.StoreUint32(&this.lastPayloadTime, current)
		case *CmdOnlySegment:
			this.HandleOption(seg.Opt)
			if seg.Cmd == SegmentCommandTerminated {
				state := this.State()
				if state == StateActive ||
					state == StateReadyToClose ||
					state == StatePeerClosed {
					this.SetState(StateTerminating)
				} else if state == StateTerminating {
					this.SetState(StateTerminated)
				}
			}
			this.sendingWorker.ProcessReceivingNext(seg.ReceivinNext)
			this.receivingWorker.ProcessSendingNext(seg.SendingNext)
		default:
		}
	}

	return 0
}

func (this *Connection) flush() {
	current := this.Elapsed()

	if this.State() == StateTerminated {
		return
	}
	if this.State() == StateActive && current-this.lastPayloadTime >= 30000 {
		this.Close()
	}

	if this.State() == StateTerminating {
		this.output.Write(&CmdOnlySegment{
			Conv: this.conv,
			Cmd:  SegmentCommandTerminated,
		})
		this.output.Flush()

		if current-this.stateBeginTime > 8000 {
			this.SetState(StateTerminated)
		}
		return
	}

	if this.State() == StateReadyToClose && current-this.stateBeginTime > 15000 {
		this.SetState(StateTerminating)
	}

	// flush acknowledges
	this.receivingWorker.Flush(current)
	this.sendingWorker.Flush(current)

	if this.sendingWorker.PingNecessary() || this.receivingWorker.PingNecessary() || current-this.lastPingTime >= 5000 {
		seg := NewCmdOnlySegment()
		seg.Conv = this.conv
		seg.Cmd = SegmentCommandPing
		seg.ReceivinNext = this.receivingWorker.nextNumber
		seg.SendingNext = this.sendingWorker.firstUnacknowledged
		if this.State() == StateReadyToClose {
			seg.Opt = SegmentOptionClose
		}
		this.output.Write(seg)
		this.lastPingTime = current
		this.sendingUpdated = false
		seg.Release()
	}

	// flash remain segments
	this.output.Flush()

}

func (this *Connection) State() State {
	return State(atomic.LoadInt32((*int32)(&this.state)))
}
