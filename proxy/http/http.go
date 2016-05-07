package http

import (
	"bufio"
	"crypto/tls"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/v2ray/v2ray-core/app/dispatcher"
	"github.com/v2ray/v2ray-core/common/alloc"
	v2io "github.com/v2ray/v2ray-core/common/io"
	"github.com/v2ray/v2ray-core/common/log"
	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/common/serial"
	"github.com/v2ray/v2ray-core/proxy"
	"github.com/v2ray/v2ray-core/transport/hub"
	"github.com/v2ray/v2ray-core/transport/ray"
)

type HttpProxyServer struct {
	sync.Mutex
	accepting        bool
	packetDispatcher dispatcher.PacketDispatcher
	config           *Config
	tcpListener      *hub.TCPHub
	listeningPort    v2net.Port
}

func NewHttpProxyServer(config *Config, packetDispatcher dispatcher.PacketDispatcher) *HttpProxyServer {
	return &HttpProxyServer{
		packetDispatcher: packetDispatcher,
		config:           config,
	}
}

func (this *HttpProxyServer) Port() v2net.Port {
	return this.listeningPort
}

func (this *HttpProxyServer) Close() {
	this.accepting = false
	if this.tcpListener != nil {
		this.Lock()
		this.tcpListener.Close()
		this.tcpListener = nil
		this.Unlock()
	}
}

func (this *HttpProxyServer) Listen(port v2net.Port) error {
	if this.accepting {
		if this.listeningPort == port {
			return nil
		} else {
			return proxy.ErrorAlreadyListening
		}
	}
	this.listeningPort = port

	var tlsConfig *tls.Config = nil
	if this.config.TlsConfig != nil {
		tlsConfig = this.config.TlsConfig.GetConfig()
	}
	tcpListener, err := hub.ListenTCP(port, this.handleConnection, tlsConfig)
	if err != nil {
		log.Error("Http: Failed listen on port ", port, ": ", err)
		return err
	}
	this.Lock()
	this.tcpListener = tcpListener
	this.Unlock()
	this.accepting = true
	return nil
}

func parseHost(rawHost string, defaultPort v2net.Port) (v2net.Destination, error) {
	port := defaultPort
	host, rawPort, err := net.SplitHostPort(rawHost)
	if err != nil {
		if addrError, ok := err.(*net.AddrError); ok && strings.Contains(addrError.Err, "missing port") {
			host = rawHost
		} else {
			return nil, err
		}
	} else {
		intPort, err := strconv.Atoi(rawPort)
		if err != nil {
			return nil, err
		}
		port = v2net.Port(intPort)
	}

	if ip := net.ParseIP(host); ip != nil {
		return v2net.TCPDestination(v2net.IPAddress(ip), port), nil
	}
	return v2net.TCPDestination(v2net.DomainAddress(host), port), nil
}

func (this *HttpProxyServer) handleConnection(conn *hub.Connection) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	request, err := http.ReadRequest(reader)
	if err != nil {
		log.Warning("Failed to read http request: ", err)
		return
	}
	log.Info("Request to Method [", request.Method, "] Host [", request.Host, "] with URL [", request.URL, "]")
	defaultPort := v2net.Port(80)
	if strings.ToLower(request.URL.Scheme) == "https" {
		defaultPort = v2net.Port(443)
	}
	host := request.Host
	if len(host) == 0 {
		host = request.URL.Host
	}
	dest, err := parseHost(host, defaultPort)
	if err != nil {
		log.Warning("Malformed proxy host (", host, "): ", err)
		return
	}
	if strings.ToUpper(request.Method) == "CONNECT" {
		this.handleConnect(request, dest, reader, conn)
	} else {
		this.handlePlainHTTP(request, dest, reader, conn)
	}
}

func (this *HttpProxyServer) handleConnect(request *http.Request, destination v2net.Destination, reader io.Reader, writer io.Writer) {
	response := &http.Response{
		Status:        "200 OK",
		StatusCode:    200,
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		Header:        http.Header(make(map[string][]string)),
		Body:          nil,
		ContentLength: 0,
		Close:         false,
	}

	buffer := alloc.NewSmallBuffer().Clear()
	response.Write(buffer)
	writer.Write(buffer.Value)
	buffer.Release()

	ray := this.packetDispatcher.DispatchToOutbound(destination)
	this.transport(reader, writer, ray)
}

func (this *HttpProxyServer) transport(input io.Reader, output io.Writer, ray ray.InboundRay) {
	var wg sync.WaitGroup
	wg.Add(2)
	defer wg.Wait()

	go func() {
		v2reader := v2io.NewAdaptiveReader(input)
		defer v2reader.Release()

		v2io.Pipe(v2reader, ray.InboundInput())
		ray.InboundInput().Close()
		wg.Done()
	}()

	go func() {
		v2writer := v2io.NewAdaptiveWriter(output)
		defer v2writer.Release()

		v2io.Pipe(ray.InboundOutput(), v2writer)
		ray.InboundOutput().Release()
		wg.Done()
	}()
}

// @VisibleForTesting
func StripHopByHopHeaders(request *http.Request) {
	// Strip hop-by-hop header basaed on RFC:
	// http://www.w3.org/Protocols/rfc2616/rfc2616-sec13.html#sec13.5.1
	// https://www.mnot.net/blog/2011/07/11/what_proxies_must_do

	request.Header.Del("Proxy-Connection")
	request.Header.Del("Proxy-Authenticate")
	request.Header.Del("Proxy-Authorization")
	request.Header.Del("TE")
	request.Header.Del("Trailers")
	request.Header.Del("Transfer-Encoding")
	request.Header.Del("Upgrade")

	// TODO: support keep-alive
	connections := request.Header.Get("Connection")
	request.Header.Set("Connection", "close")
	if len(connections) == 0 {
		return
	}
	for _, h := range strings.Split(connections, ",") {
		request.Header.Del(strings.TrimSpace(h))
	}
}

func (this *HttpProxyServer) handlePlainHTTP(request *http.Request, dest v2net.Destination, reader *bufio.Reader, writer io.Writer) {
	if len(request.URL.Host) <= 0 {
		hdr := http.Header(make(map[string][]string))
		hdr.Set("Connection", "close")
		response := &http.Response{
			Status:        "400 Bad Request",
			StatusCode:    400,
			Proto:         "HTTP/1.1",
			ProtoMajor:    1,
			ProtoMinor:    1,
			Header:        hdr,
			Body:          nil,
			ContentLength: 0,
			Close:         false,
		}

		buffer := alloc.NewSmallBuffer().Clear()
		response.Write(buffer)
		writer.Write(buffer.Value)
		buffer.Release()
		return
	}

	request.Host = request.URL.Host
	StripHopByHopHeaders(request)

	requestBuffer := alloc.NewBuffer().Clear() // Don't release this buffer as it is passed into a Packet.
	request.Write(requestBuffer)
	log.Debug("Request to remote:\n", serial.BytesLiteral(requestBuffer.Value))

	ray := this.packetDispatcher.DispatchToOutbound(dest)
	ray.InboundInput().Write(requestBuffer)
	defer ray.InboundInput().Close()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		responseReader := bufio.NewReader(NewChanReader(ray.InboundOutput()))
		response, err := http.ReadResponse(responseReader, request)
		if err != nil {
			return
		}
		responseBuffer := alloc.NewBuffer().Clear()
		defer responseBuffer.Release()
		response.Write(responseBuffer)
		writer.Write(responseBuffer.Value)
		response.Body.Close()
	}()
	wg.Wait()
}
