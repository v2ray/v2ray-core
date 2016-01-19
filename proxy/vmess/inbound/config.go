package inbound

import (
	"github.com/v2ray/v2ray-core/proxy/vmess"
)

type DetourConfig struct {
	ToTag string
}

type FeaturesConfig struct {
	Detour *DetourConfig
}

type Config struct {
	AllowedUsers []*vmess.User
	Features     *FeaturesConfig
}
