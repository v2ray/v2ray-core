package testing

import (
	"fmt"

	"github.com/v2ray/v2ray-core/proxy/registry"
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
		if err != registry.ErrNameExists {
			return name, err
		}
	}
}

func RegisterOutboundConnectionHandlerCreator(prefix string, creator registry.OutboundHandlerFactory) (string, error) {
	for {
		name := prefix + randomString()
		err := registry.RegisterOutboundHandlerCreator(name, creator)
		if err != registry.ErrNameExists {
			return name, err
		}
	}
}
