package outbound_test

import (
	"testing"

	"v2ray.com/core"
	. "v2ray.com/core/app/proxyman/outbound"
	. "v2ray.com/ext/assert"
)

func TestInterfaces(t *testing.T) {
	assert := With(t)

	assert((*Handler)(nil), Implements, (*core.OutboundHandler)(nil))
}
