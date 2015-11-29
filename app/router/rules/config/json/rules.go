package json

import (
	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/shell/point/config"
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
