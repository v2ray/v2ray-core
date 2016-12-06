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
		o.input(b.Bytes())
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

func (v *Listener) OnReceive(payload *alloc.Buffer, session *proxy.SessionInfo) {
	defer payload.Release()

	src := session.Source

	if valid := v.authenticator.Open(payload); !valid {
		log.Info("KCP|Listener: discarding invalid payload from ", src)
		return
	}
	if !v.running {
		return
	}
	v.Lock()
	defer v.Unlock()
	if !v.running {
		return
	}
	if payload.Len() < 4 {
		return
	}
	conv := serial.BytesToUint16(payload.BytesTo(2))
	cmd := Command(payload.Byte(2))
	id := ConnectionId{
		Remote: src.Address,
		Port:   src.Port,
		Conv:   conv,
	}
	conn, found := v.sessions[id]

	if !found {
		if cmd == CommandTerminate {
			return
		}
		writer := &Writer{
			id:       id,
			hub:      v.hub,
			dest:     src,
			listener: v,
		}
		remoteAddr := &net.UDPAddr{
			IP:   src.Address.IP(),
			Port: int(src.Port),
		}
		localAddr := v.hub.Addr()
		auth, err := v.config.GetAuthenticator()
		if err != nil {
			log.Error("KCP|Listener: Failed to create authenticator: ", err)
		}
		sConn := &ServerConnection{
			id:     internal.NewConnectionId(v2net.LocalHostIP, src),
			local:  localAddr,
			remote: remoteAddr,
			writer: writer,
		}
		conn = NewConnection(conv, sConn, v, auth, v.config)
		select {
		case v.awaitingConns <- conn:
		case <-time.After(time.Second * 5):
			conn.Close()
			return
		}
		v.sessions[id] = conn
	}
	conn.Input(payload.Bytes())
}

func (v *Listener) Remove(id ConnectionId) {
	if !v.running {
		return
	}
	v.Lock()
	defer v.Unlock()
	if !v.running {
		return
	}
	delete(v.sessions, id)
}

// Accept implements the Accept method in the Listener interface; it waits for the next call and returns a generic Conn.
func (v *Listener) Accept() (internet.Connection, error) {
	for {
		if !v.running {
			return nil, ErrClosedListener
		}
		select {
		case conn, open := <-v.awaitingConns:
			if !open {
				break
			}
			if v.tlsConfig != nil {
				tlsConn := tls.Server(conn, v.tlsConfig)
				return v2tls.NewConnection(tlsConn), nil
			}
			return conn, nil
		case <-time.After(time.Second):

		}
	}
}

// Close stops listening on the UDP address. Already Accepted connections are not closed.
func (v *Listener) Close() error {
	if !v.running {
		return ErrClosedListener
	}
	v.Lock()
	defer v.Unlock()

	v.running = false
	close(v.awaitingConns)
	for _, conn := range v.sessions {
		go conn.Terminate()
	}
	v.hub.Close()

	return nil
}

func (v *Listener) ActiveConnections() int {
	v.Lock()
	defer v.Unlock()

	return len(v.sessions)
}

// Addr returns the listener's network address, The Addr returned is shared by all invocations of Addr, so do not modify it.
func (v *Listener) Addr() net.Addr {
	return v.hub.Addr()
}

func (v *Listener) Put(internal.ConnectionId, net.Conn) {}

type Writer struct {
	id       ConnectionId
	dest     v2net.Destination
	hub      *udp.UDPHub
	listener *Listener
}

func (v *Writer) Write(payload []byte) (int, error) {
	return v.hub.WriteTo(payload, v.dest)
}

func (v *Writer) Close() error {
	v.listener.Remove(v.id)
	return nil
}

func ListenKCP(address v2net.Address, port v2net.Port, options internet.ListenOptions) (internet.Listener, error) {
	return NewListener(address, port, options)
}

func init() {
	internet.KCPListenFunc = ListenKCP
}
