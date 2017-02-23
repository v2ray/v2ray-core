package internet

import (
	"context"
	"errors"
	"net"

	"v2ray.com/core/common"
)

type PacketHeader interface {
	Size() int
	Write([]byte) (int, error)
}

func CreatePacketHeader(config interface{}) (PacketHeader, error) {
	header, err := common.CreateObject(context.Background(), config)
	if err != nil {
		return nil, err
	}
	if h, ok := header.(PacketHeader); ok {
		return h, nil
	}
	return nil, errors.New("Internet: Not a packet header.")
}

type ConnectionAuthenticator interface {
	Client(net.Conn) net.Conn
	Server(net.Conn) net.Conn
}

func CreateConnectionAuthenticator(config interface{}) (ConnectionAuthenticator, error) {
	auth, err := common.CreateObject(context.Background(), config)
	if err != nil {
		return nil, err
	}
	if a, ok := auth.(ConnectionAuthenticator); ok {
		return a, nil
	}
	return nil, errors.New("Internet: Not a ConnectionAuthenticator.")
}
