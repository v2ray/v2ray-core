package quic

import (
	"time"

	quic "github.com/lucas-clemente/quic-go"

	"v2ray.com/core/common/net"
	"v2ray.com/core/transport/internet"
)

type sysConn struct {
	conn   net.PacketConn
	header internet.PacketHeader
}

func (c *sysConn) ReadFrom(p []byte) (int, net.Addr, error) {
	if c.header == nil {
		return c.conn.ReadFrom(p)
	}

	overhead := int(c.header.Size())
	buffer := getBuffer()
	defer putBuffer(buffer)

	nBytes, addr, err := c.conn.ReadFrom(buffer[:len(p)+overhead])
	if err != nil {
		return 0, nil, err
	}

	copy(p, buffer[overhead:nBytes])

	return nBytes - overhead, addr, nil
}

func (c *sysConn) WriteTo(p []byte, addr net.Addr) (int, error) {
	if c.header == nil {
		return c.conn.WriteTo(p, addr)
	}

	buffer := getBuffer()
	defer putBuffer(buffer)

	overhead := int(c.header.Size())
	c.header.Serialize(buffer)
	copy(buffer[overhead:], p)
	nBytes, err := c.conn.WriteTo(buffer[:len(p)+overhead], addr)
	if err != nil {
		return 0, err
	}
	return nBytes - overhead, nil
}

func (c *sysConn) Close() error {
	return c.conn.Close()
}

func (c *sysConn) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}

func (c *sysConn) SetDeadline(t time.Time) error {
	return c.conn.SetDeadline(t)
}

func (c *sysConn) SetReadDeadline(t time.Time) error {
	return c.conn.SetReadDeadline(t)
}

func (c *sysConn) SetWriteDeadline(t time.Time) error {
	return c.conn.SetWriteDeadline(t)
}

type interConn struct {
	stream quic.Stream
	local  net.Addr
	remote net.Addr
}

func (c *interConn) Read(b []byte) (int, error) {
	return c.stream.Read(b)
}

func (c *interConn) Write(b []byte) (int, error) {
	return c.stream.Write(b)
}

func (c *interConn) Close() error {
	return c.stream.Close()
}

func (c *interConn) LocalAddr() net.Addr {
	return c.local
}

func (c *interConn) RemoteAddr() net.Addr {
	return c.remote
}

func (c *interConn) SetDeadline(t time.Time) error {
	return c.stream.SetDeadline(t)
}

func (c *interConn) SetReadDeadline(t time.Time) error {
	return c.stream.SetReadDeadline(t)
}

func (c *interConn) SetWriteDeadline(t time.Time) error {
	return c.stream.SetWriteDeadline(t)
}
