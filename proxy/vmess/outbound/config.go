package outbound

import (
	"github.com/v2ray/v2ray-core/common/protocol"
)

type Config struct {
	Receivers []*protocol.ServerSpec
}
