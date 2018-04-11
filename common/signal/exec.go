package signal

import (
	"context"
)

// Execute runs a list of tasks sequentially, returns the first error encountered or nil if all tasks pass.
func Execute(tasks ...func() error) error {
	for _, task := range tasks {
		if err := task(); err != nil {
			return err
		}
	}
	return nil
}

// ExecuteParallel executes a list of tasks asynchronously, returns the first error encountered or nil if all tasks pass.
func ExecuteParallel(ctx context.Context, tasks ...func() error) error {
	n := len(tasks)
	s := NewSemaphore(n)
	done := make(chan error, 1)

	for _, task := range tasks {
		<-s.Wait()
		go func(f func() error) {
			if err := f(); err != nil {
				select {
				case done <- err:
				default:
				}
			}
			s.Signal()
		}(task)
	}

	for i := 0; i < n; i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-done:
			return err
		case <-s.Wait():
		}
	}

	return nil
}
