package socks

import (
	"errors"
	"io"
  "log"
	"net"
  "strconv"

	"github.com/v2ray/v2ray-core"
	socksio "github.com/v2ray/v2ray-core/io/socks"
)

var (
	ErrorAuthenticationFailed = errors.New("None of the authentication methods is allowed.")
	ErrorCommandNotSupported  = errors.New("Client requested an unsupported command.")
)

// SocksServer is a SOCKS 5 proxy server
type SocksServer struct {
	accepting bool
	vPoint    *core.VPoint
}

func NewSocksServer(vp *core.VPoint) *SocksServer {
	server := new(SocksServer)
	server.vPoint = vp
	return server
}

func (server *SocksServer) Listen(port uint16) error {
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(int(port)))
	if err != nil {
		return err
	}
	server.accepting = true
	go server.AcceptConnections(listener)
	return nil
}

func (server *SocksServer) AcceptConnections(listener net.Listener) error {
	for server.accepting {
		connection, err := listener.Accept()
		if err != nil {
      log.Print(err)
			return err
		}
		go server.HandleConnection(connection)
	}
	return nil
}

func (server *SocksServer) HandleConnection(connection net.Conn) error {
	defer connection.Close()

	auth, err := socksio.ReadAuthentication(connection)
	if err != nil {
    log.Print(err)
		return err
	}
  log.Print(auth)

	if !auth.HasAuthMethod(socksio.AuthNotRequired) {
    // TODO send response with FF
    log.Print(ErrorAuthenticationFailed)
		return ErrorAuthenticationFailed
	}

	authResponse := socksio.NewAuthenticationResponse(socksio.AuthNotRequired)
	socksio.WriteAuthentication(connection, authResponse)

	request, err := socksio.ReadRequest(connection)
	if err != nil {
    log.Print(err)
		return err
	}
  
  response := socksio.NewSocks5Response()

	if request.Command == socksio.CmdBind || request.Command == socksio.CmdUdpAssociate {
		response := socksio.NewSocks5Response()
		response.Error = socksio.ErrorCommandNotSupported
		socksio.WriteResponse(connection, response)
    log.Print(ErrorCommandNotSupported)
		return ErrorCommandNotSupported
	}
  
  response.Error = socksio.ErrorSuccess
  response.Port = request.Port
  response.AddrType = request.AddrType
  switch response.AddrType {
    case socksio.AddrTypeIPv4:
    copy(response.IPv4[:], request.IPv4[:])
    case socksio.AddrTypeIPv6:
    copy(response.IPv6[:], request.IPv6[:])
    case socksio.AddrTypeDomain:
    response.Domain = request.Domain
  }
  socksio.WriteResponse(connection, response)
  

	ray := server.vPoint.NewInboundConnectionAccepted(request.Destination())
	input := ray.InboundInput()
	output := ray.InboundOutput()
	finish := make(chan bool, 2)

	go server.dumpInput(connection, input, finish)
	go server.dumpOutput(connection, output, finish)
	server.waitForFinish(finish)

	return nil
}

func (server *SocksServer) dumpInput(conn net.Conn, input chan<- []byte, finish chan<- bool) {
	for {
		buffer := make([]byte, 256)
		nBytes, err := conn.Read(buffer)
    log.Printf("Reading %d bytes, with error %v", nBytes, err)
		if err == io.EOF {
      close(input)
			finish <- true
			break
		}
		input <- buffer[:nBytes]
	}
}

func (server *SocksServer) dumpOutput(conn net.Conn, output <-chan []byte, finish chan<- bool) {
	for {
		buffer, open := <-output
		if !open {
			finish <- true
			break
		}
		nBytes, _ := conn.Write(buffer)
    log.Printf("Writing %d bytes", nBytes)
	}
}

func (server *SocksServer) waitForFinish(finish <-chan bool) {
	for i := 0; i < 2; i++ {
		<-finish
	}
}
