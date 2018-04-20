package http

import (
	"context"
	"io"
	"net/http"
	"strings"

	"v2ray.com/core/common"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/serial"
	"v2ray.com/core/common/signal"
	"v2ray.com/core/transport/internet"
	"v2ray.com/core/transport/internet/tls"
)

type Listener struct {
	server  *http.Server
	handler internet.ConnHandler
	local   net.Addr
	config  Config
}

func (l *Listener) Addr() net.Addr {
	return l.local
}

func (l *Listener) Close() error {
	return l.server.Close()
}

type flushWriter struct {
	w io.Writer
	d *signal.Done
}

func (fw flushWriter) Write(p []byte) (n int, err error) {
	if fw.d.Done() {
		return 0, io.ErrClosedPipe
	}

	n, err = fw.w.Write(p)
	if f, ok := fw.w.(http.Flusher); ok {
		f.Flush()
	}
	return
}

func (l *Listener) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	host := request.Host
	if !l.config.isValidHost(host) {
		writer.WriteHeader(404)
		return
	}
	path := l.config.getNormalizedPath()
	if !strings.HasPrefix(request.URL.Path, path) {
		writer.WriteHeader(404)
		return
	}

	writer.Header().Set("Cache-Control", "no-store")
	writer.WriteHeader(200)
	if f, ok := writer.(http.Flusher); ok {
		f.Flush()
	}
	done := signal.NewDone()
	conn := net.NewConnection(
		net.ConnectionOutput(request.Body),
		net.ConnectionInput(flushWriter{w: writer, d: done}),
		net.ConnectionOnClose(common.NewChainedClosable(done, request.Body)),
		net.ConnectionLocalAddr(l.Addr()),
		net.ConnectionRemoteAddr(l.Addr()),
	)
	l.handler(conn)
	<-done.Wait()
}

func Listen(ctx context.Context, address net.Address, port net.Port, handler internet.ConnHandler) (internet.Listener, error) {
	rawSettings := internet.TransportSettingsFromContext(ctx)
	httpSettings, ok := rawSettings.(*Config)
	if !ok {
		return nil, newError("HTTP config is not set.").AtError()
	}

	listener := &Listener{
		handler: handler,
		local: &net.TCPAddr{
			IP:   address.IP(),
			Port: int(port),
		},
		config: *httpSettings,
	}

	config := tls.ConfigFromContext(ctx)
	if config == nil {
		return nil, newError("TLS must be enabled for http transport.").AtWarning()
	}

	server := &http.Server{
		Addr:      serial.Concat(address, ":", port),
		TLSConfig: config.GetTLSConfig(tls.WithNextProto("h2")),
		Handler:   listener,
	}

	listener.server = server
	go func() {
		err := server.ListenAndServeTLS("", "")
		if err != nil {
			newError("stoping serving TLS").Base(err).WriteToLog()
		}
	}()

	return listener, nil
}

func init() {
	common.Must(internet.RegisterTransportListener(internet.TransportProtocol_HTTP, Listen))
}
