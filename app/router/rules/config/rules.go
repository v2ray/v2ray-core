package config

import (
	"github.com/v2ray/v2ray-core/app/point/config"
	v2net "github.com/v2ray/v2ray-core/common/net"
)

type Rule interface {
	Tag() config.DetourTag
	Apply(dest v2net.Destination) bool
}
