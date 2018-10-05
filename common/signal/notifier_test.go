package signal_test

import (
	"testing"

	. "v2ray.com/core/common/signal"
	//. "v2ray.com/ext/assert"
)

func TestNotifierSignal(t *testing.T) {
	//assert := With(t)

	n := NewNotifier()

	w := n.Wait()
	n.Signal()

	select {
	case <-w:
	default:
		t.Fail()
	}
}
