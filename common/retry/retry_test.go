package retry_test

import (
	"testing"
	"time"

	"v2ray.com/core/common/errors"
	. "v2ray.com/core/common/retry"
	"v2ray.com/core/testing/assert"
)

var (
	errorTestOnly = errors.New("This is a fake error.")
)

func TestNoRetry(t *testing.T) {
	assert := assert.On(t)

	startTime := time.Now().Unix()
	err := Timed(10, 100000).On(func() error {
		return nil
	})
	endTime := time.Now().Unix()

	assert.Error(err).IsNil()
	assert.Int64(endTime - startTime).AtLeast(0)
}

func TestRetryOnce(t *testing.T) {
	assert := assert.On(t)

	startTime := time.Now()
	called := 0
	err := Timed(10, 1000).On(func() error {
		if called == 0 {
			called++
			return errorTestOnly
		}
		return nil
	})
	duration := time.Since(startTime)

	assert.Error(err).IsNil()
	assert.Int64(int64(duration / time.Millisecond)).AtLeast(900)
}

func TestRetryMultiple(t *testing.T) {
	assert := assert.On(t)

	startTime := time.Now()
	called := 0
	err := Timed(10, 1000).On(func() error {
		if called < 5 {
			called++
			return errorTestOnly
		}
		return nil
	})
	duration := time.Since(startTime)

	assert.Error(err).IsNil()
	assert.Int64(int64(duration / time.Millisecond)).AtLeast(4900)
}

func TestRetryExhausted(t *testing.T) {
	assert := assert.On(t)

	startTime := time.Now()
	called := 0
	err := Timed(2, 1000).On(func() error {
		called++
		return errorTestOnly
	})
	duration := time.Since(startTime)

	assert.Error(err).Equals(ErrRetryFailed)
	assert.Int64(int64(duration / time.Millisecond)).AtLeast(1900)
}

func TestExponentialBackoff(t *testing.T) {
	assert := assert.On(t)

	startTime := time.Now()
	called := 0
	err := ExponentialBackoff(10, 100).On(func() error {
		called++
		return errorTestOnly
	})
	duration := time.Since(startTime)

	assert.Error(err).Equals(ErrRetryFailed)
	assert.Int64(int64(duration / time.Millisecond)).AtLeast(4000)
}
