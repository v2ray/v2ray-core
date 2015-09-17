package vmess

import (
	"crypto/md5"
	"io"
	"net"
	"strconv"

	"github.com/v2ray/v2ray-core"
	v2io "github.com/v2ray/v2ray-core/io"
	vmessio "github.com/v2ray/v2ray-core/io/vmess"
	"github.com/v2ray/v2ray-core/log"
	v2net "github.com/v2ray/v2ray-core/net"
)

type VMessInboundHandler struct {
	vPoint    *core.Point
	clients   core.UserSet
	accepting bool
}

func NewVMessInboundHandler(vp *core.Point, clients core.UserSet) *VMessInboundHandler {
	return &VMessInboundHandler{
		vPoint:  vp,
		clients: clients,
	}
}

func (handler *VMessInboundHandler) Listen(port uint16) error {
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(int(port)))
	if err != nil {
		return log.Error("Unable to listen tcp:%d", port)
	}
	handler.accepting = true
	go handler.AcceptConnections(listener)

	return nil
}

func (handler *VMessInboundHandler) AcceptConnections(listener net.Listener) error {
	for handler.accepting {
		connection, err := listener.Accept()
		if err != nil {
			return log.Error("Failed to accpet connection: %s", err.Error())
		}
		go handler.HandleConnection(connection)
	}
	return nil
}

func (handler *VMessInboundHandler) HandleConnection(connection net.Conn) error {
	defer connection.Close()

	reader := vmessio.NewVMessRequestReader(handler.clients)

	request, err := reader.Read(connection)
	if err != nil {
		log.Debug("Failed to parse VMess request: %v", err)
		return err
	}
	log.Debug("Received request for %s", request.Address.String())

	ray := handler.vPoint.NewInboundConnectionAccepted(request.Address)
	input := ray.InboundInput()
	output := ray.InboundOutput()

	readFinish := make(chan bool)
	writeFinish := make(chan bool)

	go handleInput(request, connection, input, readFinish)

	responseKey := md5.Sum(request.RequestKey[:])
	responseIV := md5.Sum(request.RequestIV[:])

	response := vmessio.NewVMessResponse(request)
	responseWriter, err := v2io.NewAesEncryptWriter(responseKey[:], responseIV[:], connection)
	if err != nil {
		return log.Error("Failed to create encrypt writer: %v", err)
	}
	//responseWriter.Write(response[:])

	// Optimize for small response packet
	buffer := make([]byte, 0, 1024)
	buffer = append(buffer, response[:]...)
	data, open := <-output
	if open {
		buffer = append(buffer, data...)
	}
	responseWriter.Write(buffer)

	if open {
		go handleOutput(request, responseWriter, output, writeFinish)
	} else {
		close(writeFinish)
	}

	<-writeFinish
	if tcpConn, ok := connection.(*net.TCPConn); ok {
		log.Debug("VMessIn closing write")
		tcpConn.CloseWrite()
	}
	<-readFinish

	return nil
}

func handleInput(request *vmessio.VMessRequest, reader io.Reader, input chan<- []byte, finish chan<- bool) {
	defer close(input)
	defer close(finish)

	requestReader, err := v2io.NewAesDecryptReader(request.RequestKey[:], request.RequestIV[:], reader)
	if err != nil {
		log.Error("Failed to create decrypt reader: %v", err)
		return
	}

	v2net.ReaderToChan(input, requestReader)
}

func handleOutput(request *vmessio.VMessRequest, writer io.Writer, output <-chan []byte, finish chan<- bool) {
	v2net.ChanToWriter(writer, output)
	close(finish)
}

type VMessInboundHandlerFactory struct {
}

func (factory *VMessInboundHandlerFactory) Create(vp *core.Point, rawConfig []byte) (core.InboundConnectionHandler, error) {
	config, err := loadInboundConfig(rawConfig)
	if err != nil {
		panic(log.Error("Failed to load VMess inbound config: %v", err))
	}
	allowedClients := core.NewTimedUserSet()
	for _, client := range config.AllowedClients {
		user, err := client.ToUser()
		if err != nil {
			panic(log.Error("Failed to parse user id %s: %v", client.Id, err))
		}
		allowedClients.AddUser(user)
	}
	return NewVMessInboundHandler(vp, allowedClients), nil
}

func init() {
	core.RegisterInboundConnectionHandlerFactory("vmess", &VMessInboundHandlerFactory{})
}
