// Package common contains common utilities that are shared among other packages.
// See each sub-package for detail.
package common

import (
	"errors"
)

var (
	ErrorAlreadyReleased = errors.New("Object already released.")
)

type Releasable interface {
	Release()
}
