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
}

func (m *AccessMessage) String() string {
	return strings.Join([]string{serial.ToString(m.From), string(m.Status), serial.ToString(m.To), serial.ToString(m.Reason)}, " ")
}
