package functions

import "v2ray.com/core/common"

// Task is a function that may return an error.
type Task func() error

// OnSuccess returns a Task to run a follow task if pre-condition passes, otherwise the error in pre-condition is returned.
func OnSuccess(pre func() error, followup Task) Task {
	return func() error {
		if err := pre(); err != nil {
			return err
		}
		return followup()
	}
}

// Close returns a Task to close the object.
func Close(obj interface{}) Task {
	return func() error {
		return common.Close(obj)
	}
}
