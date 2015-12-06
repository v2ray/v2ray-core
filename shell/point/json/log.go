package json

import (
	"strings"

	"github.com/v2ray/v2ray-core/common/log"
)

type LogConfig struct {
	AccessLogValue string `json:"access"`
	ErrorLogValue  string `json:"error"`
	LogLevelValue  string `json:"loglevel"`
}

func (this *LogConfig) AccessLog() string {
	return this.AccessLogValue
}

func (this *LogConfig) ErrorLog() string {
	return this.ErrorLogValue
}

func (this *LogConfig) LogLevel() log.LogLevel {
	level := strings.ToLower(this.LogLevelValue)
	switch level {
	case "debug":
		return log.DebugLevel
	case "info":
		return log.InfoLevel
	case "error":
		return log.ErrorLevel
	default:
		return log.WarningLevel
	}
}
