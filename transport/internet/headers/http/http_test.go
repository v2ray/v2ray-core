package http_test

import (
	"net"
	"testing"

	"time"

	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/serial"
	"v2ray.com/core/testing/assert"
	. "v2ray.com/core/transport/internet/headers/http"
)

func TestReaderWriter(t *testing.T) {
	assert := assert.On(t)

	cache := buf.New()
	b := buf.NewLocal(256)
	b.AppendSupplier(serial.WriteString("abcd" + ENDING))
	writer := NewHeaderWriter(b)
	err := writer.Write(cache)
	assert.Error(err).IsNil()
	assert.Int(cache.Len()).Equals(8)
	_, err = cache.Write([]byte{'e', 'f', 'g'})
	assert.Error(err).IsNil()

	reader := &HeaderReader{}
	buffer, err := reader.Read(cache)
	assert.Error(err).IsNil()
	assert.Bytes(buffer.Bytes()).Equals([]byte{'e', 'f', 'g'})
}

func TestRequestHeader(t *testing.T) {
	assert := assert.On(t)

	factory := HttpAuthenticatorFactory{}
	auth := factory.Create(&Config{
		Request: &RequestConfig{
			Uri: []string{"/"},
			Header: []*Header{
				{
					Name:  "Test",
					Value: []string{"Value"},
				},
			},
		},
	}).(HttpAuthenticator)

	cache := buf.New()
	err := auth.GetClientWriter().Write(cache)
	assert.Error(err).IsNil()

	assert.String(cache.String()).Equals("GET / HTTP/1.1\r\nTest: Value\r\n\r\n")
}

func TestConnection(t *testing.T) {
	assert := assert.On(t)

	factory := HttpAuthenticatorFactory{}
	auth := factory.Create(new(Config))

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	assert.Error(err).IsNil()

	go func() {
		conn, err := listener.Accept()
		assert.Error(err).IsNil()
		authConn := auth.Server(conn)
		b := make([]byte, 256)
		for {
			n, err := authConn.Read(b)
			assert.Error(err).IsNil()
			_, err = authConn.Write(b[:n])
			assert.Error(err).IsNil()
		}
	}()

	conn, err := net.DialTCP("tcp", nil, listener.Addr().(*net.TCPAddr))
	assert.Error(err).IsNil()

	authConn := auth.Client(conn)
	authConn.Write([]byte("Test payload"))
	authConn.Write([]byte("Test payload 2"))

	expectedResponse := "Test payloadTest payload 2"
	actualResponse := make([]byte, 256)
	deadline := time.Now().Add(time.Second * 5)
	totalBytes := 0
	for {
		n, err := authConn.Read(actualResponse[totalBytes:])
		assert.Error(err).IsNil()
		totalBytes += n
		if totalBytes >= len(expectedResponse) || time.Now().After(deadline) {
			break
		}
	}

	assert.String(string(actualResponse[:totalBytes])).Equals(expectedResponse)
}
