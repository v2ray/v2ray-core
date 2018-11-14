package task_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"v2ray.com/core/common"
	. "v2ray.com/core/common/task"
	. "v2ray.com/ext/assert"
)

func TestExecuteParallel(t *testing.T) {
	assert := With(t)

	err := Run(Parallel(func() error {
		time.Sleep(time.Millisecond * 200)
		return errors.New("test")
	}, func() error {
		time.Sleep(time.Millisecond * 500)
		return errors.New("test2")
	}))()

	assert(err.Error(), Equals, "test")
}

func TestExecuteParallelContextCancel(t *testing.T) {
	assert := With(t)

	ctx, cancel := context.WithCancel(context.Background())
	err := Run(WithContext(ctx), Parallel(func() error {
		time.Sleep(time.Millisecond * 2000)
		return errors.New("test")
	}, func() error {
		time.Sleep(time.Millisecond * 5000)
		return errors.New("test2")
	}, func() error {
		cancel()
		return nil
	}))()

	assert(err.Error(), HasSubstring, "canceled")
}

func BenchmarkExecuteOne(b *testing.B) {
	noop := func() error {
		return nil
	}
	for i := 0; i < b.N; i++ {
		common.Must(Run(Parallel(noop))())
	}
}

func BenchmarkExecuteTwo(b *testing.B) {
	noop := func() error {
		return nil
	}
	for i := 0; i < b.N; i++ {
		common.Must(Run(Parallel(noop, noop))())
	}
}

func BenchmarkExecuteContext(b *testing.B) {
	noop := func() error {
		return nil
	}
	background := context.Background()

	for i := 0; i < b.N; i++ {
		common.Must(Run(WithContext(background), Parallel(noop, noop))())
	}
}
