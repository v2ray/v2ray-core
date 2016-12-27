package internal

import (
	"fmt"
	"strings"

	"v2ray.com/core/common"
	"v2ray.com/core/common/serial"
)

type LogEntry interface {
	common.Releasable
	fmt.Stringer
}

type ErrorLog struct {
	Prefix string
	Values []interface{}
}

func (v *ErrorLog) Release() {
	for _, val := range v.Values {
		common.Release(val)
	}
}

func (v *ErrorLog) String() string {
	return v.Prefix + serial.Concat(v.Values...)
}

type AccessLog struct {
	From   interface{}
	To     interface{}
	Status string
	Reason interface{}
}

func (v *AccessLog) Release() {
	common.Release(v.From)
	common.Release(v.To)
	common.Release(v.Reason)
}

func (v *AccessLog) String() string {
	return strings.Join([]string{serial.ToString(v.From), v.Status, serial.ToString(v.To), serial.ToString(v.Reason)}, " ")
}
