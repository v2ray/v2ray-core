package vmess

import (
	"crypto/md5"
	"io"
	"net"
	"sync"

	"github.com/v2ray/v2ray-core/app"
	"github.com/v2ray/v2ray-core/common/alloc"
	v2io "github.com/v2ray/v2ray-core/common/io"
	"github.com/v2ray/v2ray-core/common/log"
	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/common/retry"
	"github.com/v2ray/v2ray-core/proxy/common/connhandler"
	"github.com/v2ray/v2ray-core/proxy/vmess/config"
	"github.com/v2ray/v2ray-core/proxy/vmess/protocol"
	"github.com/v2ray/v2ray-core/proxy/vmess/protocol/user"
)

type VMessInboundHandler struct {
	dispatcher app.PacketDispatcher
	clients    user.UserSet
	accepting  bool
	udpEnabled bool
}

func NewVMessInboundHandler(dispatcher app.PacketDispatcher, clients user.UserSet, udpEnabled bool) *VMessInboundHandler {
	return &VMessInboundHandler{
		dispatcher: dispatcher,
		clients:    clients,
		udpEnabled: udpEnabled,
	}
}

func (handler *VMessInboundHandler) Listen(port uint16) error {
	listener, err := net.ListenTCP("tcp", &net.TCPAddr{
		IP:   []byte{0, 0, 0, 0},
		Port: int(port),
		Zone: "",
	})
	if err != nil {
		log.Error("Unable to listen tcp port %d: %v", port, err)
		return err
	}
	handler.accepting = true
	go handler.AcceptConnections(listener)

	if handler.udpEnabled {
		handler.ListenUDP(port)
	}

	return nil
}

func (handler *VMessInboundHandler) AcceptConnections(listener *net.TCPListener) error {
	for handler.accepting {
		retry.Timed(100 /* times */, 100 /* ms */).On(func() error {
			connection, err := listener.AcceptTCP()
			if err != nil {
				log.Error("Failed to accpet connection: %s", err.Error())
				return err
			}
			go handler.HandleConnection(connection)
			return nil
		})

	}
	return nil
}

func (handler *VMessInboundHandler) HandleConnection(connection *net.TCPConn) error {
	defer connection.Close()

	connReader := v2net.NewTimeOutReader(16, connection)
	requestReader := protocol.NewVMessRequestReader(handler.clients)

	request, err := requestReader.Read(connReader)
	if err != nil {
		log.Access(connection.RemoteAddr().String(), "", log.AccessRejected, err.Error())
		log.Warning("VMessIn: Invalid request from (%s): %v", connection.RemoteAddr().String(), err)
		return err
	}
	log.Access(connection.RemoteAddr().String(), request.Address.String(), log.AccessAccepted, "")
	log.Debug("VMessIn: Received request for %s", request.Address.String())

	ray := handler.dispatcher.DispatchToOutbound(v2net.NewPacket(request.Destination(), nil, true))
	input := ray.InboundInput()
	output := ray.InboundOutput()
	var readFinish, writeFinish sync.Mutex
	readFinish.Lock()
	writeFinish.Lock()

	connReader.SetTimeOut(120)
	go handleInput(request, connReader, input, &readFinish)

	responseKey := md5.Sum(request.RequestKey)
	responseIV := md5.Sum(request.RequestIV)

	responseWriter, err := v2io.NewAesEncryptWriter(responseKey[:], responseIV[:], connection)
	if err != nil {
		log.Error("VMessIn: Failed to create encrypt writer: %v", err)
		return err
	}

	// Optimize for small response packet
	buffer := alloc.NewLargeBuffer().Clear()
	buffer.Append(request.ResponseHeader)

	if data, open := <-output; open {
		buffer.Append(data.Value)
		data.Release()
		responseWriter.Write(buffer.Value)
		buffer.Release()
		go handleOutput(request, responseWriter, output, &writeFinish)
		writeFinish.Lock()
	}

	connection.CloseWrite()
	readFinish.Lock()

	return nil
}

func handleInput(request *protocol.VMessRequest, reader io.Reader, input chan<- *alloc.Buffer, finish *sync.Mutex) {
	defer close(input)
	defer finish.Unlock()

	requestReader, err := v2io.NewAesDecryptReader(request.RequestKey, request.RequestIV, reader)
	if err != nil {
		log.Error("VMessIn: Failed to create decrypt reader: %v", err)
		return
	}

	v2net.ReaderToChan(input, requestReader)
}

func handleOutput(request *protocol.VMessRequest, writer io.Writer, output <-chan *alloc.Buffer, finish *sync.Mutex) {
	v2net.ChanToWriter(writer, output)
	finish.Unlock()
}

type VMessInboundHandlerFactory struct {
}

func (factory *VMessInboundHandlerFactory) Create(dispatcher app.PacketDispatcher, rawConfig interface{}) (connhandler.InboundConnectionHandler, error) {
	config := rawConfig.(config.Inbound)

	allowedClients := user.NewTimedUserSet()
	for _, user := range config.AllowedUsers() {
		allowedClients.AddUser(user)
	}

	return NewVMessInboundHandler(dispatcher, allowedClients, config.UDPEnabled()), nil
}

func init() {
	connhandler.RegisterInboundConnectionHandlerFactory("vmess", &VMessInboundHandlerFactory{})
}
