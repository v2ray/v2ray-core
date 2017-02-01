package kcp

import (
	"crypto/cipher"
	"crypto/tls"
	"io"
	"net"
	"sync"
	"time"

	"v2ray.com/core/common"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/errors"
	"v2ray.com/core/app/log"
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
	running       bool
	sessions      map[ConnectionID]*Connection
	awaitingConns chan *Connection
	hub           *udp.Hub
	tlsConfig     *tls.Config
	config        *Config
	reader        PacketReader
	header        internet.PacketHeader
	security      cipher.AEAD
}

func NewListener(address v2net.Address, port v2net.Port, options internet.ListenOptions) (*Listener, error) {
	networkSettings, err := options.Stream.GetEffectiveTransportSettings()
	if err != nil {
		log.Error("KCP|Listener: Failed to get KCP settings: ", err)
		return nil, err
	}
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
		sessions:      make(map[ConnectionID]*Connection),
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

func (v *Listener) OnReceive(payload *buf.Buffer, src v2net.Destination, originalDest v2net.Destination) {
	defer payload.Release()

	segments := v.reader.Read(payload.Bytes())
	if len(segments) == 0 {
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
		select {
		case v.awaitingConns <- conn:
		case <-time.After(time.Second * 5):
			conn.Close()
			return
		}
		v.sessions[id] = conn
	}
	conn.Input(segments)
}

func (v *Listener) Remove(id ConnectionID) {
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
				return UnreusableConnection{Conn: tlsConn}, nil
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

func ListenKCP(address v2net.Address, port v2net.Port, options internet.ListenOptions) (internet.Listener, error) {
	return NewListener(address, port, options)
}

func init() {
	common.Must(internet.RegisterTransportListener(internet.TransportProtocol_MKCP, ListenKCP))
}
