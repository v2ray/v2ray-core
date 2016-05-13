package outbound

import (
	"net"
	"sync"

	"github.com/v2ray/v2ray-core/app"
	"github.com/v2ray/v2ray-core/common/alloc"
	v2io "github.com/v2ray/v2ray-core/common/io"
	"github.com/v2ray/v2ray-core/common/log"
	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/common/protocol"
	"github.com/v2ray/v2ray-core/common/protocol/raw"
	"github.com/v2ray/v2ray-core/proxy"
	"github.com/v2ray/v2ray-core/proxy/internal"
	vmessio "github.com/v2ray/v2ray-core/proxy/vmess/io"
	"github.com/v2ray/v2ray-core/transport/dialer"
	"github.com/v2ray/v2ray-core/transport/ray"
)

type VMessOutboundHandler struct {
	receiverManager *ReceiverManager
}

func (this *VMessOutboundHandler) Dispatch(target v2net.Destination, payload *alloc.Buffer, ray ray.OutboundRay) error {
	defer ray.OutboundInput().Release()
	defer ray.OutboundOutput().Close()

	destination, vNextUser := this.receiverManager.PickReceiver()

	command := protocol.RequestCommandTCP
	if target.IsUDP() {
		command = protocol.RequestCommandUDP
	}
	request := &protocol.RequestHeader{
		Version: raw.Version,
		User:    vNextUser,
		Command: command,
		Address: target.Address(),
		Port:    target.Port(),
		Option:  protocol.RequestOptionChunkStream,
	}

	conn, err := dialer.Dial(destination)
	if err != nil {
		log.Error("Failed to open ", destination, ": ", err)
		return err
	}
	log.Info("VMessOut: Tunneling request to ", request.Address, " via ", destination)

	defer conn.Close()

	input := ray.OutboundInput()
	output := ray.OutboundOutput()

	var requestFinish, responseFinish sync.Mutex
	requestFinish.Lock()
	responseFinish.Lock()

	session := raw.NewClientSession(protocol.DefaultIDHash)

	go this.handleRequest(session, conn, request, payload, input, &requestFinish)
	go this.handleResponse(session, conn, request, destination, output, &responseFinish)

	requestFinish.Lock()
	responseFinish.Lock()
	return nil
}

func (this *VMessOutboundHandler) handleRequest(session *raw.ClientSession, conn net.Conn, request *protocol.RequestHeader, payload *alloc.Buffer, input v2io.Reader, finish *sync.Mutex) {
	defer finish.Unlock()

	writer := v2io.NewBufferedWriter(conn)
	defer writer.Release()
	session.EncodeRequestHeader(request, writer)

	bodyWriter := session.EncodeRequestBody(writer)
	var streamWriter v2io.Writer = v2io.NewAdaptiveWriter(bodyWriter)
	if request.Option.IsChunkStream() {
		streamWriter = vmessio.NewAuthChunkWriter(streamWriter)
	}
	streamWriter.Write(payload)
	writer.SetCached(false)

	v2io.Pipe(input, streamWriter)
	if request.Option.IsChunkStream() {
		streamWriter.Write(alloc.NewSmallBuffer().Clear())
	}
	streamWriter.Release()
	return
}

func (this *VMessOutboundHandler) handleResponse(session *raw.ClientSession, conn net.Conn, request *protocol.RequestHeader, dest v2net.Destination, output v2io.Writer, finish *sync.Mutex) {
	defer finish.Unlock()

	reader := v2io.NewBufferedReader(conn)
	defer reader.Release()

	header, err := session.DecodeResponseHeader(reader)
	if err != nil {
		log.Warning("VMessOut: Failed to read response: ", err)
		return
	}
	go this.handleCommand(dest, header.Command)

	reader.SetCached(false)
	decryptReader := session.DecodeResponseBody(reader)

	var bodyReader v2io.Reader
	if request.Option.IsChunkStream() {
		bodyReader = vmessio.NewAuthChunkReader(decryptReader)
	} else {
		bodyReader = v2io.NewAdaptiveReader(decryptReader)
	}

	v2io.Pipe(bodyReader, output)
	bodyReader.Release()

	return
}

func init() {
	internal.MustRegisterOutboundHandlerCreator("vmess",
		func(space app.Space, rawConfig interface{}) (proxy.OutboundHandler, error) {
			vOutConfig := rawConfig.(*Config)
			return &VMessOutboundHandler{
				receiverManager: NewReceiverManager(vOutConfig.Receivers),
			}, nil
		})
}
