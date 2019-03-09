package outbound_test

import (
	"testing"

	. "v2ray.com/core/app/proxyman/outbound"
	"v2ray.com/core/features/outbound"
)

func TestInterfaces(t *testing.T) {
	_ = (outbound.Handler)(new(Handler))
	_ = (outbound.Manager)(new(Manager))
}
