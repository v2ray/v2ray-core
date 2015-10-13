package retry

import (
	"errors"
	"time"
)

var (
	RetryFailed = errors.New("All retry attempts failed.")
)

type RetryStrategy interface {
	On(func() error) error
}

type retryer struct {
	NextDelay func(int) int
}

func (r *retryer) On(method func() error) error {
	attempt := 0
	for {
		err := method()
		if err == nil {
			return nil
		}
		delay := r.NextDelay(attempt)
		if delay < 0 {
			return RetryFailed
		}
		<-time.After(time.Duration(delay) * time.Millisecond)
	}
}

func Timed(attempts int, delay int) RetryStrategy {
	return &retryer{
		NextDelay: func(attempt int) int {
			if attempt >= attempts {
				return -1
			}
			return delay
		},
	}
}
