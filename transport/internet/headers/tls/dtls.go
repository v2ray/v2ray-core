package tls

import (
	"context"

	"v2ray.com/core/common"
	"v2ray.com/core/common/dice"
)

// DTLS writes header as DTLS. See https://tools.ietf.org/html/rfc6347
type DTLS struct {
	epoch    uint16
	sequence uint32
	length   uint16
}

// Size implements PacketHeader.
func (*DTLS) Size() int32 {
	return 1 + 2 + 2 + 3 + 2
}

// Write implements PacketHeader.
func (d *DTLS) Write(b []byte) (int, error) {
	b[0] = 23 // application data
	b[1] = 254
	b[2] = 253
	b[3] = byte(d.epoch >> 8)
	b[4] = byte(d.epoch)
	b[5] = byte(d.sequence >> 16)
	b[6] = byte(d.sequence >> 8)
	b[7] = byte(d.sequence)
	d.sequence++
	b[8] = byte(d.length >> 8)
	b[9] = byte(d.length)
	d.length += 17
	if d.length > 1024 {
		d.length -= 1024
	}
	return 10, nil
}

// New creates a new UTP header for the given config.
func New(ctx context.Context, config interface{}) (interface{}, error) {
	return &DTLS{
		epoch:    dice.RollUint16(),
		sequence: 0,
		length:   uint16(dice.Roll(1024) + 100),
	}, nil
}

func init() {
	common.Must(common.RegisterConfig((*PacketConfig)(nil), New))
}
