package kcp_test

import (
	"io"
	"testing"
	"time"

	"v2ray.com/core/common/buf"
	. "v2ray.com/core/transport/internet/kcp"
	. "v2ray.com/ext/assert"
)

type NoOpCloser int

func (NoOpCloser) Close() error {
	return nil
}

func TestConnectionReadTimeout(t *testing.T) {
	assert := With(t)

	conn := NewConnection(1, &ConnMetadata{}, &KCPPacketWriter{
		Writer: buf.DiscardBytes,
	}, NoOpCloser(0), &Config{})
	conn.SetReadDeadline(time.Now().Add(time.Second))

	b := make([]byte, 1024)
	nBytes, err := conn.Read(b)
	assert(nBytes, Equals, 0)
	assert(err, IsNotNil)

	conn.Terminate()
}

func TestConnectionInterface(t *testing.T) {
	assert := With(t)

	assert((*Connection)(nil), Implements, (*io.Writer)(nil))
	assert((*Connection)(nil), Implements, (*io.Reader)(nil))
	assert((*Connection)(nil), Implements, (*buf.Reader)(nil))
}
