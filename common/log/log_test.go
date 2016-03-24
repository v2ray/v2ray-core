package log

import (
	"bytes"
	"log"
	"testing"

	"github.com/v2ray/v2ray-core/common/serial"
    "github.com/v2ray/v2ray-core/common/platform"
	v2testing "github.com/v2ray/v2ray-core/testing"
	"github.com/v2ray/v2ray-core/testing/assert"
)

func TestLogLevelSetting(t *testing.T) {
	v2testing.Current(t)

	assert.Pointer(debugLogger).Equals(noOpLoggerInstance)
	SetLogLevel(DebugLevel)
	assert.Pointer(debugLogger).Equals(streamLoggerInstance)

	SetLogLevel(InfoLevel)
	assert.Pointer(debugLogger).Equals(noOpLoggerInstance)
	assert.Pointer(infoLogger).Equals(streamLoggerInstance)
}

func TestStreamLogger(t *testing.T) {
	v2testing.Current(t)

	buffer := bytes.NewBuffer(make([]byte, 0, 1024))
	infoLogger = &stdOutLogWriter{
		logger: log.New(buffer, "", 0),
	}
	Info("Test ", "Stream Logger", " Format")
	assert.StringLiteral(string(buffer.Bytes())).Equals("[Info]Test Stream Logger Format" + platform.LineSeparator())

	buffer.Reset()
	errorLogger = infoLogger
	Error("Test ", serial.StringLiteral("literal"), " Format")
	assert.StringLiteral(string(buffer.Bytes())).Equals("[Error]Test literal Format" + platform.LineSeparator())
}
