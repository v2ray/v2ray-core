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

// VNext is the next Point server in the connection chain.
type VNextServer struct {
	Address v2net.Address // Address of VNext server
	Users   []user.User   // User accounts for accessing VNext.
}

type VMessOutboundHandler struct {
	vPoint    *core.Point
	dest      v2net.Address
	vNextList []VNextServer
}

func NewVMessOutboundHandler(vp *core.Point, vNextList []VNextServer, dest v2net.Address) *VMessOutboundHandler {
	return &VMessOutboundHandler{
		vPoint:    vp,
		dest:      dest,
		vNextList: vNextList,
	}
}

func (handler *VMessOutboundHandler) pickVNext() (v2net.Address, user.User) {
	vNextLen := len(handler.vNextList)
	if vNextLen == 0 {
		panic("VMessOut: Zero vNext is configured.")
	}
	vNextIndex := mrand.Intn(vNextLen)
	vNext := handler.vNextList[vNextIndex]
	vNextUserLen := len(vNext.Users)
	if vNextUserLen == 0 {
		panic("VMessOut: Zero User account.")
	}
	vNextUserIndex := mrand.Intn(vNextUserLen)
	vNextUser := vNext.Users[vNextUserIndex]
	return vNext.Address, vNextUser
}

func (handler *VMessOutboundHandler) Start(ray core.OutboundRay) error {
	vNextAddress, vNextUser := handler.pickVNext()

	request := &protocol.VMessRequest{
		Version: protocol.Version,
		UserId:  vNextUser.Id,
		Command: byte(0x01),
		Address: handler.dest,
	}
	rand.Read(request.RequestIV[:])
	rand.Read(request.RequestKey[:])
	rand.Read(request.ResponseHeader[:])

	go startCommunicate(request, vNextAddress, ray)
	return nil
}

func startCommunicate(request *protocol.VMessRequest, dest v2net.Address, ray core.OutboundRay) error {
	input := ray.OutboundInput()
	output := ray.OutboundOutput()

	conn, err := net.DialTCP("tcp", nil, &net.TCPAddr{dest.IP, int(dest.Port), ""})
	if err != nil {
		log.Error("Failed to open tcp (%s): %v", dest.String(), err)
		close(output)
		return err
	}
	log.Info("VMessOut: Tunneling request for %s", request.Address.String())

	defer conn.Close()

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

	buffer, err := request.ToBytes(user.NewTimeHash(user.HMACHash{}), user.GenerateRandomInt64InRange)
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

func (factory *VMessOutboundHandlerFactory) Create(vp *core.Point, rawConfig []byte, destination v2net.Address) (core.OutboundConnectionHandler, error) {
	config, err := loadOutboundConfig(rawConfig)
	if err != nil {
		panic(log.Error("Failed to load VMess outbound config: %v", err))
	}
	servers := make([]VNextServer, 0, len(config.VNextList))
	for _, server := range config.VNextList {
		servers = append(servers, server.ToVNextServer())
	}
	return NewVMessOutboundHandler(vp, servers, destination), nil
}

func init() {
	core.RegisterOutboundConnectionHandlerFactory("vmess", &VMessOutboundHandlerFactory{})
}
