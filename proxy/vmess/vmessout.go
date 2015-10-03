package vmess

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	mrand "math/rand"
	"net"
	"sync"

	"github.com/v2ray/v2ray-core"
	v2io "github.com/v2ray/v2ray-core/common/io"
	"github.com/v2ray/v2ray-core/common/log"
	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/proxy/vmess/protocol"
	"github.com/v2ray/v2ray-core/proxy/vmess/protocol/user"
)

const (
	InfoTimeNotSync = "Please check the User ID in your vmess configuration, and make sure the time on your local and remote server are in sync."
)

// VNext is the next Point server in the connection chain.
type VNextServer struct {
	Destination v2net.Destination // Address of VNext server
	Users       []user.User       // User accounts for accessing VNext.
}

type VMessOutboundHandler struct {
	vPoint       *core.Point
	packet       v2net.Packet
	vNextList    []VNextServer
	vNextListUDP []VNextServer
}

func NewVMessOutboundHandler(vp *core.Point, vNextList, vNextListUDP []VNextServer, firstPacket v2net.Packet) *VMessOutboundHandler {
	return &VMessOutboundHandler{
		vPoint:       vp,
		packet:       firstPacket,
		vNextList:    vNextList,
		vNextListUDP: vNextListUDP,
	}
}

func pickVNext(serverList []VNextServer) (v2net.Destination, user.User) {
	vNextLen := len(serverList)
	if vNextLen == 0 {
		panic("VMessOut: Zero vNext is configured.")
	}
	vNextIndex := 0
	if vNextLen > 1 {
		vNextIndex = mrand.Intn(vNextLen)
	}

	vNext := serverList[vNextIndex]
	vNextUserLen := len(vNext.Users)
	if vNextUserLen == 0 {
		panic("VMessOut: Zero User account.")
	}
	vNextUserIndex := 0
	if vNextUserLen > 1 {
		vNextUserIndex = mrand.Intn(vNextUserLen)
	}
	vNextUser := vNext.Users[vNextUserIndex]
	return vNext.Destination, vNextUser
}

func (handler *VMessOutboundHandler) Start(ray core.OutboundRay) error {
	vNextList := handler.vNextList
	if handler.packet.Destination().IsUDP() {
		vNextList = handler.vNextListUDP
	}
	vNextAddress, vNextUser := pickVNext(vNextList)

	command := protocol.CmdTCP
	if handler.packet.Destination().IsUDP() {
		command = protocol.CmdUDP
	}
	request := &protocol.VMessRequest{
		Version: protocol.Version,
		UserId:  vNextUser.Id,
		Command: command,
		Address: handler.packet.Destination().Address(),
	}
	rand.Read(request.RequestIV[:])
	rand.Read(request.RequestKey[:])
	rand.Read(request.ResponseHeader[:])

	go startCommunicate(request, vNextAddress, ray, handler.packet)
	return nil
}

func startCommunicate(request *protocol.VMessRequest, dest v2net.Destination, ray core.OutboundRay, firstPacket v2net.Packet) error {
	conn, err := net.Dial(dest.Network(), dest.Address().String())
	if err != nil {
		log.Error("Failed to open %s: %v", dest.String(), err)
		if ray != nil {
			close(ray.OutboundOutput())
		}
		return err
	}
	log.Info("VMessOut: Tunneling request for %s", request.Address.String())

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

func handleRequest(conn net.Conn, request *protocol.VMessRequest, firstPacket v2net.Packet, input <-chan []byte, finish *sync.Mutex) {
	defer finish.Unlock()
	encryptRequestWriter, err := v2io.NewAesEncryptWriter(request.RequestKey[:], request.RequestIV[:], conn)
	if err != nil {
		log.Error("VMessOut: Failed to create encrypt writer: %v", err)
		return
	}

	buffer := make([]byte, 0, 2*1024)
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
		encryptRequestWriter.Crypt(firstChunk)
		buffer = append(buffer, firstChunk...)

		_, err = conn.Write(buffer)
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

func handleResponse(conn net.Conn, request *protocol.VMessRequest, output chan<- []byte, finish *sync.Mutex, isUDP bool) {
	defer finish.Unlock()
	defer close(output)
	responseKey := md5.Sum(request.RequestKey[:])
	responseIV := md5.Sum(request.RequestIV[:])

	decryptResponseReader, err := v2io.NewAesDecryptReader(responseKey[:], responseIV[:], conn)
	if err != nil {
		log.Error("VMessOut: Failed to create decrypt reader: %v", err)
		return
	}

	buffer := make([]byte, 2*1024)

	nBytes, err := decryptResponseReader.Read(buffer)
	if err != nil {
		//log.Error("VMessOut: Failed to read VMess response (%d bytes): %v", nBytes, err)
		return
	}
	if !bytes.Equal(buffer[:4], request.ResponseHeader[:]) {
		log.Warning("VMessOut: unexepcted response header. The connection is probably hijacked.")
		return
	}

	output <- buffer[4:nBytes]

	if !isUDP {
		v2net.ReaderToChan(output, decryptResponseReader)
	}

	return
}

type VMessOutboundHandlerFactory struct {
	servers    []VNextServer
	udpServers []VNextServer
}

func (factory *VMessOutboundHandlerFactory) Initialize(rawConfig []byte) error {
	config, err := loadOutboundConfig(rawConfig)
	if err != nil {
		panic(log.Error("Failed to load VMess outbound config: %v", err))
		return err
	}
	servers := make([]VNextServer, 0, len(config.VNextList))
	udpServers := make([]VNextServer, 0, len(config.VNextList))
	for _, server := range config.VNextList {
		if server.HasNetwork("tcp") {
			servers = append(servers, server.ToVNextServer())
		}
		if server.HasNetwork("udp") {
			udpServers = append(udpServers, server.ToVNextServer())
		}
	}
	factory.servers = servers
	factory.udpServers = udpServers
	return nil
}

func (factory *VMessOutboundHandlerFactory) Create(vp *core.Point, firstPacket v2net.Packet) (core.OutboundConnectionHandler, error) {
	return NewVMessOutboundHandler(vp, factory.servers, factory.udpServers, firstPacket), nil
}

func init() {
	core.RegisterOutboundConnectionHandlerFactory("vmess", &VMessOutboundHandlerFactory{})
}
