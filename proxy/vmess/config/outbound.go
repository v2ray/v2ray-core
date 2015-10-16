package config

import (
	v2net "github.com/v2ray/v2ray-core/common/net"
)

type OutboundTarget struct {
	Destination v2net.Destination
	Accounts    []User
}

type Outbound interface {
	Targets() []*OutboundTarget
}
