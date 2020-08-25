package signal_test

import (
	"context"
	"runtime"
	"testing"
	"time"

	. "v2ray.com/core/common/signal"
)

func TestActivityTimer(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	timer := CancelAfterInactivity(ctx, cancel, time.Second*4)
	time.Sleep(time.Second * 6)
	if ctx.Err() == nil {
		t.Error("expected some error, but got nil")
	}
	runtime.KeepAlive(timer)
}

func TestActivityTimerUpdate(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	timer := CancelAfterInactivity(ctx, cancel, time.Second*10)
	time.Sleep(time.Second * 3)
	if ctx.Err() != nil {
		t.Error("expected nil, but got ", ctx.Err().Error())
	}
	timer.SetTimeout(time.Second * 1)
	time.Sleep(time.Second * 2)
	if ctx.Err() == nil {
		t.Error("expcted some error, but got nil")
	}
	runtime.KeepAlive(timer)
}

func TestActivityTimerNonBlocking(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	timer := CancelAfterInactivity(ctx, cancel, 0)
	time.Sleep(time.Second * 1)
	select {
	case <-ctx.Done():
	default:
		t.Error("context not done")
	}
	timer.SetTimeout(0)
	timer.SetTimeout(1)
	timer.SetTimeout(2)
}

func TestActivityTimerZeroTimeout(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	timer := CancelAfterInactivity(ctx, cancel, 0)
	select {
	case <-ctx.Done():
	default:
		t.Error("context not done")
	}
	runtime.KeepAlive(timer)
}
