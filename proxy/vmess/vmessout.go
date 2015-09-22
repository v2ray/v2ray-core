package vmess

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	mrand "math/rand"
	"net"

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
	vPoint    *core.Point
	packet    v2net.Packet
	vNextList []VNextServer
}

func NewVMessOutboundHandler(vp *core.Point, vNextList []VNextServer, firstPacket v2net.Packet) *VMessOutboundHandler {
	return &VMessOutboundHandler{
		vPoint:    vp,
		packet:    firstPacket,
		vNextList: vNextList,
	}
}

func (handler *VMessOutboundHandler) pickVNext() (v2net.Destination, user.User) {
	vNextLen := len(handler.vNextList)
	if vNextLen == 0 {
		panic("VMessOut: Zero vNext is configured.")
	}
	vNextIndex := 0
	if vNextLen > 1 {
		vNextIndex = mrand.Intn(vNextLen)
	}

	vNext := handler.vNextList[vNextIndex]
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
	vNextAddress, vNextUser := handler.pickVNext()

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
	conn, err := net.DialTCP(dest.Network(), nil, &net.TCPAddr{dest.Address().IP(), int(dest.Address().Port()), ""})
	if err != nil {
		log.Error("Failed to open tcp (%s): %v", dest.String(), err)
		if ray != nil {
			close(ray.OutboundOutput())
		}
		return err
	}
	log.Info("VMessOut: Tunneling request for %s", request.Address.String())

	defer conn.Close()

	if chunk := firstPacket.Chunk(); chunk != nil {
		conn.Write(chunk)
	}

	if !firstPacket.MoreChunks() {
		if ray != nil {
			close(ray.OutboundOutput())
		}
		return nil
	}

	input := ray.OutboundInput()
	output := ray.OutboundOutput()

	requestFinish := make(chan bool)
	responseFinish := make(chan bool)

	go handleRequest(conn, request, input, requestFinish)
	go handleResponse(conn, request, output, responseFinish)

	<-requestFinish
	conn.CloseWrite()
	<-responseFinish
	return nil
}

func handleRequest(conn *net.TCPConn, request *protocol.VMessRequest, input <-chan []byte, finish chan<- bool) {
	defer close(finish)
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
	payload, open := <-input
	if open {
		encryptRequestWriter.Crypt(payload)
		buffer = append(buffer, payload...)

		_, err = conn.Write(buffer)
		if err != nil {
			log.Error("VMessOut: Failed to write VMess request: %v", err)
			return
		}

		v2net.ChanToWriter(encryptRequestWriter, input)
	}
	return
}

func handleResponse(conn *net.TCPConn, request *protocol.VMessRequest, output chan<- []byte, finish chan<- bool) {
	defer close(finish)
	defer close(output)
	responseKey := md5.Sum(request.RequestKey[:])
	responseIV := md5.Sum(request.RequestIV[:])

	decryptResponseReader, err := v2io.NewAesDecryptReader(responseKey[:], responseIV[:], conn)
	if err != nil {
		log.Error("VMessOut: Failed to create decrypt reader: %v", err)
		return
	}

	response := protocol.VMessResponse{}
	nBytes, err := decryptResponseReader.Read(response[:])
	if err != nil {
		log.Error("VMessOut: Failed to read VMess response (%d bytes): %v", nBytes, err)
		log.Error(InfoTimeNotSync)
		return
	}
	if !bytes.Equal(response[:], request.ResponseHeader[:]) {
		log.Warning("VMessOut: unexepcted response header. The connection is probably hijacked.")
		return
	}

	v2net.ReaderToChan(output, decryptResponseReader)
	return
}

type VMessOutboundHandlerFactory struct {
}

func (factory *VMessOutboundHandlerFactory) Create(vp *core.Point, rawConfig []byte, firstPacket v2net.Packet) (core.OutboundConnectionHandler, error) {
	config, err := loadOutboundConfig(rawConfig)
	if err != nil {
		panic(log.Error("Failed to load VMess outbound config: %v", err))
	}
	servers := make([]VNextServer, 0, len(config.VNextList))
	for _, server := range config.VNextList {
		servers = append(servers, server.ToVNextServer())
	}
	return NewVMessOutboundHandler(vp, servers, firstPacket), nil
}

func init() {
	core.RegisterOutboundConnectionHandlerFactory("vmess", &VMessOutboundHandlerFactory{})
}
