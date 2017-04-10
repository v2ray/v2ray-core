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

func (v *ErrorLog) String() string {
	return v.Prefix + v.Error.Error()
}

type AccessLog struct {
	From   interface{}
	To     interface{}
	Status string
	Reason interface{}
}

func (v *AccessLog) String() string {
	return strings.Join([]string{serial.ToString(v.From), v.Status, serial.ToString(v.To), serial.ToString(v.Reason)}, " ")
}
