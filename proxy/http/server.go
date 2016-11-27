package http

import (
	"bufio"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"v2ray.com/core/app"
	"v2ray.com/core/app/dispatcher"
	"v2ray.com/core/common"
	v2io "v2ray.com/core/common/io"
	"v2ray.com/core/common/loader"
	"v2ray.com/core/common/log"
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/proxy"
	"v2ray.com/core/proxy/registry"
	"v2ray.com/core/transport/internet"
	"v2ray.com/core/transport/ray"
)

// Server is a HTTP proxy server.
type Server struct {
	sync.Mutex
	accepting        bool
	packetDispatcher dispatcher.PacketDispatcher
	config           *ServerConfig
	tcpListener      *internet.TCPHub
	meta             *proxy.InboundHandlerMeta
}

func NewServer(config *ServerConfig, packetDispatcher dispatcher.PacketDispatcher, meta *proxy.InboundHandlerMeta) *Server {
	return &Server{
		packetDispatcher: packetDispatcher,
		config:           config,
		meta:             meta,
	}
}

func (v *Server) Port() v2net.Port {
	return v.meta.Port
}

func (v *Server) Close() {
	v.accepting = false
	if v.tcpListener != nil {
		v.Lock()
		v.tcpListener.Close()
		v.tcpListener = nil
		v.Unlock()
	}
}

func (v *Server) Start() error {
	if v.accepting {
		return nil
	}

	tcpListener, err := internet.ListenTCP(v.meta.Address, v.meta.Port, v.handleConnection, v.meta.StreamSettings)
	if err != nil {
		log.Error("HTTP: Failed listen on ", v.meta.Address, ":", v.meta.Port, ": ", err)
		return err
	}
	v.Lock()
	v.tcpListener = tcpListener
	v.Unlock()
	v.accepting = true
	return nil
}

func parseHost(rawHost string, defaultPort v2net.Port) (v2net.Destination, error) {
	port := defaultPort
	host, rawPort, err := net.SplitHostPort(rawHost)
	if err != nil {
		if addrError, ok := err.(*net.AddrError); ok && strings.Contains(addrError.Err, "missing port") {
			host = rawHost
		} else {
			return v2net.Destination{}, err
		}
	} else {
		intPort, err := strconv.Atoi(rawPort)
		if err != nil {
			return v2net.Destination{}, err
		}
		port = v2net.Port(intPort)
	}

	if ip := net.ParseIP(host); ip != nil {
		return v2net.TCPDestination(v2net.IPAddress(ip), port), nil
	}
	return v2net.TCPDestination(v2net.DomainAddress(host), port), nil
}

func (v *Server) handleConnection(conn internet.Connection) {
	defer conn.Close()
	timedReader := v2net.NewTimeOutReader(v.config.Timeout, conn)
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
		Source:      v2net.DestinationFromAddr(conn.RemoteAddr()),
		Destination: dest,
		Inbound:     v.meta,
	}
	if strings.ToUpper(request.Method) == "CONNECT" {
		v.handleConnect(request, session, reader, conn)
	} else {
		v.handlePlainHTTP(request, session, reader, conn)
	}
}

func (v *Server) handleConnect(request *http.Request, session *proxy.SessionInfo, reader io.Reader, writer io.Writer) {
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

	ray := v.packetDispatcher.DispatchToOutbound(session)
	v.transport(reader, writer, ray)
}

func (v *Server) transport(input io.Reader, output io.Writer, ray ray.InboundRay) {
	var wg sync.WaitGroup
	wg.Add(2)
	defer wg.Wait()

	go func() {
		v2reader := v2io.NewAdaptiveReader(input)
		defer v2reader.Release()

		if err := v2io.PipeUntilEOF(v2reader, ray.InboundInput()); err != nil {
			log.Info("HTTP: Failed to transport all TCP request: ", err)
		}
		ray.InboundInput().Close()
		wg.Done()
	}()

	go func() {
		v2writer := v2io.NewAdaptiveWriter(output)
		defer v2writer.Release()

		if err := v2io.PipeUntilEOF(ray.InboundOutput(), v2writer); err != nil {
			log.Info("HTTP: Failed to transport all TCP response: ", err)
		}
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

func (v *Server) GenerateResponse(statusCode int, status string) *http.Response {
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

func (v *Server) handlePlainHTTP(request *http.Request, session *proxy.SessionInfo, reader *bufio.Reader, writer io.Writer) {
	if len(request.URL.Host) <= 0 {
		response := v.GenerateResponse(400, "Bad Request")
		response.Write(writer)

		return
	}

	request.Host = request.URL.Host
	StripHopByHopHeaders(request)

	ray := v.packetDispatcher.DispatchToOutbound(session)
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
			response = v.GenerateResponse(503, "Service Unavailable")
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

func (v *ServerFactory) StreamCapability() v2net.NetworkList {
	return v2net.NetworkList{
		Network: []v2net.Network{v2net.Network_RawTCP},
	}
}

func (v *ServerFactory) Create(space app.Space, rawConfig interface{}, meta *proxy.InboundHandlerMeta) (proxy.InboundHandler, error) {
	if !space.HasApp(dispatcher.APP_ID) {
		return nil, common.ErrBadConfiguration
	}
	return NewServer(
		rawConfig.(*ServerConfig),
		space.GetApp(dispatcher.APP_ID).(dispatcher.PacketDispatcher),
		meta), nil
}

func init() {
	registry.MustRegisterInboundHandlerCreator(loader.GetType(new(ServerConfig)), new(ServerFactory))
}
