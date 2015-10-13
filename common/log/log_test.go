package log

import (
	"testing"

	"github.com/v2ray/v2ray-core/testing/unit"
)

func TestLogLevelSetting(t *testing.T) {
	assert := unit.Assert(t)

	assert.Pointer(debugLogger).Equals(noOpLoggerInstance)
	SetLogLevel(DebugLevel)
	assert.Pointer(debugLogger).Equals(streamLoggerInstance)

	SetLogLevel(InfoLevel)
	assert.Pointer(debugLogger).Equals(noOpLoggerInstance)
	assert.Pointer(infoLogger).Equals(streamLoggerInstance)
}
