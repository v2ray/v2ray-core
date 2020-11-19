// +build !confonly

package xtls

import (
	xtls "github.com/xtls/go"

	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/net"
)

//go:generate go run v2ray.com/core/common/errors/errorgen

var (
	_ buf.Writer = (*Conn)(nil)
)

type Conn struct {
	*xtls.Conn
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

// Client initiates a XTLS client handshake on the given connection.
func Client(c net.Conn, config *xtls.Config) net.Conn {
	xtlsConn := xtls.Client(c, config)
	return &Conn{Conn: xtlsConn}
}

// Server initiates a XTLS server handshake on the given connection.
func Server(c net.Conn, config *xtls.Config) net.Conn {
	xtlsConn := xtls.Server(c, config)
	return &Conn{Conn: xtlsConn}
}
