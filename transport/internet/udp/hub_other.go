// +build !linux

package udp

import (
	v2net "github.com/v2ray/v2ray-core/common/net"
)

func SetOriginalDestOptions(fd int) error {
	return nil
}

func RetrieveOriginalDest(oob []byte) v2net.Destination {
	return nil
}
