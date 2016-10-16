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
	v2tls "v2ray.com/core/transport/internet/tls"
	"v2ray.com/core/transport/internet/udp"
)

// Listener defines a server listening for connections
type Listener struct {
	sync.Mutex
	running       bool
	authenticator internet.Authenticator
	sessions      map[string]*Connection
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

	auth, err := kcpSettings.GetAuthenticator()
	if err != nil {
		return nil, err
	}
	l := &Listener{
		authenticator: auth,
		sessions:      make(map[string]*Connection),
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
	hub, err := udp.ListenUDP(address, port, udp.ListenOption{Callback: l.OnReceive})
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
	sourceId := src.NetAddr() + "|" + serial.Uint16ToString(conv)
	conn, found := this.sessions[sourceId]
	if !found {
		if cmd == CommandTerminate {
			return
		}
		log.Debug("KCP|Listener: Creating session with id(", sourceId, ") from ", src)
		writer := &Writer{
			id:       sourceId,
			hub:      this.hub,
			dest:     src,
			listener: this,
		}
		srcAddr := &net.UDPAddr{
			IP:   src.Address.IP(),
			Port: int(src.Port),
		}
		auth, err := this.config.GetAuthenticator()
		if err != nil {
			log.Error("KCP|Listener: Failed to create authenticator: ", err)
		}
		conn = NewConnection(conv, writer, this.Addr().(*net.UDPAddr), srcAddr, auth, this.config)
		select {
		case this.awaitingConns <- conn:
		case <-time.After(time.Second * 5):
			conn.Close()
			return
		}
		this.sessions[sourceId] = conn
	}
	conn.Input(payload.Value)
}

func (this *Listener) Remove(dest string) {
	if !this.running {
		return
	}
	this.Lock()
	defer this.Unlock()
	if !this.running {
		return
	}
	log.Debug("KCP|Listener: Removing session ", dest)
	delete(this.sessions, dest)
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

type Writer struct {
	id       string
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
