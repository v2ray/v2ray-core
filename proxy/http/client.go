// +build !confonly

package http

import (
	"bufio"
	"context"
	"encoding/base64"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"v2ray.com/core"
	"v2ray.com/core/common"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/protocol"
	"v2ray.com/core/common/retry"
	"v2ray.com/core/common/session"
	"v2ray.com/core/common/signal"
	"v2ray.com/core/common/task"
	"v2ray.com/core/features/policy"
	"v2ray.com/core/transport"
	"v2ray.com/core/transport/internet"
)

// Client is a inbound handler for HTTP protocol
type Client struct {
	serverPicker  protocol.ServerPicker
	policyManager policy.Manager
}

// NewClient create a new http client based on the given config.
func NewClient(ctx context.Context, config *ClientConfig) (*Client, error) {
	serverList := protocol.NewServerList()
	for _, rec := range config.Server {
		s, err := protocol.NewServerSpecFromPB(*rec)
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
	destination := outbound.Target

	if destination.Network == net.Network_UDP {
		return newError("UDP is not supported by HTTP outbound")
	}

	var server *protocol.ServerSpec
	var conn internet.Connection

	if err := retry.ExponentialBackoff(5, 100).On(func() error {
		server = c.serverPicker.PickServer()
		dest := server.Destination()
		rawConn, err := dialer.Dial(ctx, dest)
		if err != nil {
			return err
		}
		conn = rawConn

		return nil
	}); err != nil {
		return newError("failed to find an available destination").Base(err)
	}

	defer func() {
		if err := conn.Close(); err != nil {
			newError("failed to closed connection").Base(err).WriteToLog(session.ExportIDToError(ctx))
		}
	}()

	p := c.policyManager.ForLevel(0)

	user := server.PickUser()
	if user != nil {
		p = c.policyManager.ForLevel(user.Level)
	}

	conn = setUpHTTPTunnel(conn, &destination, user)

	ctx, cancel := context.WithCancel(ctx)
	timer := signal.CancelAfterInactivity(ctx, cancel, p.Timeouts.ConnectionIdle)

	requestFunc := func() error {
		defer timer.SetTimeout(p.Timeouts.DownlinkOnly)
		return buf.Copy(link.Reader, buf.NewWriter(conn), buf.UpdateActivity(timer))
	}
	responseFunc := func() error {
		defer timer.SetTimeout(p.Timeouts.UplinkOnly)
		bc := bufio.NewReader(conn)
		resp, err := http.ReadResponse(bc, nil)
		if err != nil {
			return err
		}
		if resp.StatusCode != http.StatusOK {
			return newError(resp.Status)
		}
		return buf.Copy(buf.NewReader(bc), link.Writer, buf.UpdateActivity(timer))
	}

	var responseDonePost = task.OnSuccess(responseFunc, task.Close(link.Writer))
	if err := task.Run(ctx, requestFunc, responseDonePost); err != nil {
		return newError("connection ends").Base(err)
	}

	return nil
}

// setUpHTTPTunnel will create a socket tunnel via HTTP CONNECT method
func setUpHTTPTunnel(conn internet.Connection, destination *net.Destination, user *protocol.MemoryUser) *tunConn {
	var headers []string
	destNetAddr := destination.NetAddr()
	headers = append(headers, "CONNECT "+destNetAddr+" HTTP/1.1")
	headers = append(headers, "Host: "+destNetAddr)
	if user != nil && user.Account != nil {
		account := user.Account.(*Account)
		auth := account.GetUsername() + ":" + account.GetPassword()
		headers = append(headers, "Proxy-Authorization: Basic "+base64.StdEncoding.EncodeToString([]byte(auth)))
	}
	headers = append(headers, "Proxy-Connection: Keep-Alive")

	b := buf.New()
	b.WriteString(strings.Join(headers, "\r\n") + "\r\n\r\n")
	return newTunConn(conn, b, 5 * time.Millisecond)
}

// tunConn is a connection that writes header before content,
// the header will be written during the next Write call or after
// specified delay.
type tunConn struct {
	internet.Connection
	header *buf.Buffer
	once   sync.Once
	timer  *time.Timer
}

func newTunConn(conn internet.Connection, header *buf.Buffer, delay time.Duration) *tunConn {
	tc := &tunConn{
		Connection: conn,
		header:     header,
	}
	if delay > 0 {
		tc.timer = time.AfterFunc(delay, func() {
			tc.Write([]byte{})
		})
	}
	return tc
}

func (c *tunConn) Write(b []byte) (n int, err error) {
	// fallback to normal write if header is sent
	if c.header == nil {
		return c.Connection.Write(b)
	}
	// Prevent timer and writer race condition
	c.once.Do(func() {
		if c.timer != nil {
			c.timer.Stop()
			c.timer = nil
		}
		lenheader := c.header.Len()
		// Concate header and b
		common.Must2(c.header.Write(b))
		// Write buffer
		var nc int64
		nc, err = io.Copy(c.Connection, c.header)
		c.header.Release()
		c.header = nil
		n = int(nc) - int(lenheader)
		if n < 0 { n = 0 }
		b = b[n:]
	})
	// Write Trailing bytes
	if len(b) > 0 && err == nil {
		var nw int
		nw, err = c.Connection.Write(b)
		n += nw
	}
	return n, err
}

func init() {
	common.Must(common.RegisterConfig((*ClientConfig)(nil), func(ctx context.Context, config interface{}) (interface{}, error) {
		return NewClient(ctx, config.(*ClientConfig))
	}))
}
