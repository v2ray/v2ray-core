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

	finish := make(chan bool, 2)
	go vconn.DumpInput(conn, input, finish)
	go vconn.DumpOutput(conn, output, finish)
	go vconn.CloseConn(conn, finish)
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

func (vconn *FreedomConnection) CloseConn(conn net.Conn, finish <-chan bool) {
	<-finish
	<-finish
	conn.Close()
}
