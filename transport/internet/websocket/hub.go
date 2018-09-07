package websocket

import (
	"context"
	"crypto/tls"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"

	"v2ray.com/core/common"
	"v2ray.com/core/common/net"
	http_proto "v2ray.com/core/common/protocol/http"
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
		newError("failed to convert to WebSocket connection").Base(err).WriteToLog()
		return
	}

	forwardedAddrs := http_proto.ParseXForwardedFor(request.Header)
	remoteAddr := conn.RemoteAddr()
	if len(forwardedAddrs) > 0 && forwardedAddrs[0].Family().Either(net.AddressFamilyIPv4, net.AddressFamilyIPv6) {
		remoteAddr.(*net.TCPAddr).IP = forwardedAddrs[0].IP()
	}

	h.ln.addConn(newConnection(conn, remoteAddr))
}

type Listener struct {
	sync.Mutex
	listener  net.Listener
	tlsConfig *tls.Config
	config    *Config
	addConn   internet.ConnHandler
}

func ListenWS(ctx context.Context, address net.Address, port net.Port, addConn internet.ConnHandler) (internet.Listener, error) {
	networkSettings := internet.StreamSettingsFromContext(ctx)
	wsSettings := networkSettings.ProtocolSettings.(*Config)

	l := &Listener{
		config:  wsSettings,
		addConn: addConn,
	}
	if config := v2tls.ConfigFromContext(ctx); config != nil {
		l.tlsConfig = config.GetTLSConfig()
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
			path: ln.config.GetNormalizedPath(),
			ln:   ln,
		})
		if err != nil {
			newError("failed to serve http for WebSocket").Base(err).AtWarning().WriteToLog()
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
	common.Must(internet.RegisterTransportListener(protocolName, ListenWS))
}
