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

type Config struct {
	AllowedUsers []*proto.User
	Features     *FeaturesConfig
}
