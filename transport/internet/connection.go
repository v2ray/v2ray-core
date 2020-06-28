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
	ReadCounter  stats.Counter
	WriteCounter stats.Counter
}

func (c *StatCouterConnection) Read(b []byte) (int, error) {
	nBytes, err := c.Connection.Read(b)
	if c.ReadCounter != nil {
		c.ReadCounter.Add(int64(nBytes))
	}

	return nBytes, err
}

func (c *StatCouterConnection) Write(b []byte) (int, error) {
	nBytes, err := c.Connection.Write(b)
	if c.WriteCounter != nil {
		c.WriteCounter.Add(int64(nBytes))
	}
	return nBytes, err
}
