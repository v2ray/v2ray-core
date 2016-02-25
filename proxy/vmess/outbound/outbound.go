package outbound

import (
	"crypto/md5"
	"crypto/rand"
	"io"
	"net"
	"sync"
	"time"

	"github.com/v2ray/v2ray-core/app"
	"github.com/v2ray/v2ray-core/common/alloc"
	v2crypto "github.com/v2ray/v2ray-core/common/crypto"
	v2io "github.com/v2ray/v2ray-core/common/io"
	"github.com/v2ray/v2ray-core/common/log"
	v2net "github.com/v2ray/v2ray-core/common/net"
	proto "github.com/v2ray/v2ray-core/common/protocol"
	"github.com/v2ray/v2ray-core/proxy"
	"github.com/v2ray/v2ray-core/proxy/internal"
	vmessio "github.com/v2ray/v2ray-core/proxy/vmess/io"
	"github.com/v2ray/v2ray-core/proxy/vmess/protocol"
	"github.com/v2ray/v2ray-core/transport/ray"
)

type VMessOutboundHandler struct {
	receiverManager *ReceiverManager
}

func (this *VMessOutboundHandler) Dispatch(firstPacket v2net.Packet, ray ray.OutboundRay) error {
	vNextAddress, vNextUser := this.receiverManager.PickReceiver()

	command := protocol.CmdTCP
	if firstPacket.Destination().IsUDP() {
		command = protocol.CmdUDP
	}
	request := &protocol.VMessRequest{
		Version: protocol.Version,
		User:    vNextUser,
		Command: command,
		Address: firstPacket.Destination().Address(),
		Port:    firstPacket.Destination().Port(),
	}
	if command == protocol.CmdUDP {
		request.Option |= protocol.OptionChunk
	}

	buffer := alloc.NewSmallBuffer()
	defer buffer.Release()                      // Buffer is released after communication finishes.
	io.ReadFull(rand.Reader, buffer.Value[:33]) // 16 + 16 + 1
	request.RequestIV = buffer.Value[:16]
	request.RequestKey = buffer.Value[16:32]
	request.ResponseHeader = buffer.Value[32]

	return this.startCommunicate(request, vNextAddress, ray, firstPacket)
}

func (this *VMessOutboundHandler) startCommunicate(request *protocol.VMessRequest, dest v2net.Destination, ray ray.OutboundRay, firstPacket v2net.Packet) error {
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

	go this.handleRequest(conn, request, firstPacket, input, &requestFinish)
	go this.handleResponse(conn, request, dest, output, &responseFinish)

	requestFinish.Lock()
	conn.CloseWrite()
	responseFinish.Lock()
	return nil
}

func (this *VMessOutboundHandler) handleRequest(conn net.Conn, request *protocol.VMessRequest, firstPacket v2net.Packet, input <-chan *alloc.Buffer, finish *sync.Mutex) {
	defer finish.Unlock()
	aesStream, err := v2crypto.NewAesEncryptionStream(request.RequestKey[:], request.RequestIV[:])
	if err != nil {
		log.Error("VMessOut: Failed to create AES encryption stream: ", err)
		return
	}
	encryptRequestWriter := v2crypto.NewCryptionWriter(aesStream, conn)

	buffer := alloc.NewBuffer().Clear()
	defer buffer.Release()
	buffer, err = request.ToBytes(proto.NewTimestampGenerator(proto.Timestamp(time.Now().Unix()), 30), buffer)
	if err != nil {
		log.Error("VMessOut: Failed to serialize VMess request: ", err)
		return
	}

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

	if request.IsChunkStream() {
		vmessio.Authenticate(firstChunk)
	}

	aesStream.XORKeyStream(firstChunk.Value, firstChunk.Value)
	buffer.Append(firstChunk.Value)
	firstChunk.Release()

	_, err = conn.Write(buffer.Value)
	if err != nil {
		log.Error("VMessOut: Failed to write VMess request: ", err)
		return
	}

	if moreChunks {
		var streamWriter v2io.Writer = v2io.NewAdaptiveWriter(encryptRequestWriter)
		if request.IsChunkStream() {
			streamWriter = vmessio.NewAuthChunkWriter(streamWriter)
		}
		v2io.ChanToWriter(streamWriter, input)
	}
	return
}

func headerMatch(request *protocol.VMessRequest, responseHeader byte) bool {
	return request.ResponseHeader == responseHeader
}

func (this *VMessOutboundHandler) handleResponse(conn net.Conn, request *protocol.VMessRequest, dest v2net.Destination, output chan<- *alloc.Buffer, finish *sync.Mutex) {
	defer finish.Unlock()
	defer close(output)
	responseKey := md5.Sum(request.RequestKey[:])
	responseIV := md5.Sum(request.RequestIV[:])

	aesStream, err := v2crypto.NewAesDecryptionStream(responseKey[:], responseIV[:])
	if err != nil {
		log.Error("VMessOut: Failed to create AES encryption stream: ", err)
		return
	}
	decryptResponseReader := v2crypto.NewCryptionReader(aesStream, conn)

	buffer := alloc.NewSmallBuffer()
	defer buffer.Release()
	_, err = io.ReadFull(decryptResponseReader, buffer.Value[:4])

	if err != nil {
		log.Error("VMessOut: Failed to read VMess response (", buffer.Len(), " bytes): ", err)
		return
	}
	if !headerMatch(request, buffer.Value[0]) {
		log.Warning("VMessOut: unexepcted response header. The connection is probably hijacked.")
		return
	}

	if buffer.Value[2] != 0 {
		command := buffer.Value[2]
		dataLen := int(buffer.Value[3])
		_, err := io.ReadFull(decryptResponseReader, buffer.Value[:dataLen])
		if err != nil {
			log.Error("VMessOut: Failed to read response command: ", err)
			return
		}
		data := buffer.Value[:dataLen]
		go this.handleCommand(dest, command, data)
	}

	var reader v2io.Reader
	if request.IsChunkStream() {
		reader = vmessio.NewAuthChunkReader(decryptResponseReader)
	} else {
		reader = v2io.NewAdaptiveReader(decryptResponseReader)
	}

	v2io.ReaderToChan(output, reader)

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
