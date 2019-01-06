package stats_test

import (
	"context"
	"testing"

	. "v2ray.com/core/app/stats"
	"v2ray.com/core/common"
	"v2ray.com/core/features/stats"
	. "v2ray.com/ext/assert"
)

func TestInternface(t *testing.T) {
	assert := With(t)

	assert((*Manager)(nil), Implements, (*stats.Manager)(nil))
}

func TestStatsCounter(t *testing.T) {
	assert := With(t)

	raw, err := common.CreateObject(context.Background(), &Config{})
	assert(err, IsNil)

	m := raw.(stats.Manager)
	c, err := m.RegisterCounter("test.counter")
	assert(err, IsNil)

	assert(c.Add(1), Equals, int64(1))
	assert(c.Set(0), Equals, int64(1))
	assert(c.Value(), Equals, int64(0))
}
