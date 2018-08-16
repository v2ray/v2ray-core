package signal_test

import (
	"context"
	"runtime"
	"testing"
	"time"

	. "v2ray.com/core/common/signal"
	. "v2ray.com/ext/assert"
)

func TestActivityTimer(t *testing.T) {
	assert := With(t)

	ctx, cancel := context.WithCancel(context.Background())
	timer := CancelAfterInactivity(ctx, cancel, time.Second*5)
	time.Sleep(time.Second * 6)
	assert(ctx.Err(), IsNotNil)
	runtime.KeepAlive(timer)
}

func TestActivityTimerUpdate(t *testing.T) {
	assert := With(t)

	ctx, cancel := context.WithCancel(context.Background())
	timer := CancelAfterInactivity(ctx, cancel, time.Second*10)
	time.Sleep(time.Second * 3)
	assert(ctx.Err(), IsNil)
	timer.SetTimeout(time.Second * 1)
	time.Sleep(time.Second * 2)
	assert(ctx.Err(), IsNotNil)
	runtime.KeepAlive(timer)
}

func TestActivityTimerNonBlocking(t *testing.T) {
	assert := With(t)

	ctx, cancel := context.WithCancel(context.Background())
	timer := CancelAfterInactivity(ctx, cancel, 0)
	time.Sleep(time.Second * 1)
	assert(ctx, HasDone)
	timer.SetTimeout(0)
	timer.SetTimeout(1)
	timer.SetTimeout(2)
}

func TestActivityTimerZeroTimeout(t *testing.T) {
	assert := With(t)

	ctx, cancel := context.WithCancel(context.Background())
	timer := CancelAfterInactivity(ctx, cancel, 0)
	assert(ctx, HasDone)
	runtime.KeepAlive(timer)
}
