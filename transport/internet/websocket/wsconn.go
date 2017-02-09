package websocket

import (
	"io"
	"net"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"v2ray.com/core/common/errors"
)

type wsconn struct {
	wsc    *websocket.Conn
	reader io.Reader
	rlock  sync.Mutex
	wlock  sync.Mutex
}

func (c *wsconn) Read(b []byte) (int, error) {
	c.rlock.Lock()
	defer c.rlock.Unlock()

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

func (c *wsconn) getReader() (io.Reader, error) {
	if c.reader != nil {
		return c.reader, nil
	}

	_, reader, err := c.wsc.NextReader()
	if err != nil {
		return nil, err
	}
	c.reader = reader
	return reader, nil
}

func (c *wsconn) Write(b []byte) (int, error) {
	c.wlock.Lock()
	defer c.wlock.Unlock()

	if err := c.wsc.WriteMessage(websocket.BinaryMessage, b); err != nil {
		return 0, err
	}
	return len(b), nil
}

func (c *wsconn) Close() error {
	c.wsc.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""), time.Now().Add(time.Second*5))
	return c.wsc.Close()
}

func (c *wsconn) LocalAddr() net.Addr {
	return c.wsc.LocalAddr()
}

func (c *wsconn) RemoteAddr() net.Addr {
	return c.wsc.RemoteAddr()
}

func (c *wsconn) SetDeadline(t time.Time) error {
	if err := c.SetReadDeadline(t); err != nil {
		return err
	}
	return c.SetWriteDeadline(t)
}

func (c *wsconn) SetReadDeadline(t time.Time) error {
	return c.wsc.SetReadDeadline(t)
}

func (c *wsconn) SetWriteDeadline(t time.Time) error {
	return c.wsc.SetWriteDeadline(t)
}
