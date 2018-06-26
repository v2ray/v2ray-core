package signal_test

import (
	"testing"

	. "v2ray.com/core/common/signal"
	//. "v2ray.com/ext/assert"
)

func TestNotifierSignal(t *testing.T) {
	//assert := With(t)

	var n Notifier

	w := n.Wait()
	n.Signal()

	select {
	case <-w:
	default:
		t.Fail()
	}
}
