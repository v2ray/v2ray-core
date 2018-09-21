package http

import (
	"context"
	gotls "crypto/tls"
	"net/http"
	"net/url"
	"sync"

	"golang.org/x/net/http2"
	"v2ray.com/core/common"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/net"
	"v2ray.com/core/transport/internet"
	"v2ray.com/core/transport/internet/tls"
	"v2ray.com/core/transport/pipe"
)

var (
	globalDialerMap    map[net.Destination]*http.Client
	globalDailerAccess sync.Mutex
)

func getHTTPClient(ctx context.Context, dest net.Destination) (*http.Client, error) {
	globalDailerAccess.Lock()
	defer globalDailerAccess.Unlock()

	if globalDialerMap == nil {
		globalDialerMap = make(map[net.Destination]*http.Client)
	}

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

			pconn, err := internet.DialSystem(context.Background(), net.TCPDestination(address, port))
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
	rawSettings := internet.StreamSettingsFromContext(ctx)
	httpSettings, ok := rawSettings.ProtocolSettings.(*Config)
	if !ok {
		return nil, newError("HTTP config is not set.").AtError()
	}

	client, err := getHTTPClient(ctx, dest)
	if err != nil {
		return nil, err
	}

	opts := pipe.OptionsFromContext(ctx)
	preader, pwriter := pipe.New(opts...)
	breader := &buf.BufferedReader{Reader: preader}
	request := &http.Request{
		Method: "PUT",
		Host:   httpSettings.getRandomHost(),
		Body:   breader,
		URL: &url.URL{
			Scheme: "https",
			Host:   dest.NetAddr(),
			Path:   httpSettings.getNormalizedPath(),
		},
		Proto:      "HTTP/2",
		ProtoMajor: 2,
		ProtoMinor: 0,
		Header:     make(http.Header),
	}
	// Disable any compression method from server.
	request.Header.Set("Accept-Encoding", "identity")

	response, err := client.Do(request)
	if err != nil {
		return nil, newError("failed to dial to ", dest).Base(err).AtWarning()
	}
	if response.StatusCode != 200 {
		return nil, newError("unexpected status", response.StatusCode).AtWarning()
	}

	bwriter := buf.NewBufferedWriter(pwriter)
	common.Must(bwriter.SetBuffered(false))
	return net.NewConnection(
		net.ConnectionOutput(response.Body),
		net.ConnectionInput(bwriter),
		net.ConnectionOnClose(common.NewChainedClosable(breader, bwriter, response.Body)),
	), nil
}

func init() {
	common.Must(internet.RegisterTransportDialer(protocolName, Dial))
}
