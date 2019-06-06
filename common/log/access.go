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
	Detour interface{}
}

// AccessMessageChan access log channel
var AccessMessageChan = make(chan AccessMessage, 100)

// AccessMessageChanChecker if AccessMessageChan full, cleanup
func AccessMessageChanChecker() {
	if (len(AccessMessageChan) > 50) {
		Loop:
		for {
			select {
			case <-AccessMessageChan:
			default:
				break Loop
			}
		}
		//println("cleanup AccessMessageChan", len(AccessMessageChan))
	}
}

func (m *AccessMessage) String() string {
	builder := strings.Builder{}
	builder.WriteString(serial.ToString(m.From))
	builder.WriteByte(' ')
	builder.WriteString(string(m.Status))
	builder.WriteByte(' ')
	builder.WriteString(serial.ToString(m.To))
	builder.WriteByte(' ')
	builder.WriteString(serial.ToString(m.Detour))
	builder.WriteByte(' ')
	builder.WriteString(serial.ToString(m.Reason))
	return builder.String()
}
