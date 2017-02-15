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
		return nil, errors.New("HTTP|Server: No space in context.")
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

func (s *Server) Process(ctx context.Context, network v2net.Network, conn internet.Connection, dispatcher dispatcher.Interface) error {
	conn.SetReusable(false)

	conn.SetReadDeadline(time.Now().Add(time.Second * 8))
	reader := bufio.NewReaderSize(conn, 2048)

	request, err := http.ReadRequest(reader)
	if err != nil {
		if errors.Cause(err) != io.EOF {
			log.Warning("HTTP: Failed to read http request: ", err)
		}
		return err
	}
	log.Info("HTTP: Request to Method [", request.Method, "] Host [", request.Host, "] with URL [", request.URL, "]")
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
		log.Warning("HTTP: Malformed proxy host (", host, "): ", err)
		return err
	}
	log.Access(conn.RemoteAddr(), request.URL, log.AccessAccepted, "")

	if strings.ToUpper(request.Method) == "CONNECT" {
		return s.handleConnect(ctx, request, reader, conn, dest, dispatcher)
	} else {
		return s.handlePlainHTTP(ctx, request, reader, conn, dest, dispatcher)
	}
}

func (s *Server) handleConnect(ctx context.Context, request *http.Request, reader io.Reader, writer io.Writer, dest v2net.Destination, dispatcher dispatcher.Interface) error {
	response := &http.Response{
		Status:        "200 OK",
		StatusCode:    200,
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		Header:        http.Header(make(map[string][]string)),
		Body:          nil,
		ContentLength: 0,
		Close:         false,
	}
	if err := response.Write(writer); err != nil {
		log.Warning("HTTP|Server: failed to write back OK response: ", err)
		return err
	}

	ctx, cancel := context.WithCancel(ctx)
	timeout := time.Second * time.Duration(s.config.Timeout)
	if timeout == 0 {
		timeout = time.Minute * 2
	}
	timer := signal.CancelAfterInactivity(ctx, cancel, timeout)
	ray, err := dispatcher.Dispatch(ctx, dest)
	if err != nil {
		return err
	}

	requestDone := signal.ExecuteAsync(func() error {
		defer ray.InboundInput().Close()

		v2reader := buf.NewReader(reader)
		if err := buf.PipeUntilEOF(timer, v2reader, ray.InboundInput()); err != nil {
			return err
		}
		return nil
	})

	responseDone := signal.ExecuteAsync(func() error {
		v2writer := buf.NewWriter(writer)
		if err := buf.PipeUntilEOF(timer, ray.InboundOutput(), v2writer); err != nil {
			return err
		}
		return nil
	})

	if err := signal.ErrorOrFinish2(ctx, requestDone, responseDone); err != nil {
		log.Info("HTTP|Server: Connection ends with: ", err)
		ray.InboundInput().CloseError()
		ray.InboundOutput().CloseError()
		return err
	}

	runtime.KeepAlive(timer)

	return nil
}

// @VisibleForTesting
func StripHopByHopHeaders(request *http.Request) {
	// Strip hop-by-hop header basaed on RFC:
	// http://www.w3.org/Protocols/rfc2616/rfc2616-sec13.html#sec13.5.1
	// https://www.mnot.net/blog/2011/07/11/what_proxies_must_do

	request.Header.Del("Proxy-Connection")
	request.Header.Del("Proxy-Authenticate")
	request.Header.Del("Proxy-Authorization")
	request.Header.Del("TE")
	request.Header.Del("Trailers")
	request.Header.Del("Transfer-Encoding")
	request.Header.Del("Upgrade")

	// TODO: support keep-alive
	connections := request.Header.Get("Connection")
	request.Header.Set("Connection", "close")
	if len(connections) == 0 {
		return
	}
	for _, h := range strings.Split(connections, ",") {
		request.Header.Del(strings.TrimSpace(h))
	}
}

func generateResponse(statusCode int, status string) *http.Response {
	hdr := http.Header(make(map[string][]string))
	hdr.Set("Connection", "close")
	return &http.Response{
		Status:        status,
		StatusCode:    statusCode,
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		Header:        hdr,
		Body:          nil,
		ContentLength: 0,
		Close:         true,
	}
}

func (s *Server) handlePlainHTTP(ctx context.Context, request *http.Request, reader io.Reader, writer io.Writer, dest v2net.Destination, dispatcher dispatcher.Interface) error {
	if len(request.URL.Host) <= 0 {
		response := generateResponse(400, "Bad Request")
		return response.Write(writer)
	}

	request.Host = request.URL.Host
	StripHopByHopHeaders(request)

	ray, err := dispatcher.Dispatch(ctx, dest)
	if err != nil {
		return err
	}
	input := ray.InboundInput()
	output := ray.InboundOutput()
	defer input.Close()

	requestDone := signal.ExecuteAsync(func() error {
		requestWriter := buf.NewBufferedWriter(buf.ToBytesWriter(ray.InboundInput()))
		err := request.Write(requestWriter)
		if err != nil {
			return err
		}
		if err := requestWriter.Flush(); err != nil {
			return err
		}
		return nil
	})

	responseDone := signal.ExecuteAsync(func() error {
		responseReader := bufio.NewReader(buf.ToBytesReader(ray.InboundOutput()))
		response, err := http.ReadResponse(responseReader, request)
		if err != nil {
			log.Warning("HTTP: Failed to read response: ", err)
			response = generateResponse(503, "Service Unavailable")
		}
		responseWriter := buf.NewBufferedWriter(writer)
		if err := response.Write(responseWriter); err != nil {
			return err
		}

		if err := responseWriter.Flush(); err != nil {
			return err
		}
		return nil
	})

	if err := signal.ErrorOrFinish2(ctx, requestDone, responseDone); err != nil {
		log.Info("HTTP|Server: Connecton ending with ", err)
		input.CloseError()
		output.CloseError()
		return err
	}

	return nil
}

func init() {
	common.Must(common.RegisterConfig((*ServerConfig)(nil), func(ctx context.Context, config interface{}) (interface{}, error) {
		return NewServer(ctx, config.(*ServerConfig))
	}))
}
