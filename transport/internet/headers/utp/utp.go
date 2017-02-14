package utp

import (
	"context"

	"v2ray.com/core/common"
	"v2ray.com/core/common/dice"
	"v2ray.com/core/common/serial"
)

type UTP struct {
	header       byte
	extension    byte
	connectionId uint16
}

func (v *UTP) Size() int {
	return 4
}

func (v *UTP) Write(b []byte) (int, error) {
	serial.Uint16ToBytes(v.connectionId, b[:0])
	b[2] = v.header
	b[3] = v.extension
	return 4, nil
}

func NewUTP(ctx context.Context, config interface{}) (interface{}, error) {
	return &UTP{
		header:       1,
		extension:    0,
		connectionId: dice.RandomUint16(),
	}, nil
}

func init() {
	common.Must(common.RegisterConfig((*Config)(nil), NewUTP))
}
