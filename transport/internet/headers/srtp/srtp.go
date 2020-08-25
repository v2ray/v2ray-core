package srtp

import (
	"context"
	"encoding/binary"

	"v2ray.com/core/common"
	"v2ray.com/core/common/dice"
)

type SRTP struct {
	header uint16
	number uint16
}

func (*SRTP) Size() int32 {
	return 4
}

// Serialize implements PacketHeader.
func (s *SRTP) Serialize(b []byte) {
	s.number++
	binary.BigEndian.PutUint16(b, s.header)
	binary.BigEndian.PutUint16(b[2:], s.number)
}

// New returns a new SRTP instance based on the given config.
func New(ctx context.Context, config interface{}) (interface{}, error) {
	return &SRTP{
		header: 0xB5E8,
		number: dice.RollUint16(),
	}, nil
}

func init() {
	common.Must(common.RegisterConfig((*Config)(nil), New))
}
