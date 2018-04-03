package ray

import (
	"io"
	"net"
	"time"

	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/signal"
)

type ConnectionOption func(*connection)

func ConnLocalAddr(addr net.Addr) ConnectionOption {
	return func(c *connection) {
		c.localAddr = addr
	}
}

func ConnRemoteAddr(addr net.Addr) ConnectionOption {
	return func(c *connection) {
		c.remoteAddr = addr
	}
}

func ConnCloseSignal(s *signal.Notifier) ConnectionOption {
	return func(c *connection) {
		c.closeSignal = s
	}
}

type connection struct {
	input       InputStream
	output      OutputStream
	closed      bool
	localAddr   net.Addr
	remoteAddr  net.Addr
	closeSignal *signal.Notifier

	reader *buf.BufferedReader
}

var zeroAddr net.Addr = &net.TCPAddr{IP: []byte{0, 0, 0, 0}}

// NewConnection wraps a Ray into net.Conn.
func NewConnection(input InputStream, output OutputStream, options ...ConnectionOption) net.Conn {
	c := &connection{
		input:      input,
		output:     output,
		localAddr:  zeroAddr,
		remoteAddr: zeroAddr,
		reader:     buf.NewBufferedReader(input),
	}

	for _, opt := range options {
		opt(c)
	}

	return c
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
	mb := buf.NewMultiBufferCap(int32(l)/buf.Size + 1)
	mb.Write(b)
	return l, c.output.WriteMultiBuffer(mb)
}

func (c *connection) WriteMultiBuffer(mb buf.MultiBuffer) error {
	if c.closed {
		return io.ErrClosedPipe
	}

	return c.output.WriteMultiBuffer(mb)
}

// Close implements net.Conn.Close().
func (c *connection) Close() error {
	c.closed = true
	c.output.Close()
	c.input.CloseError()
	if c.closeSignal != nil {
		c.closeSignal.Signal()
	}
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

// SetWriteDeadline implements net.Conn.SetWriteDeadline().
func (c *connection) SetWriteDeadline(t time.Time) error {
	return nil
}
