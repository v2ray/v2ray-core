package internal

import (
	"fmt"

	"github.com/v2ray/v2ray-core/common"
	"github.com/v2ray/v2ray-core/common/alloc"
	"github.com/v2ray/v2ray-core/common/serial"
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
		return fmt.Sprint(value)
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
	b := alloc.NewSmallBuffer().Clear()
	defer b.Release()

	b.AppendString(this.Prefix)

	for _, value := range this.Values {
		b.AppendString(InterfaceToString(value))
	}
	return b.String()
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
	b := alloc.NewSmallBuffer().Clear()
	defer b.Release()

	b.AppendString(InterfaceToString(this.From)).AppendString(" ")
	b.AppendString(this.Status).AppendString(" ")
	b.AppendString(InterfaceToString(this.To)).AppendString(" ")
	b.AppendString(InterfaceToString(this.Reason))
	return b.String()
}
