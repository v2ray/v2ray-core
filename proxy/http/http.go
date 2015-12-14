package http

import (
	"bufio"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/v2ray/v2ray-core/app"
	"github.com/v2ray/v2ray-core/common/alloc"
	"github.com/v2ray/v2ray-core/common/log"
	v2net "github.com/v2ray/v2ray-core/common/net"
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

func parseHost(rawHost string) (v2net.Address, error) {
	port := v2net.Port(80)
	host, rawPort, err := net.SplitHostPort(rawHost)
	if err != nil {
		if addrError, ok := err.(*net.AddrError); ok && strings.Contains(addrError.Err, "missing port") {
			host = rawHost
			port = v2net.Port(80)
		} else {
			return nil, err
		}
	}
	intPort, err := strconv.Atoi(rawPort)
	if err != nil {
		return nil, err
	}
	port = v2net.Port(intPort)
	if ip := net.ParseIP(host); ip != nil {
		return v2net.IPAddress(ip, port), nil
	}
	return v2net.DomainAddress(host, port), nil
}

func (this *HttpProxyServer) handleConnection(conn *net.TCPConn) {
	httpReader := bufio.NewReader(conn)
	request, err := http.ReadRequest(httpReader)
	if err != nil {
		log.Warning("Malformed HTTP request: %v", err)
		return
	}
	if strings.ToUpper(request.Method) == "CONNECT" {
		address, err := parseHost(request.Host)
		if err != nil {
			log.Warning("Malformed proxy host: %v", err)
			return
		}
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
		conn.Write(buffer.Value)

		packet := v2net.NewPacket(v2net.NewTCPDestination(address), nil, true)
		this.space.PacketDispatcher().DispatchToOutbound(packet)
	}
}
