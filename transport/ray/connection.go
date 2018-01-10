package ray

import (
	"io"
	"net"
	"time"

	"v2ray.com/core/common/buf"
)

type connection struct {
	stream     Ray
	closed     bool
	localAddr  net.Addr
	remoteAddr net.Addr

	reader *buf.BufferedReader
	writer buf.Writer
}

// NewConnection wraps a Ray into net.Conn.
func NewConnection(stream InboundRay, localAddr net.Addr, remoteAddr net.Addr) net.Conn {
	return &connection{
		stream:     stream,
		localAddr:  localAddr,
		remoteAddr: remoteAddr,
		reader:     buf.NewBufferedReader(stream.InboundOutput()),
		writer:     stream.InboundInput(),
	}
}

// Read implements net.Conn.Read().
func (c *connection) Read(b []byte) (int, error) {
	if c.closed {
		return 0, io.EOF
	}
	return c.reader.Read(b)
}

// ReadMultiBuffer implements buf.Reader.
func (c *connection) ReadMultiBuffer() (buf.MultiBuffer, error) {
	return c.reader.ReadMultiBuffer()
}

// Write implements net.Conn.Write().
func (c *connection) Write(b []byte) (int, error) {
	if c.closed {
		return 0, io.ErrClosedPipe
	}

	l := len(b)
	mb := buf.NewMultiBufferCap(l/buf.Size + 1)
	mb.Write(b)
	return l, c.writer.WriteMultiBuffer(mb)
}

func (c *connection) WriteMultiBuffer(mb buf.MultiBuffer) error {
	if c.closed {
		return io.ErrClosedPipe
	}

	return c.writer.WriteMultiBuffer(mb)
}

// Close implements net.Conn.Close().
func (c *connection) Close() error {
	c.closed = true
	c.stream.InboundInput().Close()
	c.stream.InboundOutput().CloseError()
	return nil
}

// LocalAddr implements net.Conn.LocalAddr().
func (c *connection) LocalAddr() net.Addr {
	return c.localAddr
}

// RemoteAddr implements net.Conn.RemoteAddr().
func (c *connection) RemoteAddr() net.Addr {
	return c.remoteAddr
}

// SetDeadline implements net.Conn.SetDeadline().
func (c *connection) SetDeadline(t time.Time) error {
	return nil
}

// SetReadDeadline implements net.Conn.SetReadDeadline().
func (c *connection) SetReadDeadline(t time.Time) error {
	return nil
}

// SetWriteDeadline implement net.Conn.SetWriteDeadline().
func (c *connection) SetWriteDeadline(t time.Time) error {
	return nil
}
