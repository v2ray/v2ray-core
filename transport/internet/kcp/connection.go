package kcp

import (
	"errors"
	"io"
	"net"
	"sync"
	"time"

	"github.com/v2ray/v2ray-core/common/alloc"
	"github.com/v2ray/v2ray-core/common/log"
	"github.com/v2ray/v2ray-core/common/signal"
)

var (
	errTimeout          = errors.New("i/o timeout")
	errBrokenPipe       = errors.New("broken pipe")
	errClosedListener   = errors.New("Listener closed.")
	errClosedConnection = errors.New("Connection closed.")
)

const (
	headerSize = 2
)

type Command byte

var (
	CommandData      Command = 0
	CommandTerminate Command = 1
)

type Option byte

var (
	OptionClose Option = 1
)

type ConnState byte

var (
	ConnStateActive       ConnState = 0
	ConnStateReadyToClose ConnState = 1
	ConnStatePeerClosed   ConnState = 2
	ConnStateClosed       ConnState = 4
)

func nowMillisec() int64 {
	now := time.Now()
	return now.Unix()*1000 + int64(now.Nanosecond()/1000000)
}

// Connection is a KCP connection over UDP.
type Connection struct {
	sync.RWMutex
	state         ConnState
	kcp           *KCP // the core ARQ
	kcpAccess     sync.Mutex
	block         Authenticator
	needUpdate    bool
	local, remote net.Addr
	rd            time.Time // read deadline
	wd            time.Time // write deadline
	chReadEvent   chan struct{}
	writer        io.WriteCloser
	since         int64
	terminateOnce signal.Once
}

// NewConnection create a new KCP connection between local and remote.
func NewConnection(conv uint32, writerCloser io.WriteCloser, local *net.UDPAddr, remote *net.UDPAddr, block Authenticator) *Connection {
	conn := new(Connection)
	conn.local = local
	conn.chReadEvent = make(chan struct{}, 1)
	conn.remote = remote
	conn.block = block
	conn.writer = writerCloser
	conn.since = nowMillisec()

	mtu := uint32(effectiveConfig.Mtu - block.HeaderSize() - headerSize)
	conn.kcp = NewKCP(conv, mtu, conn.output)
	conn.kcp.WndSize(effectiveConfig.GetSendingWindowSize(), effectiveConfig.GetReceivingWindowSize())
	conn.kcp.NoDelay(1, effectiveConfig.Tti, 2, effectiveConfig.Congestion)
	conn.kcp.current = conn.Elapsed()

	go conn.updateTask()

	return conn
}

func (this *Connection) Elapsed() uint32 {
	return uint32(nowMillisec() - this.since)
}

// Read implements the Conn Read method.
func (this *Connection) Read(b []byte) (int, error) {
	if this == nil || this.state == ConnStateReadyToClose || this.state == ConnStateClosed {
		return 0, io.EOF
	}

	for {
		this.RLock()
		if this.state == ConnStateReadyToClose || this.state == ConnStateClosed {
			this.RUnlock()
			return 0, io.EOF
		}

		if !this.rd.IsZero() && this.rd.Before(time.Now()) {
			this.RUnlock()
			return 0, errTimeout
		}
		this.RUnlock()

		this.kcpAccess.Lock()
		nBytes := this.kcp.Recv(b)
		this.kcpAccess.Unlock()
		if nBytes > 0 {
			return nBytes, nil
		}
		select {
		case <-this.chReadEvent:
		case <-time.After(time.Second):
		}
	}
}

// Write implements the Conn Write method.
func (this *Connection) Write(b []byte) (int, error) {
	if this == nil ||
		this.state == ConnStateReadyToClose ||
		this.state == ConnStatePeerClosed ||
		this.state == ConnStateClosed {
		return 0, io.ErrClosedPipe
	}

	for {
		this.RLock()
		if this.state == ConnStateReadyToClose ||
			this.state == ConnStatePeerClosed ||
			this.state == ConnStateClosed {
			this.RUnlock()
			return 0, io.ErrClosedPipe
		}
		this.RUnlock()

		this.kcpAccess.Lock()
		if this.kcp.WaitSnd() < int(this.kcp.snd_wnd) {
			nBytes := len(b)
			this.kcp.Send(b)
			this.kcp.current = this.Elapsed()
			this.kcp.flush()
			this.kcpAccess.Unlock()
			return nBytes, nil
		}
		this.kcpAccess.Unlock()

		if !this.wd.IsZero() && this.wd.Before(time.Now()) {
			return 0, errTimeout
		}

		// Sending windows is 1024 for the moment. This amount is not gonna sent in 1 sec.
		time.Sleep(time.Second)
	}
}

func (this *Connection) Terminate() {
	if this == nil || this.state == ConnStateClosed {
		return
	}
	this.Lock()
	defer this.Unlock()
	if this.state == ConnStateClosed {
		return
	}

	this.state = ConnStateClosed
	this.writer.Close()
}

func (this *Connection) NotifyTermination() {
	for i := 0; i < 16; i++ {
		this.RLock()
		if this.state == ConnStateClosed {
			this.RUnlock()
			break
		}
		this.RUnlock()
		buffer := alloc.NewSmallBuffer().Clear()
		buffer.AppendBytes(byte(CommandTerminate), byte(OptionClose), byte(0), byte(0), byte(0), byte(0))
		this.outputBuffer(buffer)

		time.Sleep(time.Second)

	}
	this.Terminate()
}

func (this *Connection) ForceTimeout() {
	if this == nil {
		return
	}
	for i := 0; i < 5; i++ {
		if this.state == ConnStateClosed {
			return
		}
		time.Sleep(time.Minute)
	}
	go this.terminateOnce.Do(this.NotifyTermination)
}

// Close closes the connection.
func (this *Connection) Close() error {
	if this == nil || this.state == ConnStateClosed || this.state == ConnStateReadyToClose {
		return errClosedConnection
	}
	log.Debug("KCP|Connection: Closing connection to ", this.remote)
	this.Lock()
	defer this.Unlock()

	if this.state == ConnStateActive {
		this.state = ConnStateReadyToClose
		if this.kcp.WaitSnd() == 0 {
			go this.terminateOnce.Do(this.NotifyTermination)
		} else {
			go this.ForceTimeout()
		}
	}

	if this.state == ConnStatePeerClosed {
		go this.Terminate()
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
	if this == nil || this.state != ConnStateActive {
		return errClosedConnection
	}
	this.Lock()
	defer this.Unlock()
	this.rd = t
	this.wd = t
	return nil
}

// SetReadDeadline implements the Conn SetReadDeadline method.
func (this *Connection) SetReadDeadline(t time.Time) error {
	if this == nil || this.state != ConnStateActive {
		return errClosedConnection
	}
	this.Lock()
	defer this.Unlock()
	this.rd = t
	return nil
}

// SetWriteDeadline implements the Conn SetWriteDeadline method.
func (this *Connection) SetWriteDeadline(t time.Time) error {
	if this == nil || this.state != ConnStateActive {
		return errClosedConnection
	}
	this.Lock()
	defer this.Unlock()
	this.wd = t
	return nil
}

func (this *Connection) outputBuffer(payload *alloc.Buffer) {
	defer payload.Release()
	if this == nil {
		return
	}

	this.RLock()
	defer this.RUnlock()
	if this.state == ConnStatePeerClosed || this.state == ConnStateClosed {
		return
	}
	this.block.Seal(payload)

	this.writer.Write(payload.Value)
}

func (this *Connection) output(payload []byte) {
	if this == nil || this.state == ConnStateClosed {
		return
	}

	if this.state == ConnStateReadyToClose && this.kcp.WaitSnd() == 0 {
		go this.terminateOnce.Do(this.NotifyTermination)
	}

	if len(payload) < IKCP_OVERHEAD {
		return
	}

	buffer := alloc.NewBuffer().Clear().Append(payload)
	cmd := CommandData
	opt := Option(0)
	if this.state == ConnStateReadyToClose {
		opt = OptionClose
	}
	buffer.Prepend([]byte{byte(cmd), byte(opt)})
	this.outputBuffer(buffer)
}

// kcp update, input loop
func (this *Connection) updateTask() {
	for this.state != ConnStateClosed {
		current := this.Elapsed()
		this.kcpAccess.Lock()
		this.kcp.Update(current)
		interval := this.kcp.Check(this.Elapsed())
		this.kcpAccess.Unlock()
		sleep := interval - current
		if sleep < 10 {
			sleep = 10
		}
		time.Sleep(time.Duration(sleep) * time.Millisecond)
	}
}

func (this *Connection) notifyReadEvent() {
	select {
	case this.chReadEvent <- struct{}{}:
	default:
	}
}

func (this *Connection) MarkPeerClose() {
	this.Lock()
	defer this.Unlock()
	if this.state == ConnStateReadyToClose {
		this.state = ConnStateClosed
		go this.Terminate()
		return
	}
	if this.state == ConnStateActive {
		this.state = ConnStatePeerClosed
	}
}

func (this *Connection) kcpInput(data []byte) {
	cmd := Command(data[0])
	opt := Option(data[1])
	if cmd == CommandTerminate {
		go this.Terminate()
		return
	}
	if opt == OptionClose {
		go this.MarkPeerClose()
	}
	this.kcpAccess.Lock()
	this.kcp.current = this.Elapsed()
	this.kcp.Input(data[2:])

	this.kcpAccess.Unlock()
	this.notifyReadEvent()
}

func (this *Connection) FetchInputFrom(conn net.Conn) {
	go func() {
		for {
			payload := alloc.NewBuffer()
			nBytes, err := conn.Read(payload.Value)
			if err != nil {
				return
			}
			payload.Slice(0, nBytes)
			if this.block.Open(payload) {
				this.kcpInput(payload.Value)
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
