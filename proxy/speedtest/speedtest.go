package speedtest

//go:generate go run $GOPATH/src/v2ray.com/core/common/errors/errorgen/main.go -pkg speedtest -path Proxy,SpeedTest

import (
	"bufio"
	"context"
	"io"
	"net/http"
	"strconv"
	"strings"

	"v2ray.com/core"
	"v2ray.com/core/common"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/serial"
	"v2ray.com/core/common/session"
	"v2ray.com/core/proxy"
)

var rndBytes = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!\"#$%&\\'()*+,-")

type rndBytesReader struct{}

func (rndBytesReader) Read(b []byte) (int, error) {
	totalBytes := 0
	for totalBytes < len(b) {
		nBytes := copy(b[totalBytes:], rndBytes[:])
		totalBytes += nBytes
	}
	return totalBytes, nil
}

func (rndBytesReader) Close() error {
	return nil
}

type SpeedTestHandler struct{}

func New(ctx context.Context, config *Config) (*SpeedTestHandler, error) {
	return &SpeedTestHandler{}, nil
}

type noOpCloser struct {
	io.Reader
}

func (c *noOpCloser) Close() error {
	return nil
}

func defaultResponse() *http.Response {
	response := &http.Response{
		Status:        "Not Found",
		StatusCode:    404,
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		Header:        http.Header(make(map[string][]string)),
		Body:          nil,
		ContentLength: 0,
		Close:         true,
	}
	response.Header.Set("Content-Type", "text/plain; charset=UTF-8")
	return response
}

func (h *SpeedTestHandler) Process(ctx context.Context, link *core.Link, dialer proxy.Dialer) error {
	reader := link.Reader
	writer := link.Writer

	defer func() {
		common.Close(writer)
	}()

	bufReader := bufio.NewReader(&buf.BufferedReader{
		Reader: reader,
		Direct: true,
	})

	bufWriter := buf.NewBufferedWriter(writer)
	common.Must(bufWriter.SetBuffered(false))

	request, err := http.ReadRequest(bufReader)
	if err != nil {
		return newError("failed to read speedtest request").Base(err)
	}

	path := strings.ToLower(request.URL.Path)
	if len(path) > 0 && path[0] == '/' {
		path = path[1:]
	}

	switch path {
	case "hello":
		respBody := "hello 2.5 2017-08-15.1314.4ae12d5"
		response := &http.Response{
			Status:        "OK",
			StatusCode:    200,
			Proto:         "HTTP/1.1",
			ProtoMajor:    1,
			ProtoMinor:    1,
			Header:        http.Header(make(map[string][]string)),
			Body:          &noOpCloser{strings.NewReader(respBody)},
			ContentLength: int64(len(respBody)),
			Close:         true,
		}
		response.Header.Set("Content-Type", "text/plain; charset=UTF-8")
		return response.Write(bufWriter)
	case "upload":
		switch strings.ToUpper(request.Method) {
		case "POST":
			var sc buf.SizeCounter
			buf.Copy(buf.NewReader(request.Body), buf.Discard, buf.CountSize(&sc)) // nolint: errcheck

			response := &http.Response{
				Status:     "OK",
				StatusCode: 200,
				Proto:      "HTTP/1.1",
				ProtoMajor: 1,
				ProtoMinor: 1,
				Header:     http.Header(make(map[string][]string)),
				Body:       &noOpCloser{strings.NewReader(serial.Concat("size=", sc.Size))},
				Close:      true,
			}
			response.Header.Set("Content-Type", "text/plain; charset=UTF-8")
			return response.Write(bufWriter)
		case "OPTIONS":
			response := &http.Response{
				Status:        "OK",
				StatusCode:    200,
				Proto:         "HTTP/1.1",
				ProtoMajor:    1,
				ProtoMinor:    1,
				Header:        http.Header(make(map[string][]string)),
				Body:          nil,
				ContentLength: 0,
				Close:         true,
			}

			response.Header.Set("Content-Type", "text/plain; charset=UTF-8")
			response.Header.Set("Connection", "Close")
			response.Header.Set("Access-Control-Allow-Methods", "OPTIONS, POST")
			response.Header.Set("Access-Control-Allow-Headers", "content-type")
			response.Header.Set("Access-Control-Allow-Origin", "http://www.speedtest.net")
			return response.Write(bufWriter)
		default:
			newError("unknown method for upload: ", request.Method).WriteToLog(session.ExportIDToError(ctx))
			return defaultResponse().Write(bufWriter)
		}
	case "download":
		query := request.URL.Query()
		sizeStr := query.Get("size")
		size, err := strconv.Atoi(sizeStr)
		if err != nil {
			return defaultResponse().Write(bufWriter)
		}
		response := &http.Response{
			Status:        "OK",
			StatusCode:    200,
			Proto:         "HTTP/1.1",
			ProtoMajor:    1,
			ProtoMinor:    1,
			Header:        http.Header(make(map[string][]string)),
			Body:          rndBytesReader{},
			ContentLength: int64(size),
			Close:         true,
		}
		response.Header.Set("Content-Type", "text/plain; charset=UTF-8")
		return response.Write(bufWriter)
	default:
		newError("unknown path: ", path).WriteToLog(session.ExportIDToError(ctx))
		return defaultResponse().Write(bufWriter)
	}
}

func init() {
	common.Must(common.RegisterConfig((*Config)(nil), func(ctx context.Context, config interface{}) (interface{}, error) {
		return New(ctx, config.(*Config))
	}))
}
