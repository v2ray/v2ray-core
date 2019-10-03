package router_test

import (
	"context"
	"testing"
	"time"

	"v2ray.com/core/app/proxyman/outbound"
	. "v2ray.com/core/app/router"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/transport"
)

// mock proxy/outbound/handler
type mockHandler struct {
	tag     string
	timeout time.Duration
}

func (h *mockHandler) Tag() string {
	return h.tag
}

func (h *mockHandler) Start() error {
	return nil
}

func (h *mockHandler) Close() error {
	return nil
}

func (h *mockHandler) Dispatch(ctx context.Context, link *transport.Link) {
	mockHTTPResponse := `HTTP/1.1 200 OK
Date: Mon, 27 Jul 2080 12:28:53 GMT
Server: MockServer/0.0.1
Content-Length: 53
Content-Type: text/html
Connection: Closed

<html>
<body>
<h1>Hello, World!</h1>
</body>
</html>
`
	link.Reader.ReadMultiBuffer()
	if h.timeout != 0 {
		time.Sleep(h.timeout)
	}
	link.Writer.WriteMultiBuffer(buf.MergeBytes(buf.MultiBuffer{}, []byte(mockHTTPResponse)))
}

func TestRandomStrategy(t *testing.T) {
	strategy := RandomStrategy{}
	if strategy.PickOutbound(nil, []string{"test"}) != "test" {
		t.Error("Random strategy test fail")
	}
}

func TestOptimalStrategy(t *testing.T) {
	ctx := context.Background()
	obm, _ := outbound.New(ctx, nil)
	obm.AddHandler(ctx, &mockHandler{tag: "test1", timeout: time.Millisecond * 100})
	obm.AddHandler(ctx, &mockHandler{tag: "test2"})
	strategy := NewOptimalStrategy(&OptimalStrategyConfig{URL: "http://test.com"})

	tag := strategy.PickOutbound(obm, []string{"test1", "test2"})
	if tag != "test1" {
		t.Error("Should pick first tag on start")
	}
	// waiting outbound first round test
	time.Sleep(time.Second * 1)
	tag = strategy.PickOutbound(obm, []string{"test1", "test2"})
	if tag != "test2" {
		t.Error("Should pick fastest tag")
	}
}
