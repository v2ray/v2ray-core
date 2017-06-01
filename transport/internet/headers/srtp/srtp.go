package srtp

import (
	"context"

	"v2ray.com/core/common"
	"v2ray.com/core/common/dice"
	"v2ray.com/core/common/serial"
)

type SRTP struct {
	header uint16
	number uint16
}

func (*SRTP) Size() int {
	return 4
}

// Write implements io.Writer.
func (s *SRTP) Write(b []byte) (int, error) {
	s.number++
	serial.Uint16ToBytes(s.number, b[:0])
	serial.Uint16ToBytes(s.number, b[:2])
	return 4, nil
}

func NewSRTP(ctx context.Context, config interface{}) (interface{}, error) {
	return &SRTP{
		header: 0xB5E8,
		number: dice.RollUint16(),
	}, nil
}

func init() {
	common.Must(common.RegisterConfig((*Config)(nil), NewSRTP))
}
