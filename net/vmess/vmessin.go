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
	clients   *core.UserSet
	accepting bool
}

func NewVMessInboundHandler(vp *core.Point, clients *core.UserSet) *VMessInboundHandler {
	handler := new(VMessInboundHandler)
	handler.vPoint = vp
	handler.clients = clients
	return handler
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
		return err
	}
	log.Debug("Received request for %s", request.Address.String())

	response := vmessio.NewVMessResponse(request)
	nBytes, err := connection.Write(response[:])
	log.Debug("Writing VMess response %v", response)
	if err != nil {
		return log.Error("Failed to write VMess response (%d bytes): %v", nBytes, err)
	}

	requestKey := request.RequestKey[:]
	requestIV := request.RequestIV[:]
	responseKey := md5.Sum(requestKey)
	responseIV := md5.Sum(requestIV)

	requestReader, err := v2io.NewAesDecryptReader(requestKey, requestIV, connection)
	if err != nil {
		return log.Error("Failed to create decrypt reader: %v", err)
	}

	responseWriter, err := v2io.NewAesEncryptWriter(responseKey[:], responseIV[:], connection)
	if err != nil {
		return log.Error("Failed to create encrypt writer: %v", err)
	}

	ray := handler.vPoint.NewInboundConnectionAccepted(request.Address)
	input := ray.InboundInput()
	output := ray.InboundOutput()
	finish := make(chan bool, 2)

	go handler.dumpInput(requestReader, input, finish)
	go handler.dumpOutput(responseWriter, output, finish)
	handler.waitForFinish(finish)

	return nil
}

func (handler *VMessInboundHandler) dumpInput(reader io.Reader, input chan<- []byte, finish chan<- bool) {
	v2net.ReaderToChan(input, reader)
	close(input)
	finish <- true
}

func (handler *VMessInboundHandler) dumpOutput(writer io.Writer, output <-chan []byte, finish chan<- bool) {
	v2net.ChanToWriter(writer, output)
	finish <- true
}

func (handler *VMessInboundHandler) waitForFinish(finish <-chan bool) {
	<-finish
	<-finish
}

type VMessInboundHandlerFactory struct {
}

func (factory *VMessInboundHandlerFactory) Create(vp *core.Point, rawConfig []byte) (core.InboundConnectionHandler, error) {
	config, err := loadInboundConfig(rawConfig)
	if err != nil {
		panic(log.Error("Failed to load VMess inbound config: %v", err))
	}
	allowedClients := core.NewUserSet()
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
