package signal_test

import (
	"testing"
	"time"

	"v2ray.com/core/common"
	. "v2ray.com/core/common/signal"
	. "v2ray.com/ext/assert"
)

func TestPeriodicTaskStop(t *testing.T) {
	assert := With(t)

	value := 0
	task := &PeriodicTask{
		Interval: time.Second * 2,
		Execute: func() error {
			value++
			return nil
		},
	}
	common.Must(task.Start())
	time.Sleep(time.Second * 5)
	common.Must(task.Close())
	assert(value, Equals, 3)
	time.Sleep(time.Second * 4)
	assert(value, Equals, 3)
}
