// Package errors is a drop-in replacement for Golang lib 'errors'.
package errors // import "v2ray.com/core/common/errors"

import (
	"os"
	"reflect"
	"strings"

	"v2ray.com/core/common/log"
	"v2ray.com/core/common/serial"
)

type hasInnerError interface {
	// Inner returns the underlying error of this one.
	Inner() error
}

type hasSeverity interface {
	Severity() log.Severity
}

// Error is an error object with underlying error.
type Error struct {
	pathObj  interface{}
	prefix   []interface{}
	message  []interface{}
	inner    error
	severity log.Severity
}

func (err *Error) WithPathObj(obj interface{}) *Error {
	err.pathObj = obj
	return err
}

func (err *Error) pkgPath() string {
	if err.pathObj == nil {
		return ""
	}
	return reflect.TypeOf(err.pathObj).PkgPath()
}

// Error implements error.Error().
func (v *Error) Error() string {
	builder := strings.Builder{}
	for _, prefix := range v.prefix {
		builder.WriteByte('[')
		builder.WriteString(serial.ToString(prefix))
		builder.WriteString("] ")
	}

	path := v.pkgPath()
	if len(path) > 0 {
		builder.WriteString(path)
		builder.WriteString(": ")
	}

	msg := serial.Concat(v.message...)
	builder.WriteString(msg)

	if v.inner != nil {
		builder.WriteString(" > ")
		builder.WriteString(v.inner.Error())
	}

	return builder.String()
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

func (v *Error) atSeverity(s log.Severity) *Error {
	v.severity = s
	return v
}

func (v *Error) Severity() log.Severity {
	if v.inner == nil {
		return v.severity
	}

	if s, ok := v.inner.(hasSeverity); ok {
		as := s.Severity()
		if as < v.severity {
			return as
		}
	}

	return v.severity
}

// AtDebug sets the severity to debug.
func (v *Error) AtDebug() *Error {
	return v.atSeverity(log.Severity_Debug)
}

// AtInfo sets the severity to info.
func (v *Error) AtInfo() *Error {
	return v.atSeverity(log.Severity_Info)
}

// AtWarning sets the severity to warning.
func (v *Error) AtWarning() *Error {
	return v.atSeverity(log.Severity_Warning)
}

// AtError sets the severity to error.
func (v *Error) AtError() *Error {
	return v.atSeverity(log.Severity_Error)
}

// String returns the string representation of this error.
func (v *Error) String() string {
	return v.Error()
}

// WriteToLog writes current error into log.
func (v *Error) WriteToLog(opts ...ExportOption) {
	var holder ExportOptionHolder

	for _, opt := range opts {
		opt(&holder)
	}

	if holder.SessionID > 0 {
		v.prefix = append(v.prefix, holder.SessionID)
	}

	log.Record(&log.GeneralMessage{
		Severity: GetSeverity(v),
		Content:  v,
	})
}

type ExportOptionHolder struct {
	SessionID uint32
}

type ExportOption func(*ExportOptionHolder)

// New returns a new error object with message formed from given arguments.
func New(msg ...interface{}) *Error {
	return &Error{
		message:  msg,
		severity: log.Severity_Info,
	}
}

// Cause returns the root cause of this error.
func Cause(err error) error {
	if err == nil {
		return nil
	}
L:
	for {
		switch inner := err.(type) {
		case hasInnerError:
			if inner.Inner() == nil {
				break L
			}
			err = inner.Inner()
		case *os.PathError:
			if inner.Err == nil {
				break L
			}
			err = inner.Err
		case *os.SyscallError:
			if inner.Err == nil {
				break L
			}
			err = inner.Err
		default:
			break L
		}
	}
	return err
}

// GetSeverity returns the actual severity of the error, including inner errors.
func GetSeverity(err error) log.Severity {
	if s, ok := err.(hasSeverity); ok {
		return s.Severity()
	}
	return log.Severity_Info
}
