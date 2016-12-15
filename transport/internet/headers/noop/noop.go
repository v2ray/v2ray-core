package noop

import (
	"net"

	"v2ray.com/core/common/serial"
	"v2ray.com/core/transport/internet"
)

type NoOpHeader struct{}

func (v NoOpHeader) Size() int {
	return 0
}
func (v NoOpHeader) Write([]byte) (int, error) {
	return 0, nil
}

type NoOpHeaderFactory struct{}

func (v NoOpHeaderFactory) Create(config interface{}) internet.PacketHeader {
	return NoOpHeader{}
}

type NoOpConnectionHeader struct{}

func (NoOpConnectionHeader) Client(conn net.Conn) net.Conn {
	return conn
}

func (NoOpConnectionHeader) Server(conn net.Conn) net.Conn {
	return conn
}

type NoOpConnectionHeaderFactory struct{}

func (NoOpConnectionHeaderFactory) Create(config interface{}) internet.ConnectionAuthenticator {
	return NoOpConnectionHeader{}
}

func init() {
	internet.RegisterPacketHeader(serial.GetMessageType(new(Config)), NoOpHeaderFactory{})
	internet.RegisterConnectionAuthenticator(serial.GetMessageType(new(Config)), NoOpConnectionHeaderFactory{})
}
