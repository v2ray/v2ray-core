package vmess

import (
	"crypto/md5"
	"crypto/rand"
	mrand "math/rand"
	"net"
	"sync"

	"github.com/v2ray/v2ray-core/common/alloc"
	v2crypto "github.com/v2ray/v2ray-core/common/crypto"
	"github.com/v2ray/v2ray-core/common/log"
	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/proxy/common/connhandler"
	"github.com/v2ray/v2ray-core/proxy/vmess/config"
	"github.com/v2ray/v2ray-core/proxy/vmess/protocol"
	"github.com/v2ray/v2ray-core/proxy/vmess/protocol/user"
	"github.com/v2ray/v2ray-core/transport/ray"
)

type VMessOutboundHandler struct {
	vNextList []*config.OutboundTarget
}

func NewVMessOutboundHandler(vNextList []*config.OutboundTarget) *VMessOutboundHandler {
	return &VMessOutboundHandler{
		vNextList: vNextList,
	}
}

func pickVNext(serverList []*config.OutboundTarget) (v2net.Destination, config.User) {
	vNextLen := len(serverList)
	if vNextLen == 0 {
		panic("VMessOut: Zero vNext is configured.")
	}
	vNextIndex := 0
	if vNextLen > 1 {
		vNextIndex = mrand.Intn(vNextLen)
	}

	vNext := serverList[vNextIndex]
	vNextUserLen := len(vNext.Accounts)
	if vNextUserLen == 0 {
		panic("VMessOut: Zero User account.")
	}
	vNextUserIndex := 0
	if vNextUserLen > 1 {
		vNextUserIndex = mrand.Intn(vNextUserLen)
	}
	vNextUser := vNext.Accounts[vNextUserIndex]
	return vNext.Destination, vNextUser
}

func (this *VMessOutboundHandler) Dispatch(firstPacket v2net.Packet, ray ray.OutboundRay) error {
	vNextList := this.vNextList
	vNextAddress, vNextUser := pickVNext(vNextList)

	command := protocol.CmdTCP
	if firstPacket.Destination().IsUDP() {
		command = protocol.CmdUDP
	}
	request := &protocol.VMessRequest{
		Version: protocol.Version,
		User:    vNextUser,
		Command: command,
		Address: firstPacket.Destination().Address(),
	}

	buffer := alloc.NewSmallBuffer()
	defer buffer.Release()
	v2net.ReadAllBytes(rand.Reader, buffer.Value[:36]) // 16 + 16 + 4
	request.RequestIV = buffer.Value[:16]
	request.RequestKey = buffer.Value[16:32]
	request.ResponseHeader = buffer.Value[32:36]

	return startCommunicate(request, vNextAddress, ray, firstPacket)
}

func startCommunicate(request *protocol.VMessRequest, dest v2net.Destination, ray ray.OutboundRay, firstPacket v2net.Packet) error {
	conn, err := net.Dial(dest.Network(), dest.Address().String())
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
	go handleResponse(conn, request, output, &responseFinish, dest.IsUDP())

	requestFinish.Lock()
	if tcpConn, ok := conn.(*net.TCPConn); ok {
		tcpConn.CloseWrite()
	}
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
	buffer, err = request.ToBytes(user.NewTimeHash(user.HMACHash{}), user.GenerateRandomInt64InRange, buffer)
	if err != nil {
		log.Error("VMessOut: Failed to serialize VMess request: %v", err)
		return
	}

	// Send first packet of payload together with request, in favor of small requests.
	firstChunk := firstPacket.Chunk()
	moreChunks := firstPacket.MoreChunks()

	if firstChunk == nil && moreChunks {
		firstChunk, moreChunks = <-input
	}

	if firstChunk != nil {
		aesStream.XORKeyStream(firstChunk.Value, firstChunk.Value)
		buffer.Append(firstChunk.Value)
		firstChunk.Release()

		_, err = conn.Write(buffer.Value)
		buffer.Release()
		if err != nil {
			log.Error("VMessOut: Failed to write VMess request: %v", err)
			return
		}
	}

	if moreChunks {
		v2net.ChanToWriter(encryptRequestWriter, input)
	}
	return
}

func headerMatch(request *protocol.VMessRequest, responseHeader []byte) bool {
	return ((request.ResponseHeader[0] ^ request.ResponseHeader[1]) == responseHeader[0]) &&
		((request.ResponseHeader[2] ^ request.ResponseHeader[3]) == responseHeader[1])
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
		return
	}
	if buffer.Len() < 4 || !headerMatch(request, buffer.Value[:2]) {
		log.Warning("VMessOut: unexepcted response header. The connection is probably hijacked.")
		return
	}
	log.Info("VMessOut received %d bytes from %s", buffer.Len()-4, conn.RemoteAddr().String())

	buffer.SliceFrom(4)
	output <- buffer

	if !isUDP {
		v2net.ReaderToChan(output, decryptResponseReader)
	}

	return
}

type VMessOutboundHandlerFactory struct {
}

func (this *VMessOutboundHandlerFactory) Create(rawConfig interface{}) (connhandler.OutboundConnectionHandler, error) {
	vOutConfig := rawConfig.(config.Outbound)
	return NewVMessOutboundHandler(vOutConfig.Targets()), nil
}

func init() {
	connhandler.RegisterOutboundConnectionHandlerFactory("vmess", &VMessOutboundHandlerFactory{})
}
