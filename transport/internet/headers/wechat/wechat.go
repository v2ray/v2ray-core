package wechat

import (
	"v2ray.com/core/common"
	"v2ray.com/core/common/dice"
	"v2ray.com/core/common/serial"
	"v2ray.com/core/transport/internet"
)

type VideoChat struct {
	sn int
}

func (vc *VideoChat) Size() int {
	return 13
}

func (vc *VideoChat) Write(b []byte) (int, error) {
	vc.sn++
	b = append(b[:0], 0xa1, 0x08)
	b = serial.IntToBytes(vc.sn, b)
	b = append(b, 0x10, 0x11, 0x18, 0x30, 0x22, 0x30)
	return 13, nil
}

type VideoChatFactory struct{}

func (VideoChatFactory) Create(rawSettings interface{}) internet.PacketHeader {
	return &VideoChat{
		sn: dice.Roll(65535),
	}
}

func init() {
	common.Must(internet.RegisterPacketHeader(serial.GetMessageType(new(VideoConfig)), VideoChatFactory{}))
}
