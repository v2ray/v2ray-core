package vmess

import (
	"crypto/md5"
	"io"
	"net"

	"github.com/v2ray/v2ray-core"
	v2io "github.com/v2ray/v2ray-core/io"
	vmessio "github.com/v2ray/v2ray-core/io/vmess"
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
	allowedClients *core.VUserSet
}

func NewVMessInboundHandlerFactory(clients []core.VUser) *VMessInboundHandlerFactory {
	factory := new(VMessInboundHandlerFactory)
	factory.allowedClients = core.NewVUserSet()
	for _, user := range clients {
		factory.allowedClients.AddUser(user)
	}
	return factory
}

func (factory *VMessInboundHandlerFactory) Create(vp *core.VPoint) *VMessInboundHandler {
	return NewVMessInboundHandler(vp, factory.allowedClients)
}
