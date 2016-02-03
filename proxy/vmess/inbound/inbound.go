package inbound

import (
	"crypto/md5"
	"io"
	"sync"

	"github.com/v2ray/v2ray-core/app"
	"github.com/v2ray/v2ray-core/app/dispatcher"
	"github.com/v2ray/v2ray-core/app/proxyman"
	"github.com/v2ray/v2ray-core/common/alloc"
	v2crypto "github.com/v2ray/v2ray-core/common/crypto"
	v2io "github.com/v2ray/v2ray-core/common/io"
	"github.com/v2ray/v2ray-core/common/log"
	v2net "github.com/v2ray/v2ray-core/common/net"
	proto "github.com/v2ray/v2ray-core/common/protocol"
	"github.com/v2ray/v2ray-core/common/serial"
	"github.com/v2ray/v2ray-core/proxy"
	"github.com/v2ray/v2ray-core/proxy/internal"
	vmessio "github.com/v2ray/v2ray-core/proxy/vmess/io"
	"github.com/v2ray/v2ray-core/proxy/vmess/protocol"
	"github.com/v2ray/v2ray-core/transport/hub"
)

// Inbound connection handler that handles messages in VMess format.
type VMessInboundHandler struct {
	sync.Mutex
	packetDispatcher      dispatcher.PacketDispatcher
	inboundHandlerManager proxyman.InboundHandlerManager
	clients               protocol.UserSet
	user                  *proto.User
	accepting             bool
	listener              *hub.TCPHub
	features              *FeaturesConfig
	listeningPort         v2net.Port
}

func (this *VMessInboundHandler) Port() v2net.Port {
	return this.listeningPort
}

func (this *VMessInboundHandler) Close() {
	this.accepting = false
	if this.listener != nil {
		this.Lock()
		this.listener.Close()
		this.listener = nil
		this.Unlock()
	}
}

func (this *VMessInboundHandler) GetUser() *proto.User {
	return this.user
}

func (this *VMessInboundHandler) Listen(port v2net.Port) error {
	if this.accepting {
		if this.listeningPort == port {
			return nil
		} else {
			return proxy.ErrorAlreadyListening
		}
	}
	this.listeningPort = port

	tcpListener, err := hub.ListenTCP(port, this.HandleConnection)
	if err != nil {
		log.Error("Unable to listen tcp port ", port, ": ", err)
		return err
	}
	this.accepting = true
	this.Lock()
	this.listener = tcpListener
	this.Unlock()
	return nil
}

func (this *VMessInboundHandler) HandleConnection(connection *hub.TCPConn) {
	defer connection.Close()

	connReader := v2net.NewTimeOutReader(16, connection)
	requestReader := protocol.NewVMessRequestReader(this.clients)

	request, err := requestReader.Read(connReader)
	if err != nil {
		log.Access(connection.RemoteAddr(), serial.StringLiteral(""), log.AccessRejected, serial.StringLiteral(err.Error()))
		log.Warning("VMessIn: Invalid request from ", connection.RemoteAddr(), ": ", err)
		return
	}
	log.Access(connection.RemoteAddr(), request.Address, log.AccessAccepted, serial.StringLiteral(""))
	log.Debug("VMessIn: Received request for ", request.Address)

	ray := this.packetDispatcher.DispatchToOutbound(v2net.NewPacket(request.Destination(), nil, true))
	input := ray.InboundInput()
	output := ray.InboundOutput()
	var readFinish, writeFinish sync.Mutex
	readFinish.Lock()
	writeFinish.Lock()

	userSettings := proto.GetUserSettings(request.User.Level)
	connReader.SetTimeOut(userSettings.PayloadReadTimeout)
	go handleInput(request, connReader, input, &readFinish)

	responseKey := md5.Sum(request.RequestKey)
	responseIV := md5.Sum(request.RequestIV)

	aesStream, err := v2crypto.NewAesEncryptionStream(responseKey[:], responseIV[:])
	if err != nil {
		log.Error("VMessIn: Failed to create AES decryption stream: ", err)
		close(input)
		return
	}

	responseWriter := v2crypto.NewCryptionWriter(aesStream, connection)

	// Optimize for small response packet
	buffer := alloc.NewLargeBuffer().Clear()
	defer buffer.Release()
	buffer.AppendBytes(request.ResponseHeader, byte(0))
	this.generateCommand(buffer)

	if data, open := <-output; open {
		if request.IsChunkStream() {
			vmessio.Authenticate(data)
		}
		buffer.Append(data.Value)
		data.Release()
		responseWriter.Write(buffer.Value)
		go func(finish *sync.Mutex) {
			var writer v2io.Writer = v2io.NewAdaptiveWriter(responseWriter)
			if request.IsChunkStream() {
				writer = vmessio.NewAuthChunkWriter(writer)
			}
			v2io.ChanToWriter(writer, output)
			finish.Unlock()
		}(&writeFinish)
		writeFinish.Lock()
	}

	connection.CloseWrite()
	readFinish.Lock()
}

func handleInput(request *protocol.VMessRequest, reader io.Reader, input chan<- *alloc.Buffer, finish *sync.Mutex) {
	defer close(input)
	defer finish.Unlock()

	aesStream, err := v2crypto.NewAesDecryptionStream(request.RequestKey, request.RequestIV)
	if err != nil {
		log.Error("VMessIn: Failed to create AES decryption stream: ", err)
		return
	}
	descriptionReader := v2crypto.NewCryptionReader(aesStream, reader)
	var requestReader v2io.Reader
	if request.IsChunkStream() {
		requestReader = vmessio.NewAuthChunkReader(descriptionReader)
	} else {
		requestReader = v2io.NewAdaptiveReader(descriptionReader)
	}
	v2io.ReaderToChan(input, requestReader)
}

func init() {
	internal.MustRegisterInboundHandlerCreator("vmess",
		func(space app.Space, rawConfig interface{}) (proxy.InboundHandler, error) {
			if !space.HasApp(dispatcher.APP_ID) {
				return nil, internal.ErrorBadConfiguration
			}
			config := rawConfig.(*Config)

			allowedClients := protocol.NewTimedUserSet()
			for _, user := range config.AllowedUsers {
				allowedClients.AddUser(user)
			}

			handler := &VMessInboundHandler{
				packetDispatcher: space.GetApp(dispatcher.APP_ID).(dispatcher.PacketDispatcher),
				clients:          allowedClients,
				features:         config.Features,
				user:             config.AllowedUsers[0],
			}

			if space.HasApp(proxyman.APP_ID_INBOUND_MANAGER) {
				handler.inboundHandlerManager = space.GetApp(proxyman.APP_ID_INBOUND_MANAGER).(proxyman.InboundHandlerManager)
			}

			return handler, nil
		})
}
