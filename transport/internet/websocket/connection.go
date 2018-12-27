package websocket

import (
	"io"
	"net"
	"time"

	"github.com/gorilla/websocket"

	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/errors"
	"v2ray.com/core/common/serial"
)

var (
	_ buf.Writer = (*connection)(nil)
)

// connection is a wrapper for net.Conn over WebSocket connection.
type connection struct {
	conn       *websocket.Conn
	reader     io.Reader
	remoteAddr net.Addr
}

func newConnection(conn *websocket.Conn, remoteAddr net.Addr) *connection {
	return &connection{
		conn:       conn,
		remoteAddr: remoteAddr,
	}
}

// Read implements net.Conn.Read()
func (c *connection) Read(b []byte) (int, error) {
	for {
		reader, err := c.getReader()
		if err != nil {
			return 0, err
		}

		nBytes, err := reader.Read(b)
		if errors.Cause(err) == io.EOF {
			c.reader = nil
			continue
		}
		return nBytes, err
	}
}

func (c *connection) getReader() (io.Reader, error) {
	if c.reader != nil {
		return c.reader, nil
	}

	_, reader, err := c.conn.NextReader()
	if err != nil {
		return nil, err
	}
	c.reader = reader
	return reader, nil
}

// Write implements io.Writer.
func (c *connection) Write(b []byte) (int, error) {
	if err := c.conn.WriteMessage(websocket.BinaryMessage, b); err != nil {
		return 0, err
	}
	return len(b), nil
}

func (c *connection) WriteMultiBuffer(mb buf.MultiBuffer) error {
	mb = buf.Compact(mb)

	for {
		mb2, b := buf.SplitFirst(mb)
		mb = mb2
		if b == nil {
			break
		}

		if err := buf.WriteAllBytes(c, b.Bytes()); err != nil {
			buf.ReleaseMulti(mb)
			return err
		}
	}

	return nil
}

func (c *connection) Close() error {
	var errors []interface{}
	if err := c.conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""), time.Now().Add(time.Second*5)); err != nil {
		errors = append(errors, err)
	}
	if err := c.conn.Close(); err != nil {
		errors = append(errors, err)
	}
	if len(errors) > 0 {
		return newError("failed to close connection").Base(newError(serial.Concat(errors...)))
	}
	return nil
}

func (c *connection) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}

func (c *connection) RemoteAddr() net.Addr {
	return c.remoteAddr
}

func (c *connection) SetDeadline(t time.Time) error {
	if err := c.SetReadDeadline(t); err != nil {
		return err
	}
	return c.SetWriteDeadline(t)
}

func (c *connection) SetReadDeadline(t time.Time) error {
	return c.conn.SetReadDeadline(t)
}

func (c *connection) SetWriteDeadline(t time.Time) error {
	return c.conn.SetWriteDeadline(t)
}
