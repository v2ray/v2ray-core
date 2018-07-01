package pubsub_test

import (
	"testing"

	. "v2ray.com/core/common/signal/pubsub"
	. "v2ray.com/ext/assert"
)

func TestPubsub(t *testing.T) {
	assert := With(t)

	service := NewService()

	sub := service.Subscribe("a")
	service.Publish("a", 1)

	select {
	case v := <-sub.Wait():
		assert(v.(int), Equals, 1)
	default:
		t.Fail()
	}

	sub.Close()
	service.Publish("a", 2)

	select {
	case <-sub.Wait():
		t.Fail()
	default:
	}

	service.Cleanup()
}
