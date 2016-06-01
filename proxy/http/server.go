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

	"github.com/v2ray/v2ray-core/app"
	"github.com/v2ray/v2ray-core/app/dispatcher"
	"github.com/v2ray/v2ray-core/common/alloc"
	v2io "github.com/v2ray/v2ray-core/common/io"
	"github.com/v2ray/v2ray-core/common/log"
	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/proxy"
	"github.com/v2ray/v2ray-core/proxy/internal"
	"github.com/v2ray/v2ray-core/transport/hub"
	"github.com/v2ray/v2ray-core/transport/ray"
)

// Server is a HTTP proxy server.
type Server struct {
	sync.Mutex
	accepting        bool
	packetDispatcher dispatcher.PacketDispatcher
	config           *Config
	tcpListener      *hub.TCPHub
	listeningPort    v2net.Port
	listeningAddress v2net.Address
}

func NewServer(config *Config, packetDispatcher dispatcher.PacketDispatcher) *Server {
	return &Server{
		packetDispatcher: packetDispatcher,
		config:           config,
	}
}

func (this *Server) Port() v2net.Port {
	return this.listeningPort
}

func (this *Server) Close() {
	this.accepting = false
	if this.tcpListener != nil {
		this.Lock()
		this.tcpListener.Close()
		this.tcpListener = nil
		this.Unlock()
	}
}

func (this *Server) Listen(address v2net.Address, port v2net.Port) error {
	if this.accepting {
		if this.listeningPort == port && this.listeningAddress.Equals(address) {
			return nil
		} else {
			return proxy.ErrorAlreadyListening
		}
	}
	this.listeningPort = port
	this.listeningAddress = address

	var tlsConfig *tls.Config
	if this.config.TLSConfig != nil {
		tlsConfig = this.config.TLSConfig.GetConfig()
	}
	tcpListener, err := hub.ListenTCP(address, port, this.handleConnection, tlsConfig)
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

func (this *Server) handleConnection(conn *hub.Connection) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	request, err := http.ReadRequest(reader)
	if err != nil {
		log.Warning("HTTP: Failed to read http request: ", err)
		return
	}
	log.Info("HTTP: Request to Method [", request.Method, "] Host [", request.Host, "] with URL [", request.URL, "]")
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
		log.Warning("HTTP: Malformed proxy host (", host, "): ", err)
		return
	}
	if strings.ToUpper(request.Method) == "CONNECT" {
		this.handleConnect(request, dest, reader, conn)
	} else {
		this.handlePlainHTTP(request, dest, reader, conn)
	}
}

func (this *Server) handleConnect(request *http.Request, destination v2net.Destination, reader io.Reader, writer io.Writer) {
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

func (this *Server) transport(input io.Reader, output io.Writer, ray ray.InboundRay) {
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

func (this *Server) GenerateResponse(statusCode int, status string) *http.Response {
	hdr := http.Header(make(map[string][]string))
	hdr.Set("Connection", "close")
	return &http.Response{
		Status:        status,
		StatusCode:    statusCode,
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		Header:        hdr,
		Body:          nil,
		ContentLength: 0,
		Close:         false,
	}
}

func (this *Server) handlePlainHTTP(request *http.Request, dest v2net.Destination, reader *bufio.Reader, writer io.Writer) {
	if len(request.URL.Host) <= 0 {
		response := this.GenerateResponse(400, "Bad Request")

		buffer := alloc.NewSmallBuffer().Clear()
		response.Write(buffer)
		writer.Write(buffer.Value)
		buffer.Release()
		return
	}

	request.Host = request.URL.Host
	StripHopByHopHeaders(request)

	ray := this.packetDispatcher.DispatchToOutbound(dest)
	defer ray.InboundInput().Close()
	defer ray.InboundOutput().Release()

	var finish sync.WaitGroup
	finish.Add(1)
	go func() {
		defer finish.Done()
		requestWriter := v2io.NewBufferedWriter(v2io.NewChainWriter(ray.InboundInput()))
		err := request.Write(requestWriter)
		if err != nil {
			log.Warning("HTTP: Failed to write request: ", err)
			return
		}
		requestWriter.Flush()
	}()

	finish.Add(1)
	go func() {
		defer finish.Done()
		responseReader := bufio.NewReader(v2io.NewChanReader(ray.InboundOutput()))
		response, err := http.ReadResponse(responseReader, request)
		if err != nil {
			log.Warning("HTTP: Failed to read response: ", err)
			response = this.GenerateResponse(503, "Service Unavailable")
		}
		responseWriter := v2io.NewBufferedWriter(writer)
		err = response.Write(responseWriter)
		if err != nil {
			log.Warning("HTTP: Failed to write response: ", err)
			return
		}
		responseWriter.Flush()
	}()
	finish.Wait()
}

func init() {
	internal.MustRegisterInboundHandlerCreator("http",
		func(space app.Space, rawConfig interface{}) (proxy.InboundHandler, error) {
			if !space.HasApp(dispatcher.APP_ID) {
				return nil, internal.ErrorBadConfiguration
			}
			return NewServer(
				rawConfig.(*Config),
				space.GetApp(dispatcher.APP_ID).(dispatcher.PacketDispatcher)), nil
		})
}
