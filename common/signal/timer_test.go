package signal_test

import (
	"context"
	"runtime"
	"testing"
	"time"

	. "v2ray.com/core/common/signal"
	"v2ray.com/core/testing/assert"
)

func TestActivityTimer(t *testing.T) {
	assert := assert.On(t)

	ctx, timer := CancelAfterInactivity(context.Background(), time.Second*5)
	time.Sleep(time.Second * 6)
	assert.Error(ctx.Err()).IsNotNil()
	runtime.KeepAlive(timer)
}

func TestActivityTimerUpdate(t *testing.T) {
	assert := assert.On(t)

	ctx, timer := CancelAfterInactivity(context.Background(), time.Second*10)
	time.Sleep(time.Second * 3)
	assert.Error(ctx.Err()).IsNil()
	timer.SetTimeout(time.Second * 1)
	time.Sleep(time.Second * 2)
	assert.Error(ctx.Err()).IsNotNil()
	runtime.KeepAlive(timer)
}
