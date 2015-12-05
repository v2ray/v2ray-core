package config

import (
	v2net "github.com/v2ray/v2ray-core/common/net"
)

type Receiver struct {
	Address  v2net.Address
	Accounts []User
}

type Outbound interface {
	Receivers() []*Receiver
}
