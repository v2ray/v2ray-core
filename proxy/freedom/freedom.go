package freedom

import (
	"net"

	"github.com/v2ray/v2ray-core"
	"github.com/v2ray/v2ray-core/common/log"
	v2net "github.com/v2ray/v2ray-core/common/net"
)

type FreedomConnection struct {
	dest v2net.Destination
}

func NewFreedomConnection(dest v2net.Destination) *FreedomConnection {
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

	go dumpInput(conn, input, writeFinish)
	go dumpOutput(conn, output, readFinish)
	go closeConn(conn, readFinish, writeFinish)
	return nil
}

func dumpInput(conn net.Conn, input <-chan []byte, finish chan<- bool) {
	v2net.ChanToWriter(conn, input)
	close(finish)
}

func dumpOutput(conn net.Conn, output chan<- []byte, finish chan<- bool) {
	v2net.ReaderToChan(output, conn)
	close(output)
	close(finish)
}

func closeConn(conn net.Conn, readFinish <-chan bool, writeFinish <-chan bool) {
	<-writeFinish
	if tcpConn, ok := conn.(*net.TCPConn); ok {
		tcpConn.CloseWrite()
	}
	<-readFinish
	conn.Close()
}
