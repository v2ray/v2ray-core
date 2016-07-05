package kcp

import (
	"errors"
	"io"
	"net"
	"sync"
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

type State int

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

// Connection is a KCP connection over UDP.
type Connection struct {
	sync.RWMutex
	block         Authenticator
	local, remote net.Addr
	wd            time.Time // write deadline
	writer        io.WriteCloser
	since         int64

	conv             uint16
	state            State
	stateBeginTime   uint32
	lastIncomingTime uint32
	lastPayloadTime  uint32
	sendingUpdated   bool
	lastPingTime     uint32

	mss                        uint32
	rx_rttvar, rx_srtt, rx_rto uint32
	interval                   uint32

	receivingWorker *ReceivingWorker
	sendingWorker   *SendingWorker

	fastresend        uint32
	congestionControl bool
	output            *BufferedSegmentWriter
}

// NewConnection create a new KCP connection between local and remote.
func NewConnection(conv uint16, writerCloser io.WriteCloser, local *net.UDPAddr, remote *net.UDPAddr, block Authenticator) *Connection {
	log.Debug("KCP|Connection: creating connection ", conv)

	conn := new(Connection)
	conn.local = local
	conn.remote = remote
	conn.block = block
	conn.writer = writerCloser
	conn.since = nowMillisec()

	authWriter := &AuthenticationWriter{
		Authenticator: block,
		Writer:        writerCloser,
	}
	conn.conv = conv
	conn.output = NewSegmentWriter(authWriter)

	conn.mss = authWriter.Mtu() - DataSegmentOverhead
	conn.rx_rto = 100
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

	state := this.State()
	if state == StateTerminating || state == StateTerminated {
		return 0, io.EOF
	}
	return this.receivingWorker.Read(b)
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
	this.Lock()
	this.state = state
	this.stateBeginTime = this.Elapsed()
	this.Unlock()

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

	state := this.State()
	if state == StateReadyToClose ||
		state == StateTerminating ||
		state == StateTerminated {
		return errClosedConnection
	}
	log.Debug("KCP|Connection: Closing connection to ", this.remote)

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
	this.receivingWorker.SetReadDeadline(t)
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
	log.Info("Terminating connection to ", this.RemoteAddr())

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

// https://tools.ietf.org/html/rfc6298
func (this *Connection) update_ack(rtt int32) {
	this.Lock()
	defer this.Unlock()

	if this.rx_srtt == 0 {
		this.rx_srtt = uint32(rtt)
		this.rx_rttvar = uint32(rtt) / 2
	} else {
		delta := rtt - int32(this.rx_srtt)
		if delta < 0 {
			delta = -delta
		}
		this.rx_rttvar = (3*this.rx_rttvar + uint32(delta)) / 4
		this.rx_srtt = (7*this.rx_srtt + uint32(rtt)) / 8
		if this.rx_srtt < this.interval {
			this.rx_srtt = this.interval
		}
	}
	var rto uint32
	if this.interval < 4*this.rx_rttvar {
		rto = this.rx_srtt + 4*this.rx_rttvar
	} else {
		rto = this.rx_srtt + this.interval
	}

	if rto > 10000 {
		rto = 10000
	}
	this.rx_rto = rto * 3 / 2
}

// Input when you received a low level packet (eg. UDP packet), call it
func (kcp *Connection) Input(data []byte) int {
	current := kcp.Elapsed()
	kcp.lastIncomingTime = current

	var seg Segment
	for {
		seg, data = ReadSegment(data)
		if seg == nil {
			break
		}

		switch seg := seg.(type) {
		case *DataSegment:
			kcp.HandleOption(seg.Opt)
			kcp.receivingWorker.ProcessSegment(seg)
			kcp.lastPayloadTime = current
		case *AckSegment:
			kcp.HandleOption(seg.Opt)
			kcp.sendingWorker.ProcessSegment(current, seg)
			kcp.lastPayloadTime = current
		case *CmdOnlySegment:
			kcp.HandleOption(seg.Opt)
			if seg.Cmd == SegmentCommandTerminated {
				if kcp.state == StateActive ||
					kcp.state == StateReadyToClose ||
					kcp.state == StatePeerClosed {
					kcp.SetState(StateTerminating)
				} else if kcp.state == StateTerminating {
					kcp.SetState(StateTerminated)
				}
			}
			kcp.sendingWorker.ProcessReceivingNext(seg.ReceivinNext)
			kcp.receivingWorker.ProcessSendingNext(seg.SendingNext)
		default:
		}
	}

	return 0
}

func (this *Connection) flush() {
	current := this.Elapsed()
	state := this.State()

	if state == StateTerminated {
		return
	}
	if state == StateActive && current-this.lastPayloadTime >= 30000 {
		this.Close()
	}

	if state == StateTerminating {
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

	if state == StateReadyToClose && current-this.stateBeginTime > 15000 {
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
		if state == StateReadyToClose {
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
	this.RLock()
	defer this.RUnlock()

	return this.state
}
