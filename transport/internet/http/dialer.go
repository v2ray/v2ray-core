package http

import (
	"context"
	gotls "crypto/tls"
	"io"
	"net/http"
	"net/url"
	"sync"

	"golang.org/x/net/http2"

	"v2ray.com/core/common"
	"v2ray.com/core/common/net"
	"v2ray.com/core/transport/internet"
	"v2ray.com/core/transport/internet/tls"
)

var (
	globalDialerMap    = make(map[net.Destination]*http.Client)
	globalDailerAccess sync.Mutex
)

func getHTTPClient(ctx context.Context, dest net.Destination) (*http.Client, error) {
	globalDailerAccess.Lock()
	defer globalDailerAccess.Unlock()

	if client, found := globalDialerMap[dest]; found {
		return client, nil
	}

	config := tls.ConfigFromContext(ctx)
	if config == nil {
		return nil, newError("TLS must be enabled for http transport.").AtWarning()
	}

	transport := &http2.Transport{
		DialTLS: func(network string, addr string, tlsConfig *gotls.Config) (net.Conn, error) {
			rawHost, rawPort, err := net.SplitHostPort(addr)
			if err != nil {
				return nil, err
			}
			if len(rawPort) == 0 {
				rawPort = "443"
			}
			port, err := net.PortFromString(rawPort)
			if err != nil {
				return nil, err
			}
			address := net.ParseAddress(rawHost)

			pconn, err := internet.DialSystem(context.Background(), nil, net.TCPDestination(address, port))
			if err != nil {
				return nil, err
			}
			return gotls.Client(pconn, tlsConfig), nil
		},
		TLSClientConfig: config.GetTLSConfig(tls.WithDestination(dest), tls.WithNextProto("h2")),
	}

	client := &http.Client{
		Transport: transport,
	}

	globalDialerMap[dest] = client
	return client, nil
}

// Dial dials a new TCP connection to the given destination.
func Dial(ctx context.Context, dest net.Destination) (internet.Connection, error) {
	client, err := getHTTPClient(ctx, dest)
	if err != nil {
		return nil, err
	}

	preader, pwriter := io.Pipe()
	request := &http.Request{
		Method: "PUT",
		Host:   "www.v2ray.com",
		Body:   preader,
		URL: &url.URL{
			Scheme: "https",
			Host:   dest.NetAddr(),
			Path:   "/",
		},
		Proto:      "HTTP/2",
		ProtoMajor: 2,
		ProtoMinor: 0,
	}
	response, err := client.Do(request)
	if err != nil {
		return nil, newError("failed to dial to ", dest).Base(err).AtWarning()
	}
	if response.StatusCode != 200 {
		return nil, newError("unexpected status", response.StatusCode).AtWarning()
	}

	return &Connection{
		Reader: response.Body,
		Writer: pwriter,
		Closer: common.NewChainedClosable(preader, pwriter, response.Body),
		Local: &net.TCPAddr{
			IP:   []byte{0, 0, 0, 0},
			Port: 0,
		},
		Remote: &net.TCPAddr{
			IP:   []byte{0, 0, 0, 0},
			Port: 0,
		},
	}, nil
}

func init() {
	common.Must(internet.RegisterTransportDialer(internet.TransportProtocol_HTTP, Dial))
}
