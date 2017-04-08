// Package common contains common utilities that are shared among other packages.
// See each sub-package for detail.
package common

//go:generate go run $GOPATH/src/v2ray.com/core/tools/generrorgen/main.go -pkg common -path Common

// Must panics if err is not nil.
func Must(err error) {
	if err != nil {
		panic(err)
	}
}
