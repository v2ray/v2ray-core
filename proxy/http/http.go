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
	"github.com/v2ray/v2ray-core/common/alloc"
	"github.com/v2ray/v2ray-core/common/log"
	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/common/retry"
	"github.com/v2ray/v2ray-core/transport/ray"
)

type HttpProxyServer struct {
	sync.Mutex
	accepting   bool
	space       app.Space
	config      Config
	tcpListener *net.TCPListener
}

func NewHttpProxyServer(space app.Space, config Config) *HttpProxyServer {
	return &HttpProxyServer{
		space:  space,
		config: config,
	}
}

func (this *HttpProxyServer) Close() {
	this.accepting = false
	if this.tcpListener != nil {
		this.tcpListener.Close()
		this.Lock()
		this.tcpListener = nil
		this.Unlock()
	}
}

func (this *HttpProxyServer) Listen(port v2net.Port) error {
	tcpListener, err := net.ListenTCP("tcp", &net.TCPAddr{
		Port: int(port.Value()),
		IP:   []byte{0, 0, 0, 0},
	})
	if err != nil {
		return err
	}
	this.tcpListener = tcpListener
	this.accepting = true
	go this.accept()
	return nil
}

func (this *HttpProxyServer) accept() {
	for this.accepting {
		retry.Timed(100 /* times */, 100 /* ms */).On(func() error {
			if !this.accepting {
				return nil
			}
			this.Lock()
			defer this.Unlock()
			if this.tcpListener != nil {
				tcpConn, err := this.tcpListener.AcceptTCP()
				if err != nil {
					log.Error("Failed to accept HTTP connection: %v", err)
					return err
				}
				go this.handleConnection(tcpConn)
			}
			return nil
		})
	}
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

func (this *HttpProxyServer) handleConnection(conn *net.TCPConn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	request, err := http.ReadRequest(reader)
	if err != nil {
		return
	}
	log.Info("Request to Method [%s] Host [%s] with URL [%s]", request.Method, request.Host, request.URL.String())
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
		log.Warning("Malformed proxy host (%s): %v", host, err)
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

	packet := v2net.NewPacket(destination, nil, true)
	ray := this.space.PacketDispatcher().DispatchToOutbound(packet)
	this.transport(reader, writer, ray)
}

func (this *HttpProxyServer) transport(input io.Reader, output io.Writer, ray ray.InboundRay) {
	var wg sync.WaitGroup
	wg.Add(2)
	defer wg.Wait()

	go func() {
		v2net.ReaderToChan(ray.InboundInput(), input)
		close(ray.InboundInput())
		wg.Done()
	}()

	go func() {
		v2net.ChanToWriter(output, ray.InboundOutput())
		wg.Done()
	}()
}

func stripHopByHopHeaders(request *http.Request) {
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
	stripHopByHopHeaders(request)

	requestBuffer := alloc.NewBuffer().Clear()
	request.Write(requestBuffer)
	log.Info("Request to remote:\n%s", string(requestBuffer.Value))

	packet := v2net.NewPacket(dest, requestBuffer, true)
	ray := this.space.PacketDispatcher().DispatchToOutbound(packet)
	defer close(ray.InboundInput())

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		responseReader := bufio.NewReader(NewChanReader(ray.InboundOutput()))
		responseBuffer := alloc.NewBuffer()
		defer responseBuffer.Release()
		response, err := http.ReadResponse(responseReader, request)
		if err != nil {
			return
		}
		responseBuffer.Clear()
		response.Write(responseBuffer)
		writer.Write(responseBuffer.Value)
		response.Body.Close()
	}()
	wg.Wait()
}
