package inbound

import (
	proto "github.com/v2ray/v2ray-core/common/protocol"
)

type DetourConfig struct {
	ToTag string
}

type FeaturesConfig struct {
	Detour *DetourConfig
}

type DefaultConfig struct {
	AlterIDs uint16
	Level    proto.UserLevel
}

type Config struct {
	AllowedUsers []*proto.User
	Features     *FeaturesConfig
	Defaults     *DefaultConfig
	DetourConfig *DetourConfig
}
