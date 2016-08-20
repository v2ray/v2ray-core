package outbound

import (
	"v2ray.com/core/common/protocol"
)

type Config struct {
	Receivers []*protocol.ServerSpec
}
