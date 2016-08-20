package dns

import (
	"net"

	v2net "v2ray.com/core/common/net"
)

type Config struct {
	Hosts       map[string]net.IP
	NameServers []v2net.Destination
}
