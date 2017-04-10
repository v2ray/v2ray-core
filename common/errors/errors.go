// Package errors is a drop-in replacement for Golang lib 'errors'.
package errors

import (
	"strings"

	"v2ray.com/core/common/serial"
)

// Severity describes how severe the error is.
type Severity int

const (
	SeverityDebug Severity = iota
	SeverityInfo
	SeverityWarning
	SeverityError
)

type hasInnerError interface {
	// Inner returns the underlying error of this one.
	Inner() error
}

type hasSeverity interface {
	Severity() Severity
}

// Error is an error object with underlying error.
type Error struct {
	message  []interface{}
	inner    error
	severity Severity
	path     []string
}

// Error implements error.Error().
func (v *Error) Error() string {
	msg := serial.Concat(v.message...)
	if v.inner != nil {
		msg += " > " + v.inner.Error()
	}
	if len(v.path) > 0 {
		msg = strings.Join(v.path, "|") + ": " + msg
	}
	return msg
}

// Inner implements hasInnerError.Inner()
func (v *Error) Inner() error {
	if v.inner == nil {
		return nil
	}
	return v.inner
}

func (v *Error) Base(err error) *Error {
	v.inner = err
	return v
}

func (v *Error) atSeverity(s Severity) *Error {
	v.severity = s
	return v
}

func (v *Error) Severity() Severity {
	if v.inner == nil {
		return v.severity
	}

	if s, ok := v.inner.(hasSeverity); ok {
		as := s.Severity()
		if as > v.severity {
			return as
		}
	}

	return v.severity
}

// AtDebug sets the severity to debug.
func (v *Error) AtDebug() *Error {
	return v.atSeverity(SeverityDebug)
}

// AtInfo sets the severity to info.
func (v *Error) AtInfo() *Error {
	return v.atSeverity(SeverityInfo)
}

// AtWarning sets the severity to warning.
func (v *Error) AtWarning() *Error {
	return v.atSeverity(SeverityWarning)
}

// AtError sets the severity to error.
func (v *Error) AtError() *Error {
	return v.atSeverity(SeverityError)
}

// Path sets the path to the location where this error happens.
func (v *Error) Path(path ...string) *Error {
	v.path = path
	return v
}

// New returns a new error object with message formed from given arguments.
func New(msg ...interface{}) *Error {
	return &Error{
		message:  msg,
		severity: SeverityInfo,
	}
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

func GetSeverity(err error) Severity {
	if s, ok := err.(hasSeverity); ok {
		return s.Severity()
	}
	return SeverityInfo
}
