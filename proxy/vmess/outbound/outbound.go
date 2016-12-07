package outbound

import (
	"sync"

	"v2ray.com/core/app"
	"v2ray.com/core/common/alloc"
	v2io "v2ray.com/core/common/io"
	"v2ray.com/core/common/loader"
	"v2ray.com/core/common/log"
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/common/protocol"
	"v2ray.com/core/common/retry"
	"v2ray.com/core/proxy"
	"v2ray.com/core/proxy/registry"
	"v2ray.com/core/proxy/vmess"
	"v2ray.com/core/proxy/vmess/encoding"
	"v2ray.com/core/transport/internet"
	"v2ray.com/core/transport/ray"
)

type VMessOutboundHandler struct {
	serverList   *protocol.ServerList
	serverPicker protocol.ServerPicker
	meta         *proxy.OutboundHandlerMeta
}

func (v *VMessOutboundHandler) Dispatch(target v2net.Destination, payload *alloc.Buffer, ray ray.OutboundRay) {
	defer ray.OutboundInput().Release()
	defer ray.OutboundOutput().Close()

	var rec *protocol.ServerSpec
	var conn internet.Connection

	err := retry.ExponentialBackoff(5, 100).On(func() error {
		rec = v.serverPicker.PickServer()
		rawConn, err := internet.Dial(v.meta.Address, rec.Destination(), v.meta.GetDialerOptions())
		if err != nil {
			return err
		}
		conn = rawConn

		return nil
	})
	if err != nil {
		log.Warning("VMess|Outbound: Failed to find an available destination:", err)
		return
	}
	log.Info("VMess|Outbound: Tunneling request to ", target, " via ", rec.Destination())

	command := protocol.RequestCommandTCP
	if target.Network == v2net.Network_UDP {
		command = protocol.RequestCommandUDP
	}
	request := &protocol.RequestHeader{
		Version: encoding.Version,
		User:    rec.PickUser(),
		Command: command,
		Address: target.Address,
		Port:    target.Port,
		Option:  protocol.RequestOptionChunkStream,
	}

	rawAccount, err := request.User.GetTypedAccount()
	if err != nil {
		log.Warning("VMess|Outbound: Failed to get user account: ", err)
	}
	account := rawAccount.(*vmess.InternalAccount)
	request.Security = account.Security

	defer conn.Close()

	conn.SetReusable(true)
	if conn.Reusable() { // Conn reuse may be disabled on transportation layer
		request.Option.Set(protocol.RequestOptionConnectionReuse)
	}

	input := ray.OutboundInput()
	output := ray.OutboundOutput()

	var requestFinish, responseFinish sync.Mutex
	requestFinish.Lock()
	responseFinish.Lock()

	session := encoding.NewClientSession(protocol.DefaultIDHash)

	go v.handleRequest(session, conn, request, payload, input, &requestFinish)
	go v.handleResponse(session, conn, request, rec.Destination(), output, &responseFinish)

	requestFinish.Lock()
	responseFinish.Lock()
	return
}

func (v *VMessOutboundHandler) handleRequest(session *encoding.ClientSession, conn internet.Connection, request *protocol.RequestHeader, payload *alloc.Buffer, input v2io.Reader, finish *sync.Mutex) {
	defer finish.Unlock()

	writer := v2io.NewBufferedWriter(conn)
	defer writer.Release()
	session.EncodeRequestHeader(request, writer)

	bodyWriter := session.EncodeRequestBody(request, writer)
	defer bodyWriter.Release()

	if !payload.IsEmpty() {
		if err := bodyWriter.Write(payload); err != nil {
			log.Info("VMess|Outbound: Failed to write payload. Disabling connection reuse.", err)
			conn.SetReusable(false)
		}
		payload.Release()
	}
	writer.SetCached(false)

	if err := v2io.PipeUntilEOF(input, bodyWriter); err != nil {
		conn.SetReusable(false)
	}

	if request.Option.Has(protocol.RequestOptionChunkStream) {
		err := bodyWriter.Write(alloc.NewLocalBuffer(32))
		if err != nil {
			conn.SetReusable(false)
		}
	}
	return
}

func (v *VMessOutboundHandler) handleResponse(session *encoding.ClientSession, conn internet.Connection, request *protocol.RequestHeader, dest v2net.Destination, output v2io.Writer, finish *sync.Mutex) {
	defer finish.Unlock()

	reader := v2io.NewBufferedReader(conn)
	defer reader.Release()

	header, err := session.DecodeResponseHeader(reader)
	if err != nil {
		conn.SetReusable(false)
		log.Warning("VMess|Outbound: Failed to read response from ", request.Destination(), ": ", err)
		return
	}
	go v.handleCommand(dest, header.Command)

	if !header.Option.Has(protocol.ResponseOptionConnectionReuse) {
		conn.SetReusable(false)
	}

	reader.SetCached(false)
	bodyReader := session.DecodeResponseBody(request, reader)
	defer bodyReader.Release()

	if err := v2io.PipeUntilEOF(bodyReader, output); err != nil {
		conn.SetReusable(false)
	}

	return
}

type Factory struct{}

func (v *Factory) StreamCapability() v2net.NetworkList {
	return v2net.NetworkList{
		Network: []v2net.Network{v2net.Network_TCP, v2net.Network_KCP, v2net.Network_WebSocket},
	}
}

func (v *Factory) Create(space app.Space, rawConfig interface{}, meta *proxy.OutboundHandlerMeta) (proxy.OutboundHandler, error) {
	vOutConfig := rawConfig.(*Config)

	serverList := protocol.NewServerList()
	for _, rec := range vOutConfig.Receiver {
		serverList.AddServer(protocol.NewServerSpecFromPB(*rec))
	}
	handler := &VMessOutboundHandler{
		serverList:   serverList,
		serverPicker: protocol.NewRoundRobinServerPicker(serverList),
		meta:         meta,
	}

	return handler, nil
}

func init() {
	registry.MustRegisterOutboundHandlerCreator(loader.GetType(new(Config)), new(Factory))
}
