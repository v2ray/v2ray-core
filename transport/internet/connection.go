package internet

import (
	"net"

	"v2ray.com/core/features/stats"
)

type Connection interface {
	net.Conn
}

type StatCouterConnection struct {
	Connection
	Uplink   stats.Counter
	Downlink stats.Counter
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
