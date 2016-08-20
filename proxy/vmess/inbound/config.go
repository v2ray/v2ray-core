package inbound

import (
	"v2ray.com/core/common/protocol"
)

type DetourConfig struct {
	ToTag string
}

type FeaturesConfig struct {
	Detour *DetourConfig
}

type DefaultConfig struct {
	AlterIDs uint16
	Level    protocol.UserLevel
}

type Config struct {
	AllowedUsers []*protocol.User
	Features     *FeaturesConfig
	Defaults     *DefaultConfig
	DetourConfig *DetourConfig
}
