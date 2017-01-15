package internet

import (
	"errors"
	"net"

	"context"

	"v2ray.com/core/common"
)

type ConnectionAuthenticator interface {
	Client(net.Conn) net.Conn
	Server(net.Conn) net.Conn
}

func CreateConnectionAuthenticator(config interface{}) (ConnectionAuthenticator, error) {
	auth, err := common.CreateObject(context.Background(), config)
	if err != nil {
		return nil, err
	}
	switch a := auth.(type) {
	case ConnectionAuthenticator:
		return a, nil
	default:
		return nil, errors.New("Internet: Not a ConnectionAuthenticator.")
	}
}
