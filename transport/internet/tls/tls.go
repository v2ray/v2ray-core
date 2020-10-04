// +build !confonly

package tls

import (
	"crypto/tls"

	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/net"
)

//go:generate go run v2ray.com/core/common/errors/errorgen

var (
	_ buf.Writer = (*Conn)(nil)
)

type Conn struct {
	*tls.Conn
}

func (c *Conn) WriteMultiBuffer(mb buf.MultiBuffer) error {
	mb = buf.Compact(mb)
	mb, err := buf.WriteMultiBuffer(c, mb)
	buf.ReleaseMulti(mb)
	return err
}

func (c *Conn) HandshakeAddress() net.Address {
	if err := c.Handshake(); err != nil {
		return nil
	}
	state := c.ConnectionState()
	if state.ServerName == "" {
		return nil
	}
	return net.ParseAddress(state.ServerName)
}

// Client initiates a TLS client handshake on the given connection.
func Client(c net.Conn, config *tls.Config) net.Conn {
	tlsConn := tls.Client(c, config)
	return &Conn{Conn: tlsConn}
}

/*
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
*/

// Server initiates a TLS server handshake on the given connection.
func Server(c net.Conn, config *tls.Config) net.Conn {
	tlsConn := tls.Server(c, config)
	return &Conn{Conn: tlsConn}
}
