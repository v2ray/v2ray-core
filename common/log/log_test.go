package log_test

import (
	"testing"

	"v2ray.com/core/common/log"
	"v2ray.com/core/common/net"
	. "v2ray.com/ext/assert"
)

type testLogger struct {
	value string
}

func (l *testLogger) Handle(msg log.Message) {
	l.value = msg.String()
}

func TestLogRecord(t *testing.T) {
	assert := With(t)

	var logger testLogger
	log.RegisterHandler(&logger)

	ip := "8.8.8.8"
	log.Record(&log.GeneralMessage{
		Severity: log.Severity_Error,
		Content:  net.ParseAddress(ip),
	})

	assert(logger.value, Equals, "[Error]: "+ip)
}
