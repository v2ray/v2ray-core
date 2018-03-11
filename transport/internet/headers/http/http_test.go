package http_test

import (
	"context"
	"testing"
	"time"

	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/net"
	"v2ray.com/core/common/serial"
	. "v2ray.com/core/transport/internet/headers/http"
	. "v2ray.com/ext/assert"
)

func TestReaderWriter(t *testing.T) {
	assert := With(t)

	cache := buf.New()
	b := buf.NewSize(256)
	b.AppendSupplier(serial.WriteString("abcd" + ENDING))
	writer := NewHeaderWriter(b)
	err := writer.Write(cache)
	assert(err, IsNil)
	assert(cache.Len(), Equals, 8)
	_, err = cache.Write([]byte{'e', 'f', 'g'})
	assert(err, IsNil)

	reader := &HeaderReader{}
	buffer, err := reader.Read(cache)
	assert(err, IsNil)
	assert(buffer.Bytes(), Equals, []byte{'e', 'f', 'g'})
}

func TestRequestHeader(t *testing.T) {
	assert := With(t)

	auth, err := NewHttpAuthenticator(context.Background(), &Config{
		Request: &RequestConfig{
			Uri: []string{"/"},
			Header: []*Header{
				{
					Name:  "Test",
					Value: []string{"Value"},
				},
			},
		},
	})
	assert(err, IsNil)

	cache := buf.New()
	err = auth.GetClientWriter().Write(cache)
	assert(err, IsNil)

	assert(cache.String(), Equals, "GET / HTTP/1.1\r\nTest: Value\r\n\r\n")
}

func TestConnection(t *testing.T) {
	assert := With(t)

	auth, err := NewHttpAuthenticator(context.Background(), &Config{
		Request: &RequestConfig{
			Method: &Method{Value: "Post"},
			Uri:    []string{"/testpath"},
			Header: []*Header{
				{
					Name:  "Host",
					Value: []string{"www.v2ray.com", "www.google.com"},
				},
				{
					Name:  "User-Agent",
					Value: []string{"Test-Agent"},
				},
			},
		},
		Response: &ResponseConfig{
			Version: &Version{
				Value: "1.1",
			},
			Status: &Status{
				Code:   "404",
				Reason: "Not Found",
			},
		},
	})
	assert(err, IsNil)

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	assert(err, IsNil)

	go func() {
		conn, err := listener.Accept()
		assert(err, IsNil)
		authConn := auth.Server(conn)
		b := make([]byte, 256)
		for {
			n, err := authConn.Read(b)
			assert(err, IsNil)
			_, err = authConn.Write(b[:n])
			assert(err, IsNil)
		}
	}()

	conn, err := net.DialTCP("tcp", nil, listener.Addr().(*net.TCPAddr))
	assert(err, IsNil)

	authConn := auth.Client(conn)
	authConn.Write([]byte("Test payload"))
	authConn.Write([]byte("Test payload 2"))

	expectedResponse := "Test payloadTest payload 2"
	actualResponse := make([]byte, 256)
	deadline := time.Now().Add(time.Second * 5)
	totalBytes := 0
	for {
		n, err := authConn.Read(actualResponse[totalBytes:])
		assert(err, IsNil)
		totalBytes += n
		if totalBytes >= len(expectedResponse) || time.Now().After(deadline) {
			break
		}
	}

	assert(string(actualResponse[:totalBytes]), Equals, expectedResponse)
}
