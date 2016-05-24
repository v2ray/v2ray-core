package internal

import (
	"fmt"

	"github.com/v2ray/v2ray-core/common"
	"github.com/v2ray/v2ray-core/common/alloc"
)

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
		switch typedVal := value.(type) {
		case string:
			b.AppendString(typedVal)
		case *string:
			b.AppendString(*typedVal)
		case fmt.Stringer:
			b.AppendString(typedVal.String())
		case error:
			b.AppendString(typedVal.Error())
		default:
			b.AppendString(fmt.Sprint(value))
		}
	}
	return b.String()
}

type AccessLog struct {
	From   fmt.Stringer
	To     fmt.Stringer
	Status string
	Reason fmt.Stringer
}

func (this *AccessLog) Release() {
	this.From = nil
	this.To = nil
	this.Reason = nil
}

func (this *AccessLog) String() string {
	b := alloc.NewSmallBuffer().Clear()
	defer b.Release()

	return b.AppendString(this.From.String()).AppendString(" ").AppendString(this.Status).AppendString(" ").AppendString(this.To.String()).AppendString(" ").AppendString(this.Reason.String()).String()
}
