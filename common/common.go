// Package common contains common utilities that are shared among other packages.
// See each sub-package for detail.
package common

// Must panics if err is not nil.
func Must(err error) {
	if err != nil {
		panic(err)
	}
}
