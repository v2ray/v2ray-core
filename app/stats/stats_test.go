package stats_test

import (
	"testing"

	"v2ray.com/core"
	. "v2ray.com/core/app/stats"
	. "v2ray.com/ext/assert"
)

func TestInternface(t *testing.T) {
	assert := With(t)

	assert((*Manager)(nil), Implements, (*core.StatManager)(nil))
}
