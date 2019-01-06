package wechat

import (
	"context"
	"encoding/binary"

	"v2ray.com/core/common"
	"v2ray.com/core/common/dice"
)

type VideoChat struct {
	sn uint32
}

func (vc *VideoChat) Size() int32 {
	return 13
}

// Serialize implements PacketHeader.
func (vc *VideoChat) Serialize(b []byte) {
	vc.sn++
	b[0] = 0xa1
	b[1] = 0x08
	binary.BigEndian.PutUint32(b[2:], vc.sn) // b[2:6]
	b[6] = 0x00
	b[7] = 0x10
	b[8] = 0x11
	b[9] = 0x18
	b[10] = 0x30
	b[11] = 0x22
	b[12] = 0x30
}

// NewVideoChat returns a new VideoChat instance based on given config.
func NewVideoChat(ctx context.Context, config interface{}) (interface{}, error) {
	return &VideoChat{
		sn: uint32(dice.RollUint16()),
	}, nil
}

func init() {
	common.Must(common.RegisterConfig((*VideoConfig)(nil), NewVideoChat))
}
