package task_test

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	"v2ray.com/core/common"
	. "v2ray.com/core/common/task"
)

func TestExecuteParallel(t *testing.T) {
	err := Run(context.Background(),
		func() error {
			time.Sleep(time.Millisecond * 200)
			return errors.New("test")
		}, func() error {
			time.Sleep(time.Millisecond * 500)
			return errors.New("test2")
		})

	if r := cmp.Diff(err.Error(), "test"); r != "" {
		t.Error(r)
	}
}

func TestExecuteParallelContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	err := Run(ctx, func() error {
		time.Sleep(time.Millisecond * 2000)
		return errors.New("test")
	}, func() error {
		time.Sleep(time.Millisecond * 5000)
		return errors.New("test2")
	}, func() error {
		cancel()
		return nil
	})

	errStr := err.Error()
	if !strings.Contains(errStr, "canceled") {
		t.Error("expected error string to contain 'canceled', but actually not: ", errStr)
	}
}

func BenchmarkExecuteOne(b *testing.B) {
	noop := func() error {
		return nil
	}
	for i := 0; i < b.N; i++ {
		common.Must(Run(context.Background(), noop))
	}
}

func BenchmarkExecuteTwo(b *testing.B) {
	noop := func() error {
		return nil
	}
	for i := 0; i < b.N; i++ {
		common.Must(Run(context.Background(), noop, noop))
	}
}
