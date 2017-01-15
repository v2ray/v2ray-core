package internet

import (
	"context"
	"errors"

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
	switch h := header.(type) {
	case PacketHeader:
		return h, nil
	default:
		return nil, errors.New("Internet: Not a packet header.")
	}
}
