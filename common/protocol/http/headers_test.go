package http_test

import (
	"bufio"
	"net/http"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"

	"v2ray.com/core/common"
	"v2ray.com/core/common/net"
	. "v2ray.com/core/common/protocol/http"
)

func TestParseXForwardedFor(t *testing.T) {
	header := http.Header{}
	header.Add("X-Forwarded-For", "129.78.138.66, 129.78.64.103")
	addrs := ParseXForwardedFor(header)
	if r := cmp.Diff(addrs, []net.Address{net.ParseAddress("129.78.138.66"), net.ParseAddress("129.78.64.103")}); r != "" {
		t.Error(r)
	}
}

func TestHopByHopHeadersRemoving(t *testing.T) {
	rawRequest := `GET /pkg/net/http/ HTTP/1.1
Host: golang.org
Connection: keep-alive,Foo, Bar
Foo: foo
Bar: bar
Proxy-Connection: keep-alive
Proxy-Authenticate: abc
Accept-Encoding: gzip
Accept-Charset: ISO-8859-1,UTF-8;q=0.7,*;q=0.7
Cache-Control: no-cache
Accept-Language: de,en;q=0.7,en-us;q=0.3

`
	b := bufio.NewReader(strings.NewReader(rawRequest))
	req, err := http.ReadRequest(b)
	common.Must(err)
	headers := []struct {
		Key   string
		Value string
	}{
		{
			Key:   "Foo",
			Value: "foo",
		},
		{
			Key:   "Bar",
			Value: "bar",
		},
		{
			Key:   "Connection",
			Value: "keep-alive,Foo, Bar",
		},
		{
			Key:   "Proxy-Connection",
			Value: "keep-alive",
		},
		{
			Key:   "Proxy-Authenticate",
			Value: "abc",
		},
	}
	for _, header := range headers {
		if v := req.Header.Get(header.Key); v != header.Value {
			t.Error("header ", header.Key, " = ", v, " want ", header.Value)
		}
	}

	RemoveHopByHopHeaders(req.Header)

	for _, header := range []string{"Connection", "Foo", "Bar", "Proxy-Connection", "Proxy-Authenticate"} {
		if v := req.Header.Get(header); v != "" {
			t.Error("header ", header, " = ", v)
		}
	}
}

func TestParseHost(t *testing.T) {
	testCases := []struct {
		RawHost     string
		DefaultPort net.Port
		Destination net.Destination
		Error       bool
	}{
		{
			RawHost:     "v2ray.com:80",
			DefaultPort: 443,
			Destination: net.TCPDestination(net.DomainAddress("v2ray.com"), 80),
		},
		{
			RawHost:     "tls.v2ray.com",
			DefaultPort: 443,
			Destination: net.TCPDestination(net.DomainAddress("tls.v2ray.com"), 443),
		},
		{
			RawHost:     "[2401:1bc0:51f0:ec08::1]:80",
			DefaultPort: 443,
			Destination: net.TCPDestination(net.ParseAddress("[2401:1bc0:51f0:ec08::1]"), 80),
		},
	}

	for _, testCase := range testCases {
		dest, err := ParseHost(testCase.RawHost, testCase.DefaultPort)
		if testCase.Error {
			if err == nil {
				t.Error("for test case: ", testCase.RawHost, " expected error, but actually nil")
			}
		} else {
			if dest != testCase.Destination {
				t.Error("for test case: ", testCase.RawHost, " expected host: ", testCase.Destination.String(), " but got ", dest.String())
			}
		}
	}
}
