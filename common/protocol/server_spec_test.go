package protocol_test

import (
	"testing"
	"time"

	. "v2ray.com/core/common/protocol"
	. "v2ray.com/ext/assert"
)

func TestAlwaysValidStrategy(t *testing.T) {
	assert := With(t)

	strategy := AlwaysValid()
	assert(strategy.IsValid(), IsTrue)
	strategy.Invalidate()
	assert(strategy.IsValid(), IsTrue)
}

func TestTimeoutValidStrategy(t *testing.T) {
	assert := With(t)

	strategy := BeforeTime(time.Now().Add(2 * time.Second))
	assert(strategy.IsValid(), IsTrue)
	time.Sleep(3 * time.Second)
	assert(strategy.IsValid(), IsFalse)

	strategy = BeforeTime(time.Now().Add(2 * time.Second))
	strategy.Invalidate()
	assert(strategy.IsValid(), IsFalse)
}
