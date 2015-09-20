package freedom

import (
	"net"

	"github.com/v2ray/v2ray-core"
	"github.com/v2ray/v2ray-core/common/log"
	v2net "github.com/v2ray/v2ray-core/common/net"
)

type FreedomConnection struct {
	dest *v2net.Destination
}

func NewFreedomConnection(dest *v2net.Destination) *FreedomConnection {
	return &FreedomConnection{
		dest: dest,
	}
}

func (vconn *FreedomConnection) Start(ray core.OutboundRay) error {
	input := ray.OutboundInput()
	output := ray.OutboundOutput()
	conn, err := net.Dial(vconn.dest.Network(), vconn.dest.Address().String())
	log.Info("Freedom: Opening connection to %s", vconn.dest.String())
	if err != nil {
		close(output)
		return log.Error("Freedom: Failed to open connection: %s : %v", vconn.dest.String(), err)
	}

	readFinish := make(chan bool)
	writeFinish := make(chan bool)

	go vconn.DumpInput(conn, input, writeFinish)
	go vconn.DumpOutput(conn, output, readFinish)
	go vconn.CloseConn(conn, readFinish, writeFinish)
	return nil
}

func (vconn *FreedomConnection) DumpInput(conn net.Conn, input <-chan []byte, finish chan<- bool) {
	v2net.ChanToWriter(conn, input)
	finish <- true
}

func (vconn *FreedomConnection) DumpOutput(conn net.Conn, output chan<- []byte, finish chan<- bool) {
	v2net.ReaderToChan(output, conn)
	close(output)
	finish <- true
}

func (vconn *FreedomConnection) CloseConn(conn net.Conn, readFinish <-chan bool, writeFinish <-chan bool) {
	<-writeFinish
	if tcpConn, ok := conn.(*net.TCPConn); ok {
		tcpConn.CloseWrite()
	}
	<-readFinish
	conn.Close()
}
