package signal

import (
	"context"
)

func executeAndFulfill(f func() error, done chan<- error) {
	err := f()
	if err != nil {
		done <- err
	}
	close(done)
}

func ExecuteAsync(f func() error) <-chan error {
	done := make(chan error, 1)
	go executeAndFulfill(f, done)
	return done
}

func ErrorOrFinish1(ctx context.Context, c <-chan error) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-c:
		return err
	}
}

func ErrorOrFinish2(ctx context.Context, c1, c2 <-chan error) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case err, failed := <-c1:
		if failed {
			return err
		}
		return ErrorOrFinish1(ctx, c2)
	case err, failed := <-c2:
		if failed {
			return err
		}
		return ErrorOrFinish1(ctx, c1)
	}
}
