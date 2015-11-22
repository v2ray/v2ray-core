package json

import (
	"github.com/v2ray/v2ray-core/app/point/config"
	v2net "github.com/v2ray/v2ray-core/common/net"
)

type Rule struct {
	Type        string `json:"type"`
	OutboundTag string `json:"outboundTag"`
}

func (this *Rule) Tag() config.DetourTag {
	return config.DetourTag(this.OutboundTag)
}

func (this *Rule) Apply(dest v2net.Destination) bool {
	return false
}
