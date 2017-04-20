package tls

import (
	"crypto/tls"
	"net"

	"v2ray.com/core/common/buf"
)

//go:generate go run $GOPATH/src/v2ray.com/core/tools/generrorgen/main.go -pkg tls -path Transport,Internet,TLS

type conn struct {
	net.Conn

	mergingReader buf.Reader
	mergingWriter buf.Writer
}

func (c *conn) ReadMultiBuffer() (buf.MultiBuffer, error) {
	if c.mergingReader == nil {
		c.mergingReader = buf.NewMergingReaderSize(c.Conn, 16*1024)
	}
	return c.mergingReader.Read()
}

func (c *conn) WriteMultiBuffer(mb buf.MultiBuffer) (int, error) {
	if c.mergingWriter == nil {
		c.mergingWriter = buf.NewMergingWriter(c.Conn)
	}
	return mb.Len(), c.mergingWriter.Write(mb)
}

func Client(c net.Conn, config *tls.Config) net.Conn {
	tlsConn := tls.Client(c, config)
	return &conn{Conn: tlsConn}
}

func Server(c net.Conn, config *tls.Config) net.Conn {
	tlsConn := tls.Server(c, config)
	return &conn{Conn: tlsConn}
}
