package errors

import (
	"fmt"

	"v2ray.com/core/common/serial"
)

type HasInnerError interface {
	Inner() error
}

type Error struct {
	message string
	inner   error
}

func (v *Error) Error() string {
	return v.message
}

func (v *Error) Inner() error {
	if v.inner == nil {
		return nil
	}
	return v.inner
}

func New(msg ...interface{}) error {
	return &Error{
		message: serial.Concat(msg),
	}
}

func Base(err error) ErrorBuilder {
	return ErrorBuilder{
		error: err,
	}
}

func Format(format string, values ...interface{}) error {
	return New(fmt.Sprintf(format, values...))
}

// Cause returns the root cause of this error.
func Cause(err error) error {
	if err == nil {
		return nil
	}
	for {
		inner, ok := err.(HasInnerError)
		if !ok {
			break
		}
		if inner.Inner() == nil {
			break
		}
		err = inner.Inner()
	}
	return err
}

type ErrorBuilder struct {
	error
}

func (v ErrorBuilder) Message(msg ...interface{}) error {
	if v.error == nil {
		return nil
	}

	return &Error{
		message: serial.ToString(msg) + " > " + v.error.Error(),
		inner:   v.error,
	}
}

func (v ErrorBuilder) Format(format string, values ...interface{}) error {
	if v.error == nil {
		return nil
	}
	return &Error{
		message: fmt.Sprintf(format, values...) + " > " + v.error.Error(),
		inner:   v.error,
	}
}
