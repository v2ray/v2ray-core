package log

import (
	"strings"

	"v2ray.com/core/common/serial"
)

type AccessStatus string

const (
	AccessAccepted = AccessStatus("accepted")
	AccessRejected = AccessStatus("rejected")
)

type AccessMessage struct {
	From   interface{}
	To     interface{}
	Status AccessStatus
	Reason interface{}
	Email  string
}

func (m *AccessMessage) String() string {
	builder := strings.Builder{}
	builder.WriteString(serial.ToString(m.From))
	builder.WriteByte(' ')
	builder.WriteString(string(m.Status))
	builder.WriteByte(' ')
	builder.WriteString(serial.ToString(m.To))
	builder.WriteByte(' ')
	builder.WriteString(serial.ToString(m.Reason))

	if len(m.Email) > 0 {
		builder.WriteString("email:")
		builder.WriteString(m.Email)
		builder.WriteByte(' ')
	}
	return builder.String()
}
