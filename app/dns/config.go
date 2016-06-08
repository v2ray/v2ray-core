package dns

import (
	"net"

	v2net "github.com/v2ray/v2ray-core/common/net"
)

type Config struct {
	Hosts       map[string]net.IP
	NameServers []v2net.Destination
}
