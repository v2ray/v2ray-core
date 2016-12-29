package signal

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

func ErrorOrFinish2(c1, c2 <-chan error) error {
	select {
	case err, failed := <-c1:
		if failed {
			return err
		}
		return <-c2
	case err, failed := <-c2:
		if failed {
			return err
		}
		return <-c1
	}
}
