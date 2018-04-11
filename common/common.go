// Package common contains common utilities that are shared among other packages.
// See each sub-package for detail.
package common

//go:generate go run $GOPATH/src/v2ray.com/core/common/errors/errorgen/main.go -pkg common -path Common

// Must panics if err is not nil.
func Must(err error) {
	if err != nil {
		panic(err)
	}
}

// Must2 panics if the second parameter is not nil, otherwise returns the first parameter.
func Must2(v interface{}, err error) interface{} {
	Must(err)
	return v
}

// Error2 returns the err from the 2nd parameter.
func Error2(v interface{}, err error) error {
	return err
}
