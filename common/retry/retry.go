package retry

import (
	"errors"
	"time"
)

var (
	errorRetryFailed = errors.New("All retry attempts failed.")
)

// Strategy is a way to retry on a specific function.
type Strategy interface {
  // On performs a retry on a specific function, until it doesn't return any error.
	On(func() error) error
}

type retryer struct {
	NextDelay func(int) int
}

// On implements Strategy.On.
func (r *retryer) On(method func() error) error {
	attempt := 0
	for {
		err := method()
		if err == nil {
			return nil
		}
		delay := r.NextDelay(attempt)
		if delay < 0 {
			return errorRetryFailed
		}
		<-time.After(time.Duration(delay) * time.Millisecond)
		attempt++
	}
}

// Timed returns a retry strategy with fixed interval.
func Timed(attempts int, delay int) Strategy {
	return &retryer{
		NextDelay: func(attempt int) int {
			if attempt >= attempts {
				return -1
			}
			return delay
		},
	}
}
