package http

import (
	"context"
	"io"
	"net/http"

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
}

func (l *Listener) Addr() net.Addr {
	return l.local
}

func (l *Listener) Close() error {
	return l.server.Shutdown(context.Background())
}

type flushWriter struct {
	w io.Writer
}

func (fw flushWriter) Write(p []byte) (n int, err error) {
	n, err = fw.w.Write(p)
	if f, ok := fw.w.(http.Flusher); ok {
		f.Flush()
	}
	return
}

func (l *Listener) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(200)
	if f, ok := writer.(http.Flusher); ok {
		f.Flush()
	}
	done := signal.NewDone()
	l.handler(&Connection{
		Reader: request.Body,
		Writer: flushWriter{writer},
		Closer: common.NewChainedClosable(request.Body, done),
		Local:  l.Addr(),
		Remote: l.Addr(),
	})
	<-done.C()
}

func Listen(ctx context.Context, address net.Address, port net.Port, handler internet.ConnHandler) (internet.Listener, error) {
	listener := &Listener{
		handler: handler,
		local: &net.TCPAddr{
			IP:   address.IP(),
			Port: int(port),
		},
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
	go server.ListenAndServeTLS("", "")

	return listener, nil
}

func init() {
	common.Must(internet.RegisterTransportListener(internet.TransportProtocol_HTTP, Listen))
}
