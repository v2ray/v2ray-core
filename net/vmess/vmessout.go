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

// VNext is the next VPoint server in the connection chain.
type VNextServer struct {
	Address v2net.VAddress // Address of VNext server
	Users   []core.VUser   // User accounts for accessing VNext.
}

type VMessOutboundHandler struct {
	vPoint    *core.VPoint
	dest      v2net.VAddress
	vNextList []VNextServer
}

func NewVMessOutboundHandler(vp *core.VPoint, vNextList []VNextServer, dest v2net.VAddress) *VMessOutboundHandler {
	handler := new(VMessOutboundHandler)
	handler.vPoint = vp
	handler.dest = dest
	handler.vNextList = vNextList
	return handler
}

func (handler *VMessOutboundHandler) pickVNext() (v2net.VAddress, core.VUser) {
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

func (handler *VMessOutboundHandler) Start(ray core.OutboundVRay) error {
	vNextAddress, vNextUser := handler.pickVNext()

	request := new(vmessio.VMessRequest)
	request.Version = vmessio.Version
	request.UserId = vNextUser.Id
	rand.Read(request.RequestIV[:])
	rand.Read(request.RequestKey[:])
	rand.Read(request.ResponseHeader[:])
	request.Command = byte(0x01)
	request.Address = handler.dest

	conn, err := net.Dial("tcp", vNextAddress.String())
	if err != nil {
		return err
	}
	defer conn.Close()

	requestWriter := vmessio.NewVMessRequestWriter()
	requestWriter.Write(conn, request)

	requestKey := request.RequestKey[:]
	requestIV := request.RequestIV[:]
	responseKey := md5.Sum(requestKey)
	responseIV := md5.Sum(requestIV)

	encryptRequestWriter, err := v2io.NewAesEncryptWriter(requestKey, requestIV, conn)
	if err != nil {
		return err
	}
	responseReader, err := v2io.NewAesDecryptReader(responseKey[:], responseIV[:], conn)
	if err != nil {
		return err
	}

	input := ray.OutboundInput()
	output := ray.OutboundOutput()
	finish := make(chan bool, 2)

	go handler.dumpInput(encryptRequestWriter, input, finish)
	go handler.dumpOutput(responseReader, output, finish)
	handler.waitForFinish(finish)
	return nil
}

func (handler *VMessOutboundHandler) dumpOutput(reader io.Reader, output chan<- []byte, finish chan<- bool) {
	for {
		buffer := make([]byte, BufferSize)
		nBytes, err := reader.Read(buffer)
		if err == io.EOF {
			close(output)
			finish <- true
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
			break
		}
		writer.Write(buffer)
	}
}

func (handler *VMessOutboundHandler) waitForFinish(finish <-chan bool) {
	for i := 0; i < 2; i++ {
		<-finish
	}
}

type VMessOutboundHandlerFactory struct {
}

func (factory *VMessOutboundHandlerFactory) Create(vp *core.VPoint, rawConfig []byte, destination v2net.VAddress) *VMessOutboundHandler {
	config, err := loadOutboundConfig(rawConfig)
	if err != nil {
		panic(log.Error("Failed to load VMess outbound config: %v", err))
	}
	servers := make([]VNextServer, 0, len(config.VNextList))
	for _, server := range config.VNextList {
		servers = append(servers, server.ToVNextServer())
	}
	return NewVMessOutboundHandler(vp, servers, destination)
}
