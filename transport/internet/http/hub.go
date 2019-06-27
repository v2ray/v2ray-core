// +build !confonly

package http

import (
	"context"
	"io"
	"net/http"
	"strings"
	"time"

	"v2ray.com/core/common"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/serial"
	"v2ray.com/core/common/session"
	"v2ray.com/core/common/signal/done"
	"v2ray.com/core/transport/internet"
	"v2ray.com/core/transport/internet/tls"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
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
	d *done.Instance
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

	remoteAddr := l.Addr()
	dest, err := net.ParseDestination(request.RemoteAddr)
	if err != nil {
		newError("failed to parse request remote addr: ", request.RemoteAddr).Base(err).WriteToLog()
	} else {
		remoteAddr = &net.TCPAddr{
			IP:   dest.Address.IP(),
			Port: int(dest.Port),
		}
	}

	done := done.New()
	conn := net.NewConnection(
		net.ConnectionOutput(request.Body),
		net.ConnectionInput(flushWriter{w: writer, d: done}),
		net.ConnectionOnClose(common.ChainedClosable{done, request.Body}),
		net.ConnectionLocalAddr(l.Addr()),
		net.ConnectionRemoteAddr(remoteAddr),
	)
	l.handler(conn)
	<-done.Wait()
}

func Listen(ctx context.Context, address net.Address, port net.Port, streamSettings *internet.MemoryStreamConfig, handler internet.ConnHandler) (internet.Listener, error) {
	httpSettings := streamSettings.ProtocolSettings.(*Config)
	listener := &Listener{
		handler: handler,
		local: &net.TCPAddr{
			IP:   address.IP(),
			Port: int(port),
		},
		config: *httpSettings,
	}

	config := tls.ConfigFromStreamSettings(streamSettings)
	if config == nil {
		// return nil, newError("TLS must be enabled for http transport.").AtWarning()
		h2s:=&http2.Server{}

		server := &http.Server{
			Addr:              serial.Concat(address, ":", port),
			// TLSConfig:         config.GetTLSConfig(tls.WithNextProto("h2")),
			Handler:            h2c.NewHandler(listener,h2s),
			ReadHeaderTimeout: time.Second * 4,
		}

		listener.server = server
		go func() {
			tcpListener, err := internet.ListenSystem(ctx, &net.TCPAddr{
				IP:   address.IP(),
				Port: int(port),
			}, streamSettings.SocketSettings)
			if err != nil {
				newError("failed to listen on", address, ":", port).Base(err).WriteToLog(session.ExportIDToError(ctx))
				return
			}

			err = server.Serve(tcpListener)
			if err != nil {
				newError("stoping serving H2C").Base(err).WriteToLog(session.ExportIDToError(ctx))
			}
		}()

		return listener, nil
	} else {
		server := &http.Server{
			Addr:              serial.Concat(address, ":", port),
			TLSConfig:         config.GetTLSConfig(tls.WithNextProto("h2")),
			Handler:           listener,
			ReadHeaderTimeout: time.Second * 4,
		}

		listener.server = server
		go func() {
			tcpListener, err := internet.ListenSystem(ctx, &net.TCPAddr{
				IP:   address.IP(),
				Port: int(port),
			}, streamSettings.SocketSettings)
			if err != nil {
				newError("failed to listen on", address, ":", port).Base(err).WriteToLog(session.ExportIDToError(ctx))
				return
			}

			err = server.ServeTLS(tcpListener, "", "")
			if err != nil {
				newError("stoping serving TLS").Base(err).WriteToLog(session.ExportIDToError(ctx))
			}
		}()

		return listener, nil
	}
}

func init() {
	common.Must(internet.RegisterTransportListener(protocolName, Listen))
}
