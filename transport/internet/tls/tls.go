package tls

import (
	"crypto/tls"
	"net"

	"v2ray.com/core/common/buf"
)

//go:generate go run $GOPATH/src/v2ray.com/core/common/errors/errorgen/main.go -pkg tls -path Transport,Internet,TLS

var (
	_ buf.Writer = (*conn)(nil)
)

type conn struct {
	net.Conn

	mergingWriter *buf.BufferedWriter
}

func (c *conn) WriteMultiBuffer(mb buf.MultiBuffer) error {
	if c.mergingWriter == nil {
		c.mergingWriter = buf.NewBufferedWriter(buf.NewWriter(c.Conn))
	}
	if err := c.mergingWriter.WriteMultiBuffer(mb); err != nil {
		return err
	}
	return c.mergingWriter.Flush()
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
