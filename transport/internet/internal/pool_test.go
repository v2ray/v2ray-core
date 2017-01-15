package internal_test

import (
	"net"
	"testing"
	"time"

	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/testing/assert"
	. "v2ray.com/core/transport/internet/internal"
)

type TestConnection struct {
	id     string
	closed bool
}

func (o *TestConnection) Read([]byte) (int, error) {
	return 0, nil
}

func (o *TestConnection) Write([]byte) (int, error) {
	return 0, nil
}

func (o *TestConnection) Close() error {
	o.closed = true
	return nil
}

func (o *TestConnection) LocalAddr() net.Addr {
	return nil
}

func (o *TestConnection) RemoteAddr() net.Addr {
	return nil
}

func (o *TestConnection) SetDeadline(t time.Time) error {
	return nil
}

func (o *TestConnection) SetReadDeadline(t time.Time) error {
	return nil
}

func (o *TestConnection) SetWriteDeadline(t time.Time) error {
	return nil
}

func TestConnectionCache(t *testing.T) {
	assert := assert.On(t)

	pool := NewConnectionPool()
	conn := pool.Get(NewConnectionID(v2net.LocalHostIP, v2net.TCPDestination(v2net.LocalHostIP, v2net.Port(80))))
	assert.Pointer(conn).IsNil()

	pool.Put(NewConnectionID(v2net.LocalHostIP, v2net.TCPDestination(v2net.LocalHostIP, v2net.Port(80))), &TestConnection{id: "test"})
	conn = pool.Get(NewConnectionID(v2net.LocalHostIP, v2net.TCPDestination(v2net.LocalHostIP, v2net.Port(80))))
	assert.String(conn.(*TestConnection).id).Equals("test")
}

func TestConnectionRecycle(t *testing.T) {
	assert := assert.On(t)

	pool := NewConnectionPool()
	c := &TestConnection{id: "test"}
	pool.Put(NewConnectionID(v2net.LocalHostIP, v2net.TCPDestination(v2net.LocalHostIP, v2net.Port(80))), c)
	time.Sleep(6 * time.Second)
	assert.Bool(c.closed).IsTrue()
	conn := pool.Get(NewConnectionID(v2net.LocalHostIP, v2net.TCPDestination(v2net.LocalHostIP, v2net.Port(80))))
	assert.Pointer(conn).IsNil()
}
