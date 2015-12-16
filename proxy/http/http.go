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
	"github.com/v2ray/v2ray-core/transport/ray"
)

type HttpProxyServer struct {
	accepting bool
	space     app.Space
	config    Config
}

func NewHttpProxyServer(space app.Space, config Config) *HttpProxyServer {
	return &HttpProxyServer{
		space:  space,
		config: config,
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
	go this.accept(tcpListener)
	return nil
}

func (this *HttpProxyServer) accept(listener *net.TCPListener) {
	this.accepting = true
	for this.accepting {
		tcpConn, err := listener.AcceptTCP()
		if err != nil {
			log.Error("Failed to accept HTTP connection: %v", err)
			continue
		}
		go this.handleConnection(tcpConn)
	}
}

func parseHost(rawHost string, defaultPort v2net.Port) (v2net.Address, error) {
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
		return v2net.IPAddress(ip, port), nil
	}
	return v2net.DomainAddress(host, port), nil
}

func (this *HttpProxyServer) handleConnection(conn *net.TCPConn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for true {
		request, err := http.ReadRequest(reader)
		if err != nil {
			break
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
		address, err := parseHost(host, defaultPort)
		if err != nil {
			log.Warning("Malformed proxy host (%s): %v", host, err)
		}
		if strings.ToUpper(request.Method) == "CONNECT" {
			this.handleConnect(request, address, reader, conn)
		} else if len(request.URL.Host) > 0 {
			request.Host = request.URL.Host
			request.Header.Set("Connection", "keep-alive")
			request.Header.Del("Proxy-Connection")
			buffer := alloc.NewBuffer().Clear()
			request.Write(buffer)
			log.Info("Request to remote: %s", string(buffer.Value))
			packet := v2net.NewPacket(v2net.NewTCPDestination(address), buffer, true)
			ray := this.space.PacketDispatcher().DispatchToOutbound(packet)
			defer close(ray.InboundInput())

			responseReader := bufio.NewReader(NewChanReader(ray.InboundOutput()))
			response, err := http.ReadResponse(responseReader, request)
			if err != nil {
				return
			}

			responseBuffer := alloc.NewBuffer().Clear()
			response.Write(responseBuffer)
			conn.Write(responseBuffer.Value)
			responseBuffer.Release()
		} else {
			response := &http.Response{
				Status:        "400 Bad Request",
				StatusCode:    400,
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
			conn.Write(buffer.Value)
			buffer.Release()
		}
	}
}

func (this *HttpProxyServer) handleConnect(request *http.Request, address v2net.Address, reader io.Reader, writer io.Writer) {
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

	packet := v2net.NewPacket(v2net.NewTCPDestination(address), nil, true)
	ray := this.space.PacketDispatcher().DispatchToOutbound(packet)
	this.transport(reader, writer, ray)
}

func (this *HttpProxyServer) transport(input io.Reader, output io.Writer, ray ray.InboundRay) {
	var outputFinish sync.Mutex
	outputFinish.Lock()

	go func() {
		v2net.ReaderToChan(ray.InboundInput(), input)
		close(ray.InboundInput())
	}()

	go func() {
		v2net.ChanToWriter(output, ray.InboundOutput())
		outputFinish.Unlock()
	}()

	outputFinish.Lock()
}
