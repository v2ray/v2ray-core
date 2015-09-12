package freedom

import (
	"io"
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
	for {
		data, open := <-input
		if !open {
			finish <- true
			log.Debug("Freedom finishing input.")
			break
		}
		nBytes, err := conn.Write(data)
		log.Debug("Freedom wrote %d bytes with error %v", nBytes, err)
	}
}

func (vconn *FreedomConnection) DumpOutput(conn net.Conn, output chan<- []byte, finish chan<- bool) {
	for {
		buffer := make([]byte, 512)
		nBytes, err := conn.Read(buffer)
		log.Debug("Freedom reading %d bytes with error %v", nBytes, err)
		if err == io.EOF {
			close(output)
			finish <- true
			log.Debug("Freedom finishing output.")
			break
		}
		output <- buffer[:nBytes]
	}
}

func (vconn *FreedomConnection) CloseConn(conn net.Conn, finish <-chan bool) {
	<-finish
	<-finish
	conn.Close()
}
