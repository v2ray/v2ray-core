package http

import (
	"bufio"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/v2ray/v2ray-core/app"
	"github.com/v2ray/v2ray-core/app/dispatcher"
	v2io "github.com/v2ray/v2ray-core/common/io"
	"github.com/v2ray/v2ray-core/common/log"
	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/proxy"
	"github.com/v2ray/v2ray-core/proxy/internal"
	"github.com/v2ray/v2ray-core/transport/internet"
	"github.com/v2ray/v2ray-core/transport/ray"
)

// Server is a HTTP proxy server.
type Server struct {
	sync.Mutex
	accepting        bool
	packetDispatcher dispatcher.PacketDispatcher
	config           *Config
	tcpListener      *internet.TCPHub
	meta             *proxy.InboundHandlerMeta
}

func NewServer(config *Config, packetDispatcher dispatcher.PacketDispatcher, meta *proxy.InboundHandlerMeta) *Server {
	return &Server{
		packetDispatcher: packetDispatcher,
		config:           config,
		meta:             meta,
	}
}

func (this *Server) Port() v2net.Port {
	return this.meta.Port
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

func (this *Server) Start() error {
	if this.accepting {
		return nil
	}

	tcpListener, err := internet.ListenTCP(this.meta.Address, this.meta.Port, this.handleConnection, this.meta.StreamSettings)
	if err != nil {
		log.Error("HTTP: Failed listen on ", this.meta.Address, ":", this.meta.Port, ": ", err)
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

func (this *Server) handleConnection(conn internet.Connection) {
	defer conn.Close()
	timedReader := v2net.NewTimeOutReader(this.config.Timeout, conn)
	reader := bufio.NewReaderSize(timedReader, 2048)

	request, err := http.ReadRequest(reader)
	if err != nil {
		if err != io.EOF {
			log.Warning("HTTP: Failed to read http request: ", err)
		}
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
	log.Access(conn.RemoteAddr(), request.URL, log.AccessAccepted, "")
	session := &proxy.SessionInfo{
		Source:      v2net.TCPDestinationFromAddr(conn.RemoteAddr().(*net.TCPAddr)),
		Destination: dest,
	}
	if strings.ToUpper(request.Method) == "CONNECT" {
		this.handleConnect(request, session, reader, conn)
	} else {
		this.handlePlainHTTP(request, session, reader, conn)
	}
}

func (this *Server) handleConnect(request *http.Request, session *proxy.SessionInfo, reader io.Reader, writer io.Writer) {
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
	response.Write(writer)

	ray := this.packetDispatcher.DispatchToOutbound(this.meta, session)
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

func (this *Server) handlePlainHTTP(request *http.Request, session *proxy.SessionInfo, reader *bufio.Reader, writer io.Writer) {
	if len(request.URL.Host) <= 0 {
		response := this.GenerateResponse(400, "Bad Request")
		response.Write(writer)

		return
	}

	request.Host = request.URL.Host
	StripHopByHopHeaders(request)

	ray := this.packetDispatcher.DispatchToOutbound(this.meta, session)
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

type ServerFactory struct{}

func (this *ServerFactory) StreamCapability() internet.StreamConnectionType {
	return internet.StreamConnectionTypeRawTCP
}

func (this *ServerFactory) Create(space app.Space, rawConfig interface{}, meta *proxy.InboundHandlerMeta) (proxy.InboundHandler, error) {
	if !space.HasApp(dispatcher.APP_ID) {
		return nil, internal.ErrBadConfiguration
	}
	return NewServer(
		rawConfig.(*Config),
		space.GetApp(dispatcher.APP_ID).(dispatcher.PacketDispatcher),
		meta), nil
}

func init() {
	internal.MustRegisterInboundHandlerCreator("http", new(ServerFactory))
}
