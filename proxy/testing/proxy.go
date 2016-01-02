package testing

import (
	"fmt"

	"github.com/v2ray/v2ray-core/proxy/internal"
)

var count = 0

func randomString() string {
	count++
	return fmt.Sprintf("-%d", count)
}

func RegisterInboundConnectionHandlerCreator(prefix string, creator internal.InboundConnectionHandlerCreator) (string, error) {
	for {
		name := prefix + randomString()
		err := internal.RegisterInboundConnectionHandlerFactory(name, creator)
		if err != internal.ErrorNameExists {
			return name, err
		}
	}
}

func RegisterOutboundConnectionHandlerCreator(prefix string, creator internal.OutboundConnectionHandlerCreator) (string, error) {
	for {
		name := prefix + randomString()
		err := internal.RegisterOutboundConnectionHandlerFactory(name, creator)
		if err != internal.ErrorNameExists {
			return name, err
		}
	}
}
