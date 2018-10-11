package outbound_test

import (
	"testing"

	. "v2ray.com/core/app/proxyman/outbound"
	"v2ray.com/core/features/outbound"
	. "v2ray.com/ext/assert"
)

func TestInterfaces(t *testing.T) {
	assert := With(t)

	assert((*Handler)(nil), Implements, (*outbound.Handler)(nil))
	assert((*Manager)(nil), Implements, (*outbound.HandlerManager)(nil))
}
