// Package errors is a drop-in replacement for Golang lib 'errors'.
package errors // import "v2ray.com/core/common/errors"

import (
	"context"
	"os"
	"strings"

	"v2ray.com/core/common/log"
	"v2ray.com/core/common/serial"
	"v2ray.com/core/common/session"
)

type hasInnerError interface {
	// Inner returns the underlying error of this one.
	Inner() error
}

type hasSeverity interface {
	Severity() log.Severity
}

type hasContext interface {
	Context() context.Context
}

// Error is an error object with underlying error.
type Error struct {
	message  []interface{}
	inner    error
	severity log.Severity
	path     []string
	ctx      context.Context
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

func (v *Error) WithContext(ctx context.Context) *Error {
	v.ctx = ctx
	return v
}

// Context returns the context that associated with the Error.
func (v *Error) Context() context.Context {
	if v.ctx != nil {
		return v.ctx
	}

	if v.inner == nil {
		return nil
	}

	if c, ok := v.inner.(hasContext); ok {
		return c.Context()
	}

	return nil
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

// Path sets the path to the location where this error happens.
func (v *Error) Path(path ...string) *Error {
	v.path = path
	return v
}

// String returns the string representation of this error.
func (v *Error) String() string {
	return v.Error()
}

// WriteToLog writes current error into log.
func (v *Error) WriteToLog() {
	ctx := v.Context()
	var sid session.ID
	if ctx != nil {
		sid = session.IDFromContext(ctx)
	}
	var c interface{} = v
	if sid > 0 {
		c = sessionLog{
			id:      sid,
			content: v,
		}
	}
	log.Record(&log.GeneralMessage{
		Severity: GetSeverity(v),
		Content:  c,
	})
}

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

type sessionLog struct {
	id      session.ID
	content interface{}
}

func (s sessionLog) String() string {
	return serial.Concat("[", s.id, "] ", s.content)
}
