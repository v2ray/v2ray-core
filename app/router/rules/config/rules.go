package config

import (
	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/shell/point/config"
)

type Rule interface {
	Tag() config.DetourTag
	Apply(dest v2net.Destination) bool
}
