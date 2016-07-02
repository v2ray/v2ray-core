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
	headerSize uint32 = 2
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
	wd            time.Time // write deadline
	chReadEvent   chan struct{}
	writer        io.WriteCloser
	since         int64
	terminateOnce signal.Once
}

// NewConnection create a new KCP connection between local and remote.
func NewConnection(conv uint16, writerCloser io.WriteCloser, local *net.UDPAddr, remote *net.UDPAddr, block Authenticator) *Connection {
	conn := new(Connection)
	conn.local = local
	conn.chReadEvent = make(chan struct{}, 1)
	conn.remote = remote
	conn.block = block
	conn.writer = writerCloser
	conn.since = nowMillisec()

	authWriter := &AuthenticationWriter{
		Authenticator: block,
		Writer:        writerCloser,
	}
	conn.kcp = NewKCP(conv, authWriter)
	conn.kcp.NoDelay(effectiveConfig.Tti, 2, effectiveConfig.Congestion)
	conn.kcp.current = conn.Elapsed()

	go conn.updateTask()

	return conn
}

func (this *Connection) Elapsed() uint32 {
	return uint32(nowMillisec() - this.since)
}

// Read implements the Conn Read method.
func (this *Connection) Read(b []byte) (int, error) {
	if this == nil || this.kcp.state == StateTerminating || this.kcp.state == StateTerminated {
		return 0, io.EOF
	}
	return this.kcp.receivingWorker.Read(b)
}

// Write implements the Conn Write method.
func (this *Connection) Write(b []byte) (int, error) {
	if this == nil || this.kcp.state != StateActive {
		return 0, io.ErrClosedPipe
	}
	totalWritten := 0

	for {
		this.RLock()
		if this == nil || this.kcp.state != StateActive {
			this.RUnlock()
			return totalWritten, io.ErrClosedPipe
		}
		this.RUnlock()

		this.kcpAccess.Lock()
		nBytes := this.kcp.Send(b[totalWritten:])
		if nBytes > 0 {
			totalWritten += nBytes
			if totalWritten == len(b) {
				this.kcpAccess.Unlock()
				return totalWritten, nil
			}
		}

		this.kcpAccess.Unlock()

		if !this.wd.IsZero() && this.wd.Before(time.Now()) {
			return totalWritten, errTimeout
		}

		// Sending windows is 1024 for the moment. This amount is not gonna sent in 1 sec.
		time.Sleep(time.Second)
	}
}

// Close closes the connection.
func (this *Connection) Close() error {
	if this == nil ||
		this.kcp.state == StateReadyToClose ||
		this.kcp.state == StateTerminating ||
		this.kcp.state == StateTerminated {
		return errClosedConnection
	}
	log.Debug("KCP|Connection: Closing connection to ", this.remote)
	this.Lock()
	defer this.Unlock()

	this.kcpAccess.Lock()
	this.kcp.OnClose()
	this.kcpAccess.Unlock()

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
	if this == nil || this.kcp.state != StateActive {
		return errClosedConnection
	}
	this.kcpAccess.Lock()
	defer this.kcpAccess.Unlock()
	this.kcp.receivingWorker.SetReadDeadline(t)
	return nil
}

// SetWriteDeadline implements the Conn SetWriteDeadline method.
func (this *Connection) SetWriteDeadline(t time.Time) error {
	if this == nil || this.kcp.state != StateActive {
		return errClosedConnection
	}
	this.Lock()
	defer this.Unlock()
	this.wd = t
	return nil
}

// kcp update, input loop
func (this *Connection) updateTask() {
	for this.kcp.state != StateTerminated {
		current := this.Elapsed()
		this.kcpAccess.Lock()
		this.kcp.Update(current)
		this.kcpAccess.Unlock()

		interval := time.Duration(effectiveConfig.Tti) * time.Millisecond
		if this.kcp.state == StateTerminating {
			interval = time.Second
		}
		time.Sleep(interval)
	}
	this.Terminate()
}

func (this *Connection) notifyReadEvent() {
	select {
	case this.chReadEvent <- struct{}{}:
	default:
	}
}

func (this *Connection) kcpInput(data []byte) {
	this.kcpAccess.Lock()
	this.kcp.current = this.Elapsed()
	this.kcp.Input(data)

	this.kcpAccess.Unlock()
	this.notifyReadEvent()
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

func (this *Connection) Terminate() {
	if this == nil || this.writer == nil {
		return
	}
	log.Info("Terminating connection to ", this.RemoteAddr())

	this.writer.Close()
}
