package task

import (
	"sync"
)

type Task interface {
	Execute() error
}

type ParallelExecutor struct {
	sync.Mutex
	tasks  sync.WaitGroup
	errors []error
}

func (pe *ParallelExecutor) track(err error) {
	if err == nil {
		return
	}

	pe.Lock()
	pe.errors = append(pe.errors, err)
	pe.Unlock()
}

func (pe *ParallelExecutor) Execute(task Task) {
	pe.tasks.Add(1)
	go func() {
		pe.track(task.Execute())
		pe.tasks.Done()
	}()
}

func (pe *ParallelExecutor) Wait() {
	pe.tasks.Wait()
}

func (pe *ParallelExecutor) Errors() []error {
	return pe.errors
}
