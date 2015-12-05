package log

import (
	"bytes"
	"log"
	"testing"

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
	Info("Test %s Format", "Stream Logger")
	assert.Bytes(buffer.Bytes()).Equals([]byte("[Info]Test Stream Logger Format\n"))

	buffer.Reset()
	errorLogger = infoLogger
	Error("Test No Format")
	assert.Bytes(buffer.Bytes()).Equals([]byte("[Error]Test No Format\n"))
}
