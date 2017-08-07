package http

import (
	"bufio"
	"context"
	"io"
	"net"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"time"

	"v2ray.com/core/app"
	"v2ray.com/core/app/dispatcher"
	"v2ray.com/core/app/log"
	"v2ray.com/core/common"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/errors"
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/common/signal"
	"v2ray.com/core/transport/internet"
)

// Server is a HTTP proxy server.
type Server struct {
	config *ServerConfig
}

// NewServer creates a new HTTP inbound handler.
func NewServer(ctx context.Context, config *ServerConfig) (*Server, error) {
	space := app.SpaceFromContext(ctx)
	if space == nil {
		return nil, newError("no space in context.")
	}
	s := &Server{
		config: config,
	}
	return s, nil
}

func (*Server) Network() v2net.NetworkList {
	return v2net.NetworkList{
		Network: []v2net.Network{v2net.Network_TCP},
	}
}

func parseHost(rawHost string, defaultPort v2net.Port) (v2net.Destination, error) {
	port := defaultPort
	host, rawPort, err := net.SplitHostPort(rawHost)
	if err != nil {
		if addrError, ok := err.(*net.AddrError); ok && strings.Contains(addrError.Err, "missing port") {
			host = rawHost
		} else {
			return v2net.Destination{}, err
		}
	} else {
		intPort, err := strconv.Atoi(rawPort)
		if err != nil {
			return v2net.Destination{}, err
		}
		port = v2net.Port(intPort)
	}

	if ip := net.ParseIP(host); ip != nil {
		return v2net.TCPDestination(v2net.IPAddress(ip), port), nil
	}
	return v2net.TCPDestination(v2net.DomainAddress(host), port), nil
}

func isTimeout(err error) bool {
	nerr, ok := err.(net.Error)
	return ok && nerr.Timeout()
}

func (s *Server) Process(ctx context.Context, network v2net.Network, conn internet.Connection, dispatcher dispatcher.Interface) error {
	reader := bufio.NewReaderSize(conn, 2048)

Start:
	conn.SetReadDeadline(time.Now().Add(time.Second * 16))

	request, err := http.ReadRequest(reader)
	if err != nil {
		trace := newError("failed to read http request").Base(err)
		if errors.Cause(err) != io.EOF && !isTimeout(errors.Cause(err)) {
			trace.AtWarning()
		}
		return trace
	}
	log.Trace(newError("request to Method [", request.Method, "] Host [", request.Host, "] with URL [", request.URL, "]"))
	conn.SetReadDeadline(time.Time{})

	defaultPort := v2net.Port(80)
	if strings.ToLower(request.URL.Scheme) == "https" {
		defaultPort = v2net.Port(443)
	}
	host := request.Host
	if len(host) == 0 {
		host = request.URL.Host
	}
	dest, err := parseHost(host, defaultPort)
	if err != nil {
		return newError("malformed proxy host: ", host).AtWarning().Base(err)
	}
	log.Access(conn.RemoteAddr(), request.URL, log.AccessAccepted, "")

	if strings.ToUpper(request.Method) == "CONNECT" {
		return s.handleConnect(ctx, request, reader, conn, dest, dispatcher)
	}

	keepAlive := (strings.TrimSpace(strings.ToLower(request.Header.Get("Proxy-Connection"))) == "keep-alive")

	err = s.handlePlainHTTP(ctx, request, reader, conn, dest, dispatcher)
	if err == errWaitAnother {
		if keepAlive {
			goto Start
		}
		err = nil
	}

	return err
}

func (s *Server) handleConnect(ctx context.Context, request *http.Request, reader io.Reader, writer io.Writer, dest v2net.Destination, dispatcher dispatcher.Interface) error {
	_, err := writer.Write([]byte("HTTP/1.1 200 Connection established\r\n\r\n"))
	if err != nil {
		return newError("failed to write back OK response").Base(err)
	}

	timeout := time.Second * time.Duration(s.config.Timeout)
	if timeout == 0 {
		timeout = time.Minute * 2
	}
	ctx, timer := signal.CancelAfterInactivity(ctx, timeout)
	ray, err := dispatcher.Dispatch(ctx, dest)
	if err != nil {
		return err
	}

	requestDone := signal.ExecuteAsync(func() error {
		defer ray.InboundInput().Close()

		v2reader := buf.NewReader(reader)
		if err := buf.Copy(v2reader, ray.InboundInput(), buf.UpdateActivity(timer)); err != nil {
			return err
		}
		return nil
	})

	responseDone := signal.ExecuteAsync(func() error {
		v2writer := buf.NewWriter(writer)
		if err := buf.Copy(ray.InboundOutput(), v2writer, buf.UpdateActivity(timer)); err != nil {
			return err
		}
		return nil
	})

	if err := signal.ErrorOrFinish2(ctx, requestDone, responseDone); err != nil {
		ray.InboundInput().CloseError()
		ray.InboundOutput().CloseError()
		return newError("connection ends").Base(err)
	}

	runtime.KeepAlive(timer)

	return nil
}

// @VisibleForTesting
func StripHopByHopHeaders(header http.Header) {
	// Strip hop-by-hop header basaed on RFC:
	// http://www.w3.org/Protocols/rfc2616/rfc2616-sec13.html#sec13.5.1
	// https://www.mnot.net/blog/2011/07/11/what_proxies_must_do

	header.Del("Proxy-Connection")
	header.Del("Proxy-Authenticate")
	header.Del("Proxy-Authorization")
	header.Del("TE")
	header.Del("Trailers")
	header.Del("Transfer-Encoding")
	header.Del("Upgrade")

	connections := header.Get("Connection")
	header.Del("Connection")
	if len(connections) == 0 {
		return
	}
	for _, h := range strings.Split(connections, ",") {
		header.Del(strings.TrimSpace(h))
	}
}

var errWaitAnother = newError("keep alive")

func (s *Server) handlePlainHTTP(ctx context.Context, request *http.Request, reader io.Reader, writer io.Writer, dest v2net.Destination, dispatcher dispatcher.Interface) error {
	if len(request.URL.Host) <= 0 {
		response := &http.Response{
			Status:        "Bad Request",
			StatusCode:    400,
			Proto:         "HTTP/1.1",
			ProtoMajor:    1,
			ProtoMinor:    1,
			Header:        http.Header(make(map[string][]string)),
			Body:          nil,
			ContentLength: 0,
			Close:         true,
		}
		response.Header.Set("Proxy-Connection", "close")
		response.Header.Set("Connection", "close")
		return response.Write(writer)
	}

	request.Host = request.URL.Host
	StripHopByHopHeaders(request.Header)

	ray, err := dispatcher.Dispatch(ctx, dest)
	if err != nil {
		return err
	}
	input := ray.InboundInput()
	output := ray.InboundOutput()
	defer input.Close()

	var result error = errWaitAnother

	requestDone := signal.ExecuteAsync(func() error {
		request.Header.Set("Connection", "close")

		requestWriter := buf.ToBytesWriter(ray.InboundInput())
		if err := request.Write(requestWriter); err != nil {
			return err
		}
		return nil
	})

	responseDone := signal.ExecuteAsync(func() error {
		responseReader := bufio.NewReaderSize(buf.ToBytesReader(ray.InboundOutput()), 2048)
		response, err := http.ReadResponse(responseReader, request)
		if err == nil {
			StripHopByHopHeaders(response.Header)
			if response.ContentLength >= 0 {
				response.Header.Set("Proxy-Connection", "keep-alive")
				response.Header.Set("Connection", "keep-alive")
				response.Header.Set("Keep-Alive", "timeout=4")
				response.Close = false
			} else {
				response.Close = true
				result = nil
			}
		} else {
			log.Trace(newError("failed to read response from ", request.Host).Base(err).AtWarning())
			response = &http.Response{
				Status:        "Service Unavailable",
				StatusCode:    503,
				Proto:         "HTTP/1.1",
				ProtoMajor:    1,
				ProtoMinor:    1,
				Header:        http.Header(make(map[string][]string)),
				Body:          nil,
				ContentLength: 0,
				Close:         true,
			}
			response.Header.Set("Connection", "close")
			response.Header.Set("Proxy-Connection", "close")
		}
		if err := response.Write(writer); err != nil {
			return newError("failed to write response").Base(err)
		}
		return nil
	})

	if err := signal.ErrorOrFinish2(ctx, requestDone, responseDone); err != nil {
		input.CloseError()
		output.CloseError()
		return newError("connection ends").Base(err)
	}

	return result
}

func init() {
	common.Must(common.RegisterConfig((*ServerConfig)(nil), func(ctx context.Context, config interface{}) (interface{}, error) {
		return NewServer(ctx, config.(*ServerConfig))
	}))
}
