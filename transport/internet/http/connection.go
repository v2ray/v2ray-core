package http

import (
	"io"
	"time"

	"v2ray.com/core/common"
	"v2ray.com/core/common/net"
)

type Connection struct {
	Reader io.Reader
	Writer io.Writer
	Closer common.Closable
	Local  net.Addr
	Remote net.Addr
}

func (c *Connection) Read(b []byte) (int, error) {
	return c.Reader.Read(b)
}

func (c *Connection) Write(b []byte) (int, error) {
	return c.Writer.Write(b)
}

func (c *Connection) Close() error {
	return c.Closer.Close()
}

func (c *Connection) LocalAddr() net.Addr {
	return c.Local
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.Remote
}

func (c *Connection) SetDeadline(t time.Time) error {
	return nil
}

func (c *Connection) SetReadDeadline(t time.Time) error {
	return nil
}

func (c *Connection) SetWriteDeadline(t time.Time) error {
	return nil
}
