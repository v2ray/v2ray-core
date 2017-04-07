package websocket

import (
	"context"
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
	v2tls "v2ray.com/core/transport/internet/tls"
)

var (
	ErrClosedListener = errors.New("Listener is closed.")
)

type requestHandler struct {
	path string
	ln   *Listener
}

func (h *requestHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if request.URL.Path != h.path {
		writer.WriteHeader(http.StatusNotFound)
		return
	}
	conn, err := converttovws(writer, request)
	if err != nil {
		log.Trace(errors.New("WebSocket|Listener: Failed to convert to WebSocket connection: ", err))
		return
	}

	select {
	case <-h.ln.ctx.Done():
		conn.Close()
	case h.ln.conns <- internet.Connection(conn):
	case <-time.After(time.Second * 5):
		conn.Close()
	}
}

type Listener struct {
	sync.Mutex
	ctx       context.Context
	listener  net.Listener
	tlsConfig *tls.Config
	config    *Config
	conns     chan<- internet.Connection
}

func ListenWS(ctx context.Context, address v2net.Address, port v2net.Port, conns chan<- internet.Connection) (internet.Listener, error) {
	networkSettings := internet.TransportSettingsFromContext(ctx)
	wsSettings := networkSettings.(*Config)

	l := &Listener{
		ctx:    ctx,
		config: wsSettings,
		conns:  conns,
	}
	if securitySettings := internet.SecuritySettingsFromContext(ctx); securitySettings != nil {
		tlsConfig, ok := securitySettings.(*v2tls.Config)
		if ok {
			l.tlsConfig = tlsConfig.GetTLSConfig()
		}
	}

	err := l.listenws(address, port)

	return l, err
}

func (ln *Listener) listenws(address v2net.Address, port v2net.Port) error {
	netAddr := address.String() + ":" + strconv.Itoa(int(port.Value()))
	var listener net.Listener
	if ln.tlsConfig == nil {
		l, err := net.Listen("tcp", netAddr)
		if err != nil {
			return errors.New("failed to listen TCP ", netAddr).Base(err).Path("WebSocket", "Listener")
		}
		listener = l
	} else {
		l, err := tls.Listen("tcp", netAddr, ln.tlsConfig)
		if err != nil {
			return errors.New("failed to listen TLS ", netAddr).Base(err).Path("WebSocket", "Listener")
		}
		listener = l
	}
	ln.listener = listener

	go func() {
		http.Serve(listener, &requestHandler{
			path: ln.config.GetNormailzedPath(),
			ln:   ln,
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

func (ln *Listener) Addr() net.Addr {
	return ln.listener.Addr()
}

func (ln *Listener) Close() error {
	return ln.listener.Close()
}

func init() {
	common.Must(internet.RegisterTransportListener(internet.TransportProtocol_WebSocket, ListenWS))
}
