package signal_test

import (
	"context"
	"errors"
	"testing"
	"time"

	. "v2ray.com/core/common/signal"
	. "v2ray.com/ext/assert"
)

func TestExecuteParallel(t *testing.T) {
	assert := With(t)

	err := ExecuteParallel(context.Background(), func() error {
		time.Sleep(time.Millisecond * 200)
		return errors.New("test")
	}, func() error {
		time.Sleep(time.Millisecond * 500)
		return errors.New("test2")
	})

	assert(err.Error(), Equals, "test")
}

func TestExecuteParallelContextCancel(t *testing.T) {
	assert := With(t)

	ctx, cancel := context.WithCancel(context.Background())
	err := ExecuteParallel(ctx, func() error {
		time.Sleep(time.Millisecond * 2000)
		return errors.New("test")
	}, func() error {
		time.Sleep(time.Millisecond * 5000)
		return errors.New("test2")
	}, func() error {
		cancel()
		return nil
	})

	assert(err.Error(), HasSubstring, "canceled")
}
