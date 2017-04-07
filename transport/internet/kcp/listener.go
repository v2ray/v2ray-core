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
	v2tls "v2ray.com/core/transport/internet/tls"
	"v2ray.com/core/transport/internet/udp"
)

type ConnectionID struct {
	Remote v2net.Address
	Port   v2net.Port
	Conv   uint16
}

type ServerConnection struct {
	local  net.Addr
	remote net.Addr
	writer PacketWriter
	closer io.Closer
}

func (c *ServerConnection) Overhead() int {
	return c.writer.Overhead()
}

func (*ServerConnection) Read([]byte) (int, error) {
	panic("KCP|ServerConnection: Read should not be called.")
}

func (c *ServerConnection) Write(b []byte) (int, error) {
	return c.writer.Write(b)
}

func (c *ServerConnection) Close() error {
	return c.closer.Close()
}

func (*ServerConnection) Reset(input func([]Segment)) {
}

func (c *ServerConnection) LocalAddr() net.Addr {
	return c.local
}

func (c *ServerConnection) RemoteAddr() net.Addr {
	return c.remote
}

func (*ServerConnection) SetDeadline(time.Time) error {
	return nil
}

func (*ServerConnection) SetReadDeadline(time.Time) error {
	return nil
}

func (*ServerConnection) SetWriteDeadline(time.Time) error {
	return nil
}

// Listener defines a server listening for connections
type Listener struct {
	sync.Mutex
	ctx       context.Context
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

	header, err := kcpSettings.GetPackerHeader()
	if err != nil {
		return nil, errors.New("KCP|Listener: Failed to create packet header.").Base(err)
	}
	security, err := kcpSettings.GetSecurity()
	if err != nil {
		return nil, errors.New("KCP|Listener: Failed to create security.").Base(err)
	}
	l := &Listener{
		header:   header,
		security: security,
		reader: &KCPPacketReader{
			Header:   header,
			Security: security,
		},
		sessions: make(map[ConnectionID]*Connection),
		ctx:      ctx,
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
	log.Trace(errors.New("KCP|Listener: listening on ", address, ":", port))
	return l, nil
}

func (v *Listener) OnReceive(payload *buf.Buffer, src v2net.Destination, originalDest v2net.Destination) {
	defer payload.Release()

	segments := v.reader.Read(payload.Bytes())
	if len(segments) == 0 {
		log.Trace(errors.New("KCP|Listener: discarding invalid payload from ", src))
		return
	}

	v.Lock()
	defer v.Unlock()

	select {
	case <-v.ctx.Done():
		return
	default:
	}

	if v.hub == nil {
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
			local:  localAddr,
			remote: remoteAddr,
			writer: &KCPPacketWriter{
				Header:   v.header,
				Writer:   writer,
				Security: v.security,
			},
			closer: writer,
		}
		conn = NewConnection(conv, sConn, v.config)
		var netConn internet.Connection = conn
		if v.tlsConfig != nil {
			tlsConn := tls.Server(conn, v.tlsConfig)
			netConn = tlsConn
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
	case <-v.ctx.Done():
		return
	default:
		v.Lock()
		delete(v.sessions, id)
		v.Unlock()
	}
}

// Close stops listening on the UDP address. Already Accepted connections are not closed.
func (v *Listener) Close() error {
	v.hub.Close()

	v.Lock()
	defer v.Unlock()

	for _, conn := range v.sessions {
		go conn.Terminate()
	}

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
