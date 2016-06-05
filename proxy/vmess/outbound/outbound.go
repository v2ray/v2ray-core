package outbound

import (
	"io"
	"sync"

	"github.com/v2ray/v2ray-core/app"
	"github.com/v2ray/v2ray-core/common/alloc"
	v2io "github.com/v2ray/v2ray-core/common/io"
	"github.com/v2ray/v2ray-core/common/log"
	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/common/protocol"
	"github.com/v2ray/v2ray-core/common/protocol/raw"
	"github.com/v2ray/v2ray-core/common/retry"
	"github.com/v2ray/v2ray-core/proxy"
	"github.com/v2ray/v2ray-core/proxy/internal"
	vmessio "github.com/v2ray/v2ray-core/proxy/vmess/io"
	"github.com/v2ray/v2ray-core/transport"
	"github.com/v2ray/v2ray-core/transport/hub"
	"github.com/v2ray/v2ray-core/transport/ray"
)

type VMessOutboundHandler struct {
	receiverManager *ReceiverManager
	meta            *proxy.OutboundHandlerMeta
}

func (this *VMessOutboundHandler) Dispatch(target v2net.Destination, payload *alloc.Buffer, ray ray.OutboundRay) error {
	defer ray.OutboundInput().Release()
	defer ray.OutboundOutput().Close()

	var rec *Receiver
	var conn *hub.Connection

	err := retry.Timed(5, 100).On(func() error {
		rec = this.receiverManager.PickReceiver()
		rawConn, err := hub.Dial(this.meta.Address, rec.Destination)
		if err != nil {
			return err
		}
		conn = rawConn

		return nil
	})
	if err != nil {
		log.Error("Failed to find an available destination:", err)
		return err
	}
	log.Info("VMessOut: Tunneling request to ", target, " via ", rec.Destination)

	command := protocol.RequestCommandTCP
	if target.IsUDP() {
		command = protocol.RequestCommandUDP
	}
	request := &protocol.RequestHeader{
		Version: raw.Version,
		User:    rec.PickUser(),
		Command: command,
		Address: target.Address(),
		Port:    target.Port(),
		Option:  protocol.RequestOptionChunkStream,
	}

	defer conn.Close()

	if transport.IsConnectionReusable() {
		request.Option.Set(protocol.RequestOptionConnectionReuse)
		conn.SetReusable(true)
	}

	input := ray.OutboundInput()
	output := ray.OutboundOutput()

	var requestFinish, responseFinish sync.Mutex
	requestFinish.Lock()
	responseFinish.Lock()

	session := raw.NewClientSession(protocol.DefaultIDHash)

	go this.handleRequest(session, conn, request, payload, input, &requestFinish)
	go this.handleResponse(session, conn, request, rec.Destination, output, &responseFinish)

	requestFinish.Lock()
	responseFinish.Lock()
	return nil
}

func (this *VMessOutboundHandler) handleRequest(session *raw.ClientSession, conn *hub.Connection, request *protocol.RequestHeader, payload *alloc.Buffer, input v2io.Reader, finish *sync.Mutex) {
	defer finish.Unlock()

	writer := v2io.NewBufferedWriter(conn)
	defer writer.Release()
	session.EncodeRequestHeader(request, writer)

	bodyWriter := session.EncodeRequestBody(writer)
	var streamWriter v2io.Writer = v2io.NewAdaptiveWriter(bodyWriter)
	if request.Option.Has(protocol.RequestOptionChunkStream) {
		streamWriter = vmessio.NewAuthChunkWriter(streamWriter)
	}
	if err := streamWriter.Write(payload); err != nil {
		conn.SetReusable(false)
	}
	writer.SetCached(false)

	err := v2io.Pipe(input, streamWriter)
	if err != io.EOF {
		conn.SetReusable(false)
	}

	if request.Option.Has(protocol.RequestOptionChunkStream) {
		err := streamWriter.Write(alloc.NewSmallBuffer().Clear())
		if err != nil {
			conn.SetReusable(false)
		}
	}
	streamWriter.Release()
	return
}

func (this *VMessOutboundHandler) handleResponse(session *raw.ClientSession, conn *hub.Connection, request *protocol.RequestHeader, dest v2net.Destination, output v2io.Writer, finish *sync.Mutex) {
	defer finish.Unlock()

	reader := v2io.NewBufferedReader(conn)
	defer reader.Release()

	header, err := session.DecodeResponseHeader(reader)
	if err != nil {
		conn.SetReusable(false)
		log.Warning("VMessOut: Failed to read response: ", err)
		return
	}
	go this.handleCommand(dest, header.Command)

	if !header.Option.Has(protocol.ResponseOptionConnectionReuse) {
		conn.SetReusable(false)
	}

	reader.SetCached(false)
	decryptReader := session.DecodeResponseBody(reader)

	var bodyReader v2io.Reader
	if request.Option.Has(protocol.RequestOptionChunkStream) {
		bodyReader = vmessio.NewAuthChunkReader(decryptReader)
	} else {
		bodyReader = v2io.NewAdaptiveReader(decryptReader)
	}

	err = v2io.Pipe(bodyReader, output)
	if err != io.EOF {
		conn.SetReusable(false)
	}

	bodyReader.Release()

	return
}

func init() {
	internal.MustRegisterOutboundHandlerCreator("vmess",
		func(space app.Space, rawConfig interface{}, meta *proxy.OutboundHandlerMeta) (proxy.OutboundHandler, error) {
			vOutConfig := rawConfig.(*Config)
			return &VMessOutboundHandler{
				receiverManager: NewReceiverManager(vOutConfig.Receivers),
				meta:            meta,
			}, nil
		})
}
