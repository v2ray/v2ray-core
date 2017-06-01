package internal

import (
	"fmt"
	"strings"

	"v2ray.com/core/common/serial"
)

type LogEntry interface {
	fmt.Stringer
}

type ErrorLog struct {
	Prefix string
	Error  error
}

func (l *ErrorLog) String() string {
	return l.Prefix + l.Error.Error()
}

type AccessLog struct {
	From   interface{}
	To     interface{}
	Status string
	Reason interface{}
}

func (l *AccessLog) String() string {
	return strings.Join([]string{serial.ToString(l.From), l.Status, serial.ToString(l.To), serial.ToString(l.Reason)}, " ")
}
