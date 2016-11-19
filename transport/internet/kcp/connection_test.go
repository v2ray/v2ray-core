package kcp_test

import (
	"testing"
	"time"

	"v2ray.com/core/testing/assert"
	. "v2ray.com/core/transport/internet/kcp"
)

type NoOpWriteCloser struct{}

func (this *NoOpWriteCloser) Write(b []byte) (int, error) {
	return len(b), nil
}

func (this *NoOpWriteCloser) Close() error {
	return nil
}

func TestConnectionReadTimeout(t *testing.T) {
	assert := assert.On(t)

	conn := NewConnection(1, &NoOpWriteCloser{}, nil, nil, NewSimpleAuthenticator(), &Config{})
	conn.SetReadDeadline(time.Now().Add(time.Second))

	b := make([]byte, 1024)
	nBytes, err := conn.Read(b)
	assert.Int(nBytes).Equals(0)
	assert.Error(err).IsNotNil()

	conn.Terminate()
}
