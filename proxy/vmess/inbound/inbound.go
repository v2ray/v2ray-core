package inbound

import (
	"io"
	"sync"

	"v2ray.com/core/app"
	"v2ray.com/core/app/dispatcher"
	"v2ray.com/core/app/proxyman"
	"v2ray.com/core/common"
	"v2ray.com/core/common/alloc"
	"v2ray.com/core/common/errors"
	v2io "v2ray.com/core/common/io"
	"v2ray.com/core/common/loader"
	"v2ray.com/core/common/log"
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/common/protocol"
	"v2ray.com/core/common/uuid"
	"v2ray.com/core/proxy"
	"v2ray.com/core/proxy/registry"
	"v2ray.com/core/proxy/vmess"
	"v2ray.com/core/proxy/vmess/encoding"
	vmessio "v2ray.com/core/proxy/vmess/io"
	"v2ray.com/core/transport/internet"
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
				Account: loader.NewTypedSettings(account),
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

func (v *VMessInboundHandler) HandleConnection(connection internet.Connection) {
	defer connection.Close()

	if !v.accepting {
		return
	}

	connReader := v2net.NewTimeOutReader(8, connection)
	defer connReader.Release()

	reader := v2io.NewBufferedReader(connReader)
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
			log.Info("VMessIn: Invalid request from ", connection.RemoteAddr(), ": ", err)
		}
		connection.SetReusable(false)
		return
	}
	log.Access(connection.RemoteAddr(), request.Destination(), log.AccessAccepted, "")
	log.Info("VMessIn: Received request for ", request.Destination())

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

	var readFinish sync.Mutex
	readFinish.Lock()

	userSettings := request.User.GetSettings()
	connReader.SetTimeOut(userSettings.PayloadReadTimeout)
	reader.SetCached(false)

	go func() {
		bodyReader := session.DecodeRequestBody(reader)
		var requestReader v2io.Reader
		if request.Option.Has(protocol.RequestOptionChunkStream) {
			requestReader = vmessio.NewAuthChunkReader(bodyReader)
		} else {
			requestReader = v2io.NewAdaptiveReader(bodyReader)
		}
		if err := v2io.PipeUntilEOF(requestReader, input); err != nil {
			connection.SetReusable(false)
		}

		requestReader.Release()
		input.Close()
		readFinish.Unlock()
	}()

	writer := v2io.NewBufferedWriter(connection)
	defer writer.Release()

	response := &protocol.ResponseHeader{
		Command: v.generateCommand(request),
	}

	if connection.Reusable() {
		response.Option.Set(protocol.ResponseOptionConnectionReuse)
	}

	session.EncodeResponseHeader(response, writer)

	bodyWriter := session.EncodeResponseBody(writer)
	var v2writer v2io.Writer = v2io.NewAdaptiveWriter(bodyWriter)
	if request.Option.Has(protocol.RequestOptionChunkStream) {
		v2writer = vmessio.NewAuthChunkWriter(v2writer)
	}

	// Optimize for small response packet
	if data, err := output.Read(); err == nil {
		if err := v2writer.Write(data); err != nil {
			connection.SetReusable(false)
		}

		writer.SetCached(false)

		if err := v2io.PipeUntilEOF(output, v2writer); err != nil {
			connection.SetReusable(false)
		}

	}
	output.Release()
	if request.Option.Has(protocol.RequestOptionChunkStream) {
		if err := v2writer.Write(alloc.NewLocalBuffer(32).Clear()); err != nil {
			connection.SetReusable(false)
		}
	}
	writer.Flush()
	v2writer.Release()

	readFinish.Lock()
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
	registry.MustRegisterInboundHandlerCreator(loader.GetType(new(Config)), new(Factory))
}
