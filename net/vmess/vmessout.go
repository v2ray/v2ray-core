package vmess

import (
	"crypto/md5"
	"crypto/rand"
	mrand "math/rand"
	"net"

	"github.com/v2ray/v2ray-core"
	v2hash "github.com/v2ray/v2ray-core/hash"
	v2io "github.com/v2ray/v2ray-core/io"
	vmessio "github.com/v2ray/v2ray-core/io/vmess"
	"github.com/v2ray/v2ray-core/log"
	v2math "github.com/v2ray/v2ray-core/math"
	v2net "github.com/v2ray/v2ray-core/net"
)

// VNext is the next Point server in the connection chain.
type VNextServer struct {
	Address v2net.Address // Address of VNext server
	Users   []core.User   // User accounts for accessing VNext.
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

func (handler *VMessOutboundHandler) pickVNext() (v2net.Address, core.User) {
	vNextLen := len(handler.vNextList)
	if vNextLen == 0 {
		panic("Zero vNext is configured.")
	}
	vNextIndex := mrand.Intn(vNextLen)
	vNext := handler.vNextList[vNextIndex]
	vNextUserLen := len(vNext.Users)
	if vNextUserLen == 0 {
		panic("Zero User account.")
	}
	vNextUserIndex := mrand.Intn(vNextUserLen)
	vNextUser := vNext.Users[vNextUserIndex]
	return vNext.Address, vNextUser
}

func (handler *VMessOutboundHandler) Start(ray core.OutboundRay) error {
	vNextAddress, vNextUser := handler.pickVNext()

	request := &vmessio.VMessRequest{
		Version: vmessio.Version,
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

func startCommunicate(request *vmessio.VMessRequest, dest v2net.Address, ray core.OutboundRay) error {
	input := ray.OutboundInput()
	output := ray.OutboundOutput()

	conn, err := net.DialTCP("tcp", nil, &net.TCPAddr{dest.IP, int(dest.Port), ""})
	log.Debug("VMessOutbound dialing tcp: %s", dest.String())
	if err != nil {
		log.Error("Failed to open tcp (%s): %v", dest.String(), err)
		close(output)
		return err
	}

	log.Debug("VMessOut: Tunneling request for %s", request.Address.String())

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

func handleRequest(conn *net.TCPConn, request *vmessio.VMessRequest, input <-chan []byte, finish chan<- bool) error {
	defer close(finish)
	requestWriter := vmessio.NewVMessRequestWriter(v2hash.NewTimeHash(v2hash.HMACHash{}), v2math.GenerateRandomInt64InRange)
	err := requestWriter.Write(conn, request)
	if err != nil {
		log.Error("Failed to write VMess request: %v", err)
		return err
	}

	encryptRequestWriter, err := v2io.NewAesEncryptWriter(request.RequestKey[:], request.RequestIV[:], conn)
	if err != nil {
		log.Error("Failed to create encrypt writer: %v", err)
		return err
	}

	v2net.ChanToWriter(encryptRequestWriter, input)
	return nil
}

func handleResponse(conn *net.TCPConn, request *vmessio.VMessRequest, output chan<- []byte, finish chan<- bool) error {
	defer close(finish)
	defer close(output)
	responseKey := md5.Sum(request.RequestKey[:])
	responseIV := md5.Sum(request.RequestIV[:])

	decryptResponseReader, err := v2io.NewAesDecryptReader(responseKey[:], responseIV[:], conn)
	if err != nil {
		log.Error("Failed to create decrypt reader: %v", err)
		return err
	}

	response := vmessio.VMessResponse{}
	nBytes, err := decryptResponseReader.Read(response[:])
	if err != nil {
		log.Error("Failed to read VMess response (%d bytes): %v", nBytes, err)
		return err
	}
	log.Debug("Got response %v", response)
	// TODO: check response

	v2net.ReaderToChan(output, decryptResponseReader)
	return nil
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
