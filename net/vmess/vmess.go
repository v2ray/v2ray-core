package vemss

import (
	"net"
)

type VMessHandler struct {
}

func (*VMessHandler) Listen(port uint8) error {
	_, err := net.Listen("tcp", ":"+string(port))
	if err != nil {
		return err
	}

	return nil
}
