package websocket

import (
	"crypto/tls"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"v2ray.com/core/app/log"
	"v2ray.com/core/common"
	"v2ray.com/core/common/errors"
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/transport/internet"
	"v2ray.com/core/transport/internet/internal"
	v2tls "v2ray.com/core/transport/internet/tls"
)

var (
	ErrClosedListener = errors.New("Listener is closed.")
)

type ConnectionWithError struct {
	conn net.Conn
	err  error
}

type requestHandler struct {
	path  string
	conns chan *ConnectionWithError
}

func (h *requestHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if request.URL.Path != h.path {
		writer.WriteHeader(http.StatusNotFound)
		return
	}
	conn, err := converttovws(writer, request)
	if err != nil {
		log.Info("WebSocket|Listener: Failed to convert to WebSocket connection: ", err)
		return
	}

	select {
	case h.conns <- &ConnectionWithError{conn: conn}:
	default:
		conn.Close()
	}
}

type Listener struct {
	sync.Mutex
	closed        chan bool
	awaitingConns chan *ConnectionWithError
	listener      net.Listener
	tlsConfig     *tls.Config
	config        *Config
}

func ListenWS(address v2net.Address, port v2net.Port, options internet.ListenOptions) (internet.Listener, error) {
	networkSettings, err := options.Stream.GetEffectiveTransportSettings()
	if err != nil {
		return nil, err
	}
	wsSettings := networkSettings.(*Config)

	l := &Listener{
		closed:        make(chan bool),
		awaitingConns: make(chan *ConnectionWithError, 32),
		config:        wsSettings,
	}
	if options.Stream != nil && options.Stream.HasSecuritySettings() {
		securitySettings, err := options.Stream.GetEffectiveSecuritySettings()
		if err != nil {
			return nil, errors.Base(err).Message("WebSocket: Failed to create apply TLS config.")
		}
		tlsConfig, ok := securitySettings.(*v2tls.Config)
		if ok {
			l.tlsConfig = tlsConfig.GetTLSConfig()
		}
	}

	err = l.listenws(address, port)

	return l, err
}

func (ln *Listener) listenws(address v2net.Address, port v2net.Port) error {
	netAddr := address.String() + ":" + strconv.Itoa(int(port.Value()))
	var listener net.Listener
	if ln.tlsConfig == nil {
		l, err := net.Listen("tcp", netAddr)
		if err != nil {
			return errors.Base(err).Message("WebSocket|Listener: Failed to listen TCP ", netAddr)
		}
		listener = l
	} else {
		l, err := tls.Listen("tcp", netAddr, ln.tlsConfig)
		if err != nil {
			return errors.Base(err).Message("WebSocket|Listener: Failed to listen TLS ", netAddr)
		}
		listener = l
	}
	ln.listener = listener

	go func() {
		http.Serve(listener, &requestHandler{
			path:  ln.config.GetNormailzedPath(),
			conns: ln.awaitingConns,
		})
	}()

	return nil
}

func converttovws(w http.ResponseWriter, r *http.Request) (*connection, error) {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  32 * 1024,
		WriteBufferSize: 32 * 1024,
	}
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		return nil, err
	}

	return &connection{wsc: conn}, nil
}

func (ln *Listener) Accept() (internet.Connection, error) {
	for {
		select {
		case <-ln.closed:
			return nil, ErrClosedListener
		case connErr, open := <-ln.awaitingConns:
			if !open {
				return nil, ErrClosedListener
			}
			if connErr.err != nil {
				return nil, connErr.err
			}
			return internal.NewConnection(internal.ConnectionID{}, connErr.conn, ln, internal.ReuseConnection(ln.config.IsConnectionReuse())), nil
		case <-time.After(time.Second * 2):
		}
	}
}

func (ln *Listener) Put(id internal.ConnectionID, conn net.Conn) {
	ln.Lock()
	defer ln.Unlock()
	select {
	case <-ln.closed:
		return
	default:
	}
	select {
	case ln.awaitingConns <- &ConnectionWithError{conn: conn}:
	default:
		conn.Close()
	}
}

func (ln *Listener) Addr() net.Addr {
	return ln.listener.Addr()
}

func (ln *Listener) Close() error {
	ln.Lock()
	defer ln.Unlock()
	select {
	case <-ln.closed:
		return ErrClosedListener
	default:
	}
	close(ln.closed)
	ln.listener.Close()
	close(ln.awaitingConns)
	for connErr := range ln.awaitingConns {
		if connErr.conn != nil {
			connErr.conn.Close()
		}
	}
	return nil
}

func init() {
	common.Must(internet.RegisterTransportListener(internet.TransportProtocol_WebSocket, ListenWS))
}
