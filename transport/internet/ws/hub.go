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
		log.Warning("WS:WSListener->listenws->(HandleFunc,lambda 2)! Accepting websocket")
		con, err := wsl.converttovws(w, r)
		if err != nil {
			log.Warning("WS:WSListener->listenws->(HandleFunc,lambda 2)!" + err.Error())
			return
		}

		select {
		case wsl.awaitingConns <- &ConnectionWithError{
			conn: con,
			err:  err,
		}:
			log.Warning("WS:WSListener->listenws->(HandleFunc,lambda 2)! transferd websocket")
		default:
			if con != nil {
				con.Close()
			}
		}
		//con.retloc.Wait()
		return

	})

	errchan := make(chan error)

	listenerfunc := func() error {
		return http.ListenAndServe(address.String()+":"+strconv.Itoa(int(port.Value())), nil)
	}

	if effectiveConfig.Pto == "wss" {
		listenerfunc = func() error {
			return http.ListenAndServeTLS(address.String()+":"+strconv.Itoa(int(port.Value())), effectiveConfig.Cert, effectiveConfig.PrivKey, nil)
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
		log.Error("WS:WSListener->listenws->ListenAndServe!" + err.Error())
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
			log.Info("WSListener: conn accepted")
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

	log.Warning("WSListener: Yet to support close listening HTTP service")

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
