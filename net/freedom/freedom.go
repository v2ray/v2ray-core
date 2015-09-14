package freedom

import (
	"net"

	"github.com/v2ray/v2ray-core"
	"github.com/v2ray/v2ray-core/log"
	v2net "github.com/v2ray/v2ray-core/net"
)

type FreedomConnection struct {
	dest v2net.Address
}

func NewFreedomConnection(dest v2net.Address) *FreedomConnection {
	conn := new(FreedomConnection)
	conn.dest = dest
	return conn
}

func (vconn *FreedomConnection) Start(ray core.OutboundRay) error {
	input := ray.OutboundInput()
	output := ray.OutboundOutput()
	conn, err := net.Dial("tcp", vconn.dest.String())
	if err != nil {
		return log.Error("Failed to open tcp: %s", vconn.dest.String())
	}
	log.Debug("Sending outbound tcp: %s", vconn.dest.String())

	readFinish := make(chan bool)
	writeFinish := make(chan bool)

	go vconn.DumpInput(conn, input, writeFinish)
	go vconn.DumpOutput(conn, output, readFinish)
	go vconn.CloseConn(conn, readFinish, writeFinish)
	return nil
}

func (vconn *FreedomConnection) DumpInput(conn net.Conn, input <-chan []byte, finish chan<- bool) {
	v2net.ChanToWriter(conn, input)
	log.Debug("Freedom closing input")
	finish <- true
}

func (vconn *FreedomConnection) DumpOutput(conn net.Conn, output chan<- []byte, finish chan<- bool) {
	v2net.ReaderToChan(output, conn)
	close(output)
	log.Debug("Freedom closing output")
	finish <- true
}

func (vconn *FreedomConnection) CloseConn(conn net.Conn, readFinish <-chan bool, writeFinish <-chan bool) {
	<-writeFinish
	if tcpConn, ok := conn.(*net.TCPConn); ok {
		log.Debug("Closing freedom write.")
		tcpConn.CloseWrite()
	}
	<-readFinish
	conn.Close()
}
