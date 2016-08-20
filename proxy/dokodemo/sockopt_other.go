// +build !linux

package dokodemo

import (
	v2net "v2ray.com/core/common/net"
	"v2ray.com/core/transport/internet"
)

func GetOriginalDestination(conn internet.Connection) v2net.Destination {
	return nil
}
