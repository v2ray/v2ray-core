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
	v2net "github.com/v2ray/v2ray-core/net"
)

type VMessOutboundHandler struct {
	vPoint *core.VPoint
	dest   v2net.VAddress
}

func NewVMessOutboundHandler(vp *core.VPoint, dest v2net.VAddress) *VMessOutboundHandler {
	handler := new(VMessOutboundHandler)
	handler.vPoint = vp
	handler.dest = dest
	return handler
}

func (handler *VMessOutboundHandler) pickVNext() (v2net.VAddress, core.VUser) {
	vNextLen := len(handler.vPoint.Config.VNextList)
	if vNextLen == 0 {
		panic("Zero vNext is configured.")
	}
	vNextIndex := mrand.Intn(vNextLen)
	vNext := handler.vPoint.Config.VNextList[vNextIndex]
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
	request.SetVersion(vmessio.Version)
	copy(request.UserHash(), vNextUser.Id.Hash([]byte("ASK")))
	rand.Read(request.RequestIV())
	rand.Read(request.RequestKey())
	rand.Read(request.ResponseHeader())
	request.SetCommand(byte(0x01))
	request.SetPort(handler.dest.Port)

	address := handler.dest
	switch {
	case address.IsIPv4():
		request.SetIPv4(address.IP)
	case address.IsIPv6():
		request.SetIPv6(address.IP)
	case address.IsDomain():
		request.SetDomain(address.Domain)
	}

	conn, err := net.Dial("tcp", vNextAddress.String())
	if err != nil {
		return err
	}
	defer conn.Close()

	requestWriter := vmessio.NewVMessRequestWriter(handler.vPoint.UserSet)
	requestWriter.Write(conn, request)

	requestKey := request.RequestKey()
	requestIV := request.RequestIV()
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

func (factory *VMessOutboundHandlerFactory) Create(vp *core.VPoint, destination v2net.VAddress) *VMessOutboundHandler {
	return NewVMessOutboundHandler(vp, destination)
}
