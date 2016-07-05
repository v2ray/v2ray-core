package kcp

import (
	"net"
	"sync"
	"time"

	"github.com/v2ray/v2ray-core/common/alloc"
	"github.com/v2ray/v2ray-core/common/log"
	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/common/serial"
	"github.com/v2ray/v2ray-core/transport/internet"
	"github.com/v2ray/v2ray-core/transport/internet/udp"
)

// Listener defines a server listening for connections
type Listener struct {
	sync.Mutex
	running       bool
	block         Authenticator
	sessions      map[string]*Connection
	awaitingConns chan *Connection
	hub           *udp.UDPHub
	localAddr     *net.UDPAddr
}

func NewListener(address v2net.Address, port v2net.Port) (*Listener, error) {
	l := &Listener{
		block:         NewSimpleAuthenticator(),
		sessions:      make(map[string]*Connection),
		awaitingConns: make(chan *Connection, 64),
		localAddr: &net.UDPAddr{
			IP:   address.IP(),
			Port: int(port),
		},
		running: true,
	}
	hub, err := udp.ListenUDP(address, port, l.OnReceive)
	if err != nil {
		return nil, err
	}
	l.hub = hub
	log.Info("KCP|Listener: listening on ", address, ":", port)
	return l, nil
}

func (this *Listener) OnReceive(payload *alloc.Buffer, src v2net.Destination) {
	defer payload.Release()

	if valid := this.block.Open(payload); !valid {
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
	srcAddrStr := src.NetAddr()
	conn, found := this.sessions[srcAddrStr]
	if !found {
		conv := serial.BytesToUint16(payload.Value)
		writer := &Writer{
			hub:      this.hub,
			dest:     src,
			listener: this,
		}
		srcAddr := &net.UDPAddr{
			IP:   src.Address().IP(),
			Port: int(src.Port()),
		}
		conn = NewConnection(conv, writer, this.localAddr, srcAddr, this.block)
		select {
		case this.awaitingConns <- conn:
		case <-time.After(time.Second * 5):
			conn.Close()
			return
		}
		this.sessions[srcAddrStr] = conn
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
	delete(this.sessions, dest)
}

// Accept implements the Accept method in the Listener interface; it waits for the next call and returns a generic Conn.
func (this *Listener) Accept() (internet.Connection, error) {
	for {
		if !this.running {
			return nil, errClosedListener
		}
		select {
		case conn := <-this.awaitingConns:
			return conn, nil
		case <-time.After(time.Second):

		}
	}
}

// Close stops listening on the UDP address. Already Accepted connections are not closed.
func (this *Listener) Close() error {
	if !this.running {
		return errClosedListener
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

// Addr returns the listener's network address, The Addr returned is shared by all invocations of Addr, so do not modify it.
func (this *Listener) Addr() net.Addr {
	return this.localAddr
}

type Writer struct {
	dest     v2net.Destination
	hub      *udp.UDPHub
	listener *Listener
}

func (this *Writer) Write(payload []byte) (int, error) {
	return this.hub.WriteTo(payload, this.dest)
}

func (this *Writer) Close() error {
	this.listener.Remove(this.dest.NetAddr())
	return nil
}

func ListenKCP(address v2net.Address, port v2net.Port) (internet.Listener, error) {
	return NewListener(address, port)
}

func init() {
	internet.KCPListenFunc = ListenKCP
}
