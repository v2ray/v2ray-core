package inbound

import (
	"io"
	"sync"

	"github.com/v2ray/v2ray-core/app"
	"github.com/v2ray/v2ray-core/app/dispatcher"
	"github.com/v2ray/v2ray-core/app/proxyman"
	"github.com/v2ray/v2ray-core/common/alloc"
	v2io "github.com/v2ray/v2ray-core/common/io"
	"github.com/v2ray/v2ray-core/common/log"
	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/common/protocol"
	"github.com/v2ray/v2ray-core/common/protocol/raw"
	"github.com/v2ray/v2ray-core/common/uuid"
	"github.com/v2ray/v2ray-core/proxy"
	"github.com/v2ray/v2ray-core/proxy/internal"
	vmessio "github.com/v2ray/v2ray-core/proxy/vmess/io"
	"github.com/v2ray/v2ray-core/transport"
	"github.com/v2ray/v2ray-core/transport/hub"
)

type userByEmail struct {
	sync.RWMutex
	cache           map[string]*protocol.User
	defaultLevel    protocol.UserLevel
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
		defaultAlterIDs: config.AlterIDs,
	}
}

func (this *userByEmail) Get(email string) (*protocol.User, bool) {
	var user *protocol.User
	var found bool
	this.RLock()
	user, found = this.cache[email]
	this.RUnlock()
	if !found {
		this.Lock()
		user, found = this.cache[email]
		if !found {
			id := protocol.NewID(uuid.New())
			alterIDs := protocol.NewAlterIDs(id, this.defaultAlterIDs)
			account := &protocol.VMessAccount{
				ID:       id,
				AlterIDs: alterIDs,
			}
			user = protocol.NewUser(account, this.defaultLevel, email)
			this.cache[email] = user
		}
		this.Unlock()
	}
	return user, found
}

// Inbound connection handler that handles messages in VMess format.
type VMessInboundHandler struct {
	sync.Mutex
	packetDispatcher      dispatcher.PacketDispatcher
	inboundHandlerManager proxyman.InboundHandlerManager
	clients               protocol.UserValidator
	usersByEmail          *userByEmail
	accepting             bool
	listener              *hub.TCPHub
	detours               *DetourConfig
	meta                  *proxy.InboundHandlerMeta
}

func (this *VMessInboundHandler) Port() v2net.Port {
	return this.meta.Port
}

func (this *VMessInboundHandler) Close() {
	this.accepting = false
	if this.listener != nil {
		this.Lock()
		this.listener.Close()
		this.listener = nil
		this.clients.Release()
		this.clients = nil
		this.Unlock()
	}
}

func (this *VMessInboundHandler) GetUser(email string) *protocol.User {
	user, existing := this.usersByEmail.Get(email)
	if !existing {
		this.clients.Add(user)
	}
	return user
}

func (this *VMessInboundHandler) Start() error {
	if this.accepting {
		return nil
	}

	tcpListener, err := hub.ListenTCP(this.meta.Address, this.meta.Port, this.HandleConnection, nil)
	if err != nil {
		log.Error("Unable to listen tcp ", this.meta.Address, ":", this.meta.Port, ": ", err)
		return err
	}
	this.accepting = true
	this.Lock()
	this.listener = tcpListener
	this.Unlock()
	return nil
}

func (this *VMessInboundHandler) HandleConnection(connection *hub.Connection) {
	defer connection.Close()

	connReader := v2net.NewTimeOutReader(8, connection)
	defer connReader.Release()

	reader := v2io.NewBufferedReader(connReader)
	defer reader.Release()

	session := raw.NewServerSession(this.clients)
	defer session.Release()

	request, err := session.DecodeRequestHeader(reader)
	if err != nil {
		if err != io.EOF {
			log.Access(connection.RemoteAddr(), "", log.AccessRejected, err)
			log.Warning("VMessIn: Invalid request from ", connection.RemoteAddr(), ": ", err)
		}
		connection.SetReusable(false)
		return
	}
	log.Access(connection.RemoteAddr(), request.Destination(), log.AccessAccepted, "")
	log.Debug("VMessIn: Received request for ", request.Destination())

	if request.Option.Has(protocol.RequestOptionConnectionReuse) {
		connection.SetReusable(true)
	}

	ray := this.packetDispatcher.DispatchToOutbound(request.Destination())
	input := ray.InboundInput()
	output := ray.InboundOutput()
	defer input.Close()
	defer output.Release()

	var readFinish sync.Mutex
	readFinish.Lock()

	userSettings := protocol.GetUserSettings(request.User.Level)
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
		err := v2io.Pipe(requestReader, input)
		if err != io.EOF {
			connection.SetReusable(false)
		}

		requestReader.Release()
		input.Close()
		readFinish.Unlock()
	}()

	writer := v2io.NewBufferedWriter(connection)
	defer writer.Release()

	response := &protocol.ResponseHeader{
		Command: this.generateCommand(request),
	}

	if request.Option.Has(protocol.RequestOptionConnectionReuse) && transport.IsConnectionReusable() {
		response.Option.Set(protocol.ResponseOptionConnectionReuse)
	}

	session.EncodeResponseHeader(response, writer)

	bodyWriter := session.EncodeResponseBody(writer)

	// Optimize for small response packet
	if data, err := output.Read(); err == nil {
		var v2writer v2io.Writer = v2io.NewAdaptiveWriter(bodyWriter)
		if request.Option.Has(protocol.RequestOptionChunkStream) {
			v2writer = vmessio.NewAuthChunkWriter(v2writer)
		}

		if err := v2writer.Write(data); err != nil {
			connection.SetReusable(false)
		}

		writer.SetCached(false)

		err = v2io.Pipe(output, v2writer)
		if err != io.EOF {
			connection.SetReusable(false)
		}

		output.Release()
		if request.Option.Has(protocol.RequestOptionChunkStream) {
			if err := v2writer.Write(alloc.NewSmallBuffer().Clear()); err != nil {
				connection.SetReusable(false)
			}
		}
		v2writer.Release()
	}

	readFinish.Lock()
}

func init() {
	internal.MustRegisterInboundHandlerCreator("vmess",
		func(space app.Space, rawConfig interface{}, meta *proxy.InboundHandlerMeta) (proxy.InboundHandler, error) {
			if !space.HasApp(dispatcher.APP_ID) {
				return nil, internal.ErrorBadConfiguration
			}
			config := rawConfig.(*Config)

			allowedClients := protocol.NewTimedUserValidator(protocol.DefaultIDHash)
			for _, user := range config.AllowedUsers {
				allowedClients.Add(user)
			}

			handler := &VMessInboundHandler{
				packetDispatcher: space.GetApp(dispatcher.APP_ID).(dispatcher.PacketDispatcher),
				clients:          allowedClients,
				detours:          config.DetourConfig,
				usersByEmail:     NewUserByEmail(config.AllowedUsers, config.Defaults),
				meta:             meta,
			}

			if space.HasApp(proxyman.APP_ID_INBOUND_MANAGER) {
				handler.inboundHandlerManager = space.GetApp(proxyman.APP_ID_INBOUND_MANAGER).(proxyman.InboundHandlerManager)
			}

			return handler, nil
		})
}
