package internet

import (
	"net"
)

type Connection interface {
	net.Conn
}

type addInt64 interface {
	Add(int64) int64
}

type StatCouterConnection struct {
	Connection
	Uplink   addInt64
	Downlink addInt64
}

func (c *StatCouterConnection) Read(b []byte) (int, error) {
	nBytes, err := c.Connection.Read(b)
	if c.Uplink != nil {
		c.Uplink.Add(int64(nBytes))
	}

	return nBytes, err
}

func (c *StatCouterConnection) Write(b []byte) (int, error) {
	nBytes, err := c.Connection.Write(b)
	if c.Downlink != nil {
		c.Downlink.Add(int64(nBytes))
	}
	return nBytes, err
}
