package http_test

import (
	"bytes"
	"context"
	"crypto/rand"
	"testing"
	"time"

	"v2ray.com/core/common"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/net"
	. "v2ray.com/core/transport/internet/headers/http"
	. "v2ray.com/ext/assert"
)

func TestReaderWriter(t *testing.T) {
	assert := With(t)

	cache := buf.New()
	b := buf.New()
	common.Must2(b.WriteString("abcd" + ENDING))
	writer := NewHeaderWriter(b)
	err := writer.Write(cache)
	common.Must(err)
	assert(cache.Len(), Equals, int32(8))
	_, err = cache.Write([]byte{'e', 'f', 'g'})
	common.Must(err)

	reader := &HeaderReader{}
	buffer, err := reader.Read(cache)
	common.Must(err)
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
	common.Must(err)

	cache := buf.New()
	err = auth.GetClientWriter().Write(cache)
	common.Must(err)

	assert(cache.String(), Equals, "GET / HTTP/1.1\r\nTest: Value\r\n\r\n")
}

func TestLongRequestHeader(t *testing.T) {
	payload := make([]byte, buf.Size+2)
	common.Must2(rand.Read(payload[:buf.Size-2]))
	copy(payload[buf.Size-2:], []byte(ENDING))
	payload = append(payload, []byte("abcd")...)

	reader := HeaderReader{}
	b, err := reader.Read(bytes.NewReader(payload))
	common.Must(err)
	if b.String() != "abcd" {
		t.Error("expect content abcd, but actually ", b.String())
	}
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
	common.Must(err)

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	common.Must(err)

	go func() {
		conn, err := listener.Accept()
		common.Must(err)
		authConn := auth.Server(conn)
		b := make([]byte, 256)
		for {
			n, err := authConn.Read(b)
			common.Must(err)
			_, err = authConn.Write(b[:n])
			common.Must(err)
		}
	}()

	conn, err := net.DialTCP("tcp", nil, listener.Addr().(*net.TCPAddr))
	common.Must(err)

	authConn := auth.Client(conn)
	authConn.Write([]byte("Test payload"))
	authConn.Write([]byte("Test payload 2"))

	expectedResponse := "Test payloadTest payload 2"
	actualResponse := make([]byte, 256)
	deadline := time.Now().Add(time.Second * 5)
	totalBytes := 0
	for {
		n, err := authConn.Read(actualResponse[totalBytes:])
		common.Must(err)
		totalBytes += n
		if totalBytes >= len(expectedResponse) || time.Now().After(deadline) {
			break
		}
	}

	assert(string(actualResponse[:totalBytes]), Equals, expectedResponse)
}
