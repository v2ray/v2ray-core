package inbound

import (
	"crypto/md5"
	"io"
	"net"
	"sync"

	"github.com/v2ray/v2ray-core/app"
	"github.com/v2ray/v2ray-core/common/alloc"
	v2crypto "github.com/v2ray/v2ray-core/common/crypto"
	"github.com/v2ray/v2ray-core/common/log"
	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/common/retry"
	"github.com/v2ray/v2ray-core/proxy/common/connhandler"
	"github.com/v2ray/v2ray-core/proxy/vmess"
	"github.com/v2ray/v2ray-core/proxy/vmess/protocol"
	"github.com/v2ray/v2ray-core/proxy/vmess/protocol/user"
)

// Inbound connection handler that handles messages in VMess format.
type VMessInboundHandler struct {
	space     *app.Space
	clients   user.UserSet
	accepting bool
}

func NewVMessInboundHandler(space *app.Space, clients user.UserSet) *VMessInboundHandler {
	return &VMessInboundHandler{
		space:   space,
		clients: clients,
	}
}

func (this *VMessInboundHandler) Listen(port v2net.Port) error {
	listener, err := net.ListenTCP("tcp", &net.TCPAddr{
		IP:   []byte{0, 0, 0, 0},
		Port: int(port),
		Zone: "",
	})
	if err != nil {
		log.Error("Unable to listen tcp port %d: %v", port, err)
		return err
	}
	this.accepting = true
	go this.AcceptConnections(listener)
	return nil
}

func (this *VMessInboundHandler) AcceptConnections(listener *net.TCPListener) error {
	for this.accepting {
		retry.Timed(100 /* times */, 100 /* ms */).On(func() error {
			connection, err := listener.AcceptTCP()
			if err != nil {
				log.Error("Failed to accpet connection: %s", err.Error())
				return err
			}
			go this.HandleConnection(connection)
			return nil
		})

	}
	return nil
}

func (this *VMessInboundHandler) HandleConnection(connection *net.TCPConn) error {
	defer connection.Close()

	connReader := v2net.NewTimeOutReader(16, connection)
	requestReader := protocol.NewVMessRequestReader(this.clients)

	request, err := requestReader.Read(connReader)
	if err != nil {
		log.Access(connection.RemoteAddr().String(), "", log.AccessRejected, err.Error())
		log.Warning("VMessIn: Invalid request from (%s): %v", connection.RemoteAddr().String(), err)
		return err
	}
	log.Access(connection.RemoteAddr().String(), request.Address.String(), log.AccessAccepted, "")
	log.Debug("VMessIn: Received request for %s", request.Address.String())

	ray := this.space.PacketDispatcher().DispatchToOutbound(v2net.NewPacket(request.Destination(), nil, true))
	input := ray.InboundInput()
	output := ray.InboundOutput()
	var readFinish, writeFinish sync.Mutex
	readFinish.Lock()
	writeFinish.Lock()

	userSettings := vmess.GetUserSettings(request.User.Level())
	connReader.SetTimeOut(userSettings.PayloadReadTimeout)
	go handleInput(request, connReader, input, &readFinish)

	responseKey := md5.Sum(request.RequestKey)
	responseIV := md5.Sum(request.RequestIV)

	aesStream, err := v2crypto.NewAesEncryptionStream(responseKey[:], responseIV[:])
	if err != nil {
		log.Error("VMessIn: Failed to create AES decryption stream: %v", err)
		return err
	}

	responseWriter := v2crypto.NewCryptionWriter(aesStream, connection)

	// Optimize for small response packet
	buffer := alloc.NewLargeBuffer().Clear()
	buffer.AppendBytes(request.ResponseHeader[0] ^ request.ResponseHeader[1])
	buffer.AppendBytes(request.ResponseHeader[2] ^ request.ResponseHeader[3])
	buffer.AppendBytes(byte(0), byte(0))

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

	aesStream, err := v2crypto.NewAesDecryptionStream(request.RequestKey, request.RequestIV)
	if err != nil {
		log.Error("VMessIn: Failed to create AES decryption stream: %v", err)
		return
	}
	requestReader := v2crypto.NewCryptionReader(aesStream, reader)
	v2net.ReaderToChan(input, requestReader)
}

func handleOutput(request *protocol.VMessRequest, writer io.Writer, output <-chan *alloc.Buffer, finish *sync.Mutex) {
	v2net.ChanToWriter(writer, output)
	finish.Unlock()
}

type VMessInboundHandlerFactory struct {
}

func (this *VMessInboundHandlerFactory) Create(space *app.Space, rawConfig interface{}) (connhandler.InboundConnectionHandler, error) {
	config := rawConfig.(Config)

	allowedClients := user.NewTimedUserSet()
	for _, user := range config.AllowedUsers() {
		allowedClients.AddUser(user)
	}

	return NewVMessInboundHandler(space, allowedClients), nil
}

func init() {
	connhandler.RegisterInboundConnectionHandlerFactory("vmess", &VMessInboundHandlerFactory{})
}
