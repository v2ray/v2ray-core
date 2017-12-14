package kcp

import (
	"context"
	"crypto/cipher"
	"crypto/tls"
	"sync"

	"v2ray.com/core/app/log"
	"v2ray.com/core/common"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/net"
	"v2ray.com/core/transport/internet"
	v2tls "v2ray.com/core/transport/internet/tls"
	"v2ray.com/core/transport/internet/udp"
)

type ConnectionID struct {
	Remote net.Address
	Port   net.Port
	Conv   uint16
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
	addConn   internet.AddConnection
}

func NewListener(ctx context.Context, address net.Address, port net.Port, addConn internet.AddConnection) (*Listener, error) {
	networkSettings := internet.TransportSettingsFromContext(ctx)
	kcpSettings := networkSettings.(*Config)

	header, err := kcpSettings.GetPackerHeader()
	if err != nil {
		return nil, newError("failed to create packet header").Base(err).AtError()
	}
	security, err := kcpSettings.GetSecurity()
	if err != nil {
		return nil, newError("failed to create security").Base(err).AtError()
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
		addConn:  addConn,
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
	log.Trace(newError("listening on ", address, ":", port))
	return l, nil
}

func (v *Listener) OnReceive(payload *buf.Buffer, src net.Destination, originalDest net.Destination) {
	defer payload.Release()

	segments := v.reader.Read(payload.Bytes())
	if len(segments) == 0 {
		log.Trace(newError("discarding invalid payload from ", src))
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
		conn = NewConnection(ConnMetadata{
			LocalAddr:    localAddr,
			RemoteAddr:   remoteAddr,
			Conversation: conv,
		}, &KCPPacketWriter{
			Header:   v.header,
			Security: v.security,
			Writer:   writer,
		}, writer, v.config)
		var netConn internet.Connection = conn
		if v.tlsConfig != nil {
			tlsConn := tls.Server(conn, v.tlsConfig)
			netConn = tlsConn
		}

		if !v.addConn(context.Background(), netConn) {
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
	dest     net.Destination
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

func ListenKCP(ctx context.Context, address net.Address, port net.Port, addConn internet.AddConnection) (internet.Listener, error) {
	return NewListener(ctx, address, port, addConn)
}

func init() {
	common.Must(internet.RegisterTransportListener(internet.TransportProtocol_MKCP, ListenKCP))
}
