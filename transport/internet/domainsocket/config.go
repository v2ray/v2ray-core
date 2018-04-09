package domainsocket

import "net"

func (c *Config) GetUnixAddr() (*net.UnixAddr, error) {
	path := c.Path
	if len(path) == 0 {
		return nil, newError("empty domain socket path")
	}
	if c.Abstract && path[0] != '\x00' {
		path = "\x00" + path
	}
	return &net.UnixAddr{
		Name: path,
		Net:  "unix",
	}, nil
}
