package retry

import (
	"time"
	"v2ray.com/core/common/errors"
)

var (
	ErrRetryFailed = errors.New("All retry attempts failed.")
)

// Strategy is a way to retry on a specific function.
type Strategy interface {
	// On performs a retry on a specific function, until it doesn't return any error.
	On(func() error) error
}

type retryer struct {
	totalAttempt int
	nextDelay    func() uint32
}

// On implements Strategy.On.
func (r *retryer) On(method func() error) error {
	attempt := 0
	for attempt < r.totalAttempt {
		err := method()
		if err == nil {
			return nil
		}
		delay := r.nextDelay()
		<-time.After(time.Duration(delay) * time.Millisecond)
		attempt++
	}
	return ErrRetryFailed
}

// Timed returns a retry strategy with fixed interval.
func Timed(attempts int, delay uint32) Strategy {
	return &retryer{
		totalAttempt: attempts,
		nextDelay: func() uint32 {
			return delay
		},
	}
}

func ExponentialBackoff(attempts int, delay uint32) Strategy {
	nextDelay := uint32(0)
	return &retryer{
		totalAttempt: attempts,
		nextDelay: func() uint32 {
			r := nextDelay
			nextDelay += delay
			return r
		},
	}
}
