// +build !confonly

package http

import (
	"bufio"
	"context"
	"encoding/base64"
	"io"
	"net/http"
	"strings"
	"time"

	"v2ray.com/core"
	"v2ray.com/core/common"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/errors"
	"v2ray.com/core/common/log"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/protocol"
	http_proto "v2ray.com/core/common/protocol/http"
	"v2ray.com/core/common/session"
	"v2ray.com/core/common/signal"
	"v2ray.com/core/common/task"
	"v2ray.com/core/features/policy"
	"v2ray.com/core/features/routing"
	"v2ray.com/core/transport/internet"
)

// Server is an HTTP proxy server.
type Server struct {
	config        *ServerConfig
	policyManager policy.Manager
}

// NewServer creates a new HTTP inbound handler.
func NewServer(ctx context.Context, config *ServerConfig) (*Server, error) {
	v := core.MustFromContext(ctx)
	s := &Server{
		config:        config,
		policyManager: v.GetFeature(policy.ManagerType()).(policy.Manager),
	}

	return s, nil
}

func (s *Server) policy() policy.Session {
	config := s.config
	p := s.policyManager.ForLevel(config.UserLevel)
	if config.Timeout > 0 && config.UserLevel == 0 {
		p.Timeouts.ConnectionIdle = time.Duration(config.Timeout) * time.Second
	}
	return p
}

// Network implements proxy.Inbound.
func (*Server) Network() []net.Network {
	return []net.Network{net.Network_TCP}
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

func (s *Server) Process(ctx context.Context, network net.Network, conn internet.Connection, dispatcher routing.Dispatcher) error {
	if inbound := session.InboundFromContext(ctx); inbound != nil {
		inbound.User = &protocol.MemoryUser{
			Level: s.config.UserLevel,
		}
	}

	reader := bufio.NewReaderSize(readerOnly{conn}, buf.Size)

Start:
	if err := conn.SetReadDeadline(time.Now().Add(s.policy().Timeouts.Handshake)); err != nil {
		newError("failed to set read deadline").Base(err).WriteToLog(session.ExportIDToError(ctx))
	}

	request, err := http.ReadRequest(reader)
	if err != nil {
		trace := newError("failed to read http request").Base(err)
		if errors.Cause(err) != io.EOF && !isTimeout(errors.Cause(err)) {
			trace.AtWarning() // nolint: errcheck
		}
		return trace
	}

	if len(s.config.Accounts) > 0 {
		user, pass, ok := parseBasicAuth(request.Header.Get("Proxy-Authorization"))
		if !ok || !s.config.HasAccount(user, pass) {
			return common.Error2(conn.Write([]byte("HTTP/1.1 407 Proxy Authentication Required\r\nProxy-Authenticate: Basic realm=\"proxy\"\r\n\r\n")))
		}
	}

	newError("request to Method [", request.Method, "] Host [", request.Host, "] with URL [", request.URL, "]").WriteToLog(session.ExportIDToError(ctx))
	if err := conn.SetReadDeadline(time.Time{}); err != nil {
		newError("failed to clear read deadline").Base(err).WriteToLog(session.ExportIDToError(ctx))
	}

	defaultPort := net.Port(80)
	if strings.EqualFold(request.URL.Scheme, "https") {
		defaultPort = net.Port(443)
	}
	host := request.Host
	if host == "" {
		host = request.URL.Host
	}
	dest, err := http_proto.ParseHost(host, defaultPort)
	if err != nil {
		return newError("malformed proxy host: ", host).AtWarning().Base(err)
	}
	log.Record(&log.AccessMessage{
		From:   conn.RemoteAddr(),
		To:     request.URL,
		Status: log.AccessAccepted,
	})

	if strings.EqualFold(request.Method, "CONNECT") {
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

func (s *Server) handleConnect(ctx context.Context, request *http.Request, reader *bufio.Reader, conn internet.Connection, dest net.Destination, dispatcher routing.Dispatcher) error {
	_, err := conn.Write([]byte("HTTP/1.1 200 Connection established\r\n\r\n"))
	if err != nil {
		return newError("failed to write back OK response").Base(err)
	}

	plcy := s.policy()
	ctx, cancel := context.WithCancel(ctx)
	timer := signal.CancelAfterInactivity(ctx, cancel, plcy.Timeouts.ConnectionIdle)

	ctx = policy.ContextWithBufferPolicy(ctx, plcy.Buffer)
	link, err := dispatcher.Dispatch(ctx, dest)
	if err != nil {
		return err
	}

	if reader.Buffered() > 0 {
		payload, err := buf.ReadFrom(io.LimitReader(reader, int64(reader.Buffered())))
		if err != nil {
			return err
		}
		if err := link.Writer.WriteMultiBuffer(payload); err != nil {
			return err
		}
		reader = nil
	}

	requestDone := func() error {
		defer timer.SetTimeout(plcy.Timeouts.DownlinkOnly)

		return buf.Copy(buf.NewReader(conn), link.Writer, buf.UpdateActivity(timer))
	}

	responseDone := func() error {
		defer timer.SetTimeout(plcy.Timeouts.UplinkOnly)

		v2writer := buf.NewWriter(conn)
		if err := buf.Copy(link.Reader, v2writer, buf.UpdateActivity(timer)); err != nil {
			return err
		}

		return nil
	}

	var closeWriter = task.OnSuccess(requestDone, task.Close(link.Writer))
	if err := task.Run(ctx, closeWriter, responseDone); err != nil {
		common.Interrupt(link.Reader)
		common.Interrupt(link.Writer)
		return newError("connection ends").Base(err)
	}

	return nil
}

var errWaitAnother = newError("keep alive")

func (s *Server) handlePlainHTTP(ctx context.Context, request *http.Request, writer io.Writer, dest net.Destination, dispatcher routing.Dispatcher) error {
	if !s.config.AllowTransparent && request.URL.Host == "" {
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
	if request.Header.Get("User-Agent") == "" {
		request.Header.Set("User-Agent", "")
	}

	content := &session.Content{
		Protocol: "http/1.1",
	}

	content.SetAttribute(":method", strings.ToUpper(request.Method))
	content.SetAttribute(":path", request.URL.Path)
	for key := range request.Header {
		value := request.Header.Get(key)
		content.SetAttribute(strings.ToLower(key), value)
	}

	ctx = session.ContextWithContent(ctx, content)

	link, err := dispatcher.Dispatch(ctx, dest)
	if err != nil {
		return err
	}

	// Plain HTTP request is not a stream. The request always finishes before response. Hense request has to be closed later.
	defer common.Close(link.Writer) // nolint: errcheck
	var result error = errWaitAnother

	requestDone := func() error {
		request.Header.Set("Connection", "close")

		requestWriter := buf.NewBufferedWriter(link.Writer)
		common.Must(requestWriter.SetBuffered(false))
		if err := request.Write(requestWriter); err != nil {
			return newError("failed to write whole request").Base(err).AtWarning()
		}
		return nil
	}

	responseDone := func() error {
		responseReader := bufio.NewReaderSize(&buf.BufferedReader{Reader: link.Reader}, buf.Size)
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
			newError("failed to read response from ", request.Host).Base(err).AtWarning().WriteToLog(session.ExportIDToError(ctx))
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
	}

	if err := task.Run(ctx, requestDone, responseDone); err != nil {
		common.Interrupt(link.Reader)
		common.Interrupt(link.Writer)
		return newError("connection ends").Base(err)
	}

	return result
}

func init() {
	common.Must(common.RegisterConfig((*ServerConfig)(nil), func(ctx context.Context, config interface{}) (interface{}, error) {
		return NewServer(ctx, config.(*ServerConfig))
	}))
}
