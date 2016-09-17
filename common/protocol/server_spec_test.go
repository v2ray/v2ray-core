package protocol_test

import (
	"testing"
	"time"

	. "v2ray.com/core/common/protocol"
	"v2ray.com/core/testing/assert"
)

func TestAlwaysValidStrategy(t *testing.T) {
	assert := assert.On(t)

	strategy := AlwaysValid()
	assert.Bool(strategy.IsValid()).IsTrue()
	strategy.Invalidate()
	assert.Bool(strategy.IsValid()).IsTrue()
}

func TestTimeoutValidStrategy(t *testing.T) {
	assert := assert.On(t)

	strategy := BeforeTime(time.Now().Add(2 * time.Second))
	assert.Bool(strategy.IsValid()).IsTrue()
	time.Sleep(3 * time.Second)
	assert.Bool(strategy.IsValid()).IsFalse()

	strategy = BeforeTime(time.Now().Add(2 * time.Second))
	strategy.Invalidate()
	assert.Bool(strategy.IsValid()).IsFalse()
}
