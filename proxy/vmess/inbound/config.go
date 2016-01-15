package inbound

import (
	"github.com/v2ray/v2ray-core/proxy/vmess"
)

type Config struct {
	AllowedUsers []*vmess.User
}
