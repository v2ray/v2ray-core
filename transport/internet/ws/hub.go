package ws

import (
	"errors"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/v2ray/v2ray-core/common/log"
	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/transport/internet"
)

var (
	ErrClosedListener = errors.New("Listener is closed.")
)

type ConnectionWithError struct {
	conn net.Conn
	err  error
}

type WSListener struct {
	sync.Mutex
	acccepting    bool
	awaitingConns chan *ConnectionWithError
	listener      *StoppableListener
}

func ListenWS(address v2net.Address, port v2net.Port) (internet.Listener, error) {

	l := &WSListener{
		acccepting:    true,
		awaitingConns: make(chan *ConnectionWithError, 32),
	}

	err := l.listenws(address, port)

	return l, err
}

func (wsl *WSListener) listenws(address v2net.Address, port v2net.Port) error {

	http.HandleFunc("/"+effectiveConfig.Path, func(w http.ResponseWriter, r *http.Request) {
		con, err := wsl.converttovws(w, r)
		if err != nil {
			log.Warning("WebSocket|Listener: Failed to convert connection: ", err)
			return
		}

		select {
		case wsl.awaitingConns <- &ConnectionWithError{
			conn: con,
		}:
		default:
			if con != nil {
				con.Close()
			}
		}
		return
	})

	errchan := make(chan error)

	listenerfunc := func() error {
		ol, err := net.Listen("tcp", address.String()+":"+strconv.Itoa(int(port.Value())))
		if err != nil {
			return err
		}
		wsl.listener, err = NewStoppableListener(ol)
		if err != nil {
			return err
		}
		return http.Serve(wsl.listener, nil)
	}

	if effectiveConfig.Pto == "wss" {
		listenerfunc = func() error {
			var err error
			wsl.listener, err = getstopableTLSlistener(effectiveConfig.Cert, effectiveConfig.PrivKey, address.String()+":"+strconv.Itoa(int(port.Value())))
			if err != nil {
				return err
			}
			return http.Serve(wsl.listener, nil)
		}
	}

	go func() {
		err := listenerfunc()
		errchan <- err
		return
	}()

	var err error
	select {
	case err = <-errchan:
	case <-time.After(time.Second * 2):
		//Should this listen fail after 2 sec, it could gone untracked.
	}

	if err != nil {
		log.Error("WebSocket|Listener: Failed to serve: ", err)
	}

	return err

}

func (wsl *WSListener) converttovws(w http.ResponseWriter, r *http.Request) (*wsconn, error) {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  65536,
		WriteBufferSize: 65536,
	}
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		return nil, err
	}

	wrapedConn := &wsconn{wsc: conn, connClosing: false}
	wrapedConn.setup()
	return wrapedConn, nil
}

func (this *WSListener) Accept() (internet.Connection, error) {
	for this.acccepting {
		select {
		case connErr, open := <-this.awaitingConns:
			if !open {
				return nil, ErrClosedListener
			}
			if connErr.err != nil {
				return nil, connErr.err
			}
			return NewConnection("", connErr.conn.(*wsconn), this), nil
		case <-time.After(time.Second * 2):
		}
	}
	return nil, ErrClosedListener
}

func (this *WSListener) Recycle(dest string, conn *wsconn) {
	this.Lock()
	defer this.Unlock()
	if !this.acccepting {
		return
	}
	select {
	case this.awaitingConns <- &ConnectionWithError{conn: conn}:
	default:
		conn.Close()
	}
}

func (this *WSListener) Addr() net.Addr {
	return nil
}

func (this *WSListener) Close() error {
	this.Lock()
	defer this.Unlock()
	this.acccepting = false

	this.listener.Stop()

	close(this.awaitingConns)
	for connErr := range this.awaitingConns {
		if connErr.conn != nil {
			go connErr.conn.Close()
		}
	}
	return nil
}

func init() {
	internet.WSListenFunc = ListenWS
}
