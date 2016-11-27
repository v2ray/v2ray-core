package kcp

import (
	"crypto/tls"
	"net"
	"sync"
	"time"

	"v2ray.com/core/common/alloc"
	"v2ray.com/core/common/log"
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/common/serial"
	"v2ray.com/core/proxy"
	"v2ray.com/core/transport/internet"
	"v2ray.com/core/transport/internet/internal"
	v2tls "v2ray.com/core/transport/internet/tls"
	"v2ray.com/core/transport/internet/udp"
)

type ConnectionId struct {
	Remote v2net.Address
	Port   v2net.Port
	Conv   uint16
}

type ServerConnection struct {
	id     internal.ConnectionId
	writer *Writer
	local  net.Addr
	remote net.Addr
	auth   internet.Authenticator
	input  func([]byte)
}

func (o *ServerConnection) Read([]byte) (int, error) {
	panic("KCP|ServerConnection: Read should not be called.")
}

func (o *ServerConnection) Write(b []byte) (int, error) {
	return o.writer.Write(b)
}

func (o *ServerConnection) Close() error {
	return o.writer.Close()
}

func (o *ServerConnection) Reset(auth internet.Authenticator, input func([]byte)) {
	o.auth = auth
	o.input = input
}

func (o *ServerConnection) Input(b *alloc.Buffer) {
	defer b.Release()

	if o.auth.Open(b) {
		o.input(b.Value)
	}
}

func (o *ServerConnection) LocalAddr() net.Addr {
	return o.local
}

func (o *ServerConnection) RemoteAddr() net.Addr {
	return o.remote
}

func (o *ServerConnection) SetDeadline(time.Time) error {
	return nil
}

func (o *ServerConnection) SetReadDeadline(time.Time) error {
	return nil
}

func (o *ServerConnection) SetWriteDeadline(time.Time) error {
	return nil
}

func (o *ServerConnection) Id() internal.ConnectionId {
	return o.id
}

// Listener defines a server listening for connections
type Listener struct {
	sync.Mutex
	running       bool
	authenticator internet.Authenticator
	sessions      map[ConnectionId]*Connection
	awaitingConns chan *Connection
	hub           *udp.UDPHub
	tlsConfig     *tls.Config
	config        *Config
}

func NewListener(address v2net.Address, port v2net.Port, options internet.ListenOptions) (*Listener, error) {
	networkSettings, err := options.Stream.GetEffectiveNetworkSettings()
	if err != nil {
		log.Error("KCP|Listener: Failed to get KCP settings: ", err)
		return nil, err
	}
	kcpSettings := networkSettings.(*Config)
	kcpSettings.ConnectionReuse = &ConnectionReuse{Enable: false}

	auth, err := kcpSettings.GetAuthenticator()
	if err != nil {
		return nil, err
	}
	l := &Listener{
		authenticator: auth,
		sessions:      make(map[ConnectionId]*Connection),
		awaitingConns: make(chan *Connection, 64),
		running:       true,
		config:        kcpSettings,
	}
	if options.Stream != nil && options.Stream.HasSecuritySettings() {
		securitySettings, err := options.Stream.GetEffectiveSecuritySettings()
		if err != nil {
			log.Error("KCP|Listener: Failed to get security settings: ", err)
			return nil, err
		}
		switch securitySettings := securitySettings.(type) {
		case *v2tls.Config:
			l.tlsConfig = securitySettings.GetTLSConfig()
		}
	}
	hub, err := udp.ListenUDP(address, port, udp.ListenOption{Callback: l.OnReceive, Concurrency: 2})
	if err != nil {
		return nil, err
	}
	l.hub = hub
	log.Info("KCP|Listener: listening on ", address, ":", port)
	return l, nil
}

func (this *Listener) OnReceive(payload *alloc.Buffer, session *proxy.SessionInfo) {
	defer payload.Release()

	src := session.Source

	if valid := this.authenticator.Open(payload); !valid {
		log.Info("KCP|Listener: discarding invalid payload from ", src)
		return
	}
	if !this.running {
		return
	}
	this.Lock()
	defer this.Unlock()
	if !this.running {
		return
	}
	if payload.Len() < 4 {
		return
	}
	conv := serial.BytesToUint16(payload.Value)
	cmd := Command(payload.Value[2])
	id := ConnectionId{
		Remote: src.Address,
		Port:   src.Port,
		Conv:   conv,
	}
	conn, found := this.sessions[id]

	if !found {
		if cmd == CommandTerminate {
			return
		}
		writer := &Writer{
			id:       id,
			hub:      this.hub,
			dest:     src,
			listener: this,
		}
		remoteAddr := &net.UDPAddr{
			IP:   src.Address.IP(),
			Port: int(src.Port),
		}
		localAddr := this.hub.Addr()
		auth, err := this.config.GetAuthenticator()
		if err != nil {
			log.Error("KCP|Listener: Failed to create authenticator: ", err)
		}
		sConn := &ServerConnection{
			id:     internal.NewConnectionId(v2net.LocalHostIP, src),
			local:  localAddr,
			remote: remoteAddr,
			writer: writer,
		}
		conn = NewConnection(conv, sConn, this, auth, this.config)
		select {
		case this.awaitingConns <- conn:
		case <-time.After(time.Second * 5):
			conn.Close()
			return
		}
		this.sessions[id] = conn
	}
	conn.Input(payload.Value)
}

func (this *Listener) Remove(id ConnectionId) {
	if !this.running {
		return
	}
	this.Lock()
	defer this.Unlock()
	if !this.running {
		return
	}
	delete(this.sessions, id)
}

// Accept implements the Accept method in the Listener interface; it waits for the next call and returns a generic Conn.
func (this *Listener) Accept() (internet.Connection, error) {
	for {
		if !this.running {
			return nil, ErrClosedListener
		}
		select {
		case conn, open := <-this.awaitingConns:
			if !open {
				break
			}
			if this.tlsConfig != nil {
				tlsConn := tls.Server(conn, this.tlsConfig)
				return v2tls.NewConnection(tlsConn), nil
			}
			return conn, nil
		case <-time.After(time.Second):

		}
	}
}

// Close stops listening on the UDP address. Already Accepted connections are not closed.
func (this *Listener) Close() error {
	if !this.running {
		return ErrClosedListener
	}
	this.Lock()
	defer this.Unlock()

	this.running = false
	close(this.awaitingConns)
	for _, conn := range this.sessions {
		go conn.Terminate()
	}
	this.hub.Close()

	return nil
}

func (this *Listener) ActiveConnections() int {
	this.Lock()
	defer this.Unlock()

	return len(this.sessions)
}

// Addr returns the listener's network address, The Addr returned is shared by all invocations of Addr, so do not modify it.
func (this *Listener) Addr() net.Addr {
	return this.hub.Addr()
}

func (this *Listener) Put(internal.ConnectionId, net.Conn) {}

type Writer struct {
	id       ConnectionId
	dest     v2net.Destination
	hub      *udp.UDPHub
	listener *Listener
}

func (this *Writer) Write(payload []byte) (int, error) {
	return this.hub.WriteTo(payload, this.dest)
}

func (this *Writer) Close() error {
	this.listener.Remove(this.id)
	return nil
}

func ListenKCP(address v2net.Address, port v2net.Port, options internet.ListenOptions) (internet.Listener, error) {
	return NewListener(address, port, options)
}

func init() {
	internet.KCPListenFunc = ListenKCP
}
