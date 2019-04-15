// +build !confonly

package tls

import (
	"crypto/tls"
	"os"

	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/net"

	utls "v2ray.com/core/external/github.com/refraction-networking/utls"
)

//go:generate errorgen

var (
	_ buf.Writer = (*conn)(nil)
)

type conn struct {
	*tls.Conn
}

func (c *conn) WriteMultiBuffer(mb buf.MultiBuffer) error {
	mb = buf.Compact(mb)
	mb, err := buf.WriteMultiBuffer(c, mb)
	buf.ReleaseMulti(mb)
	return err
}

func (c *conn) HandshakeAddress() net.Address {
	if err := c.Handshake(); err != nil {
		return nil
	}
	state := c.Conn.ConnectionState()
	if len(state.ServerName) == 0 {
		return nil
	}
	return net.ParseAddress(state.ServerName)
}

// Client initiates a TLS client handshake on the given connection.
func Client(c net.Conn, config *tls.Config) net.Conn {
	tlsConn := tls.Client(c, config)
	return &conn{Conn: tlsConn}
}

func copyConfig(c *tls.Config) *utls.Config {
	return &utls.Config{
		NextProtos:         c.NextProtos,
		ServerName:         c.ServerName,
		InsecureSkipVerify: c.InsecureSkipVerify,
		MinVersion:         utls.VersionTLS12,
		MaxVersion:         utls.VersionTLS12,
	}
}

func UClient(c net.Conn, config *tls.Config) net.Conn {
	uConfig := copyConfig(config)
	return utls.Client(c, uConfig)
}

// Server initiates a TLS server handshake on the given connection.
func Server(c net.Conn, config *tls.Config) net.Conn {
	tlsConn := tls.Server(c, config)
	return &conn{Conn: tlsConn}
}

func init() {
	// opt-in TLS 1.3 for Go1.12
	// TODO: remove this line when Go1.13 is released.
	_ = os.Setenv("GODEBUG", os.Getenv("GODEBUG")+",tls13=1")
}
