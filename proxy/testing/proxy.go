package testing

import (
	"fmt"

	"v2ray.com/core/common"
	"v2ray.com/core/proxy/registry"
)

var count = 0

func randomString() string {
	count++
	return fmt.Sprintf("-%d", count)
}

func RegisterInboundConnectionHandlerCreator(prefix string, creator registry.InboundHandlerFactory) (string, error) {
	for {
		name := prefix + randomString()
		err := registry.RegisterInboundHandlerCreator(name, creator)
		if err != common.ErrDuplicatedName {
			return name, err
		}
	}
}

func RegisterOutboundConnectionHandlerCreator(prefix string, creator registry.OutboundHandlerFactory) (string, error) {
	for {
		name := prefix + randomString()
		err := registry.RegisterOutboundHandlerCreator(name, creator)
		if err != common.ErrDuplicatedName {
			return name, err
		}
	}
}
