package srtp

import (
	"context"
	"math/rand"

	"v2ray.com/core/common"
	"v2ray.com/core/common/serial"
)

type SRTP struct {
	header uint16
	number uint16
}

func (v *SRTP) Size() int {
	return 4
}

func (v *SRTP) Write(b []byte) (int, error) {
	v.number++
	serial.Uint16ToBytes(v.number, b[:0])
	serial.Uint16ToBytes(v.number, b[:2])
	return 4, nil
}

func NewSRTP(ctx context.Context, config interface{}) (interface{}, error) {
	return &SRTP{
		header: 0xB5E8,
		number: uint16(rand.Intn(65536)),
	}, nil
}

func init() {
	common.Must(common.RegisterConfig((*Config)(nil), NewSRTP))
}
