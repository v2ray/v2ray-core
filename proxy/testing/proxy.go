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

func RegisterInboundConnectionHandlerCreator(prefix string, creator internal.InboundHandlerCreator) (string, error) {
	for {
		name := prefix + randomString()
		err := internal.RegisterInboundHandlerCreator(name, creator)
		if err != internal.ErrorNameExists {
			return name, err
		}
	}
}

func RegisterOutboundConnectionHandlerCreator(prefix string, creator internal.OutboundHandlerCreator) (string, error) {
	for {
		name := prefix + randomString()
		err := internal.RegisterOutboundHandlerCreator(name, creator)
		if err != internal.ErrorNameExists {
			return name, err
		}
	}
}
