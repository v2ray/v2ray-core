// Package common contains common utilities that are shared among other packages.
// See each sub-package for detail.
package common

import (
	"errors"
)

var (
	ErrorAlreadyReleased = errors.New("Object already released.")
	ErrBadConfiguration  = errors.New("Bad configuration.")
)

// Releasable interface is for those types that can release its members.
type Releasable interface {
	// Release releases all references to accelerate garbage collection.
	Release()
}
