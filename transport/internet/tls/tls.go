// +build !confonly

package tls

import (
	"crypto/tls"

	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/net"
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

// Server initiates a TLS server handshake on the given connection.
func Server(c net.Conn, config *tls.Config) net.Conn {
	tlsConn := tls.Server(c, config)
	return &conn{Conn: tlsConn}
}
