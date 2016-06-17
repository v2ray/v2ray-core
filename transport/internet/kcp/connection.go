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
	errTimeout        = errors.New("i/o timeout")
	errBrokenPipe     = errors.New("broken pipe")
	errClosedListener = errors.New("Listener closed.")
)

const (
	basePort       = 20000 // minimum port for listening
	maxPort        = 65535 // maximum port for listening
	defaultWndSize = 128   // default window size, in packet
	mtuLimit       = 4096
	rxQueueLimit   = 8192
	rxFecLimit     = 2048

	headerSize        = 2
	cmdData    uint16 = 0
	cmdClose   uint16 = 1
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

// UDPSession defines a KCP session implemented by UDP
type UDPSession struct {
	sync.Mutex
	state         ConnState
	kcp           *KCP // the core ARQ
	kcpAccess     sync.Mutex
	block         Authenticator
	needUpdate    bool
	local, remote net.Addr
	rd            time.Time // read deadline
	wd            time.Time // write deadline
	chReadEvent   chan struct{}
	chWriteEvent  chan struct{}
	writer        io.WriteCloser
	since         int64
}

// newUDPSession create a new udp session for client or server
func newUDPSession(conv uint32, writerCloser io.WriteCloser, local *net.UDPAddr, remote *net.UDPAddr, block Authenticator) *UDPSession {
	sess := new(UDPSession)
	sess.local = local
	sess.chReadEvent = make(chan struct{}, 1)
	sess.chWriteEvent = make(chan struct{}, 1)
	sess.remote = remote
	sess.block = block
	sess.writer = writerCloser
	sess.since = nowMillisec()

	mtu := uint32(effectiveConfig.Mtu - block.HeaderSize() - headerSize)
	sess.kcp = NewKCP(conv, mtu, func(buf []byte, size int) {
		if size >= IKCP_OVERHEAD {
			ext := alloc.NewBuffer().Clear().Append(buf[:size])
			cmd := cmdData
			opt := Option(0)
			if sess.state == ConnStateReadyToClose {
				opt = OptionClose
			}
			ext.Prepend([]byte{byte(cmd), byte(opt)})
			sess.output(ext)
		}
	})
	sess.kcp.WndSize(effectiveConfig.Sndwnd, effectiveConfig.Rcvwnd)
	sess.kcp.NoDelay(1, 20, 2, 1)
	sess.kcp.current = sess.Elapsed()

	go sess.updateTask()

	return sess
}

func (this *UDPSession) Elapsed() uint32 {
	return uint32(nowMillisec() - this.since)
}

// Read implements the Conn Read method.
func (s *UDPSession) Read(b []byte) (int, error) {
	if s.state == ConnStateReadyToClose || s.state == ConnStateClosed {
		return 0, io.EOF
	}

	for {
		s.Lock()
		if s.state == ConnStateReadyToClose || s.state == ConnStateClosed {
			s.Unlock()
			return 0, io.EOF
		}

		if !s.rd.IsZero() {
			if time.Now().After(s.rd) {
				s.Unlock()
				return 0, errTimeout
			}
		}

		nBytes := s.kcp.Recv(b)
		if nBytes > 0 {
			s.Unlock()
			return nBytes, nil
		}

		var timeout <-chan time.Time
		if !s.rd.IsZero() {
			delay := s.rd.Sub(time.Now())
			timeout = time.After(delay)
		}

		s.Unlock()
		select {
		case <-s.chReadEvent:
		case <-timeout:
			return 0, errTimeout
		}
	}
}

// Write implements the Conn Write method.
func (s *UDPSession) Write(b []byte) (int, error) {
	if s.state == ConnStateReadyToClose ||
		s.state == ConnStatePeerClosed ||
		s.state == ConnStateClosed {
		return 0, io.ErrClosedPipe
	}

	for {
		if s.state == ConnStateReadyToClose ||
			s.state == ConnStatePeerClosed ||
			s.state == ConnStateClosed {
			return 0, io.ErrClosedPipe
		}

		s.kcpAccess.Lock()
		if s.kcp.WaitSnd() < int(s.kcp.snd_wnd) {
			nBytes := len(b)
			s.kcp.Send(b)
			s.kcp.current = s.Elapsed()
			s.kcp.flush()
			s.kcpAccess.Unlock()
			return nBytes, nil
		}
		s.kcpAccess.Unlock()

		if !s.wd.IsZero() && s.wd.Before(time.Now()) {
			return 0, errTimeout
		}

		time.Sleep(time.Duration(s.kcp.WaitSnd()*5) * time.Millisecond)
	}
}

func (this *UDPSession) Terminate() {
	if this.state == ConnStateClosed {
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

func (this *UDPSession) NotifyTermination() {
	for i := 0; i < 16; i++ {
		this.Lock()
		if this.state == ConnStateClosed {
			this.Unlock()
			break
		}
		buffer := alloc.NewSmallBuffer().Clear()
		buffer.AppendBytes(byte(CommandTerminate), byte(OptionClose), byte(0), byte(0), byte(0), byte(0))
		this.output(buffer)
		time.Sleep(time.Second)
		this.Unlock()
	}
	this.Terminate()
}

// Close closes the connection.
func (s *UDPSession) Close() error {
	log.Debug("KCP|Connection: Closing connection to ", s.remote)
	s.Lock()
	defer s.Unlock()

	if s.state == ConnStateActive {
		s.state = ConnStateReadyToClose
		if s.kcp.WaitSnd() == 0 {
			go s.NotifyTermination()
		}
	}

	if s.state == ConnStatePeerClosed {
		go s.Terminate()
	}

	return nil
}

// LocalAddr returns the local network address. The Addr returned is shared by all invocations of LocalAddr, so do not modify it.
func (s *UDPSession) LocalAddr() net.Addr {
	return s.local
}

// RemoteAddr returns the remote network address. The Addr returned is shared by all invocations of RemoteAddr, so do not modify it.
func (s *UDPSession) RemoteAddr() net.Addr { return s.remote }

// SetDeadline sets the deadline associated with the listener. A zero time value disables the deadline.
func (s *UDPSession) SetDeadline(t time.Time) error {
	s.Lock()
	defer s.Unlock()
	s.rd = t
	s.wd = t
	return nil
}

// SetReadDeadline implements the Conn SetReadDeadline method.
func (s *UDPSession) SetReadDeadline(t time.Time) error {
	s.Lock()
	defer s.Unlock()
	s.rd = t
	return nil
}

// SetWriteDeadline implements the Conn SetWriteDeadline method.
func (s *UDPSession) SetWriteDeadline(t time.Time) error {
	s.Lock()
	defer s.Unlock()
	s.wd = t
	return nil
}

func (s *UDPSession) output(payload *alloc.Buffer) {
	defer payload.Release()

	if s.state == ConnStatePeerClosed || s.state == ConnStateClosed {
		return
	}
	s.block.Seal(payload)

	s.writer.Write(payload.Value)
}

// kcp update, input loop
func (s *UDPSession) updateTask() {
	for s.state != ConnStateClosed {
		current := s.Elapsed()
		s.kcp.Update(current)
		interval := s.kcp.Check(s.Elapsed())
		time.Sleep(time.Duration(interval) * time.Millisecond)
	}
}

func (s *UDPSession) notifyReadEvent() {
	select {
	case s.chReadEvent <- struct{}{}:
	default:
	}
}

func (s *UDPSession) notifyWriteEvent() {
	select {
	case s.chWriteEvent <- struct{}{}:
	default:
	}
}

func (this *UDPSession) MarkPeerClose() {
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

func (s *UDPSession) kcpInput(data []byte) {
	cmd := Command(data[0])
	opt := Option(data[1])
	if cmd == CommandTerminate {
		go s.Terminate()
		return
	}
	if opt == OptionClose {
		go s.MarkPeerClose()
	}
	s.kcpAccess.Lock()
	s.kcp.current = s.Elapsed()
	s.kcp.Input(data[2:])

	s.kcpAccess.Unlock()
	s.notifyReadEvent()
}

func (this *UDPSession) FetchInputFrom(conn net.Conn) {
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

func (this *UDPSession) Reusable() bool {
	return false
}

func (this *UDPSession) SetReusable(b bool) {}
