package retry

import (
	"errors"
	"testing"
	"time"

	"github.com/v2ray/v2ray-core/testing/unit"
)

var (
	TestError = errors.New("This is a fake error.")
)

func TestNoRetry(t *testing.T) {
	assert := unit.Assert(t)

	startTime := time.Now().Unix()
	Timed(10, 100000).On(func() error {
		return nil
	})
	endTime := time.Now().Unix()

	assert.Int64(endTime - startTime).AtLeast(0)
}

func TestRetryOnce(t *testing.T) {
	assert := unit.Assert(t)

	startTime := time.Now()
	called := 0
	Timed(10, 1000).On(func() error {
		if called == 0 {
			called++
			return TestError
		}
		return nil
	})
	duration := time.Since(startTime)

	assert.Int64(int64(duration / time.Millisecond)).AtLeast(900)
}

func TestRetryMultiple(t *testing.T) {
	assert := unit.Assert(t)

	startTime := time.Now()
	called := 0
	Timed(10, 1000).On(func() error {
		if called < 5 {
			called++
			return TestError
		}
		return nil
	})
	duration := time.Since(startTime)

	assert.Int64(int64(duration / time.Millisecond)).AtLeast(4900)
}
