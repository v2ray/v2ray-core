package task_test

import (
	"testing"
	"time"

	. "v2ray.com/core/common/task"
	. "v2ray.com/ext/assert"

	"v2ray.com/core/common"
)

func TestPeriodicTaskStop(t *testing.T) {
	assert := With(t)

	value := 0
	task := &Periodic{
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
	common.Must(task.Start())
	time.Sleep(time.Second * 3)
	if value != 5 {
		t.Fatal("Expected 5, but ", value)
	}
	common.Must(task.Close())
}
