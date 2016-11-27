package internal

import (
	"fmt"
	"strings"

	"v2ray.com/core/common"
	"v2ray.com/core/common/serial"
)

func InterfaceToString(value interface{}) string {
	if value == nil {
		return " "
	}
	switch value := value.(type) {
	case string:
		return value
	case *string:
		return *value
	case fmt.Stringer:
		return value.String()
	case error:
		return value.Error()
	case []byte:
		return serial.BytesToHexString(value)
	default:
		return fmt.Sprintf("%+v", value)
	}
}

type LogEntry interface {
	common.Releasable
	fmt.Stringer
}

type ErrorLog struct {
	Prefix string
	Values []interface{}
}

func (v *ErrorLog) Release() {
	for index := range v.Values {
		v.Values[index] = nil
	}
	v.Values = nil
}

func (v *ErrorLog) String() string {
	values := make([]string, len(v.Values)+1)
	values[0] = v.Prefix
	for i, value := range v.Values {
		values[i+1] = InterfaceToString(value)
	}
	return strings.Join(values, "")
}

type AccessLog struct {
	From   interface{}
	To     interface{}
	Status string
	Reason interface{}
}

func (v *AccessLog) Release() {
	v.From = nil
	v.To = nil
	v.Reason = nil
}

func (v *AccessLog) String() string {
	return strings.Join([]string{InterfaceToString(v.From), v.Status, InterfaceToString(v.To), InterfaceToString(v.Reason)}, " ")
}
