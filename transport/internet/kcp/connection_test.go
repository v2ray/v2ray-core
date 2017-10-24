package kcp_test

import (
	"net"
	"testing"
	"time"

	. "v2ray.com/ext/assert"
	. "v2ray.com/core/transport/internet/kcp"
)

type NoOpConn struct{}

func (o *NoOpConn) Overhead() int {
	return 0
}

// Write implements io.Writer.
func (o *NoOpConn) Write(b []byte) (int, error) {
	return len(b), nil
}

func (o *NoOpConn) Close() error {
	return nil
}

func (o *NoOpConn) Read([]byte) (int, error) {
	panic("Should not be called.")
}

func (o *NoOpConn) LocalAddr() net.Addr {
	return nil
}

func (o *NoOpConn) RemoteAddr() net.Addr {
	return nil
}

func (o *NoOpConn) SetDeadline(time.Time) error {
	return nil
}

func (o *NoOpConn) SetReadDeadline(time.Time) error {
	return nil
}

func (o *NoOpConn) SetWriteDeadline(time.Time) error {
	return nil
}

func (o *NoOpConn) Reset(input func([]Segment)) {}

func TestConnectionReadTimeout(t *testing.T) {
	assert := With(t)

	conn := NewConnection(1, &NoOpConn{}, &Config{})
	conn.SetReadDeadline(time.Now().Add(time.Second))

	b := make([]byte, 1024)
	nBytes, err := conn.Read(b)
	assert(nBytes, Equals, 0)
	assert(err, IsNotNil)

	conn.Terminate()
}
