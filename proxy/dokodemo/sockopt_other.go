// +build !linux

package dokodemo

import (
	"v2ray.com/core/common/net"
	"v2ray.com/core/transport/internet"
)

func GetOriginalDestination(conn internet.Connection) net.Destination {
	return net.Destination{}
}
