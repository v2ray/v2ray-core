// +build !confonly

package http

import (
	"bufio"
	"context"
	"encoding/base64"
	"io"
	"net/http"
	"net/url"
	"sync"

	"golang.org/x/net/http2"

	"v2ray.com/core"
	"v2ray.com/core/common"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/bytespool"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/protocol"
	"v2ray.com/core/common/retry"
	"v2ray.com/core/common/session"
	"v2ray.com/core/common/signal"
	"v2ray.com/core/common/task"
	"v2ray.com/core/features/policy"
	"v2ray.com/core/transport"
	"v2ray.com/core/transport/internet"
	"v2ray.com/core/transport/internet/tls"
)

type Client struct {
	serverPicker  protocol.ServerPicker
	policyManager policy.Manager
}

type h2Conn struct {
	rawConn net.Conn
	h2Conn  *http2.ClientConn
}

var (
	cachedH2Mutex sync.Mutex
	cachedH2Conns map[net.Destination]h2Conn
)

// NewClient create a new http client based on the given config.
func NewClient(ctx context.Context, config *ClientConfig) (*Client, error) {
	serverList := protocol.NewServerList()
	for _, rec := range config.Server {
		s, err := protocol.NewServerSpecFromPB(rec)
		if err != nil {
			return nil, newError("failed to get server spec").Base(err)
		}
		serverList.AddServer(s)
	}
	if serverList.Size() == 0 {
		return nil, newError("0 target server")
	}

	v := core.MustFromContext(ctx)
	return &Client{
		serverPicker:  protocol.NewRoundRobinServerPicker(serverList),
		policyManager: v.GetFeature(policy.ManagerType()).(policy.Manager),
	}, nil
}

// Process implements proxy.Outbound.Process. We first create a socket tunnel via HTTP CONNECT method, then redirect all inbound traffic to that tunnel.
func (c *Client) Process(ctx context.Context, link *transport.Link, dialer internet.Dialer) error {
	outbound := session.OutboundFromContext(ctx)
	if outbound == nil || !outbound.Target.IsValid() {
		return newError("target not specified.")
	}
	target := outbound.Target
	targetAddr := target.NetAddr()

	if target.Network == net.Network_UDP {
		return newError("UDP is not supported by HTTP outbound")
	}

	var user *protocol.MemoryUser
	var conn internet.Connection

	mbuf, _ := link.Reader.ReadMultiBuffer()
	len := mbuf.Len()
	firstPayload := bytespool.Alloc(len)
	mbuf, _ = buf.SplitBytes(mbuf, firstPayload)
	firstPayload = firstPayload[:len]

	buf.ReleaseMulti(mbuf)
	defer bytespool.Free(firstPayload)

	if err := retry.ExponentialBackoff(5, 100).On(func() error {
		server := c.serverPicker.PickServer()
		dest := server.Destination()
		user = server.PickUser()

		netConn, err := setUpHTTPTunnel(ctx, dest, targetAddr, user, dialer, firstPayload)
		if netConn != nil {
			conn = internet.Connection(netConn)
		}
		return err
	}); err != nil {
		return newError("failed to find an available destination").Base(err)
	}

	defer func() {
		if err := conn.Close(); err != nil {
			newError("failed to closed connection").Base(err).WriteToLog(session.ExportIDToError(ctx))
		}
	}()

	p := c.policyManager.ForLevel(0)
	if user != nil {
		p = c.policyManager.ForLevel(user.Level)
	}

	ctx, cancel := context.WithCancel(ctx)
	timer := signal.CancelAfterInactivity(ctx, cancel, p.Timeouts.ConnectionIdle)

	requestFunc := func() error {
		defer timer.SetTimeout(p.Timeouts.DownlinkOnly)
		return buf.Copy(link.Reader, buf.NewWriter(conn), buf.UpdateActivity(timer))
	}
	responseFunc := func() error {
		defer timer.SetTimeout(p.Timeouts.UplinkOnly)
		return buf.Copy(buf.NewReader(conn), link.Writer, buf.UpdateActivity(timer))
	}

	var responseDonePost = task.OnSuccess(responseFunc, task.Close(link.Writer))
	if err := task.Run(ctx, requestFunc, responseDonePost); err != nil {
		return newError("connection ends").Base(err)
	}

	return nil
}

// setUpHTTPTunnel will create a socket tunnel via HTTP CONNECT method
func setUpHTTPTunnel(ctx context.Context, dest net.Destination, target string, user *protocol.MemoryUser, dialer internet.Dialer, firstPayload []byte) (net.Conn, error) {
	req := &http.Request{
		Method: http.MethodConnect,
		URL:    &url.URL{Host: target},
		Header: make(http.Header),
		Host:   target,
	}

	if user != nil && user.Account != nil {
		account := user.Account.(*Account)
		auth := account.GetUsername() + ":" + account.GetPassword()
		req.Header.Set("Proxy-Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(auth)))
	}

	connectHTTP1 := func(rawConn net.Conn) (net.Conn, error) {
		req.Header.Set("Proxy-Connection", "Keep-Alive")

		err := req.Write(rawConn)
		if err != nil {
			rawConn.Close()
			return nil, err
		}

		if _, err := rawConn.Write(firstPayload); err != nil {
			rawConn.Close()
			return nil, err
		}

		resp, err := http.ReadResponse(bufio.NewReader(rawConn), req)
		if err != nil {
			rawConn.Close()
			return nil, err
		}

		if resp.StatusCode != http.StatusOK {
			rawConn.Close()
			return nil, newError("Proxy responded with non 200 code: " + resp.Status)
		}
		return rawConn, nil
	}

	connectHTTP2 := func(rawConn net.Conn, h2clientConn *http2.ClientConn) (net.Conn, error) {
		pr, pw := io.Pipe()
		req.Body = pr

		var pErr error
		var wg sync.WaitGroup
		wg.Add(1)

		go func() {
			_, pErr = pw.Write(firstPayload)
			wg.Done()
		}()

		resp, err := h2clientConn.RoundTrip(req)
		if err != nil {
			rawConn.Close()
			return nil, err
		}

		wg.Wait()
		if pErr != nil {
			rawConn.Close()
			return nil, pErr
		}

		if resp.StatusCode != http.StatusOK {
			rawConn.Close()
			return nil, newError("Proxy responded with non 200 code: " + resp.Status)
		}
		return newHTTP2Conn(rawConn, pw, resp.Body), nil
	}

	cachedH2Mutex.Lock()
	cachedConn, cachedConnFound := cachedH2Conns[dest]
	cachedH2Mutex.Unlock()

	if cachedConnFound {
		rc, cc := cachedConn.rawConn, cachedConn.h2Conn
		if cc.CanTakeNewRequest() {
			proxyConn, err := connectHTTP2(rc, cc)
			if err != nil {
				return nil, err
			}

			return proxyConn, nil
		}
	}

	rawConn, err := dialer.Dial(ctx, dest)
	if err != nil {
		return nil, err
	}

	iConn := rawConn
	if statConn, ok := iConn.(*internet.StatCouterConnection); ok {
		iConn = statConn.Connection
	}

	nextProto := ""
	if tlsConn, ok := iConn.(*tls.Conn); ok {
		if err := tlsConn.Handshake(); err != nil {
			rawConn.Close()
			return nil, err
		}
		nextProto = tlsConn.ConnectionState().NegotiatedProtocol
	}

	switch nextProto {
	case "", "http/1.1":
		return connectHTTP1(rawConn)
	case "h2":
		t := http2.Transport{}
		h2clientConn, err := t.NewClientConn(rawConn)
		if err != nil {
			rawConn.Close()
			return nil, err
		}

		proxyConn, err := connectHTTP2(rawConn, h2clientConn)
		if err != nil {
			rawConn.Close()
			return nil, err
		}

		cachedH2Mutex.Lock()
		if cachedH2Conns == nil {
			cachedH2Conns = make(map[net.Destination]h2Conn)
		}

		cachedH2Conns[dest] = h2Conn{
			rawConn: rawConn,
			h2Conn:  h2clientConn,
		}
		cachedH2Mutex.Unlock()

		return proxyConn, err
	default:
		return nil, newError("negotiated unsupported application layer protocol: " + nextProto)
	}
}

func newHTTP2Conn(c net.Conn, pipedReqBody *io.PipeWriter, respBody io.ReadCloser) net.Conn {
	return &http2Conn{Conn: c, in: pipedReqBody, out: respBody}
}

type http2Conn struct {
	net.Conn
	in  *io.PipeWriter
	out io.ReadCloser
}

func (h *http2Conn) Read(p []byte) (n int, err error) {
	return h.out.Read(p)
}

func (h *http2Conn) Write(p []byte) (n int, err error) {
	return h.in.Write(p)
}

func (h *http2Conn) Close() error {
	h.in.Close()
	return h.out.Close()
}

func init() {
	common.Must(common.RegisterConfig((*ClientConfig)(nil), func(ctx context.Context, config interface{}) (interface{}, error) {
		return NewClient(ctx, config.(*ClientConfig))
	}))
}
