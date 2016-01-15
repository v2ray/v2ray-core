package outbound

import (
	"crypto/md5"
	"crypto/rand"
	"net"
	"sync"
	"time"

	"github.com/v2ray/v2ray-core/app"
	"github.com/v2ray/v2ray-core/common/alloc"
	v2crypto "github.com/v2ray/v2ray-core/common/crypto"
	"github.com/v2ray/v2ray-core/common/log"
	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/proxy"
	"github.com/v2ray/v2ray-core/proxy/internal"
	"github.com/v2ray/v2ray-core/proxy/vmess/protocol"
	"github.com/v2ray/v2ray-core/transport/ray"
)

type VMessOutboundHandler struct {
	receiverManager *ReceiverManager
	space           app.Space
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

	buffer := alloc.NewSmallBuffer()
	defer buffer.Release()                             // Buffer is released after communication finishes.
	v2net.ReadAllBytes(rand.Reader, buffer.Value[:33]) // 16 + 16 + 1
	buffer.Value[33] = 0
	buffer.Value[34] = 0
	buffer.Value[35] = 0
	request.RequestIV = buffer.Value[:16]
	request.RequestKey = buffer.Value[16:32]
	request.ResponseHeader = buffer.Value[32:36]

	return startCommunicate(request, vNextAddress, ray, firstPacket)
}

func startCommunicate(request *protocol.VMessRequest, dest v2net.Destination, ray ray.OutboundRay, firstPacket v2net.Packet) error {
	var destIp net.IP
	if dest.Address().IsIPv4() || dest.Address().IsIPv6() {
		destIp = dest.Address().IP()
	} else {
		ips, err := net.LookupIP(dest.Address().Domain())
		if err != nil {
			return err
		}
		destIp = ips[0]
	}
	conn, err := net.DialTCP("tcp", nil, &net.TCPAddr{
		IP:   destIp,
		Port: int(dest.Port()),
	})
	if err != nil {
		log.Error("Failed to open %s: %v", dest.String(), err)
		if ray != nil {
			close(ray.OutboundOutput())
		}
		return err
	}
	log.Info("VMessOut: Tunneling request to %s via %s", request.Address.String(), dest.String())

	defer conn.Close()

	input := ray.OutboundInput()
	output := ray.OutboundOutput()

	var requestFinish, responseFinish sync.Mutex
	requestFinish.Lock()
	responseFinish.Lock()

	go handleRequest(conn, request, firstPacket, input, &requestFinish)
	go handleResponse(conn, request, output, &responseFinish, (request.Command == protocol.CmdUDP))

	requestFinish.Lock()
	conn.CloseWrite()
	responseFinish.Lock()
	return nil
}

func handleRequest(conn net.Conn, request *protocol.VMessRequest, firstPacket v2net.Packet, input <-chan *alloc.Buffer, finish *sync.Mutex) {
	defer finish.Unlock()
	aesStream, err := v2crypto.NewAesEncryptionStream(request.RequestKey[:], request.RequestIV[:])
	if err != nil {
		log.Error("VMessOut: Failed to create AES encryption stream: %v", err)
		return
	}
	encryptRequestWriter := v2crypto.NewCryptionWriter(aesStream, conn)

	buffer := alloc.NewBuffer().Clear()
	defer buffer.Release()
	buffer, err = request.ToBytes(protocol.NewRandomTimestampGenerator(protocol.Timestamp(time.Now().Unix()), 30), buffer)
	if err != nil {
		log.Error("VMessOut: Failed to serialize VMess request: %v", err)
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

	aesStream.XORKeyStream(firstChunk.Value, firstChunk.Value)
	buffer.Append(firstChunk.Value)
	firstChunk.Release()

	_, err = conn.Write(buffer.Value)
	if err != nil {
		log.Error("VMessOut: Failed to write VMess request: %v", err)
		return
	}

	if moreChunks {
		v2net.ChanToWriter(encryptRequestWriter, input)
	}
	return
}

func headerMatch(request *protocol.VMessRequest, responseHeader []byte) bool {
	return (request.ResponseHeader[0] == responseHeader[0])
}

func handleResponse(conn net.Conn, request *protocol.VMessRequest, output chan<- *alloc.Buffer, finish *sync.Mutex, isUDP bool) {
	defer finish.Unlock()
	defer close(output)
	responseKey := md5.Sum(request.RequestKey[:])
	responseIV := md5.Sum(request.RequestIV[:])

	aesStream, err := v2crypto.NewAesDecryptionStream(responseKey[:], responseIV[:])
	if err != nil {
		log.Error("VMessOut: Failed to create AES encryption stream: %v", err)
		return
	}
	decryptResponseReader := v2crypto.NewCryptionReader(aesStream, conn)

	buffer, err := v2net.ReadFrom(decryptResponseReader, nil)
	if err != nil {
		log.Error("VMessOut: Failed to read VMess response (%d bytes): %v", buffer.Len(), err)
		buffer.Release()
		return
	}
	if buffer.Len() < 4 || !headerMatch(request, buffer.Value[:2]) {
		log.Warning("VMessOut: unexepcted response header. The connection is probably hijacked.")
		return
	}
	log.Info("VMessOut received %d bytes from %s", buffer.Len()-4, conn.RemoteAddr().String())

	responseBegin := 4
	if buffer.Value[2] != 0 {
		dataLen := int(buffer.Value[3])
		if buffer.Len() < dataLen+4 { // Rare case
			diffBuffer := make([]byte, dataLen+4-buffer.Len())
			v2net.ReadAllBytes(decryptResponseReader, diffBuffer)
			buffer.Append(diffBuffer)
		}
		command := buffer.Value[2]
		data := buffer.Value[4 : 4+dataLen]
		go handleCommand(command, data)
		responseBegin = 4 + dataLen
	}

	buffer.SliceFrom(responseBegin)
	output <- buffer

	if !isUDP {
		v2net.ReaderToChan(output, decryptResponseReader)
	}

	return
}

func init() {
	internal.MustRegisterOutboundConnectionHandlerCreator("vmess",
		func(space app.Space, rawConfig interface{}) (proxy.OutboundConnectionHandler, error) {
			vOutConfig := rawConfig.(*Config)
			return &VMessOutboundHandler{
				space:           space,
				receiverManager: NewReceiverManager(vOutConfig.Receivers),
			}, nil
		})
}
