package retry_test

import (
	"testing"
	"time"

	"v2ray.com/core/common"
	"v2ray.com/core/common/errors"
	. "v2ray.com/core/common/retry"
)

var (
	errorTestOnly = errors.New("This is a fake error.")
)

func TestNoRetry(t *testing.T) {
	startTime := time.Now().Unix()
	err := Timed(10, 100000).On(func() error {
		return nil
	})
	endTime := time.Now().Unix()

	common.Must(err)
	if endTime < startTime {
		t.Error("endTime < startTime: ", startTime, " -> ", endTime)
	}
}

func TestRetryOnce(t *testing.T) {
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

	common.Must(err)
	if v := int64(duration / time.Millisecond); v < 900 {
		t.Error("duration: ", v)
	}
}

func TestRetryMultiple(t *testing.T) {
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

	common.Must(err)
	if v := int64(duration / time.Millisecond); v < 4900 {
		t.Error("duration: ", v)
	}
}

func TestRetryExhausted(t *testing.T) {
	startTime := time.Now()
	called := 0
	err := Timed(2, 1000).On(func() error {
		called++
		return errorTestOnly
	})
	duration := time.Since(startTime)

	if errors.Cause(err) != ErrRetryFailed {
		t.Error("cause: ", err)
	}

	if v := int64(duration / time.Millisecond); v < 1900 {
		t.Error("duration: ", v)
	}
}

func TestExponentialBackoff(t *testing.T) {
	startTime := time.Now()
	called := 0
	err := ExponentialBackoff(10, 100).On(func() error {
		called++
		return errorTestOnly
	})
	duration := time.Since(startTime)

	if errors.Cause(err) != ErrRetryFailed {
		t.Error("cause: ", err)
	}
	if v := int64(duration / time.Millisecond); v < 4000 {
		t.Error("duration: ", v)
	}
}
