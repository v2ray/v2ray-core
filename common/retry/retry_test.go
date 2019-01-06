package retry_test

import (
	"testing"
	"time"

	"v2ray.com/core/common/errors"
	. "v2ray.com/core/common/retry"
	. "v2ray.com/ext/assert"
)

var (
	errorTestOnly = errors.New("This is a fake error.")
)

func TestNoRetry(t *testing.T) {
	assert := With(t)

	startTime := time.Now().Unix()
	err := Timed(10, 100000).On(func() error {
		return nil
	})
	endTime := time.Now().Unix()

	assert(err, IsNil)
	assert(endTime-startTime, AtLeast, int64(0))
}

func TestRetryOnce(t *testing.T) {
	assert := With(t)

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

	assert(err, IsNil)
	assert(int64(duration/time.Millisecond), AtLeast, int64(900))
}

func TestRetryMultiple(t *testing.T) {
	assert := With(t)

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

	assert(err, IsNil)
	assert(int64(duration/time.Millisecond), AtLeast, int64(4900))
}

func TestRetryExhausted(t *testing.T) {
	assert := With(t)

	startTime := time.Now()
	called := 0
	err := Timed(2, 1000).On(func() error {
		called++
		return errorTestOnly
	})
	duration := time.Since(startTime)

	assert(errors.Cause(err), Equals, ErrRetryFailed)
	assert(int64(duration/time.Millisecond), AtLeast, int64(1900))
}

func TestExponentialBackoff(t *testing.T) {
	assert := With(t)

	startTime := time.Now()
	called := 0
	err := ExponentialBackoff(10, 100).On(func() error {
		called++
		return errorTestOnly
	})
	duration := time.Since(startTime)

	assert(errors.Cause(err), Equals, ErrRetryFailed)
	assert(int64(duration/time.Millisecond), AtLeast, int64(4000))
}
