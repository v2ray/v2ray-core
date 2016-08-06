package kcp_test

import (
	"crypto/rand"
	"io"
	"net"
	"testing"
	"time"

	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/testing/assert"
	"github.com/v2ray/v2ray-core/transport/internet"
	"github.com/v2ray/v2ray-core/transport/internet/authenticators/srtp"
	. "github.com/v2ray/v2ray-core/transport/internet/kcp"
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

	conn := NewConnection(1, &NoOpWriteCloser{}, nil, nil, NewSimpleAuthenticator())
	conn.SetReadDeadline(time.Now().Add(time.Second))

	b := make([]byte, 1024)
	nBytes, err := conn.Read(b)
	assert.Int(nBytes).Equals(0)
	assert.Error(err).IsNotNil()
}

func TestConnectionReadWrite(t *testing.T) {
	assert := assert.On(t)

	upReader, upWriter := io.Pipe()
	downReader, downWriter := io.Pipe()

	auth := internet.NewAuthenticatorChain(srtp.ObfuscatorSRTPFactory{}.Create(nil), NewSimpleAuthenticator())

	connClient := NewConnection(1, upWriter, &net.UDPAddr{IP: v2net.LocalHostIP.IP(), Port: 1}, &net.UDPAddr{IP: v2net.LocalHostIP.IP(), Port: 2}, auth)
	connClient.FetchInputFrom(downReader)

	connServer := NewConnection(1, downWriter, &net.UDPAddr{IP: v2net.LocalHostIP.IP(), Port: 2}, &net.UDPAddr{IP: v2net.LocalHostIP.IP(), Port: 1}, auth)
	connServer.FetchInputFrom(upReader)

	totalWritten := 1024 * 1024
	clientSend := make([]byte, totalWritten)
	rand.Read(clientSend)
	go func() {
		nBytes, err := connClient.Write(clientSend)
		assert.Int(nBytes).Equals(totalWritten)
		assert.Error(err).IsNil()
	}()

	serverReceived := make([]byte, totalWritten)
	totalRead := 0
	for totalRead < totalWritten {
		nBytes, err := connServer.Read(serverReceived[totalRead:])
		assert.Error(err).IsNil()
		totalRead += nBytes
	}
	assert.Bytes(serverReceived).Equals(clientSend)

	connClient.Close()
	connServer.Close()

	for connClient.State() != StateTerminated || connServer.State() != StateTerminated {
		time.Sleep(time.Second)
	}
}
