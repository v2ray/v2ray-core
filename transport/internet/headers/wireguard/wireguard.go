package wireguard

import (
	"context"

	"v2ray.com/core/common"
)

type Wireguard struct{}

func (Wireguard) Size() int32 {
	return 4
}

// Write implements io.Writer.
func (Wireguard) Write(b []byte) (int, error) {
	b = append(b[:0], 0x04, 0x00, 0x00, 0x00)
	return 4, nil
}

// NewWireguard returns a new VideoChat instance based on given config.
func NewWireguard(ctx context.Context, config interface{}) (interface{}, error) {
	return Wireguard{}, nil
}

func init() {
	common.Must(common.RegisterConfig((*WireguardConfig)(nil), NewWireguard))
}
