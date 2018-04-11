package http

import (
	"bufio"
	"context"
	"encoding/base64"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"v2ray.com/core"
	"v2ray.com/core/common"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/errors"
	"v2ray.com/core/common/log"
	"v2ray.com/core/common/net"
	http_proto "v2ray.com/core/common/protocol/http"
	"v2ray.com/core/common/signal"
	"v2ray.com/core/transport/internet"
)

// Server is an HTTP proxy server.
type Server struct {
	config *ServerConfig
	v      *core.Instance
}

// NewServer creates a new HTTP inbound handler.
func NewServer(ctx context.Context, config *ServerConfig) (*Server, error) {
	s := &Server{
		config: config,
		v:      core.MustFromContext(ctx),
	}

	return s, nil
}

func (s *Server) policy() core.Policy {
	config := s.config
	p := s.v.PolicyManager().ForLevel(config.UserLevel)
	if config.Timeout > 0 && config.UserLevel == 0 {
		p.Timeouts.ConnectionIdle = time.Duration(config.Timeout) * time.Second
	}
	return p
}

func (*Server) Network() net.NetworkList {
	return net.NetworkList{
		Network: []net.Network{net.Network_TCP},
	}
}

func parseHost(rawHost string, defaultPort net.Port) (net.Destination, error) {
	port := defaultPort
	host, rawPort, err := net.SplitHostPort(rawHost)
	if err != nil {
		if addrError, ok := err.(*net.AddrError); ok && strings.Contains(addrError.Err, "missing port") {
			host = rawHost
		} else {
			return net.Destination{}, err
		}
	} else if len(rawPort) > 0 {
		intPort, err := strconv.Atoi(rawPort)
		if err != nil {
			return net.Destination{}, err
		}
		port = net.Port(intPort)
	}

	return net.TCPDestination(net.ParseAddress(host), port), nil
}

func isTimeout(err error) bool {
	nerr, ok := errors.Cause(err).(net.Error)
	return ok && nerr.Timeout()
}

func parseBasicAuth(auth string) (username, password string, ok bool) {
	const prefix = "Basic "
	if !strings.HasPrefix(auth, prefix) {
		return
	}
	c, err := base64.StdEncoding.DecodeString(auth[len(prefix):])
	if err != nil {
		return
	}
	cs := string(c)
	s := strings.IndexByte(cs, ':')
	if s < 0 {
		return
	}
	return cs[:s], cs[s+1:], true
}

type readerOnly struct {
	io.Reader
}

func (s *Server) Process(ctx context.Context, network net.Network, conn internet.Connection, dispatcher core.Dispatcher) error {
	reader := bufio.NewReaderSize(readerOnly{conn}, buf.Size)

Start:
	conn.SetReadDeadline(time.Now().Add(s.policy().Timeouts.Handshake))

	request, err := http.ReadRequest(reader)
	if err != nil {
		trace := newError("failed to read http request").Base(err)
		if errors.Cause(err) != io.EOF && !isTimeout(errors.Cause(err)) {
			trace.AtWarning()
		}
		return trace
	}

	if len(s.config.Accounts) > 0 {
		user, pass, ok := parseBasicAuth(request.Header.Get("Proxy-Authorization"))
		if !ok || !s.config.HasAccount(user, pass) {
			return common.Error2(conn.Write([]byte("HTTP/1.1 407 Proxy Authentication Required\r\nProxy-Authenticate: Basic realm=\"proxy\"\r\n\r\n")))
		}
	}

	newError("request to Method [", request.Method, "] Host [", request.Host, "] with URL [", request.URL, "]").WithContext(ctx).WriteToLog()
	conn.SetReadDeadline(time.Time{})

	defaultPort := net.Port(80)
	if strings.ToLower(request.URL.Scheme) == "https" {
		defaultPort = net.Port(443)
	}
	host := request.Host
	if len(host) == 0 {
		host = request.URL.Host
	}
	dest, err := parseHost(host, defaultPort)
	if err != nil {
		return newError("malformed proxy host: ", host).AtWarning().Base(err)
	}
	log.Record(&log.AccessMessage{
		From:   conn.RemoteAddr(),
		To:     request.URL,
		Status: log.AccessAccepted,
	})

	if strings.ToUpper(request.Method) == "CONNECT" {
		return s.handleConnect(ctx, request, reader, conn, dest, dispatcher)
	}

	keepAlive := (strings.TrimSpace(strings.ToLower(request.Header.Get("Proxy-Connection"))) == "keep-alive")

	err = s.handlePlainHTTP(ctx, request, conn, dest, dispatcher)
	if err == errWaitAnother {
		if keepAlive {
			goto Start
		}
		err = nil
	}

	return err
}

func (s *Server) handleConnect(ctx context.Context, request *http.Request, reader *bufio.Reader, conn internet.Connection, dest net.Destination, dispatcher core.Dispatcher) error {
	_, err := conn.Write([]byte("HTTP/1.1 200 Connection established\r\n\r\n"))
	if err != nil {
		return newError("failed to write back OK response").Base(err)
	}

	ctx, cancel := context.WithCancel(ctx)
	timer := signal.CancelAfterInactivity(ctx, cancel, s.policy().Timeouts.ConnectionIdle)
	ray, err := dispatcher.Dispatch(ctx, dest)
	if err != nil {
		return err
	}

	if reader.Buffered() > 0 {
		payload, err := buf.ReadSizeToMultiBuffer(reader, int32(reader.Buffered()))
		if err != nil {
			return err
		}
		if err := ray.InboundInput().WriteMultiBuffer(payload); err != nil {
			return err
		}
		reader = nil
	}

	requestDone := signal.ExecuteAsync(func() error {
		defer ray.InboundInput().Close()
		defer timer.SetTimeout(s.policy().Timeouts.DownlinkOnly)

		v2reader := buf.NewReader(conn)
		return buf.Copy(v2reader, ray.InboundInput(), buf.UpdateActivity(timer))
	})

	responseDone := signal.ExecuteAsync(func() error {
		defer timer.SetTimeout(s.policy().Timeouts.UplinkOnly)

		v2writer := buf.NewWriter(conn)
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

	return nil
}

var errWaitAnother = newError("keep alive")

func (s *Server) handlePlainHTTP(ctx context.Context, request *http.Request, writer io.Writer, dest net.Destination, dispatcher core.Dispatcher) error {
	if !s.config.AllowTransparent && len(request.URL.Host) <= 0 {
		// RFC 2068 (HTTP/1.1) requires URL to be absolute URL in HTTP proxy.
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

	if len(request.URL.Host) > 0 {
		request.Host = request.URL.Host
	}
	http_proto.RemoveHopByHopHeaders(request.Header)

	// Prevent UA from being set to golang's default ones
	if len(request.Header.Get("User-Agent")) == 0 {
		request.Header.Set("User-Agent", "")
	}

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

		requestWriter := buf.NewBufferedWriter(ray.InboundInput())
		common.Must(requestWriter.SetBuffered(false))
		if err := request.Write(requestWriter); err != nil {
			return newError("failed to write whole request").Base(err).AtWarning()
		}
		return nil
	})

	responseDone := signal.ExecuteAsync(func() error {
		responseReader := bufio.NewReaderSize(buf.NewBufferedReader(ray.InboundOutput()), buf.Size)
		response, err := http.ReadResponse(responseReader, request)
		if err == nil {
			http_proto.RemoveHopByHopHeaders(response.Header)
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
			newError("failed to read response from ", request.Host).Base(err).AtWarning().WithContext(ctx).WriteToLog()
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
			return newError("failed to write response").Base(err).AtWarning()
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
