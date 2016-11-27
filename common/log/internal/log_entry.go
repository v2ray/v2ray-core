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

func (this *ErrorLog) Release() {
	for index := range this.Values {
		this.Values[index] = nil
	}
	this.Values = nil
}

func (this *ErrorLog) String() string {
	values := make([]string, len(this.Values)+1)
	values[0] = this.Prefix
	for i, value := range this.Values {
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

func (this *AccessLog) Release() {
	this.From = nil
	this.To = nil
	this.Reason = nil
}

func (this *AccessLog) String() string {
	return strings.Join([]string{InterfaceToString(this.From), this.Status, InterfaceToString(this.To), InterfaceToString(this.Reason)}, " ")
}
