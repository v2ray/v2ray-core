package kcp

import (
	crand "crypto/rand"
	"encoding/binary"
	"errors"
	"hash/crc32"
	"io"
	"log"
	"math/rand"
	"net"
	"sync"
	"time"

	"golang.org/x/net/ipv4"
)

var (
	errTimeout    = errors.New("i/o timeout")
	errBrokenPipe = errors.New("broken pipe")
)

const (
	basePort        = 20000 // minimum port for listening
	maxPort         = 65535 // maximum port for listening
	defaultWndSize  = 128   // default window size, in packet
	otpSize         = 16    // magic number
	crcSize         = 4     // 4bytes packet checksum
	cryptHeaderSize = otpSize + crcSize
	connTimeout     = 60 * time.Second
	mtuLimit        = 4096
	rxQueueLimit    = 8192
	rxFecLimit      = 2048
)

type (
	// UDPSession defines a KCP session implemented by UDP
	UDPSession struct {
		kcp           *KCP         // the core ARQ
		conn          *net.UDPConn // the underlying UDP socket
		block         BlockCrypt
		needUpdate    bool
		l             *Listener // point to server listener if it's a server socket
		local, remote net.Addr
		rd            time.Time // read deadline
		wd            time.Time // write deadline
		sockbuff      []byte    // kcp receiving is based on packet, I turn it into stream
		die           chan struct{}
		isClosed      bool
		mu            sync.Mutex
		chReadEvent   chan struct{}
		chWriteEvent  chan struct{}
		chTicker      chan time.Time
		chUDPOutput   chan []byte
		headerSize    int
		lastInputTs   time.Time
		ackNoDelay    bool
	}
)

// newUDPSession create a new udp session for client or server
func newUDPSession(conv uint32, l *Listener, conn *net.UDPConn, remote *net.UDPAddr, block BlockCrypt) *UDPSession {
	sess := new(UDPSession)
	sess.chTicker = make(chan time.Time, 1)
	sess.chUDPOutput = make(chan []byte, rxQueueLimit)
	sess.die = make(chan struct{})
	sess.local = conn.LocalAddr()
	sess.chReadEvent = make(chan struct{}, 1)
	sess.chWriteEvent = make(chan struct{}, 1)
	sess.remote = remote
	sess.conn = conn
	sess.l = l
	sess.block = block
	sess.lastInputTs = time.Now()

	// caculate header size
	if sess.block != nil {
		sess.headerSize += cryptHeaderSize
	}

	sess.kcp = NewKCP(conv, func(buf []byte, size int) {
		if size >= IKCP_OVERHEAD {
			ext := make([]byte, sess.headerSize+size)
			copy(ext[sess.headerSize:], buf)
			sess.chUDPOutput <- ext
		}
	})
	sess.kcp.WndSize(defaultWndSize, defaultWndSize)
	sess.kcp.SetMtu(IKCP_MTU_DEF - sess.headerSize)

	go sess.updateTask()
	go sess.outputTask()
	if l == nil { // it's a client connection
		go sess.readLoop()
	}

	return sess
}

// Read implements the Conn Read method.
func (s *UDPSession) Read(b []byte) (n int, err error) {
	for {
		s.mu.Lock()
		if len(s.sockbuff) > 0 { // copy from buffer
			n = copy(b, s.sockbuff)
			s.sockbuff = s.sockbuff[n:]
			s.mu.Unlock()
			return n, nil
		}

		if s.isClosed {
			s.mu.Unlock()
			return 0, errBrokenPipe
		}

		if !s.rd.IsZero() {
			if time.Now().After(s.rd) { // timeout
				s.mu.Unlock()
				return 0, errTimeout
			}
		}

		if n := s.kcp.PeekSize(); n > 0 { // data arrived
			if len(b) >= n {
				s.kcp.Recv(b)
			} else {
				buf := make([]byte, n)
				s.kcp.Recv(buf)
				n = copy(b, buf)
				s.sockbuff = buf[n:] // store remaining bytes into sockbuff for next read
			}
			s.mu.Unlock()
			return n, nil
		}

		var timeout <-chan time.Time
		if !s.rd.IsZero() {
			delay := s.rd.Sub(time.Now())
			timeout = time.After(delay)
		}
		s.mu.Unlock()

		// wait for read event or timeout
		select {
		case <-s.chReadEvent:
		case <-timeout:
		case <-s.die:
		}
	}
}

// Write implements the Conn Write method.
func (s *UDPSession) Write(b []byte) (n int, err error) {
	for {
		s.mu.Lock()
		if s.isClosed {
			s.mu.Unlock()
			return 0, errBrokenPipe
		}

		if !s.wd.IsZero() {
			if time.Now().After(s.wd) { // timeout
				s.mu.Unlock()
				return 0, errTimeout
			}
		}

		if s.kcp.WaitSnd() < int(s.kcp.snd_wnd) {
			n = len(b)
			max := s.kcp.mss << 8
			for {
				if len(b) <= int(max) { // in most cases
					s.kcp.Send(b)
					break
				} else {
					s.kcp.Send(b[:max])
					b = b[max:]
				}
			}
			s.kcp.current = currentMs()
			s.kcp.flush()
			s.mu.Unlock()
			return n, nil
		}

		var timeout <-chan time.Time
		if !s.wd.IsZero() {
			delay := s.wd.Sub(time.Now())
			timeout = time.After(delay)
		}
		s.mu.Unlock()

		// wait for write event or timeout
		select {
		case <-s.chWriteEvent:
		case <-timeout:
		case <-s.die:
		}
	}
}

// Close closes the connection.
func (s *UDPSession) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.isClosed {
		return errBrokenPipe
	}
	close(s.die)
	s.isClosed = true
	if s.l == nil { // client socket close
		s.conn.Close()
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
	s.mu.Lock()
	defer s.mu.Unlock()
	s.rd = t
	s.wd = t
	return nil
}

// SetReadDeadline implements the Conn SetReadDeadline method.
func (s *UDPSession) SetReadDeadline(t time.Time) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.rd = t
	return nil
}

// SetWriteDeadline implements the Conn SetWriteDeadline method.
func (s *UDPSession) SetWriteDeadline(t time.Time) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.wd = t
	return nil
}

// SetWindowSize set maximum window size
func (s *UDPSession) SetWindowSize(sndwnd, rcvwnd int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.kcp.WndSize(sndwnd, rcvwnd)
}

// SetMtu sets the maximum transmission unit
func (s *UDPSession) SetMtu(mtu int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.kcp.SetMtu(mtu - s.headerSize)
}

// SetACKNoDelay changes ack flush option, set true to flush ack immediately,
func (s *UDPSession) SetACKNoDelay(nodelay bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.ackNoDelay = nodelay
}

// SetNoDelay calls nodelay() of kcp
func (s *UDPSession) SetNoDelay(nodelay, interval, resend, nc int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.kcp.NoDelay(nodelay, interval, resend, nc)
}

// SetDSCP sets the DSCP field of IP header
func (s *UDPSession) SetDSCP(tos int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := ipv4.NewConn(s.conn).SetTOS(tos << 2); err != nil {
		log.Println("set tos:", err)
	}
}

func (s *UDPSession) outputTask() {
	// ping
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case ext := <-s.chUDPOutput:
			if s.block != nil {
				io.ReadFull(crand.Reader, ext[:otpSize]) // OTP
				checksum := crc32.ChecksumIEEE(ext[cryptHeaderSize:])
				binary.LittleEndian.PutUint32(ext[otpSize:], checksum)
				s.block.Encrypt(ext, ext)
			}

			//if rand.Intn(100) < 80 {
			n, err := s.conn.WriteTo(ext, s.remote)
			if err != nil {
				log.Println(err, n)
			}
			//}

		case <-ticker.C:
			sz := rand.Intn(IKCP_MTU_DEF - s.headerSize - IKCP_OVERHEAD)
			sz += s.headerSize + IKCP_OVERHEAD
			ping := make([]byte, sz)
			io.ReadFull(crand.Reader, ping)
			if s.block != nil {
				checksum := crc32.ChecksumIEEE(ping[cryptHeaderSize:])
				binary.LittleEndian.PutUint32(ping[otpSize:], checksum)
				s.block.Encrypt(ping, ping)
			}

			n, err := s.conn.WriteTo(ping, s.remote)
			if err != nil {
				log.Println(err, n)
			}
		case <-s.die:
			return
		}
	}
}

// kcp update, input loop
func (s *UDPSession) updateTask() {
	var tc <-chan time.Time
	if s.l == nil { // client
		ticker := time.NewTicker(10 * time.Millisecond)
		tc = ticker.C
		defer ticker.Stop()
	} else {
		tc = s.chTicker
	}

	var nextupdate uint32
	for {
		select {
		case <-tc:
			s.mu.Lock()
			current := currentMs()
			if current >= nextupdate || s.needUpdate {
				s.kcp.Update(current)
				nextupdate = s.kcp.Check(current)
			}
			if s.kcp.WaitSnd() < int(s.kcp.snd_wnd) {
				s.notifyWriteEvent()
			}
			s.needUpdate = false
			s.mu.Unlock()
		case <-s.die:
			if s.l != nil { // has listener
				s.l.chDeadlinks <- s.remote
			}
			return
		}
	}
}

// GetConv gets conversation id of a session
func (s *UDPSession) GetConv() uint32 {
	return s.kcp.conv
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

func (s *UDPSession) kcpInput(data []byte) {
	now := time.Now()
	if now.Sub(s.lastInputTs) > connTimeout {
		s.Close()
		return
	}
	s.lastInputTs = now

	s.mu.Lock()
	s.kcp.current = currentMs()
	s.kcp.Input(data)

	if s.ackNoDelay {
		s.kcp.current = currentMs()
		s.kcp.flush()
	} else {
		s.needUpdate = true
	}
	s.mu.Unlock()
	s.notifyReadEvent()
}

func (s *UDPSession) receiver(ch chan []byte) {
	for {
		data := make([]byte, mtuLimit)
		if n, _, err := s.conn.ReadFromUDP(data); err == nil && n >= s.headerSize+IKCP_OVERHEAD {
			ch <- data[:n]
		} else if err != nil {
			return
		}
	}
}

// read loop for client session
func (s *UDPSession) readLoop() {
	chPacket := make(chan []byte, rxQueueLimit)
	go s.receiver(chPacket)

	for {
		select {
		case data := <-chPacket:
			dataValid := false
			if s.block != nil {
				s.block.Decrypt(data, data)
				data = data[otpSize:]
				checksum := crc32.ChecksumIEEE(data[crcSize:])
				if checksum == binary.LittleEndian.Uint32(data) {
					data = data[crcSize:]
					dataValid = true
				}
			} else if s.block == nil {
				dataValid = true
			}

			if dataValid {
				s.kcpInput(data)
			}
		case <-s.die:
			return
		}
	}
}

type (
	// Listener defines a server listening for connections
	Listener struct {
		block       BlockCrypt
		conn        *net.UDPConn
		sessions    map[string]*UDPSession
		chAccepts   chan *UDPSession
		chDeadlinks chan net.Addr
		headerSize  int
		die         chan struct{}
	}

	packet struct {
		from *net.UDPAddr
		data []byte
	}
)

// monitor incoming data for all connections of server
func (l *Listener) monitor() {
	chPacket := make(chan packet, rxQueueLimit)
	go l.receiver(chPacket)
	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case p := <-chPacket:
			data := p.data
			from := p.from
			dataValid := false
			if l.block != nil {
				l.block.Decrypt(data, data)
				data = data[otpSize:]
				checksum := crc32.ChecksumIEEE(data[crcSize:])
				if checksum == binary.LittleEndian.Uint32(data) {
					data = data[crcSize:]
					dataValid = true
				}
			} else if l.block == nil {
				dataValid = true
			}

			if dataValid {
				addr := from.String()
				s, ok := l.sessions[addr]
				if !ok { // new session
					var conv uint32
					convValid := false

					conv = binary.LittleEndian.Uint32(data)
					convValid = true

					if convValid {
						s := newUDPSession(conv, l, l.conn, from, l.block)
						s.kcpInput(data)
						l.sessions[addr] = s
						l.chAccepts <- s
					}
				} else {
					s.kcpInput(data)
				}
			}
		case deadlink := <-l.chDeadlinks:
			delete(l.sessions, deadlink.String())
		case <-l.die:
			return
		case <-ticker.C:
			now := time.Now()
			for _, s := range l.sessions {
				select {
				case s.chTicker <- now:
				default:
				}
			}
		}
	}
}

func (l *Listener) receiver(ch chan packet) {
	for {
		data := make([]byte, mtuLimit)
		if n, from, err := l.conn.ReadFromUDP(data); err == nil && n >= l.headerSize+IKCP_OVERHEAD {
			ch <- packet{from, data[:n]}
		} else if err != nil {
			return
		}
	}
}

// Accept implements the Accept method in the Listener interface; it waits for the next call and returns a generic Conn.
func (l *Listener) Accept() (*UDPSession, error) {
	select {
	case c := <-l.chAccepts:
		return c, nil
	case <-l.die:
		return nil, errors.New("listener stopped")
	}
}

// Close stops listening on the UDP address. Already Accepted connections are not closed.
func (l *Listener) Close() error {
	if err := l.conn.Close(); err == nil {
		close(l.die)
		return nil
	} else {
		return err
	}
}

// Addr returns the listener's network address, The Addr returned is shared by all invocations of Addr, so do not modify it.
func (l *Listener) Addr() net.Addr {
	return l.conn.LocalAddr()
}

// Listen listens for incoming KCP packets addressed to the local address laddr on the network "udp",
func Listen(laddr string) (*Listener, error) {
	return ListenWithOptions(laddr, nil)
}

// ListenWithOptions listens for incoming KCP packets addressed to the local address laddr on the network "udp" with packet encryption,
// FEC = 0 means no FEC, FEC > 0 means num(FEC) as a FEC cluster
func ListenWithOptions(laddr string, block BlockCrypt) (*Listener, error) {
	udpaddr, err := net.ResolveUDPAddr("udp", laddr)
	if err != nil {
		return nil, err
	}
	conn, err := net.ListenUDP("udp", udpaddr)
	if err != nil {
		return nil, err
	}

	l := new(Listener)
	l.conn = conn
	l.sessions = make(map[string]*UDPSession)
	l.chAccepts = make(chan *UDPSession, 1024)
	l.chDeadlinks = make(chan net.Addr, 1024)
	l.die = make(chan struct{})
	l.block = block

	// caculate header size
	if l.block != nil {
		l.headerSize += cryptHeaderSize
	}

	go l.monitor()
	return l, nil
}

// Dial connects to the remote address raddr on the network "udp"
func Dial(raddr string) (*UDPSession, error) {
	return DialWithOptions(raddr, nil)
}

// DialWithOptions connects to the remote address raddr on the network "udp" with packet encryption
func DialWithOptions(raddr string, block BlockCrypt) (*UDPSession, error) {
	udpaddr, err := net.ResolveUDPAddr("udp", raddr)
	if err != nil {
		return nil, err
	}

	for {
		port := basePort + rand.Int()%(maxPort-basePort)
		if udpconn, err := net.ListenUDP("udp", &net.UDPAddr{Port: port}); err == nil {
			return newUDPSession(rand.Uint32(), nil, udpconn, udpaddr, block), nil
		}
	}
}

func currentMs() uint32 {
	return uint32(time.Now().UnixNano() / int64(time.Millisecond))
}
