package websocket

import (
	"crypto/tls"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"

	"v2ray.com/core/common/errors"
	"v2ray.com/core/common/log"
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/transport/internet"
	v2tls "v2ray.com/core/transport/internet/tls"

	"github.com/gorilla/websocket"
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
	tlsConfig     *tls.Config
	config        *Config
}

func ListenWS(address v2net.Address, port v2net.Port, options internet.ListenOptions) (internet.Listener, error) {
	networkSettings, err := options.Stream.GetEffectiveNetworkSettings()
	if err != nil {
		return nil, err
	}
	wsSettings := networkSettings.(*Config)

	l := &WSListener{
		acccepting:    true,
		awaitingConns: make(chan *ConnectionWithError, 32),
		config:        wsSettings,
	}
	if options.Stream != nil && options.Stream.HasSecuritySettings() {
		securitySettings, err := options.Stream.GetEffectiveSecuritySettings()
		if err != nil {
			log.Error("WebSocket: Failed to create apply TLS config: ", err)
			return nil, err
		}
		tlsConfig, ok := securitySettings.(*v2tls.Config)
		if ok {
			l.tlsConfig = tlsConfig.GetTLSConfig()
		}
	}

	err = l.listenws(address, port)

	return l, err
}

func (wsl *WSListener) listenws(address v2net.Address, port v2net.Port) error {
	http.HandleFunc("/"+wsl.config.Path, func(w http.ResponseWriter, r *http.Request) {
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

	if wsl.tlsConfig != nil {
		listenerfunc = func() error {
			var err error
			wsl.listener, err = getstopableTLSlistener(wsl.tlsConfig, address.String()+":"+strconv.Itoa(int(port.Value())))
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

func (v *WSListener) Accept() (internet.Connection, error) {
	for v.acccepting {
		select {
		case connErr, open := <-v.awaitingConns:
			if !open {
				return nil, ErrClosedListener
			}
			if connErr.err != nil {
				return nil, connErr.err
			}
			return NewConnection("", connErr.conn.(*wsconn), v, v.config), nil
		case <-time.After(time.Second * 2):
		}
	}
	return nil, ErrClosedListener
}

func (v *WSListener) Recycle(dest string, conn *wsconn) {
	v.Lock()
	defer v.Unlock()
	if !v.acccepting {
		return
	}
	select {
	case v.awaitingConns <- &ConnectionWithError{conn: conn}:
	default:
		conn.Close()
	}
}

func (v *WSListener) Addr() net.Addr {
	return nil
}

func (v *WSListener) Close() error {
	v.Lock()
	defer v.Unlock()
	v.acccepting = false

	v.listener.Stop()

	close(v.awaitingConns)
	for connErr := range v.awaitingConns {
		if connErr.conn != nil {
			go connErr.conn.Close()
		}
	}
	return nil
}

func init() {
	internet.WSListenFunc = ListenWS
}
