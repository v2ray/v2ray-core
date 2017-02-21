// Package errors is a drop-in replacement for Golang lib 'errors'.
package errors

import (
	"fmt"

	"v2ray.com/core/common/serial"
)

type hasInnerError interface {
	// Inner returns the underlying error of this one.
	Inner() error
}

type actionRequired interface {
	ActionRequired() bool
}

// Error is an error object with underlying error.
type Error struct {
	message        string
	inner          error
	actionRequired bool
}

// Error implements error.Error().
func (v *Error) Error() string {
	return v.message
}

// Inner implements hasInnerError.Inner()
func (v *Error) Inner() error {
	if v.inner == nil {
		return nil
	}
	return v.inner
}

func (v *Error) ActionRequired() bool {
	return v.actionRequired
}

// New returns a new error object with message formed from given arguments.
func New(msg ...interface{}) error {
	return &Error{
		message: serial.Concat(msg...),
	}
}

// Base returns an ErrorBuilder based on the given error.
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
		inner, ok := err.(hasInnerError)
		if !ok || inner.Inner() == nil {
			break
		}
		err = inner.Inner()
	}
	return err
}

func IsActionRequired(err error) bool {
	for err != nil {
		if ar, ok := err.(actionRequired); ok && ar.ActionRequired() {
			return true
		}
		inner, ok := err.(hasInnerError)
		if !ok || inner.Inner() == nil {
			break
		}
		err = inner.Inner()
	}
	return false
}

type ErrorBuilder struct {
	error
	actionRequired bool
}

func (v ErrorBuilder) RequireUserAction() ErrorBuilder {
	v.actionRequired = true
	return v
}

// Message returns an error object with given message and base error.
func (v ErrorBuilder) Message(msg ...interface{}) error {
	if v.error == nil {
		return nil
	}

	return &Error{
		message:        serial.Concat(msg...) + " > " + v.error.Error(),
		inner:          v.error,
		actionRequired: v.actionRequired,
	}
}

// Format returns an errors object with given message format and base error.
func (v ErrorBuilder) Format(format string, values ...interface{}) error {
	if v.error == nil {
		return nil
	}
	return &Error{
		message:        fmt.Sprintf(format, values...) + " > " + v.error.Error(),
		inner:          v.error,
		actionRequired: v.actionRequired,
	}
}
