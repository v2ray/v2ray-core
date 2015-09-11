package freedom

import (
	"io"
	"net"
  "log"

	"github.com/v2ray/v2ray-core"
	v2net "github.com/v2ray/v2ray-core/net"
)

type VFreeConnection struct {
	vPoint *core.VPoint
	dest   v2net.VAddress
}

func NewVFreeConnection(vp *core.VPoint, dest v2net.VAddress) *VFreeConnection {
	conn := new(VFreeConnection)
	conn.vPoint = vp
	conn.dest = dest
	return conn
}

func (vconn *VFreeConnection) Start(vRay core.OutboundVRay) error {
	input := vRay.OutboundInput()
	output := vRay.OutboundOutput()
	conn, err := net.Dial("tcp", vconn.dest.String())
	if err != nil {
    log.Print(err)
		return err
	}
  log.Print("Working on tcp:" + vconn.dest.String())

	finish := make(chan bool, 2)
	go vconn.DumpInput(conn, input, finish)
	go vconn.DumpOutput(conn, output, finish)
	go vconn.CloseConn(conn, finish)
	return nil
}

func (vconn *VFreeConnection) DumpInput(conn net.Conn, input <-chan []byte, finish chan<- bool) {
	for {
		data, open := <-input
		if !open {
			finish <- true
			break
		}
		conn.Write(data)
	}
}

func (vconn *VFreeConnection) DumpOutput(conn net.Conn, output chan<- []byte, finish chan<- bool) {
	for {
		buffer := make([]byte, 128)
		nBytes, err := conn.Read(buffer)
		if err == io.EOF {
      close(output)
			finish <- true
			break
		}
    log.Print(buffer[:nBytes])
		output <- buffer[:nBytes]
	}
}

func (vconn *VFreeConnection) CloseConn(conn net.Conn, finish <-chan bool) {
	for i := 0; i < 2; i++ {
		<-finish
	}
	conn.Close()
}
