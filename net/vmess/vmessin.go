package vmess

import (
	"crypto/md5"
	"io"
	"net"

	"github.com/v2ray/v2ray-core"
	v2io "github.com/v2ray/v2ray-core/io"
	vmessio "github.com/v2ray/v2ray-core/io/vmess"
	"github.com/v2ray/v2ray-core/log"
)

type VMessInboundHandler struct {
	vPoint    *core.VPoint
	clients   *core.VUserSet
	accepting bool
}

func NewVMessInboundHandler(vp *core.VPoint, clients *core.VUserSet) *VMessInboundHandler {
	handler := new(VMessInboundHandler)
	handler.vPoint = vp
	handler.clients = clients
	return handler
}

func (handler *VMessInboundHandler) Listen(port uint8) error {
	listener, err := net.Listen("tcp", ":"+string(port))
	if err != nil {
		return err
	}
	handler.accepting = true
	go handler.AcceptConnections(listener)

	return nil
}

func (handler *VMessInboundHandler) AcceptConnections(listener net.Listener) error {
	for handler.accepting {
		connection, err := listener.Accept()
		if err != nil {
			return err
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

	response := vmessio.NewVMessResponse(request)
	connection.Write(response[:])

	requestKey := request.RequestKey[:]
	requestIV := request.RequestIV[:]
	responseKey := md5.Sum(requestKey)
	responseIV := md5.Sum(requestIV)

	requestReader, err := v2io.NewAesDecryptReader(requestKey, requestIV, connection)
	if err != nil {
		return err
	}

	responseWriter, err := v2io.NewAesEncryptWriter(responseKey[:], responseIV[:], connection)
	if err != nil {
		return err
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
	for {
		buffer := make([]byte, BufferSize)
		nBytes, err := reader.Read(buffer)
		if err == io.EOF {
			close(input)
			finish <- true
			break
		}
		input <- buffer[:nBytes]
	}
}

func (handler *VMessInboundHandler) dumpOutput(writer io.Writer, output <-chan []byte, finish chan<- bool) {
	for {
		buffer, open := <-output
		if !open {
			finish <- true
			break
		}
		writer.Write(buffer)
	}
}

func (handler *VMessInboundHandler) waitForFinish(finish <-chan bool) {
	for i := 0; i < 2; i++ {
		<-finish
	}
}

type VMessInboundHandlerFactory struct {
}

func (factory *VMessInboundHandlerFactory) Create(vp *core.VPoint, rawConfig []byte) *VMessInboundHandler {
	config, err := loadInboundConfig(rawConfig)
	if err != nil {
		panic(log.Error("Failed to load VMess inbound config: %v", err))
	}
	allowedClients := core.NewVUserSet()
	for _, user := range config.AllowedClients {
		allowedClients.AddUser(user)
	}
	return NewVMessInboundHandler(vp, allowedClients)
}
