package outbound

import (
	"net"
	"sync"

	"github.com/v2ray/v2ray-core/app"
	"github.com/v2ray/v2ray-core/common/alloc"
	v2io "github.com/v2ray/v2ray-core/common/io"
	"github.com/v2ray/v2ray-core/common/log"
	v2net "github.com/v2ray/v2ray-core/common/net"
	proto "github.com/v2ray/v2ray-core/common/protocol"
	raw "github.com/v2ray/v2ray-core/common/protocol/raw"
	"github.com/v2ray/v2ray-core/proxy"
	"github.com/v2ray/v2ray-core/proxy/internal"
	vmessio "github.com/v2ray/v2ray-core/proxy/vmess/io"
	"github.com/v2ray/v2ray-core/transport/ray"
)

type VMessOutboundHandler struct {
	receiverManager *ReceiverManager
}

func (this *VMessOutboundHandler) Dispatch(firstPacket v2net.Packet, ray ray.OutboundRay) error {
	vNextAddress, vNextUser := this.receiverManager.PickReceiver()

	command := proto.RequestCommandTCP
	if firstPacket.Destination().IsUDP() {
		command = proto.RequestCommandUDP
	}
	request := &proto.RequestHeader{
		Version: raw.Version,
		User:    vNextUser,
		Command: command,
		Address: firstPacket.Destination().Address(),
		Port:    firstPacket.Destination().Port(),
	}
	if command == proto.RequestCommandUDP {
		request.Option |= proto.RequestOptionChunkStream
	}

	return this.startCommunicate(request, vNextAddress, ray, firstPacket)
}

func (this *VMessOutboundHandler) startCommunicate(request *proto.RequestHeader, dest v2net.Destination, ray ray.OutboundRay, firstPacket v2net.Packet) error {
	var destIP net.IP
	if dest.Address().IsIPv4() || dest.Address().IsIPv6() {
		destIP = dest.Address().IP()
	} else {
		ips, err := net.LookupIP(dest.Address().Domain())
		if err != nil {
			return err
		}
		destIP = ips[0]
	}
	conn, err := net.DialTCP("tcp", nil, &net.TCPAddr{
		IP:   destIP,
		Port: int(dest.Port()),
	})
	if err != nil {
		log.Error("Failed to open ", dest, ": ", err)
		if ray != nil {
			close(ray.OutboundOutput())
		}
		return err
	}
	log.Info("VMessOut: Tunneling request to ", request.Address, " via ", dest)

	defer conn.Close()

	input := ray.OutboundInput()
	output := ray.OutboundOutput()

	var requestFinish, responseFinish sync.Mutex
	requestFinish.Lock()
	responseFinish.Lock()

	session := raw.NewClientSession(proto.DefaultIDHash)

	go this.handleRequest(session, conn, request, firstPacket, input, &requestFinish)
	go this.handleResponse(session, conn, request, dest, output, &responseFinish)

	requestFinish.Lock()
	conn.CloseWrite()
	responseFinish.Lock()
	return nil
}

func (this *VMessOutboundHandler) handleRequest(session *raw.ClientSession, conn net.Conn, request *proto.RequestHeader, firstPacket v2net.Packet, input <-chan *alloc.Buffer, finish *sync.Mutex) {
	defer finish.Unlock()

	writer := v2io.NewBufferedWriter(conn)
	session.EncodeRequestHeader(request, writer)

	// Send first packet of payload together with request, in favor of small requests.
	firstChunk := firstPacket.Chunk()
	moreChunks := firstPacket.MoreChunks()

	for firstChunk == nil && moreChunks {
		firstChunk, moreChunks = <-input
	}

	if firstChunk == nil && !moreChunks {
		log.Warning("VMessOut: Nothing to send. Existing...")
		return
	}

	if request.Option.IsChunkStream() {
		vmessio.Authenticate(firstChunk)
	}

	bodyWriter := session.EncodeRequestBody(writer)
	bodyWriter.Write(firstChunk.Value)
	firstChunk.Release()

	writer.SetCached(false)

	if moreChunks {
		var streamWriter v2io.ReleasableWriter = v2io.NewAdaptiveWriter(bodyWriter)
		if request.Option.IsChunkStream() {
			streamWriter = vmessio.NewAuthChunkWriter(streamWriter)
		}
		v2io.ChanToWriter(streamWriter, input)
        streamWriter.Release()
	}
	return
}

func (this *VMessOutboundHandler) handleResponse(session *raw.ClientSession, conn net.Conn, request *proto.RequestHeader, dest v2net.Destination, output chan<- *alloc.Buffer, finish *sync.Mutex) {
	defer finish.Unlock()
	defer close(output)

	reader := v2io.NewBufferedReader(conn)

	header, err := session.DecodeResponseHeader(reader)
	if err != nil {
		log.Warning("VMessOut: Failed to read response: ", err)
		return
	}
	go this.handleCommand(dest, header.Command)

	reader.SetCached(false)
	decryptReader := session.DecodeResponseBody(conn)

	var bodyReader v2io.Reader
	if request.Option.IsChunkStream() {
		bodyReader = vmessio.NewAuthChunkReader(decryptReader)
	} else {
		bodyReader = v2io.NewAdaptiveReader(decryptReader)
	}

	v2io.ReaderToChan(output, bodyReader)

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
