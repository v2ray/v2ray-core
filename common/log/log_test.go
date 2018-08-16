package log_test

import (
	"testing"

	"v2ray.com/core/common/compare"
	"v2ray.com/core/common/log"
	"v2ray.com/core/common/net"
)

type testLogger struct {
	value string
}

func (l *testLogger) Handle(msg log.Message) {
	l.value = msg.String()
}

func TestLogRecord(t *testing.T) {
	var logger testLogger
	log.RegisterHandler(&logger)

	ip := "8.8.8.8"
	log.Record(&log.GeneralMessage{
		Severity: log.Severity_Error,
		Content:  net.ParseAddress(ip),
	})

	if err := compare.StringEqualWithDetail("[Error] "+ip, logger.value); err != nil {
		t.Fatal(err)
	}
}
