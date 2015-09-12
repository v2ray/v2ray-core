package vmess

import (
	"crypto/md5"
	"crypto/rand"
	"io"
	mrand "math/rand"
	"net"

	"github.com/v2ray/v2ray-core"
	v2io "github.com/v2ray/v2ray-core/io"
	vmessio "github.com/v2ray/v2ray-core/io/vmess"
	"github.com/v2ray/v2ray-core/log"
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
	handler := new(VMessOutboundHandler)
	handler.vPoint = vp
	handler.dest = dest
	handler.vNextList = vNextList
	return handler
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

	request := new(vmessio.VMessRequest)
	request.Version = vmessio.Version
	request.UserId = vNextUser.Id
	rand.Read(request.RequestIV[:])
	rand.Read(request.RequestKey[:])
	rand.Read(request.ResponseHeader[:])
	request.Command = byte(0x01)
	request.Address = handler.dest

	go handler.startCommunicate(request, vNextAddress, ray)
	return nil
}

func (handler *VMessOutboundHandler) startCommunicate(request *vmessio.VMessRequest, dest v2net.Address, ray core.OutboundRay) error {
	conn, err := net.Dial("tcp", dest.String())
	log.Debug("VMessOutbound dialing tcp: %s", dest.String())
	if err != nil {
		log.Error("Failed to open tcp (%s): %v", dest.String(), err)
		return err
	}
	defer conn.Close()

	requestWriter := vmessio.NewVMessRequestWriter()
	err = requestWriter.Write(conn, request)
	if err != nil {
		log.Error("Failed to write VMess request: %v", err)
		return err
	}

	requestKey := request.RequestKey[:]
	requestIV := request.RequestIV[:]
	responseKey := md5.Sum(requestKey)
	responseIV := md5.Sum(requestIV)

	response := vmessio.VMessResponse{}
	nBytes, err := conn.Read(response[:])
	if err != nil {
		log.Error("Failed to read VMess response (%d bytes): %v", nBytes, err)
		return err
	}
	log.Debug("Got response %v", response)
	// TODO: check response

	encryptRequestWriter, err := v2io.NewAesEncryptWriter(requestKey, requestIV, conn)
	if err != nil {
		log.Error("Failed to create encrypt writer: %v", err)
		return err
	}
	decryptResponseReader, err := v2io.NewAesDecryptReader(responseKey[:], responseIV[:], conn)
	if err != nil {
		log.Error("Failed to create decrypt reader: %v", err)
		return err
	}

	input := ray.OutboundInput()
	output := ray.OutboundOutput()
	finish := make(chan bool, 2)

	go handler.dumpInput(encryptRequestWriter, input, finish)
	go handler.dumpOutput(decryptResponseReader, output, finish)
	handler.waitForFinish(finish)
	return nil
}

func (handler *VMessOutboundHandler) dumpOutput(reader io.Reader, output chan<- []byte, finish chan<- bool) {
	for {
		buffer := make([]byte, BufferSize)
		nBytes, err := reader.Read(buffer)
		log.Debug("VMessOutbound: Reading %d bytes, with error %v", nBytes, err)
		if err == io.EOF {
			close(output)
			finish <- true
			log.Debug("VMessOutbound finishing output.")
			break
		}
		output <- buffer[:nBytes]
	}
}

func (handler *VMessOutboundHandler) dumpInput(writer io.Writer, input <-chan []byte, finish chan<- bool) {
	for {
		buffer, open := <-input
		if !open {
			finish <- true
			log.Debug("VMessOutbound finishing input.")
			break
		}
		nBytes, err := writer.Write(buffer)
		log.Debug("VMessOutbound: Wrote %d bytes with error %v", nBytes, err)
	}
}

func (handler *VMessOutboundHandler) waitForFinish(finish <-chan bool) {
	<-finish
	<-finish
	log.Debug("Finishing waiting for VMessOutbound ending.")
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
