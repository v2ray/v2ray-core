package websocket

import (
	"context"
	"crypto/tls"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"v2ray.com/core/app/log"
	"v2ray.com/core/common"
	"v2ray.com/core/common/net"
	"v2ray.com/core/transport/internet"
	v2tls "v2ray.com/core/transport/internet/tls"
)

type requestHandler struct {
	path string
	ln   *Listener
}

var upgrader = &websocket.Upgrader{
	ReadBufferSize:   4 * 1024,
	WriteBufferSize:  4 * 1024,
	HandshakeTimeout: time.Second * 8,
}

func (h *requestHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if request.URL.Path != h.path {
		writer.WriteHeader(http.StatusNotFound)
		return
	}
	conn, err := upgrader.Upgrade(writer, request, nil)
	if err != nil {
		log.Trace(newError("failed to convert to WebSocket connection").Base(err))
		return
	}

	h.ln.addConn(h.ln.ctx, newConnection(conn))
}

type Listener struct {
	sync.Mutex
	ctx       context.Context
	listener  net.Listener
	tlsConfig *tls.Config
	config    *Config
	addConn   internet.AddConnection
}

func ListenWS(ctx context.Context, address net.Address, port net.Port, addConn internet.AddConnection) (internet.Listener, error) {
	networkSettings := internet.TransportSettingsFromContext(ctx)
	wsSettings := networkSettings.(*Config)

	l := &Listener{
		ctx:     ctx,
		config:  wsSettings,
		addConn: addConn,
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

func (ln *Listener) listenws(address net.Address, port net.Port) error {
	netAddr := address.String() + ":" + strconv.Itoa(int(port.Value()))
	var listener net.Listener
	if ln.tlsConfig == nil {
		l, err := net.Listen("tcp", netAddr)
		if err != nil {
			return newError("failed to listen TCP ", netAddr).Base(err)
		}
		listener = l
	} else {
		l, err := tls.Listen("tcp", netAddr, ln.tlsConfig)
		if err != nil {
			return newError("failed to listen TLS ", netAddr).Base(err)
		}
		listener = l
	}
	ln.listener = listener

	go func() {
		err := http.Serve(listener, &requestHandler{
			path: ln.config.GetNormailzedPath(),
			ln:   ln,
		})
		if err != nil {
			log.Trace(newError("failed to serve http for WebSocket").Base(err).AtWarning())
		}
	}()

	return nil
}

// Addr implements net.Listener.Addr().
func (ln *Listener) Addr() net.Addr {
	return ln.listener.Addr()
}

// Close implements net.Listener.Close().
func (ln *Listener) Close() error {
	return ln.listener.Close()
}

func init() {
	common.Must(internet.RegisterTransportListener(internet.TransportProtocol_WebSocket, ListenWS))
}
