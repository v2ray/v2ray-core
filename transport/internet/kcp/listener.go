package kcp

import (
	"context"
	"crypto/cipher"
	"crypto/tls"
	"io"
	"net"
	"sync"
	"time"

	"v2ray.com/core/app/log"
	"v2ray.com/core/common"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/errors"
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/transport/internet"
	"v2ray.com/core/transport/internet/internal"
	v2tls "v2ray.com/core/transport/internet/tls"
	"v2ray.com/core/transport/internet/udp"
)

type ConnectionID struct {
	Remote v2net.Address
	Port   v2net.Port
	Conv   uint16
}

type ServerConnection struct {
	id     internal.ConnectionID
	local  net.Addr
	remote net.Addr
	writer PacketWriter
	closer io.Closer
}

func (o *ServerConnection) Overhead() int {
	return o.writer.Overhead()
}

func (o *ServerConnection) Read([]byte) (int, error) {
	panic("KCP|ServerConnection: Read should not be called.")
}

func (o *ServerConnection) Write(b []byte) (int, error) {
	return o.writer.Write(b)
}

func (o *ServerConnection) Close() error {
	return o.closer.Close()
}

func (o *ServerConnection) Reset(input func([]Segment)) {
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

func (o *ServerConnection) Id() internal.ConnectionID {
	return o.id
}

// Listener defines a server listening for connections
type Listener struct {
	sync.Mutex
	closed    chan bool
	sessions  map[ConnectionID]*Connection
	hub       *udp.Hub
	tlsConfig *tls.Config
	config    *Config
	reader    PacketReader
	header    internet.PacketHeader
	security  cipher.AEAD
	conns     chan<- internet.Connection
}

func NewListener(ctx context.Context, address v2net.Address, port v2net.Port, conns chan<- internet.Connection) (*Listener, error) {
	networkSettings := internet.TransportSettingsFromContext(ctx)
	kcpSettings := networkSettings.(*Config)
	kcpSettings.ConnectionReuse = &ConnectionReuse{Enable: false}

	header, err := kcpSettings.GetPackerHeader()
	if err != nil {
		return nil, errors.Base(err).Message("KCP|Listener: Failed to create packet header.")
	}
	security, err := kcpSettings.GetSecurity()
	if err != nil {
		return nil, errors.Base(err).Message("KCP|Listener: Failed to create security.")
	}
	l := &Listener{
		header:   header,
		security: security,
		reader: &KCPPacketReader{
			Header:   header,
			Security: security,
		},
		sessions: make(map[ConnectionID]*Connection),
		closed:   make(chan bool),
		config:   kcpSettings,
		conns:    conns,
	}
	securitySettings := internet.SecuritySettingsFromContext(ctx)
	if securitySettings != nil {
		switch securitySettings := securitySettings.(type) {
		case *v2tls.Config:
			l.tlsConfig = securitySettings.GetTLSConfig()
		}
	}
	hub, err := udp.ListenUDP(address, port, udp.ListenOption{Callback: l.OnReceive, Concurrency: 2})
	if err != nil {
		return nil, err
	}
	l.Lock()
	l.hub = hub
	l.Unlock()
	log.Info("KCP|Listener: listening on ", address, ":", port)
	return l, nil
}

func (v *Listener) OnReceive(payload *buf.Buffer, src v2net.Destination, originalDest v2net.Destination) {
	defer payload.Release()

	segments := v.reader.Read(payload.Bytes())
	if len(segments) == 0 {
		log.Info("KCP|Listener: discarding invalid payload from ", src)
		return
	}

	select {
	case <-v.closed:
		return
	default:
	}

	v.Lock()
	defer v.Unlock()
	if v.hub == nil {
		return
	}
	if payload.Len() < 4 {
		return
	}
	conv := segments[0].Conversation()
	cmd := segments[0].Command()

	id := ConnectionID{
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
		sConn := &ServerConnection{
			id:     internal.NewConnectionID(v2net.LocalHostIP, src),
			local:  localAddr,
			remote: remoteAddr,
			writer: &KCPPacketWriter{
				Header:   v.header,
				Writer:   writer,
				Security: v.security,
			},
			closer: writer,
		}
		conn = NewConnection(conv, sConn, v, v.config)
		var netConn internet.Connection = conn
		if v.tlsConfig != nil {
			tlsConn := tls.Server(conn, v.tlsConfig)
			netConn = UnreusableConnection{Conn: tlsConn}
		}

		select {
		case v.conns <- netConn:
		case <-time.After(time.Second * 5):
			conn.Close()
			return
		}
		v.sessions[id] = conn
	}
	conn.Input(segments)
}

func (v *Listener) Remove(id ConnectionID) {
	select {
	case <-v.closed:
		return
	default:
		v.Lock()
		delete(v.sessions, id)
		v.Unlock()
	}
}

// Close stops listening on the UDP address. Already Accepted connections are not closed.
func (v *Listener) Close() error {

	v.Lock()
	defer v.Unlock()
	select {
	case <-v.closed:
		return ErrClosedListener
	default:
	}

	close(v.closed)
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

func (v *Listener) Put(internal.ConnectionID, net.Conn) {}

type Writer struct {
	id       ConnectionID
	dest     v2net.Destination
	hub      *udp.Hub
	listener *Listener
}

func (v *Writer) Write(payload []byte) (int, error) {
	return v.hub.WriteTo(payload, v.dest)
}

func (v *Writer) Close() error {
	v.listener.Remove(v.id)
	return nil
}

func ListenKCP(ctx context.Context, address v2net.Address, port v2net.Port, conns chan<- internet.Connection) (internet.Listener, error) {
	return NewListener(ctx, address, port, conns)
}

func init() {
	common.Must(internet.RegisterTransportListener(internet.TransportProtocol_MKCP, ListenKCP))
}
