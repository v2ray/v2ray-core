package tcp

import (
	"io"
	"net"

	"github.com/v2ray/v2ray-core"
)

type VFreeConnection struct {
	network string
	address string
}

func NewVFreeConnection(network string, address string) *VFreeConnection {
	conn := new(VFreeConnection)
	conn.network = network
	conn.address = address
	return conn
}

func (vconn *VFreeConnection) Start(vRay core.OutboundVRay) error {
	input := vRay.OutboundInput()
	output := vRay.OutboundOutput()
	conn, err := net.Dial(vconn.network, vconn.address)
	if err != nil {
		return err
	}

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
			finish <- true
			break
		}
		output <- buffer[:nBytes]
	}
}

func (vconn *VFreeConnection) CloseConn(conn net.Conn, finish <-chan bool) {
	for i := 0; i < 2; i++ {
		<-finish
	}
	conn.Close()
}
