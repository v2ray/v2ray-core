package json

import (
	v2net "github.com/v2ray/v2ray-core/common/net"
)

type Rule struct {
	Type        string `json:"type"`
	OutboundTag string `json:"outboundTag"`
}

func (this *Rule) Tag() string {
	return this.OutboundTag
}

func (this *Rule) Apply(dest v2net.Destination) bool {
	return false
}
