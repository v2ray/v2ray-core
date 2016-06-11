// +build !linux

package dokodemo

import (
	v2net "github.com/v2ray/v2ray-core/common/net"
	"github.com/v2ray/v2ray-core/transport/hub"
)

func GetOriginalDestination(conn *hub.Connection) v2net.Destination {
	return nil
}
