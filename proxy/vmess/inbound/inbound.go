package inbound

import (
	"io"
	"sync"

	"v2ray.com/core/app"
	"v2ray.com/core/app/dispatcher"
	"v2ray.com/core/app/proxyman"
	"v2ray.com/core/common"
	"v2ray.com/core/common/buf"
	"v2ray.com/core/common/bufio"
	"v2ray.com/core/common/errors"
	"v2ray.com/core/common/log"
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/common/protocol"
	"v2ray.com/core/common/serial"
	"v2ray.com/core/common/signal"
	"v2ray.com/core/common/uuid"
	"v2ray.com/core/proxy"
	"v2ray.com/core/proxy/vmess"
	"v2ray.com/core/proxy/vmess/encoding"
	"v2ray.com/core/transport/internet"
	"v2ray.com/core/transport/ray"
)

type userByEmail struct {
	sync.RWMutex
	cache           map[string]*protocol.User
	defaultLevel    uint32
	defaultAlterIDs uint16
}

func NewUserByEmail(users []*protocol.User, config *DefaultConfig) *userByEmail {
	cache := make(map[string]*protocol.User)
	for _, user := range users {
		cache[user.Email] = user
	}
	return &userByEmail{
		cache:           cache,
		defaultLevel:    config.Level,
		defaultAlterIDs: uint16(config.AlterId),
	}
}

func (v *userByEmail) Get(email string) (*protocol.User, bool) {
	var user *protocol.User
	var found bool
	v.RLock()
	user, found = v.cache[email]
	v.RUnlock()
	if !found {
		v.Lock()
		user, found = v.cache[email]
		if !found {
			account := &vmess.Account{
				Id:      uuid.New().String(),
				AlterId: uint32(v.defaultAlterIDs),
			}
			user = &protocol.User{
				Level:   v.defaultLevel,
				Email:   email,
				Account: serial.ToTypedMessage(account),
			}
			v.cache[email] = user
		}
		v.Unlock()
	}
	return user, found
}

// Inbound connection handler that handles messages in VMess format.
type VMessInboundHandler struct {
	sync.RWMutex
	packetDispatcher      dispatcher.PacketDispatcher
	inboundHandlerManager proxyman.InboundHandlerManager
	clients               protocol.UserValidator
	usersByEmail          *userByEmail
	accepting             bool
	listener              *internet.TCPHub
	detours               *DetourConfig
	meta                  *proxy.InboundHandlerMeta
}

func (v *VMessInboundHandler) Port() v2net.Port {
	return v.meta.Port
}

func (v *VMessInboundHandler) Close() {
	v.accepting = false
	if v.listener != nil {
		v.Lock()
		v.listener.Close()
		v.listener = nil
		v.clients.Release()
		v.clients = nil
		v.Unlock()
	}
}

func (v *VMessInboundHandler) GetUser(email string) *protocol.User {
	v.RLock()
	defer v.RUnlock()

	if !v.accepting {
		return nil
	}

	user, existing := v.usersByEmail.Get(email)
	if !existing {
		v.clients.Add(user)
	}
	return user
}

func (v *VMessInboundHandler) Start() error {
	if v.accepting {
		return nil
	}

	tcpListener, err := internet.ListenTCP(v.meta.Address, v.meta.Port, v.HandleConnection, v.meta.StreamSettings)
	if err != nil {
		log.Error("VMess|Inbound: Unable to listen tcp ", v.meta.Address, ":", v.meta.Port, ": ", err)
		return err
	}
	v.accepting = true
	v.Lock()
	v.listener = tcpListener
	v.Unlock()
	return nil
}

func transferRequest(session *encoding.ServerSession, request *protocol.RequestHeader, input io.Reader, output ray.OutputStream) error {
	defer output.Close()

	bodyReader := session.DecodeRequestBody(request, input)
	defer bodyReader.Release()

	if err := buf.PipeUntilEOF(bodyReader, output); err != nil {
		return err
	}
	return nil
}

func transferResponse(session *encoding.ServerSession, request *protocol.RequestHeader, response *protocol.ResponseHeader, input ray.InputStream, output io.Writer) error {
	defer input.Release()
	session.EncodeResponseHeader(response, output)

	bodyWriter := session.EncodeResponseBody(request, output)

	// Optimize for small response packet
	if data, err := input.Read(); err == nil {
		if err := bodyWriter.Write(data); err != nil {
			return err
		}

		if bufferedWriter, ok := output.(*bufio.BufferedWriter); ok {
			bufferedWriter.SetBuffered(false)
		}

		if err := buf.PipeUntilEOF(input, bodyWriter); err != nil {
			return err
		}
	}

	if request.Option.Has(protocol.RequestOptionChunkStream) {
		if err := bodyWriter.Write(buf.NewLocal(8)); err != nil {
			return err
		}
	}

	return nil
}

func (v *VMessInboundHandler) HandleConnection(connection internet.Connection) {
	defer connection.Close()

	if !v.accepting {
		return
	}

	connReader := v2net.NewTimeOutReader(8, connection)
	defer connReader.Release()

	reader := bufio.NewReader(connReader)
	defer reader.Release()

	v.RLock()
	if !v.accepting {
		v.RUnlock()
		return
	}
	session := encoding.NewServerSession(v.clients)
	defer session.Release()

	request, err := session.DecodeRequestHeader(reader)
	v.RUnlock()

	if err != nil {
		if errors.Cause(err) != io.EOF {
			log.Access(connection.RemoteAddr(), "", log.AccessRejected, err)
			log.Info("VMess|Inbound: Invalid request from ", connection.RemoteAddr(), ": ", err)
		}
		connection.SetReusable(false)
		return
	}
	log.Access(connection.RemoteAddr(), request.Destination(), log.AccessAccepted, "")
	log.Info("VMess|Inbound: Received request for ", request.Destination())

	connection.SetReusable(request.Option.Has(protocol.RequestOptionConnectionReuse))

	ray := v.packetDispatcher.DispatchToOutbound(&proxy.SessionInfo{
		Source:      v2net.DestinationFromAddr(connection.RemoteAddr()),
		Destination: request.Destination(),
		User:        request.User,
		Inbound:     v.meta,
	})
	input := ray.InboundInput()
	output := ray.InboundOutput()
	defer input.Close()
	defer output.Release()

	userSettings := request.User.GetSettings()
	connReader.SetTimeOut(userSettings.PayloadReadTimeout)
	reader.SetBuffered(false)

	requestDone := signal.ExecuteAsync(func() error {
		return transferRequest(session, request, reader, input)
	})

	writer := bufio.NewWriter(connection)
	defer writer.Release()

	response := &protocol.ResponseHeader{
		Command: v.generateCommand(request),
	}

	if connection.Reusable() {
		response.Option.Set(protocol.ResponseOptionConnectionReuse)
	}

	responseDone := signal.ExecuteAsync(func() error {
		return transferResponse(session, request, response, output, writer)
	})

	if err := signal.ErrorOrFinish2(requestDone, responseDone); err != nil {
		log.Info("VMess|Inbound: Connection ending with ", err)
		connection.SetReusable(false)
		return
	}

	if err := writer.Flush(); err != nil {
		log.Info("VMess|Inbound: Failed to flush remain data: ", err)
		connection.SetReusable(false)
		return
	}
}

type Factory struct{}

func (v *Factory) StreamCapability() v2net.NetworkList {
	return v2net.NetworkList{
		Network: []v2net.Network{v2net.Network_TCP, v2net.Network_KCP, v2net.Network_WebSocket},
	}
}

func (v *Factory) Create(space app.Space, rawConfig interface{}, meta *proxy.InboundHandlerMeta) (proxy.InboundHandler, error) {
	if !space.HasApp(dispatcher.APP_ID) {
		return nil, common.ErrBadConfiguration
	}
	config := rawConfig.(*Config)

	allowedClients := vmess.NewTimedUserValidator(protocol.DefaultIDHash)
	for _, user := range config.User {
		allowedClients.Add(user)
	}

	handler := &VMessInboundHandler{
		packetDispatcher: space.GetApp(dispatcher.APP_ID).(dispatcher.PacketDispatcher),
		clients:          allowedClients,
		detours:          config.Detour,
		usersByEmail:     NewUserByEmail(config.User, config.GetDefaultValue()),
		meta:             meta,
	}

	if space.HasApp(proxyman.APP_ID_INBOUND_MANAGER) {
		handler.inboundHandlerManager = space.GetApp(proxyman.APP_ID_INBOUND_MANAGER).(proxyman.InboundHandlerManager)
	}

	return handler, nil
}

func init() {
	common.Must(proxy.RegisterInboundHandlerCreator(serial.GetMessageType(new(Config)), new(Factory)))
}
